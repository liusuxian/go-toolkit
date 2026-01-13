/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2026-01-13 10:50:22
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2026-01-13 19:09:00
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
	"github.com/liusuxian/go-toolkit/gtktype"
	"math/rand/v2"
	"slices"
	"strings"
	"time"
)

var (
	ErrVerifyCodeNotFound     = errors.New("verify code not found")      // 验证码不存在
	ErrVerifyCodeInvalid      = errors.New("verify code is invalid")     // 验证码无效
	ErrSendVerifyCodeTooOften = errors.New("send verify code too often") // 发送验证码过于频繁
)

// Cache Key
const (
	// 验证码缓存 verifycode:[phoneNumbers] 手机号 验证码（String）
	keyVerifyCode = "verifycode:%v"
	// 短信冷却时间 cooldown:[phoneNumbers] 手机号 冷却时间（String）
	keySmsCooldown = "cooldown:%v"
)

// Cache 缓存
type Cache interface {
	Get(ctx context.Context, key string, timeout ...time.Duration) (val any, err error) // 获取缓存
	Set(ctx context.Context, key string, val any, timeout ...time.Duration) (err error) // 设置缓存
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
	cache          Cache            // 缓存器
	cacheKeyPrefix string           // 缓存键前缀
	client         *dysmsapi.Client // 阿里云短信服务客户端
}

// NewAliyunSMS 创建阿里云短信服务
func NewAliyunSMS(config SMSConfig, cache Cache, cacheKeyPrefix ...string) (s *AliyunSMS, err error) {
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
func (s *AliyunSMS) GenerateVerifyCode(ctx context.Context, phoneNumbers string) (verifyCode string, err error) {
	// 防止在极短时间内为同一手机号生成多个验证码
	var (
		cooldownKey = s.cacheKeyPrefix + fmt.Sprintf(keySmsCooldown, phoneNumbers)
		cooldownVal any
	)
	if cooldownVal, err = s.cache.Get(ctx, cooldownKey); err != nil {
		err = fmt.Errorf("generate verify code failed: %w", err)
		return
	}
	if cooldownVal != nil {
		if v, ok := cooldownVal.(string); ok && v == "1" {
			err = ErrSendVerifyCodeTooOften
			return
		}
	}
	// 生成验证码
	verifyCode = fmt.Sprintf("%06d", rand.IntN(1000000))
	// 设置冷却时间
	if err = s.cache.Set(ctx, cooldownKey, 1, s.config.VerifyCodeCooldownTime); err != nil {
		err = fmt.Errorf("generate verify code failed: %w", err)
		return
	}
	// 设置验证码过期时间
	verifyCodeKey := s.cacheKeyPrefix + fmt.Sprintf(keyVerifyCode, phoneNumbers)
	if err = s.cache.Set(ctx, verifyCodeKey, verifyCode, s.config.VerifyCodeExpireTime); err != nil {
		err = fmt.Errorf("generate verify code failed: %w", err)
		return
	}
	return
}

// CheckVerifyCode 校验验证码
func (s *AliyunSMS) CheckVerifyCode(ctx context.Context, phoneNumbers, verifyCode string) (err error) {
	// 验证码白名单
	if slices.Contains(s.config.VerifyCodeWhiteList, verifyCode) {
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
	if cacheVal == nil {
		return ErrVerifyCodeNotFound
	}
	cacheVerifyCode, ok := cacheVal.(string)
	if !ok || cacheVerifyCode == "" {
		return ErrVerifyCodeNotFound
	}
	if cacheVerifyCode != verifyCode {
		return ErrVerifyCodeInvalid
	}
	return nil
}
