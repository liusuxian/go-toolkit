/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-20 15:38:07
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-23 18:07:43
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gftoken

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/crypto/gaes"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/liusuxian/go-toolkit/gf/gflogger"
	"github.com/liusuxian/go-toolkit/gf/gfresp"
	"net/http"
	"strings"
)

// Token token 结构体
type Token struct {
	// 服务名称
	ServerName string
	// 缓存模式 1:gcache 2:gredis 3:fileCache 默认1
	CacheMode int8
	// gredis 组名称
	RedisGroupName string
	// 缓存 key
	CacheKey string
	// 超时时间 默认10天（毫秒）
	Timeout int
	// 缓存刷新时间 默认为超时时间的一半（毫秒）
	MaxRefresh int
	// Token 分隔符
	TokenDelimiter string
	// Token 加密 key
	EncryptKey []byte
	// 认证失败中文提示
	AuthFailMsg string
	// 是否支持多端登录，默认 false
	MultiLogin bool
	// 是否是全局认证，兼容历史版本，已废弃
	GlobalMiddleware bool
	// 中间件类型 1 GroupMiddleware 2 BindMiddleware  3 GlobalMiddleware
	MiddlewareType uint

	// 登录路径
	LoginPath string
	// 登录验证方法 return userKey 用户标识 如果 userKey 为空，结束执行
	LoginBeforeFunc func(r *ghttp.Request) (string, any)
	// 登录返回方法
	LoginAfterFunc func(r *ghttp.Request, respData gfresp.Response)
	// 登出地址
	LogoutPath string
	// 登出验证方法 return true 继续执行，否则结束执行
	LogoutBeforeFunc func(r *ghttp.Request) bool
	// 登出返回方法
	LogoutAfterFunc func(r *ghttp.Request, respData gfresp.Response)

	// 拦截地址
	AuthPaths g.SliceStr
	// 拦截排除地址
	AuthExcludePaths g.SliceStr
	// 认证验证方法 return true 继续执行，否则结束执行
	AuthBeforeFunc func(r *ghttp.Request) bool
	// 认证返回方法
	AuthAfterFunc func(r *ghttp.Request, respData gfresp.Response)
}

// Login 登录
func (m *Token) Login(r *ghttp.Request) {
	userKey, data := m.LoginBeforeFunc(r)
	if userKey == "" {
		gflogger.Error(r.Context(), msgLog(MsgErrUserKeyEmpty))
		return
	}

	if m.MultiLogin {
		// 支持多端重复登录，返回相同token
		userCacheResp := m.getToken(r.Context(), userKey)
		if userCacheResp.Success() {
			respToken := m.EncryptToken(r.Context(), userKey, userCacheResp.GetString(KeyUuid))
			m.LoginAfterFunc(r, respToken)
			return
		}
	}

	// 生成token
	respToken := m.genToken(r.Context(), userKey, data)
	m.LoginAfterFunc(r, respToken)

}

// Logout 登出
func (m *Token) Logout(r *ghttp.Request) {
	if !m.LogoutBeforeFunc(r) {
		return
	}

	// 获取请求token
	respData := m.getRequestToken(r)
	if respData.Success() {
		// 删除token
		m.RemoveToken(r.Context(), respData.DataString())
	}

	m.LogoutAfterFunc(r, respData)
}

// AuthMiddleware 认证拦截
func (m *Token) authMiddleware(r *ghttp.Request) {
	urlPath := r.URL.Path
	if !m.AuthPath(r.Context(), urlPath) {
		// 如果不需要认证，继续
		r.Middleware.Next()
		return
	}

	// 不需要认证，直接下一步
	if !m.AuthBeforeFunc(r) {
		r.Middleware.Next()
		return
	}

	// 获取请求token
	tokenResp := m.getRequestToken(r)
	if tokenResp.Success() {
		// 验证token
		tokenResp = m.validToken(r.Context(), tokenResp.DataString())
	}

	m.AuthAfterFunc(r, tokenResp)
}

// GetTokenData 通过token获取对象
func (m *Token) GetTokenData(r *ghttp.Request) gfresp.Response {
	respData := m.getRequestToken(r)
	if respData.Success() {
		// 验证token
		respData = m.validToken(r.Context(), respData.DataString())
	}

	return respData
}

