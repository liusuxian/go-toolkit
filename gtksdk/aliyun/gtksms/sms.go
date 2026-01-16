/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2026-01-13 10:50:22
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2026-01-16 13:06:28
 * @Description:
 *
 * Copyright (c) 2026 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtksms

import (
	"context"
	"errors"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	"github.com/liusuxian/go-toolkit/gtkconv"
	"github.com/liusuxian/go-toolkit/gtktype"
	"math/rand/v2"
	"slices"
	"strings"
	"time"
)

var (
	ErrGenerateVerifyCodeFailed = errors.New("generate verify code failed") // 生成验证码失败
	ErrVerifyCodeInvalid        = errors.New("verify code is invalid")      // 验证码无效
	ErrSendVerifyCodeTooOften   = errors.New("send verify code too often")  // 发送验证码过于频繁
)

// Cache Key
const (
	// 验证码缓存 verifycode:[phoneNumbers] 手机号 验证码（String）
	keyVerifyCode = "verifycode:%v"
	// 短信冷却时间 cooldown:[phoneNumbers] 手机号 冷却时间（String）
	keySmsCooldown = "cooldown:%v"
)

// ICache 缓存接口
type ICache interface {
	// Get 获取缓存
	//   当`timeout > 0`且缓存命中时，设置/重置`key`的过期时间
	Get(ctx context.Context, key string, timeout ...time.Duration) (val any, err error)
	// BatchSet 批量设置缓存
	//   支持为每个`key`设置不同的过期时间
	//   当所有`key`使用相同过期时间时，可以使用更简洁的`SetMap`方法
	//   defaultTimeout: 可选参数，设置默认过期时间（对所有未单独设置过期时间的 key 生效）
	//   当`defaultTimeout > 0`时，所有未单独指定过期时间的`key`将使用此默认过期时间
	//   当`defaultTimeout <= 0`时，所有未单独指定过期时间的`key`将保持原有的过期时间
	BatchSet(ctx context.Context, fn func(add func(key string, val any, timeout ...time.Duration)), defaultTimeout ...time.Duration) (err error)
	// Delete 删除缓存
	Delete(ctx context.Context, keys ...string) (err error)
}

// SMSConfig 阿里云短信配置
type SMSConfig struct {
	Endpoint               string        `json:"endpoint"`                  // 短信服务节点
	AccessKeyID            string        `json:"access_key_id"`             // accessKeyID
	AccessKeySecret        string        `json:"access_key_secret"`         // accessKeySecret
	VerifyCodeExpireTime   time.Duration `json:"verify_code_expire_time"`   // 验证码过期时间，默认5分钟
	VerifyCodeCooldownTime time.Duration `json:"verify_code_cooldown_time"` // 发送验证码冷却时间，默认1分钟
	VerifyCodeWhiteList    []string      `json:"verify_code_white_list"`    // 验证码白名单，测试用
}

// AliyunSMS 阿里云短信服务
type AliyunSMS struct {
	config         SMSConfig        // 配置
	cache          ICache           // 缓存器
	cacheKeyPrefix string           // 缓存键前缀
	client         *dysmsapi.Client // 阿里云短信服务客户端
}

// NewAliyunSMS 创建阿里云短信服务
func NewAliyunSMS(config SMSConfig, cache ICache, cacheKeyPrefix ...string) (s *AliyunSMS, err error) {
	// 检查配置
	if config.AccessKeyID == "" {
		return nil, fmt.Errorf("access key id is not set")
	}
	if config.AccessKeySecret == "" {
		return nil, fmt.Errorf("access key secret is not set")
	}
	// 设置配置默认值
	if config.Endpoint == "" {
		config.Endpoint = "dysmsapi.aliyuncs.com"
	}
	if config.VerifyCodeExpireTime <= time.Duration(0) {
		config.VerifyCodeExpireTime = 5 * time.Minute
	}
	if config.VerifyCodeCooldownTime <= time.Duration(0) {
		config.VerifyCodeCooldownTime = 1 * time.Minute
	}
	// 检查缓存器
	if cache == nil {
		return nil, fmt.Errorf("cache is not set")
	}
	// 设置缓存键前缀
	var cachePrefix string
	if len(cacheKeyPrefix) > 0 {
		cachePrefix = cacheKeyPrefix[0]
		if !strings.HasSuffix(cachePrefix, ":") {
			cachePrefix += ":"
		}
	} else {
		cachePrefix = "gtksms:"
	}
	s = &AliyunSMS{
		config:         config,
		cache:          cache,
		cacheKeyPrefix: cachePrefix,
	}
	// 创建阿里云短信服务客户端
	if s.client, err = dysmsapi.NewClient(&openapi.Config{
		Endpoint:        gtktype.String(config.Endpoint),
		AccessKeyId:     gtktype.String(config.AccessKeyID),
		AccessKeySecret: gtktype.String(config.AccessKeySecret),
	}); err != nil {
		return nil, fmt.Errorf("create aliyun sms client failed: %v", err)
	}
	return
}

