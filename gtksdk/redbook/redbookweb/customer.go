/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-27 22:00:33
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-28 14:27:15
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package redbookweb

import (
	"context"
	"fmt"
	"github.com/liusuxian/go-toolkit/gtkhttp"
	"net/http"
	"net/url"
)

const (
	sendCodeURL            = "/api/cas/sendCode"
	loginWithVerifyCodeURL = "/api/cas/loginWithVerifyCode"
)

// SendCode 发送验证码
func (c *Client) SendCode(ctx context.Context, request SendCodeRequest) (response *SendCodeResponse, err error) {
	var (
		setters = []gtkhttp.RequestOption{
			gtkhttp.SetKeyValue("Accept", "application/json, text/plain, */*"),
			gtkhttp.SetKeyValue("Accept-Encoding", "gzip, deflate, br"),
			gtkhttp.SetKeyValue("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6"),
			gtkhttp.SetKeyValue("Origin", "https://creator.xiaohongshu.com"),
			gtkhttp.SetKeyValue("Referer", "https://creator.xiaohongshu.com/"),
			gtkhttp.SetKeyValue("Sec-Ch-Ua", "Not A(Brand;v=99, Microsoft Edge;v=121, Chromium;v=121"),
			gtkhttp.SetKeyValue("Sec-Ch-Ua-Mobile", "?0"),
			gtkhttp.SetKeyValue("Sec-Ch-Ua-Platform", "macOS"),
			gtkhttp.SetKeyValue("Sec-Fetch-Dest", "empty"),
			gtkhttp.SetKeyValue("Sec-Fetch-Mode", "cors"),
			gtkhttp.SetKeyValue("Sec-Fetch-Site", "same-site"),
			gtkhttp.SetKeyValue("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/121.0.0.0"),
		}
		values = url.Values{}
		req    *http.Request
	)

	values.Add("zone", request.Zone)
	values.Add("phone", request.Phone)
	params := values.Encode()
	apiUrl := fmt.Sprintf("%s?%s", c.customerFullURL(sendCodeURL), params)
	if req, err = c.requestBuilder.Build(ctx, http.MethodGet, apiUrl, setters...); err != nil {
		return
	}

	_, err = c.sendRequest(req, &response)
	return
}

// LoginWithVerifyCode 使用验证码登录
func (c *Client) LoginWithVerifyCode(ctx context.Context, request LoginWithVerifyCodeRequest) (response *LoginWithVerifyCodeResponse, respCookies []*http.Cookie, err error) {
	var (
		setters = []gtkhttp.RequestOption{
			gtkhttp.SetBody(map[string]any{
				"zone":       request.Zone,
				"mobile":     request.Mobile,
				"verifyCode": request.VerifyCode,
				"service":    "https://creator.xiaohongshu.com",
			}),
			gtkhttp.SetKeyValue("Accept", "application/json, text/plain, */*"),
			gtkhttp.SetKeyValue("Accept-Encoding", "gzip, deflate, br"),
			gtkhttp.SetKeyValue("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6"),
			gtkhttp.SetContentType("application/json; charset=utf-8"),
			gtkhttp.SetKeyValue("Origin", "https://creator.xiaohongshu.com"),
			gtkhttp.SetKeyValue("Referer", "https://creator.xiaohongshu.com/"),
			gtkhttp.SetKeyValue("Sec-Ch-Ua", "Not A(Brand;v=99, Microsoft Edge;v=121, Chromium;v=121"),
			gtkhttp.SetKeyValue("Sec-Ch-Ua-Mobile", "?0"),
			gtkhttp.SetKeyValue("Sec-Ch-Ua-Platform", "macOS"),
			gtkhttp.SetKeyValue("Sec-Fetch-Dest", "empty"),
			gtkhttp.SetKeyValue("Sec-Fetch-Mode", "cors"),
			gtkhttp.SetKeyValue("Sec-Fetch-Site", "same-site"),
			gtkhttp.SetKeyValue("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/121.0.0.0"),
		}
		req *http.Request
	)

	if req, err = c.requestBuilder.Build(ctx, http.MethodPost, c.customerFullURL(loginWithVerifyCodeURL), setters...); err != nil {
		return
	}

	respCookies, err = c.sendRequest(req, &response)
	return
}
