/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-27 21:51:42
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-27 21:57:17
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package redbookweb

import "net/http"

const (
	customerBaseURL = "https://customer.xiaohongshu.com"
	creatorBaseURL  = "https://creator.xiaohongshu.com"
	edithBaseURL    = "https://edith.xiaohongshu.com"
)

// ClientConfig 客户端配置
type ClientConfig struct {
	CustomerBaseURL string
	CreatorBaseURL  string
	EdithBaseURL    string
	HTTPClient      *http.Client
}

// DefaultConfig 默认客户端配置
func DefaultConfig() (config ClientConfig) {
	return ClientConfig{
		CustomerBaseURL: customerBaseURL,
		CreatorBaseURL:  creatorBaseURL,
		EdithBaseURL:    edithBaseURL,
		HTTPClient:      &http.Client{},
	}
}
