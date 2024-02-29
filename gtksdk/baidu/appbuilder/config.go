/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-26 22:28:26
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-26 23:52:07
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package appbuilder

import "net/http"

const (
	baseURL                        = "https://appbuilder.baidu.com"
	defaultEmptyMessagesLimit uint = 300 // 默认空消息限制
)

// ClientConfig 客户端配置
type ClientConfig struct {
	AppToken           string // 应用 token
	BaseURL            string
	HTTPClient         *http.Client
	EmptyMessagesLimit uint // 空消息限制
}

// DefaultConfig 默认客户端配置
func DefaultConfig(appToken string) (config ClientConfig) {
	return ClientConfig{
		AppToken:           appToken,
		BaseURL:            baseURL,
		HTTPClient:         &http.Client{},
		EmptyMessagesLimit: defaultEmptyMessagesLimit,
	}
}
