/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-19 22:29:06
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-22 12:39:18
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gfweixin

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/liusuxian/go-toolkit/gfcache"
	"net/url"
	"time"
)

const (
	// 微信 Accesstoken key wechat:accesstoken:[appid] appid String
	KeyWeChatAccesstoken = "wechat:accesstoken:%v"
)

// WeChatService 微信服务
type WeChatService struct {
	appid        string               // appid
	secret       string               // 密钥
	cacheAdapter gfcache.CacheAdapter // 缓存适配器
}

// NewWeChatService 创建微信服务
func NewWeChatService(appid, secret string, adapter gfcache.CacheAdapter) (s *WeChatService) {
	return &WeChatService{
		appid:        appid,
		secret:       secret,
		cacheAdapter: adapter,
	}
}

// AuthCode2Session 登录凭证校验
func (s *WeChatService) AuthCode2Session(ctx context.Context, code string) (resMap map[string]any, err error) {
	// 组装参数
	params := url.Values{}
	params.Add("appid", s.appid)
	params.Add("secret", s.secret)
	params.Add("js_code", code)
	params.Add("grant_type", "authorization_code")
	// 发起请求
	if resMap, err = Get(ctx, "https://api.weixin.qq.com/sns/jscode2session", params); err != nil {
		return
	}
	// 处理结果
	errCode := gconv.Int(resMap["errcode"])
	if errCode != 0 {
		err = gerror.NewCode(gcode.New(errCode, gconv.String(resMap["errmsg"]), nil))
		return
	}
	return
}

// GetStableAccessToken 获取稳定版接口调用凭据
func (s *WeChatService) GetStableAccessToken(ctx context.Context, forceRefresh ...bool) (accessToken string, err error) {
	var newForceRefresh bool
	if len(forceRefresh) > 0 {
		newForceRefresh = forceRefresh[0]
	}
	// 非强制刷新，读取缓存
	cacheKey := fmt.Sprintf(KeyWeChatAccesstoken, s.appid)
	if !newForceRefresh && s.cacheAdapter != nil {
		var val *gvar.Var
		if val, err = s.cacheAdapter.Get(ctx, cacheKey); err != nil {
			return
		}
		if !val.IsNil() && val.String() != "" {
			accessToken = val.String()
			return
		}
	}
	// 组装参数
	params := url.Values{}
	body := map[string]any{
		"grant_type":    "client_credential",
		"appid":         s.appid,
		"secret":        s.secret,
		"force_refresh": newForceRefresh,
	}
	// 发起请求
	var resMap map[string]any
	if resMap, err = Post(ctx, "https://api.weixin.qq.com/cgi-bin/stable_token", params, body); err != nil {
		return
	}
	// 处理结果
	errCode := gconv.Int(resMap["errcode"])
	if errCode != 0 {
		err = gerror.NewCode(gcode.New(errCode, gconv.String(resMap["errmsg"]), nil))
		return
	}
	// 缓存写入
	accessToken = gconv.String(resMap["access_token"])
	if s.cacheAdapter != nil {
		expiresIn := gconv.Int(resMap["expires_in"])
		err = s.cacheAdapter.Set(ctx, cacheKey, accessToken, time.Second*time.Duration(expiresIn-10))
	}
	return
}

// GetPhoneNumber 获取手机号
func (s *WeChatService) GetPhoneNumber(ctx context.Context, code, accessToken string) (phoneNumber string, invalidAccessToken bool, err error) {
	// 组装参数
	params := url.Values{}
	params.Add("access_token", accessToken)
	body := map[string]any{
		"code": code,
	}
	// 发起请求
	var resMap map[string]any
	if resMap, err = Post(ctx, "https://api.weixin.qq.com/wxa/business/getuserphonenumber", params, body); err != nil {
		return
	}
	// 处理结果
	errCode := gconv.Int(resMap["errcode"])
	if errCode == 40001 || errCode == 40014 {
		// access_token 无效
		invalidAccessToken = true
		return
	}
	if errCode != 0 {
		err = gerror.NewCode(gcode.New(errCode, gconv.String(resMap["errmsg"]), nil))
		return
	}
	// 获取手机号
	phoneInfo := gconv.Map(resMap["phone_info"])
	countryCode := gconv.String(phoneInfo["countryCode"])
	purePhoneNumber := gconv.String(phoneInfo["purePhoneNumber"])
	phoneNumber = fmt.Sprintf("%v-%v", countryCode, purePhoneNumber)
	return
}
