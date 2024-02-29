/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-28 00:27:28
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-28 17:16:47
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package redbookweb

import (
	"context"
	"github.com/liusuxian/go-toolkit/gtkhttp"
	"net/http"
)

const (
	customerLoginURL = "/sso/customer_login"
	loginURL         = "/api/galaxy/user/cas/login"
	userInfoURL      = "/api/galaxy/user/info"
)

// CustomerLogin 客户登录
func (c *Client) CustomerLogin(ctx context.Context, request CustomerLoginRequest, reqCookies []*http.Cookie) (response *CustomerLoginResponse, respCookies []*http.Cookie, err error) {
	var (
		setters = []gtkhttp.RequestOption{
			gtkhttp.SetBody(map[string]any{
				"login_service":     "https://creator.xiaohongshu.com",
				"set_global_domain": true,
				"subsystem_alias":   "ares",
				"ticket":            request.Ticket,
			}),
			gtkhttp.SetKeyValue("Accept", "application/json, text/plain, */*"),
			gtkhttp.SetKeyValue("Accept-Encoding", "gzip, deflate, br"),
			gtkhttp.SetKeyValue("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6"),
			gtkhttp.SetContentType("application/json; charset=utf-8"),
			gtkhttp.SetCookie(reqCookies),
			gtkhttp.SetKeyValue("Origin", "https://creator.xiaohongshu.com"),
			gtkhttp.SetKeyValue("Referer", "https://creator.xiaohongshu.com/login?selfLogout=true"),
			gtkhttp.SetKeyValue("Sec-Ch-Ua", "Not A(Brand;v=99, Microsoft Edge;v=121, Chromium;v=121"),
			gtkhttp.SetKeyValue("Sec-Ch-Ua-Mobile", "?0"),
			gtkhttp.SetKeyValue("Sec-Ch-Ua-Platform", "macOS"),
			gtkhttp.SetKeyValue("Sec-Fetch-Dest", "empty"),
			gtkhttp.SetKeyValue("Sec-Fetch-Mode", "cors"),
			gtkhttp.SetKeyValue("Sec-Fetch-Site", "same-origin"),
			gtkhttp.SetKeyValue("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/121.0.0.0"),
		}
		req *http.Request
	)

	if req, err = c.requestBuilder.Build(ctx, http.MethodPost, c.creatorFullURL(customerLoginURL), setters...); err != nil {
		return
	}

	respCookies, err = c.sendRequest(req, &response)
	return
}

// Login 登录
func (c *Client) Login(ctx context.Context, reqCookies []*http.Cookie) (response *LoginResponse, respCookies []*http.Cookie, err error) {
	var (
		setters = []gtkhttp.RequestOption{
			gtkhttp.SetKeyValue("Accept", "application/json, text/plain, */*"),
			gtkhttp.SetKeyValue("Accept-Encoding", "gzip, deflate, br"),
			gtkhttp.SetKeyValue("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6"),
			gtkhttp.SetCookie(reqCookies),
			gtkhttp.SetKeyValue("Origin", "https://creator.xiaohongshu.com"),
			gtkhttp.SetKeyValue("Referer", "https://creator.xiaohongshu.com/login?selfLogout=true"),
			gtkhttp.SetKeyValue("Sec-Ch-Ua", "Not A(Brand;v=99, Microsoft Edge;v=121, Chromium;v=121"),
			gtkhttp.SetKeyValue("Sec-Ch-Ua-Mobile", "?0"),
			gtkhttp.SetKeyValue("Sec-Ch-Ua-Platform", "macOS"),
			gtkhttp.SetKeyValue("Sec-Fetch-Dest", "empty"),
			gtkhttp.SetKeyValue("Sec-Fetch-Mode", "cors"),
			gtkhttp.SetKeyValue("Sec-Fetch-Site", "same-origin"),
			gtkhttp.SetKeyValue("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/121.0.0.0"),
		}
		req *http.Request
	)

	if req, err = c.requestBuilder.Build(ctx, http.MethodPost, c.creatorFullURL(loginURL), setters...); err != nil {
		return
	}

	respCookies, err = c.sendRequest(req, &response)
	return
}

// UserInfo 获取用户信息
func (c *Client) UserInfo(ctx context.Context, reqCookies []*http.Cookie) (response *UserInfoResponse, respCookies []*http.Cookie, err error) {
	var (
		setters = []gtkhttp.RequestOption{
			gtkhttp.SetKeyValue("Accept", "application/json, text/plain, */*"),
			gtkhttp.SetKeyValue("Accept-Encoding", "gzip, deflate, br"),
			gtkhttp.SetKeyValue("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6"),
			gtkhttp.SetCookie(reqCookies),
			gtkhttp.SetKeyValue("Origin", "https://creator.xiaohongshu.com"),
			gtkhttp.SetKeyValue("Referer", "https://creator.xiaohongshu.com/login?selfLogout=true"),
			gtkhttp.SetKeyValue("Sec-Ch-Ua", "Not A(Brand;v=99, Microsoft Edge;v=121, Chromium;v=121"),
			gtkhttp.SetKeyValue("Sec-Ch-Ua-Mobile", "?0"),
			gtkhttp.SetKeyValue("Sec-Ch-Ua-Platform", "macOS"),
			gtkhttp.SetKeyValue("Sec-Fetch-Dest", "empty"),
			gtkhttp.SetKeyValue("Sec-Fetch-Mode", "cors"),
			gtkhttp.SetKeyValue("Sec-Fetch-Site", "same-origin"),
			gtkhttp.SetKeyValue("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/121.0.0.0"),
		}
		req *http.Request
	)

	if req, err = c.requestBuilder.Build(ctx, http.MethodGet, c.creatorFullURL(userInfoURL), setters...); err != nil {
		return
	}

	respCookies, err = c.sendRequest(req, &response)
	return
}
