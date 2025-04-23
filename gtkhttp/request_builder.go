/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-07 18:30:21
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-23 18:51:27
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

// RequestBuilder 请求构建器接口
type RequestBuilder interface {
	Build(ctx context.Context, method, url string, body any, header http.Header) (req *http.Request, err error) // 构建器
}

// HTTPRequestBuilder HTTP 请求构建器
type HTTPRequestBuilder struct {
	marshaller Marshaller
}

// RequestBuilderOption HTTP 请求构建器选项
type RequestBuilderOption func(*HTTPRequestBuilder)

// WithMarshaller 设置 marshaller
func WithMarshaller(marshaller Marshaller) (option RequestBuilderOption) {
	return func(hrb *HTTPRequestBuilder) {
		hrb.marshaller = marshaller
	}
}

// NewRequestBuilder 新建 HTTP 请求构建器
func NewRequestBuilder(opts ...RequestBuilderOption) (hrb *HTTPRequestBuilder) {
	hrb = &HTTPRequestBuilder{
		marshaller: &JSONMarshaller{},
	}

	for _, opt := range opts {
		opt(hrb)
	}
	return
}

// Build 构建器
func (hrb *HTTPRequestBuilder) Build(ctx context.Context, method, url string, body any, header http.Header) (req *http.Request, err error) {
	var bodyReader io.Reader
	if body != nil {
		if v, ok := body.(io.Reader); ok {
			bodyReader = v
		} else {
			var reqBytes []byte
			if reqBytes, err = hrb.marshaller.Marshal(body); err != nil {
				return
			}
			bodyReader = bytes.NewBuffer(reqBytes)
		}
	}
	if req, err = http.NewRequestWithContext(ctx, method, url, bodyReader); err != nil {
		return
	}
	if header != nil {
		req.Header = header
	}
	return
}