// AuthPath 判断路径是否需要进行认证拦截
// return true 需要认证
func (m *Token) AuthPath(ctx context.Context, urlPath string) bool {
	// 去除后斜杠
	if strings.HasSuffix(urlPath, "/") {
		urlPath = gstr.SubStr(urlPath, 0, len(urlPath)-1)
	}
	// 分组拦截，登录接口不拦截
	if m.MiddlewareType == MiddlewareTypeGroup {
		if (m.LoginPath != "" && gstr.HasSuffix(urlPath, m.LoginPath)) ||
			(m.LogoutPath != "" && gstr.HasSuffix(urlPath, m.LogoutPath)) {
			return false
		}
	}

	// 全局处理，认证路径拦截处理
	if m.MiddlewareType == MiddlewareTypeGlobal {
		var authFlag bool
		for _, authPath := range m.AuthPaths {
			tmpPath := authPath
			if strings.HasSuffix(tmpPath, "/*") {
				tmpPath = gstr.SubStr(tmpPath, 0, len(tmpPath)-2)
			}
			if gstr.HasPrefix(urlPath, tmpPath) {
				authFlag = true
				break
			}
		}

		if !authFlag {
			// 拦截路径不匹配
			return false
		}
	}

	// 排除路径处理，到这里nextFlag为true
	for _, excludePath := range m.AuthExcludePaths {
		tmpPath := excludePath
		// 前缀匹配
		if strings.HasSuffix(tmpPath, "/*") {
			tmpPath = gstr.SubStr(tmpPath, 0, len(tmpPath)-2)
			if gstr.HasPrefix(urlPath, tmpPath) {
				// 前缀匹配不拦截
				return false
			}
		} else {
			// 全路径匹配
			if strings.HasSuffix(tmpPath, "/") {
				tmpPath = gstr.SubStr(tmpPath, 0, len(tmpPath)-1)
			}
			if urlPath == tmpPath {
				// 全路径匹配不拦截
				return false
			}
		}
	}

	return true
}

// getRequestToken 返回请求Token
func (m *Token) getRequestToken(r *ghttp.Request) gfresp.Response {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			gflogger.Warning(r.Context(), msgLog(MsgErrAuthHeader, authHeader))
			return gfresp.Unauthorized(fmt.Sprintf(MsgErrAuthHeader, authHeader), "")
		} else if parts[1] == "" {
			gflogger.Warning(r.Context(), msgLog(MsgErrAuthHeader, authHeader))
			return gfresp.Unauthorized(fmt.Sprintf(MsgErrAuthHeader, authHeader), "")
		}

		return gfresp.Succ(parts[1])
	}

	authHeader = r.Get(KeyToken).String()
	if authHeader == "" {
		return gfresp.Unauthorized(MsgErrTokenEmpty, "")
	}
	return gfresp.Succ(authHeader)

}

// genToken 生成Token
func (m *Token) genToken(ctx context.Context, userKey string, data any) gfresp.Response {
	token := m.EncryptToken(ctx, userKey, "")
	if !token.Success() {
		return token
	}

	cacheKey := m.CacheKey + userKey
	userCache := g.Map{
		KeyUserKey:     userKey,
		KeyUuid:        token.GetString(KeyUuid),
		KeyData:        data,
		KeyCreateTime:  gtime.Now().TimestampMilli(),
		KeyRefreshTime: gtime.Now().TimestampMilli() + gconv.Int64(m.MaxRefresh),
	}

	cacheResp := m.setCache(ctx, cacheKey, userCache)
	if !cacheResp.Success() {
		return cacheResp
	}

	return token
}

// validToken 验证Token
func (m *Token) validToken(ctx context.Context, token string) gfresp.Response {
	if token == "" {
		return gfresp.Unauthorized(MsgErrTokenEmpty, "")
	}

	decryptToken := m.DecryptToken(ctx, token)
	if !decryptToken.Success() {
		return decryptToken
	}

	userKey := decryptToken.GetString(KeyUserKey)
	uuid := decryptToken.GetString(KeyUuid)

	userCacheResp := m.getToken(ctx, userKey)
	if !userCacheResp.Success() {
		return userCacheResp
	}

	if uuid != userCacheResp.GetString(KeyUuid) {
		gflogger.Debug(ctx, msgLog(MsgErrAuthUuid)+", decryptToken:"+decryptToken.Json()+" cacheValue:"+gconv.String(userCacheResp.Data))
		return gfresp.Unauthorized(MsgErrAuthUuid, "")
	}

	return userCacheResp
}

// getToken 通过userKey获取Token
func (m *Token) getToken(ctx context.Context, userKey string) gfresp.Response {
	cacheKey := m.CacheKey + userKey

	userCacheResp := m.getCache(ctx, cacheKey)
	if !userCacheResp.Success() {
		return userCacheResp
	}
	userCache := gconv.Map(userCacheResp.Data)

	nowTime := gtime.Now().TimestampMilli()
	refreshTime := userCache[KeyRefreshTime]

	// 需要进行缓存超时时间刷新
	if gconv.Int64(refreshTime) == 0 || nowTime > gconv.Int64(refreshTime) {
		userCache[KeyCreateTime] = gtime.Now().TimestampMilli()
		userCache[KeyRefreshTime] = gtime.Now().TimestampMilli() + gconv.Int64(m.MaxRefresh)
		return m.setCache(ctx, cacheKey, userCache)
	}

	return gfresp.Succ(userCache)
}

