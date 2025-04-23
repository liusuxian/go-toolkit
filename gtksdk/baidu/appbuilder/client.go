/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-26 22:33:32
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-23 19:22:11
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package appbuilder

import "github.com/liusuxian/go-toolkit/gtkhttp"

// Client 客户端
type Client struct {
	httpClient *gtkhttp.HTTPClient // http 客户端
	appToken   string              // 应用 token
}

// NewClient 新建客户端
func NewClient(appToken string) (c *Client) {
	return &Client{
		httpClient: gtkhttp.NewHTTPClient("https://appbuilder.baidu.com"),
		appToken:   appToken,
	}
}
