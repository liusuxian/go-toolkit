/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2026-01-13 10:50:22
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2026-01-13 19:07:36
 * @Description:
 *
 * Copyright (c) 2026 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtksms_test

import (
	"context"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/joho/godotenv"
	"github.com/liusuxian/go-toolkit/gtkcache"
	"github.com/liusuxian/go-toolkit/gtkenv"
	"github.com/liusuxian/go-toolkit/gtkredis"
	"github.com/liusuxian/go-toolkit/gtksdk/aliyun/gtksms"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAliyunSendSms(t *testing.T) {
	var (
		assert = assert.New(t)
		ctx    = context.Background()
		r      = miniredis.RunT(t)
		err    error
	)
	err = godotenv.Load(".env")
	assert.NoError(err)
	// 创建 RedisCache
	var cache *gtkcache.RedisCache
	cache, err = gtkcache.NewRedisCache(ctx, &gtkredis.ClientConfig{
		Addr:     r.Addr(),
		Username: "default",
		Password: "",
		DB:       1,
	})
	assert.NoError(err)
	assert.NotNil(cache)
	// 创建 AliyunSMS
	var aliyunSMS *gtksms.AliyunSMS
	aliyunSMS, err = gtksms.NewAliyunSMS(gtksms.SMSConfig{
		Endpoint:        gtkenv.Get("ALIYUN_SMS_ENDPOINT"),
		AccessKeyID:     gtkenv.Get("ALIYUN_SMS_ACCESS_KEY_ID"),
		AccessKeySecret: gtkenv.Get("ALIYUN_SMS_ACCESS_KEY_SECRET"),
	}, cache)
	assert.NoError(err)
	assert.NotNil(aliyunSMS)
	// 生成验证码
	var verifyCode string
	verifyCode, err = aliyunSMS.GenerateVerifyCode(ctx, gtkenv.Get("ALIYUN_SMS_PHONE_NUMBERS"))
	assert.NoError(err)
	assert.NotEmpty(verifyCode)
	// 发送短信
	var bizId string
	bizId, err = aliyunSMS.SendSms(&gtksms.SendSmsRequest{
		PhoneNumbers:  gtkenv.Get("ALIYUN_SMS_PHONE_NUMBERS"),
		SignName:      gtkenv.Get("ALIYUN_SMS_SIGN_NAME"),
		TemplateCode:  gtkenv.Get("ALIYUN_SMS_TEMPLATE_CODE"),
		TemplateParam: fmt.Sprintf("{\"code\":\"%s\"}", verifyCode),
	})
	assert.NoError(err)
	assert.NotEmpty(bizId)
	// 校验验证码
	err = aliyunSMS.CheckVerifyCode(ctx, gtkenv.Get("ALIYUN_SMS_PHONE_NUMBERS"), verifyCode)
	assert.NoError(err)
}