// RemoveToken 删除Token
func (m *Token) RemoveToken(ctx context.Context, token string) gfresp.Response {
	decryptToken := m.DecryptToken(ctx, token)
	if !decryptToken.Success() {
		return decryptToken
	}

	cacheKey := m.CacheKey + decryptToken.GetString(KeyUserKey)
	return m.removeCache(ctx, cacheKey)
}

// EncryptToken token加密方法
func (m *Token) EncryptToken(ctx context.Context, userKey string, uuid string) gfresp.Response {
	if userKey == "" {
		return gfresp.Fail(FAIL, MsgErrUserKeyEmpty)
	}

	if uuid == "" {
		// 重新生成uuid
		newUuid, err := gmd5.Encrypt(grand.Letters(10))
		if err != nil {
			gflogger.Error(ctx, msgLog(MsgErrAuthUuid), err)
			return gfresp.Fail(ERROR, MsgErrAuthUuid)
		}
		uuid = newUuid
	}

	tokenStr := userKey + m.TokenDelimiter + uuid

	token, err := gaes.Encrypt([]byte(tokenStr), m.EncryptKey)
	if err != nil {
		gflogger.Error(ctx, msgLog(MsgErrTokenEncrypt), tokenStr, err)
		return gfresp.Fail(ERROR, MsgErrTokenEncrypt)
	}

	return gfresp.Succ(g.Map{
		KeyUserKey: userKey,
		KeyUuid:    uuid,
		KeyToken:   gbase64.EncodeToString(token),
	})
}

// DecryptToken token解密方法
func (m *Token) DecryptToken(ctx context.Context, token string) gfresp.Response {
	if token == "" {
		return gfresp.Fail(FAIL, MsgErrTokenEmpty)
	}

	token64, err := gbase64.Decode([]byte(token))
	if err != nil {
		gflogger.Error(ctx, msgLog(MsgErrTokenDecode), token, err)
		return gfresp.Fail(ERROR, MsgErrTokenDecode)
	}
	decryptToken, err2 := gaes.Decrypt(token64, m.EncryptKey)
	if err2 != nil {
		gflogger.Error(ctx, msgLog(MsgErrTokenEncrypt), token, err2)
		return gfresp.Fail(ERROR, MsgErrTokenEncrypt)
	}
	tokenArray := gstr.Split(string(decryptToken), m.TokenDelimiter)
	if len(tokenArray) < 2 {
		gflogger.Error(ctx, msgLog(MsgErrTokenLen), token)
		return gfresp.Fail(ERROR, MsgErrTokenLen)
	}

	return gfresp.Succ(g.Map{
		KeyUserKey: tokenArray[0],
		KeyUuid:    tokenArray[1],
	})
}

// InitConfig 初始化配置信息
func (m *Token) InitConfig() bool {
	if m.CacheMode == 0 {
		m.CacheMode = CacheModeCache
	}

	if m.RedisGroupName == "" {
		m.RedisGroupName = gredis.DefaultGroupName
	}

	if m.CacheKey == "" {
		m.CacheKey = DefaultCacheKey
	}

	if m.Timeout == 0 {
		m.Timeout = DefaultTimeout
	}

	if m.MaxRefresh == 0 {
		m.MaxRefresh = m.Timeout / 2
	}

	if m.TokenDelimiter == "" {
		m.TokenDelimiter = DefaultTokenDelimiter
	}

	if len(m.EncryptKey) == 0 {
		m.EncryptKey = []byte(DefaultEncryptKey)
	}

	if m.AuthFailMsg == "" {
		m.AuthFailMsg = DefaultAuthFailMsg
	}

	// 设置中间件模式，未设置说明历史版本，通过GlobalMiddleware兼容
	if m.MiddlewareType == 0 {
		if m.GlobalMiddleware {
			m.MiddlewareType = MiddlewareTypeGlobal
		} else {
			m.MiddlewareType = MiddlewareTypeBind
		}
	}

	if m.LoginAfterFunc == nil {
		m.LoginAfterFunc = func(r *ghttp.Request, respData gfresp.Response) {
			if !respData.Success() {
				respData.Resp(r)
			} else {
				gfresp.Succ(g.Map{KeyToken: respData.GetString(KeyToken)}).Resp(r)
			}
		}
	}

	if m.LogoutBeforeFunc == nil {
		m.LogoutBeforeFunc = func(r *ghttp.Request) bool {
			return true
		}
	}

	if m.LogoutAfterFunc == nil {
		m.LogoutAfterFunc = func(r *ghttp.Request, respData gfresp.Response) {
			if respData.Success() {
				gfresp.Succ(MsgLogoutSucc).Resp(r)
			} else {
				respData.Resp(r)
			}
		}
	}

	if m.AuthBeforeFunc == nil {
		m.AuthBeforeFunc = func(r *ghttp.Request) bool {
			// 静态页面不拦截
			return !r.IsFileRequest()
		}
	}
	if m.AuthAfterFunc == nil {
		m.AuthAfterFunc = func(r *ghttp.Request, respData gfresp.Response) {
			if respData.Success() {
				r.Middleware.Next()
			} else {
				var params map[string]any
				if r.Method == http.MethodGet {
					params = r.GetMap()
				} else if r.Method == http.MethodPost {
					params = r.GetMap()
				} else {
					r.Response.Writeln(MsgErrReqMethod)
					return
				}

				no := gconv.String(gtime.TimestampMilli())

				gflogger.Warning(r.Context(), fmt.Sprintf("[AUTH_%s][url:%s][params:%s][data:%s]",
					no, r.URL.Path, params, respData.Json()))
				respData.Message = m.AuthFailMsg
				respData.Resp(r)
				r.ExitAll()
			}
		}
	}

	return true
}

