/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-19 22:29:06
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-19 23:12:15
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gfweixin

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gurl"
	"net/url"
)

// AuthCode2Session 登录凭证校验
func AuthCode2Session(ctx context.Context, appid, secret, code string) (resMap map[string]any, err error) {
	// 组装参数
	values := url.Values{}
	values.Add("appid", appid)
	values.Add("secret", secret)
	values.Add("js_code", code)
	values.Add("grant_type", "authorization_code")
	params := gurl.BuildQuery(values)
	// 发起请求
	return Get(ctx, "https://api.weixin.qq.com/sns/jscode2session", params)
}