// SendSms 发送短信
func (s *AliyunSMS) SendSms(req *SendSmsRequest) (bizId string, err error) {
	var response *dysmsapi.SendSmsResponse
	if response, err = s.client.SendSms(&dysmsapi.SendSmsRequest{
		PhoneNumbers:  gtktype.String(req.PhoneNumbers),
		SignName:      gtktype.String(req.SignName),
		TemplateCode:  gtktype.String(req.TemplateCode),
		TemplateParam: gtktype.String(req.TemplateParam),
	}); err != nil {
		err = fmt.Errorf("send sms failed: %w", err)
		return
	}
	if gtktype.StringValue(response.Body.Code) != "OK" {
		err = fmt.Errorf("send sms failed: %s", gtktype.StringValue(response.Body.Message))
		return
	}
	bizId = gtktype.StringValue(response.Body.BizId)
	return
}

// GenerateVerifyCode 生成验证码
func (s *AliyunSMS) GenerateVerifyCode(ctx context.Context, phoneNumbers string) (code string, err error) {
	// 防止在极短时间内为同一手机号生成多个验证码
	var (
		cooldownKey = s.cacheKeyPrefix + fmt.Sprintf(keySmsCooldown, phoneNumbers)
		cooldownVal any
	)
	if cooldownVal, err = s.cache.Get(ctx, cooldownKey); err != nil {
		err = errors.Join(ErrGenerateVerifyCodeFailed, err)
		return
	}
	// 将 any 转换为 bool 类型
	var isCooldown bool
	if isCooldown, err = gtkconv.ToBoolE(cooldownVal); err != nil {
		err = errors.Join(ErrGenerateVerifyCodeFailed, err)
		return
	}
	if isCooldown {
		err = ErrSendVerifyCodeTooOften
		return
	}
	// 生成验证码
	code = fmt.Sprintf("%06d", rand.IntN(1000000))
	// 设置冷却时间和验证码过期时间
	if err = s.cache.BatchSet(ctx, func(add func(key string, val any, timeout ...time.Duration)) {
		add(cooldownKey, 1, s.config.VerifyCodeCooldownTime)
		add(s.cacheKeyPrefix+fmt.Sprintf(keyVerifyCode, phoneNumbers), code, s.config.VerifyCodeExpireTime)
	}); err != nil {
		err = errors.Join(ErrGenerateVerifyCodeFailed, err)
		return
	}
	return
}

// CheckVerifyCode 校验验证码
func (s *AliyunSMS) CheckVerifyCode(ctx context.Context, phoneNumbers, code string) (err error) {
	// 验证码白名单
	if slices.Contains(s.config.VerifyCodeWhiteList, code) {
		return nil
	}
	// 检查验证码
	var (
		cacheKey = s.cacheKeyPrefix + fmt.Sprintf(keyVerifyCode, phoneNumbers)
		cacheVal any
	)
	if cacheVal, err = s.cache.Get(ctx, cacheKey); err != nil {
		return fmt.Errorf("check verify code failed: %w", err)
	}
	// 将 any 转换为 string 类型
	var cacheCode string
	if cacheCode, err = gtkconv.ToStringE(cacheVal); err != nil {
		return fmt.Errorf("check verify code failed: %w", err)
	}
	if cacheCode != code {
		return ErrVerifyCodeInvalid
	}
	// 验证成功后删除验证码
	_ = s.cache.Delete(ctx, cacheKey)
	return nil
}
