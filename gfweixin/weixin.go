/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-19 22:29:06
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-19 23:33:37
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gfweixin

import (
	"context"
	"net/url"
)

// AuthCode2Session 登录凭证校验
func AuthCode2Session(ctx context.Context, appid, secret, code string) (resMap map[string]any, err error) {
	// 组装参数
	params := url.Values{}
	params.Add("appid", appid)
	params.Add("secret", secret)
	params.Add("js_code", code)
	params.Add("grant_type", "authorization_code")
	// 发起请求
	return Get(ctx, "https://api.weixin.qq.com/sns/jscode2session", params)
}

// GetStableAccessToken 获取稳定版接口调用凭据
func GetStableAccessToken(ctx context.Context, appid, secret string, forceRefresh ...bool) (resMap map[string]any, err error) {
	var newForceRefresh bool
	if len(forceRefresh) > 0 {
		newForceRefresh = forceRefresh[0]
	}
	// 组装参数
	params := url.Values{}
	body := map[string]any{
		"grant_type":    "client_credential",
		"appid":         appid,
		"secret":        secret,
		"force_refresh": newForceRefresh,
	}
	// 发起请求
	return Post(ctx, "https://api.weixin.qq.com/cgi-bin/stable_token", params, body)
}
