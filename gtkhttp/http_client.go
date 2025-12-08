/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-28 17:56:51
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-08 22:49:10
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultHTTPClientTimeout           time.Duration = 5 * time.Second        // 默认HTTP客户端请求超时时间
	defaultStreamReturnIntervalTimeout time.Duration = 10 * time.Second       // 默认流式传输返回的间隔超时时间
	defaultEmptyMessagesLimit          uint          = 300                    // 默认空消息限制
	requestIdHeaderKey                 string        = "X-Toolkit-Request-Id" // 请求ID头键
)

// HTTPDoer HTTP 请求执行器接口
type HTTPDoer interface {
	SetTimeout(timeout time.Duration)                      // 设置请求超时时间，零值表示无超时限制
	Do(req *http.Request) (resp *http.Response, err error) // 发送请求
}

// DefaultHTTPDoer 默认 HTTP 请求执行器
type DefaultHTTPDoer struct {
	client *http.Client // 底层 HTTP 客户端
}

// NewDefaultHTTPDoer 新建默认 HTTP 请求执行器
//
//	如果 timeout 为 0，则表示无超时限制
func NewDefaultHTTPDoer(timeout time.Duration) (doer *DefaultHTTPDoer) {
	if timeout <= 0 {
		timeout = 0
	}
	return &DefaultHTTPDoer{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// SetTimeout 设置请求超时时间，零值表示无超时限制（非并发安全）
func (doer *DefaultHTTPDoer) SetTimeout(timeout time.Duration) {
	if timeout <= 0 {
		timeout = 0
	}
	doer.client.Timeout = timeout
}

// Do 发送请求
func (doer *DefaultHTTPDoer) Do(req *http.Request) (resp *http.Response, err error) {
	return doer.client.Do(req)
}

// ResponseDecoder 响应数据解码器接口
type ResponseDecoder interface {
	Decode(body io.Reader, v any) (err error) // 解码响应数据
}

// DefaultResponseDecoder 默认响应数据解码器
type DefaultResponseDecoder struct{}

// Decode 解码响应数据
func (d *DefaultResponseDecoder) Decode(body io.Reader, v any) (err error) {
	if v == nil {
		return
	}

	switch o := v.(type) {
	case *string:
		return decodeString(body, o)
	default:
		return json.NewDecoder(body).Decode(v)
	}
}

// decodeString 解码字符串
func decodeString(body io.Reader, output *string) (err error) {
	var b []byte
	if b, err = io.ReadAll(body); err != nil {
		return
	}

	*output = string(b)
	return
}

// HTTPClientConfig 客户端配置
type HTTPClientConfig struct {
	BaseURL                     string          // API 的基础 URL 地址
	HTTPClient                  HTTPDoer        // HTTP 客户端实现，用于发送请求
	ResponseDecoder             ResponseDecoder // 响应数据解码器
	EmptyMessagesLimit          uint            // 空消息限制
	StreamReturnIntervalTimeout time.Duration   // 流式传输返回的间隔超时时间，用于防止流式传输返回的间隔时间过长
}

// HTTPClient 客户端
type HTTPClient struct {
	config            HTTPClientConfig                 // 客户端配置
	requestBuilder    RequestBuilder                   // 请求构建器
	createFormBuilder func(body io.Writer) FormBuilder // 表单构建器
}

// HTTPClientOption 客户端选项
type HTTPClientOption func(c *HTTPClient)

// WithTimeout 设置请求超时时间（非并发安全）
func WithTimeout(timeout time.Duration) (opt HTTPClientOption) {
	return func(c *HTTPClient) {
		c.config.HTTPClient.SetTimeout(timeout)
	}
}

// WithStreamReturnIntervalTimeout 设置流式传输返回的间隔超时时间（非并发安全）
func WithStreamReturnIntervalTimeout(timeout time.Duration) (opt HTTPClientOption) {
	return func(c *HTTPClient) {
		c.config.StreamReturnIntervalTimeout = timeout
	}
}

// Response 响应
type Response interface {
	SetHeader(http.Header)
}

// HttpHeader 请求头
type HttpHeader http.Header

// SetHeader 设置请求头
func (h *HttpHeader) SetHeader(header http.Header) {
	*h = HttpHeader(header)
}

// Header 获取请求头
func (h *HttpHeader) Header() (header http.Header) {
	return http.Header(*h)
}

// RequestID 获取请求ID
func (h *HttpHeader) RequestID() (requestID string) {
	return h.Header().Get(requestIdHeaderKey)
}

// RawResponse 原始响应
type RawResponse struct {
	io.ReadCloser
	HttpHeader
}

// NewHTTPClient 新建 HTTP 客户端
func NewHTTPClient(baseURL string, opts ...RequestBuilderOption) (c *HTTPClient) {
	return NewHTTPClientWithConfig(HTTPClientConfig{
		BaseURL:                     baseURL,
		HTTPClient:                  NewDefaultHTTPDoer(defaultHTTPClientTimeout),
		ResponseDecoder:             &DefaultResponseDecoder{},
		EmptyMessagesLimit:          defaultEmptyMessagesLimit,
		StreamReturnIntervalTimeout: defaultStreamReturnIntervalTimeout,
	}, opts...)
}

// NewHTTPClientWithConfig 通过客户端配置新建 HTTP 客户端
func NewHTTPClientWithConfig(config HTTPClientConfig, opts ...RequestBuilderOption) (c *HTTPClient) {
	return &HTTPClient{
		config:         config,
		requestBuilder: NewRequestBuilder(opts...),
		createFormBuilder: func(body io.Writer) FormBuilder {
			return NewFormBuilder(body)
		},
	}
}

// RequestOptions 请求选项
type RequestOptions struct {
	body   any
	header http.Header
}

// RequestOption 请求选项配置器
type RequestOption func(reqOpts *RequestOptions)

// WithBody 设置 HTTP 请求的主体内容
func WithBody(body any) (reqOpt RequestOption) {
	return func(reqOpts *RequestOptions) {
		reqOpts.body = body
	}
}

// WithContentType 设置 HTTP 请求头的 Content-Type 字段
func WithContentType(contentType string) (reqOpt RequestOption) {
	return func(reqOpts *RequestOptions) {
		reqOpts.header.Set("Content-Type", contentType)
	}
}

// WithCookie 设置 HTTP 请求头的 Cookie 字段
func WithCookie(cookies []*http.Cookie) (reqOpt RequestOption) {
	return func(reqOpts *RequestOptions) {
		cookieList := make([]string, 0, len(cookies))
		for _, v := range cookies {
			cookieList = append(cookieList, fmt.Sprintf("%s=%s", v.Name, v.Value))
		}
		reqOpts.header.Set("Cookie", strings.Join(cookieList, "; "))
	}
}

// WithKeyValue 设置 HTTP 请求头的键值对
func WithKeyValue(key, value string) (reqOpt RequestOption) {
	return func(reqOpts *RequestOptions) {
		reqOpts.header.Set(key, value)
	}
}

// GetFormBuilder 获取表单构建器
func (c *HTTPClient) GetFormBuilder(body io.Writer) (builder FormBuilder) {
	return c.createFormBuilder(body)
}

// NewRequest 新建请求
func (c *HTTPClient) NewRequest(ctx context.Context, method, url string, setters ...RequestOption) (req *http.Request, err error) {
	reqOpts := &RequestOptions{
		body:   nil,
		header: make(http.Header),
	}
	for _, setter := range setters {
		setter(reqOpts)
	}

	if req, err = c.requestBuilder.Build(ctx, method, url, reqOpts.body, reqOpts.header); err != nil {
		return
	}
	return
}

// SendRequest 发送请求
func (c *HTTPClient) SendRequest(req *http.Request, v Response) (err error) {
	// 设置默认请求头
	for _, v := range [][]string{
		{"Content-Type", "application/json"},
		{"Accept", "application/json"},
	} {
		if req.Header.Get(v[0]) == "" {
			req.Header.Set(v[0], v[1])
		}
	}

	var resp *http.Response
	if resp, err = c.config.HTTPClient.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	if v != nil {
		c.setRequestID(req, resp)
		v.SetHeader(resp.Header)
	}

	if isFailureStatusCode(resp) {
		return c.handleErrorResp(resp)
	}

	return c.config.ResponseDecoder.Decode(resp.Body, v)
}

// SendRequestRaw 发送请求
func (c *HTTPClient) SendRequestRaw(req *http.Request) (response RawResponse, err error) {
	var resp *http.Response
	if resp, err = c.config.HTTPClient.Do(req); err != nil {
		return
	}

	if isFailureStatusCode(resp) {
		err = c.handleErrorResp(resp)
		return
	}

	c.setRequestID(req, resp)
	response.SetHeader(resp.Header)
	response.ReadCloser = resp.Body
	return
}

// SendRequestStream 发送流式请求
func SendRequestStream[T Streamable](client *HTTPClient, req *http.Request) (stream *StreamReader[T], err error) {
	// 设置默认请求头
	for _, v := range [][]string{
		{"Content-Type", "application/json"},
		{"Accept", "text/event-stream"},
		{"Cache-Control", "no-cache"},
		{"Connection", "keep-alive"},
	} {
		if req.Header.Get(v[0]) == "" {
			req.Header.Set(v[0], v[1])
		}
	}

	var resp *http.Response
	if resp, err = client.config.HTTPClient.Do(req); err != nil {
		stream = &StreamReader[T]{}
		return
	}

	if isFailureStatusCode(resp) {
		stream = &StreamReader[T]{}
		err = client.handleErrorResp(resp)
		return
	}

	client.setRequestID(req, resp)
	stream = &StreamReader[T]{
		emptyMessagesLimit:          client.config.EmptyMessagesLimit,
		reader:                      bufio.NewReader(resp.Body),
		response:                    resp,
		streamReturnIntervalTimeout: client.config.StreamReturnIntervalTimeout,
		errAccumulator:              NewErrorAccumulator(),
		responseDecoder:             client.config.ResponseDecoder,
		startTime:                   time.Now(),
		HttpHeader:                  HttpHeader(resp.Header),
	}
	return
}

// isFailureStatusCode 是否失败状态码
func isFailureStatusCode(resp *http.Response) (ok bool) {
	return resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest
}

// FullURL 获取完整 URL
func (c *HTTPClient) FullURL(suffix string) (url string) {
	baseURL := strings.TrimRight(c.config.BaseURL, "/")
	url = fmt.Sprintf("%s%s", baseURL, suffix)
	return
}

// setRequestID 设置请求ID
func (c *HTTPClient) setRequestID(request *http.Request, response *http.Response) {
	requestInfo := GetRequestInfo(request.Context())
	if requestInfo.RequestID != "" && requestInfo.RequestID != "unknown" {
		response.Header.Set(requestIdHeaderKey, requestInfo.RequestID)
	}
}

// handleErrorResp 处理错误响应
func (c *HTTPClient) handleErrorResp(resp *http.Response) (err error) {
	// 读取响应体
	var body []byte
	if body, err = io.ReadAll(resp.Body); err != nil {
		return fmt.Errorf("error, reading response body: %w", err)
	}
	// 尝试解析为 ErrorResponse
	var errRes ErrorResponse
	if err = json.Unmarshal(body, &errRes); err == nil && errRes.Error != nil {
		errRes.Error.HTTPStatus = resp.Status
		errRes.Error.HTTPStatusCode = resp.StatusCode
		return errRes.Error
	}
	// 尝试解析为 APIError
	var apiErr *APIError
	if err = json.Unmarshal(body, &apiErr); err == nil && apiErr != nil {
		apiErr.HTTPStatus = resp.Status
		apiErr.HTTPStatusCode = resp.StatusCode
		return apiErr
	}
	// 如果都解析失败，返回包含解析错误的 RequestError
	return &RequestError{
		HTTPStatus:     resp.Status,
		HTTPStatusCode: resp.StatusCode,
		Err:            fmt.Errorf("failed to parse error response"),
		Body:           body,
	}
}
