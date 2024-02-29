/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-27 21:54:54
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-28 15:48:36
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package redbookweb

import (
	"fmt"
	"github.com/liusuxian/go-toolkit/gtkhttp"
	"github.com/pkg/errors"
	"net/http"
)

// Client 客户端
type Client struct {
	config         ClientConfig           // 客户端配置
	requestBuilder gtkhttp.RequestBuilder // 请求构建器
}

// NewClient 新建客户端
func NewClient() (c *Client) {
	return NewClientWithConfig(DefaultConfig())
}

// NewClientWithConfig 通过客户端配置新建客户端
func NewClientWithConfig(config ClientConfig) (c *Client) {
	return &Client{
		config:         config,
		requestBuilder: gtkhttp.NewRequestBuilder(),
	}
}

// sendRequest 发送请求
func (c *Client) sendRequest(req *http.Request, v any) (respCookies []*http.Cookie, err error) {
	var resp *http.Response
	if resp, err = c.config.HTTPClient.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	if gtkhttp.IsFailureStatusCode(resp) {
		err = &gtkhttp.RequestError{
			HTTPStatusCode: resp.StatusCode,
			Err:            errors.New("Request Failed"),
		}
		return
	}

	if err = gtkhttp.DecodeResponse(resp.Body, v); err != nil {
		return
	}

	respCookies = resp.Cookies()
	return
}

// customerFullURL 获取完整链接
func (c *Client) customerFullURL(suffix string) (url string) {
	return fmt.Sprintf("%s%s", c.config.CustomerBaseURL, suffix)
}

// creatorFullURL 获取完整链接
func (c *Client) creatorFullURL(suffix string) (url string) {
	return fmt.Sprintf("%s%s", c.config.CreatorBaseURL, suffix)
}

// edithFullURL 获取完整链接
func (c *Client) edithFullURL(suffix string) (url string) {
	return fmt.Sprintf("%s%s", c.config.EdithBaseURL, suffix)
}
