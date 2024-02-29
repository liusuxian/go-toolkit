/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-28 15:51:01
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-29 15:15:00
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
	noteURL = "/web_api/sns/v2/note"
)

// TODO PublishNote 发布笔记（暂未跑通）
func (c *Client) PublishNote(ctx context.Context, request PublishNoteRequest, reqCookies []*http.Cookie) (response *PublishNoteResponse, respCookies []*http.Cookie, err error) {
	var (
		setters = []gtkhttp.RequestOption{
			gtkhttp.SetBody(map[string]any{}),
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
			gtkhttp.SetKeyValue("Sec-Fetch-Site", "same-site"),
			gtkhttp.SetKeyValue("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/121.0.0.0"),
		}
		req *http.Request
	)

	if req, err = c.requestBuilder.Build(ctx, http.MethodPost, c.edithFullURL(noteURL), setters...); err != nil {
		return
	}

	respCookies, err = c.sendRequest(req, &response)
	return
}
