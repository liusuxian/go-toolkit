/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-26 20:27:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-28 14:59:22
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// RequestOptions 请求选项
type RequestOptions struct {
	body   any
	header http.Header
}

// RequestOption 请求选项配置器
type RequestOption func(reqOpts *RequestOptions)

// SetBody 设置`HTTP`请求的主体内容
func SetBody(body any) (reqOpt RequestOption) {
	return func(reqOpts *RequestOptions) {
		reqOpts.body = body
	}
}

// SetContentType 设置`HTTP`请求头的`Content-Type`字段
func SetContentType(contentType string) (reqOpt RequestOption) {
	return func(reqOpts *RequestOptions) {
		reqOpts.header.Set("Content-Type", contentType)
	}
}

// SetCookie 设置`HTTP`请求头的`Cookie`字段
func SetCookie(cookies []*http.Cookie) (reqOpt RequestOption) {
	return func(reqOpts *RequestOptions) {
		cookieList := make([]string, 0, len(cookies))
		for _, v := range cookies {
			cookieList = append(cookieList, fmt.Sprintf("%s=%s", v.Name, v.Value))
		}
		reqOpts.header.Set("Cookie", strings.Join(cookieList, "; "))
	}
}

// SetKeyValue 设置`HTTP`请求头的键值对
func SetKeyValue(key, value string) (reqOpt RequestOption) {
	return func(reqOpts *RequestOptions) {
		reqOpts.header.Set(key, value)
	}
}

// AddKeyValue 添加`HTTP`请求头的键值对
func AddKeyValue(key, value string) (reqOpt RequestOption) {
	return func(reqOpts *RequestOptions) {
		reqOpts.header.Add(key, value)
	}
}

// RequestBuilder 请求构建器接口
type RequestBuilder interface {
	Build(ctx context.Context, method, url string, setters ...RequestOption) (req *http.Request, err error) // 构建器
}

// HTTPRequestBuilder `HTTP`请求构建器
type HTTPRequestBuilder struct {
	marshaller Marshaller
}

// NewRequestBuilder 新建`HTTP`请求构建器
func NewRequestBuilder() (hrb *HTTPRequestBuilder) {
	return &HTTPRequestBuilder{
		marshaller: &JSONMarshaller{},
	}
}

// Build 构建器
func (hrb *HTTPRequestBuilder) Build(ctx context.Context, method string, url string, setters ...RequestOption) (req *http.Request, err error) {
	reqOpts := &RequestOptions{
		body:   nil,
		header: make(http.Header),
	}
	for _, setter := range setters {
		setter(reqOpts)
	}

	var bodyReader io.Reader
	if reqOpts.body != nil {
		if v, ok := reqOpts.body.(io.Reader); ok {
			bodyReader = v
		} else {
			var reqBytes []byte
			if reqBytes, err = hrb.marshaller.Marshal(reqOpts.body); err != nil {
				return
			}
			bodyReader = bytes.NewBuffer(reqBytes)
		}
	}
	if req, err = http.NewRequestWithContext(ctx, method, url, bodyReader); err != nil {
		return
	}

	req.Header = reqOpts.header
	return
}
