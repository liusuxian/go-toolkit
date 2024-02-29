/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-26 22:33:32
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-27 13:17:34
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package appbuilder

import (
	"bufio"
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
func NewClient(appToken string) (c *Client) {
	return NewClientWithConfig(DefaultConfig(appToken))
}

// NewClientWithConfig 通过客户端配置新建客户端
func NewClientWithConfig(config ClientConfig) (c *Client) {
	return &Client{
		config:         config,
		requestBuilder: gtkhttp.NewRequestBuilder(),
	}
}

// sendRequest 发送请求
func (c *Client) sendRequest(req *http.Request, v any) (err error) {
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

	return gtkhttp.DecodeResponse(resp.Body, v)
}

// sendRequestStream 发送请求
func sendRequestStream[T streamable](client *Client, req *http.Request) (stream *streamReader[T], err error) {
	var resp *http.Response
	if resp, err = client.config.HTTPClient.Do(req); err != nil {
		stream = &streamReader[T]{}
		return
	}

	if gtkhttp.IsFailureStatusCode(resp) {
		stream = &streamReader[T]{}
		err = &gtkhttp.RequestError{
			HTTPStatusCode: resp.StatusCode,
			Err:            errors.New("Request Stream Failed"),
		}
		return
	}

	stream = &streamReader[T]{
		emptyMessagesLimit: client.config.EmptyMessagesLimit,
		reader:             bufio.NewReader(resp.Body),
		response:           resp,
		unmarshaler:        &gtkhttp.JSONUnmarshaler{},
	}
	return
}

// fullURL 获取完整链接
func (c *Client) fullURL(suffix string) (url string) {
	return fmt.Sprintf("%s%s", c.config.BaseURL, suffix)
}