// Start 启动
func (m *Token) Start() error {
	if !m.InitConfig() {
		return errors.New(MsgErrInitFail)
	}

	ctx := context.Background()
	gflogger.Info(ctx, msgLog("[params:"+m.String()+"]start... "))

	s := g.Server(m.ServerName)

	// 缓存模式
	if m.CacheMode > CacheModeFile {
		gflogger.Error(ctx, msgLog(MsgErrNotSet, "CacheMode"))
		return fmt.Errorf(MsgErrNotSet, "CacheMode")
	}

	// 初始化文件缓存
	if m.CacheMode == 3 {
		m.initFileCache(ctx)
	}

	// 认证拦截器
	if m.AuthPaths == nil {
		gflogger.Error(ctx, msgLog(MsgErrNotSet, "AuthPaths"))
		return fmt.Errorf(MsgErrNotSet, "AuthPaths")
	}

	// 是否是全局拦截
	if m.MiddlewareType == MiddlewareTypeGlobal {
		s.BindMiddlewareDefault(m.authMiddleware)
	} else {
		for _, authPath := range m.AuthPaths {
			tmpPath := authPath
			if !strings.HasSuffix(authPath, "/*") {
				tmpPath += "/*"
			}
			s.BindMiddleware(tmpPath, m.authMiddleware)
		}
	}

	// 登录
	if m.LoginPath == "" {
		gflogger.Error(ctx, msgLog(MsgErrNotSet, "LoginPath"))
		return fmt.Errorf(MsgErrNotSet, "LoginPath")
	}
	if m.LoginBeforeFunc == nil {
		gflogger.Error(ctx, msgLog(MsgErrNotSet, "LoginBeforeFunc"))
		return fmt.Errorf(MsgErrNotSet, "LoginBeforeFunc")
	}
	s.BindHandler(m.LoginPath, m.Login)

	// 登出
	if m.LogoutPath == "" {
		gflogger.Error(ctx, msgLog(MsgErrNotSet, "LogoutPath"))
		return fmt.Errorf(MsgErrNotSet, "LogoutPath")
	}
	s.BindHandler(m.LogoutPath, m.Logout)

	return nil
}

// Stop 结束
func (m *Token) Stop(ctx context.Context) error {
	gflogger.Info(ctx, "[GToken]stop. ")
	return nil
}

// String token解密方法
func (m *Token) String() string {
	return gconv.String(g.Map{
		// 缓存模式 1:gcache 2:gredis 3:fileCache 默认1
		"CacheMode":        m.CacheMode,
		"RedisGroupName":   m.RedisGroupName,
		"CacheKey":         m.CacheKey,
		"Timeout":          m.Timeout,
		"TokenDelimiter":   m.TokenDelimiter,
		"AuthFailMsg":      m.AuthFailMsg,
		"MultiLogin":       m.MultiLogin,
		"MiddlewareType":   m.MiddlewareType,
		"LoginPath":        m.LoginPath,
		"LogoutPath":       m.LogoutPath,
		"AuthPaths":        gconv.String(m.AuthPaths),
		"AuthExcludePaths": gconv.String(m.AuthExcludePaths),
	})
}
