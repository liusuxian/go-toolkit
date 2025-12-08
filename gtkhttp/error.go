/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-28 17:56:51
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-09 00:21:42
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
)

var (
	ErrMethodNotSupported           = errors.New("method is not supported")                                                            // 方法不支持
	ErrCompletionStreamNotSupported = errors.New("streaming is not supported with this method, please use CreateChatCompletionStream") // 流式传输不支持
	ErrTooManyEmptyStreamMessages   = errors.New("stream has sent too many empty messages")                                            // 流式传输发送了太多空消息
	ErrStreamReturnIntervalTimeout  = errors.New("stream return interval timeout")                                                     // 流式传输返回间隔超时
)

// APIError API错误信息
type APIError struct {
	Code           any         `json:"code,omitempty"`
	Message        string      `json:"message"`
	RequestId      string      `json:"request_id,omitempty"`
	Param          *string     `json:"param,omitempty"`
	Type           string      `json:"type"`
	HTTPStatus     string      `json:"-"`
	HTTPStatusCode int         `json:"-"`
	InnerError     *InnerError `json:"innererror,omitempty"`
}

// InnerError 内部错误信息
type InnerError struct {
	Code                 string                `json:"code,omitempty"`
	ContentFilterResults *ContentFilterResults `json:"content_filter_result,omitempty"`
}

// Hate 仇恨内容过滤结果
type Hate struct {
	Filtered bool   `json:"filtered"`           // 是否过滤
	Severity string `json:"severity,omitempty"` // 严重程度
}

// SelfHarm 自残内容过滤结果
type SelfHarm struct {
	Filtered bool   `json:"filtered"`           // 是否过滤
	Severity string `json:"severity,omitempty"` // 严重程度
}

// Sexual 性内容过滤结果
type Sexual struct {
	Filtered bool   `json:"filtered"`           // 是否过滤
	Severity string `json:"severity,omitempty"` // 严重程度
}

// Violence 暴力内容过滤结果
type Violence struct {
	Filtered bool   `json:"filtered"`           // 是否过滤
	Severity string `json:"severity,omitempty"` // 严重程度
}

// JailBreak 越狱内容过滤结果
type JailBreak struct {
	Filtered bool `json:"filtered"` // 是否过滤
	Detected bool `json:"detected"` // 是否检测到
}

// Profanity 亵渎内容过滤结果
type Profanity struct {
	Filtered bool `json:"filtered"` // 是否过滤
	Detected bool `json:"detected"` // 是否检测到
}

// ContentFilterResults 内容过滤结果
type ContentFilterResults struct {
	Hate      *Hate      `json:"hate,omitempty"`      // 仇恨内容过滤结果
	SelfHarm  *SelfHarm  `json:"self_harm,omitempty"` // 自残内容过滤结果
	Sexual    *Sexual    `json:"sexual,omitempty"`    // 性内容过滤结果
	Violence  *Violence  `json:"violence,omitempty"`  // 暴力内容过滤结果
	JailBreak *JailBreak `json:"jailbreak,omitempty"` // 越狱内容过滤结果
	Profanity *Profanity `json:"profanity,omitempty"` // 亵渎内容过滤结果
}

// RequestError 请求错误
type RequestError struct {
	HTTPStatus     string // HTTP 状态描述
	HTTPStatusCode int    // HTTP 状态码
	Err            error  // 错误信息
	Body           []byte // 响应体
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error *APIError `json:"error,omitempty"` // 错误信息
}

// Error 实现 error 接口的方法
func (e *APIError) Error() (s string) {
	if e.HTTPStatusCode > 0 {
		if e.RequestId != "" {
			return fmt.Sprintf("error, status code: %d, status: %s, message: %s, request id: %s", e.HTTPStatusCode, e.HTTPStatus, e.Message, e.RequestId)
		}
		return fmt.Sprintf("error, status code: %d, status: %s, message: %s", e.HTTPStatusCode, e.HTTPStatus, e.Message)
	}
	return e.Message
}

// UnmarshalJSON 反序列化JSON
func (e *APIError) UnmarshalJSON(data []byte) (err error) {
	var rawMap map[string]json.RawMessage
	if err = json.Unmarshal(data, &rawMap); err != nil {
		return
	}

	if err = json.Unmarshal(rawMap["message"], &e.Message); err != nil {
		var messages []string
		if err = json.Unmarshal(rawMap["message"], &messages); err != nil {
			return
		}
		e.Message = strings.Join(messages, ", ")
	}

	if _, ok := rawMap["type"]; ok {
		if err = json.Unmarshal(rawMap["type"], &e.Type); err != nil {
			return
		}
	}

	if _, ok := rawMap["innererror"]; ok {
		if err = json.Unmarshal(rawMap["innererror"], &e.InnerError); err != nil {
			return
		}
	}

	if _, ok := rawMap["request_id"]; ok {
		if err = json.Unmarshal(rawMap["request_id"], &e.RequestId); err != nil {
			return
		}
	}

	if _, ok := rawMap["param"]; ok {
		if err = json.Unmarshal(rawMap["param"], &e.Param); err != nil {
			return
		}
	}

	if _, ok := rawMap["code"]; !ok {
		return nil
	}

	var intCode int
	if err = json.Unmarshal(rawMap["code"], &intCode); err == nil {
		e.Code = intCode
		return nil
	}

	return json.Unmarshal(rawMap["code"], &e.Code)
}

// Error 实现 error 接口的方法
func (e *RequestError) Error() (s string) {
	return fmt.Sprintf(
		"error, status code: %d, status: %s, message: %s, body: %s",
		e.HTTPStatusCode, e.HTTPStatus, e.Err, e.Body,
	)
}

// Unwrap 解包错误
func (e *RequestError) Unwrap() (err error) {
	return e.Err
}

// IsMethodNotSupportedError 判断是否是方法不支持错误
func IsMethodNotSupportedError(err error) (is bool) {
	return errors.Is(err, ErrMethodNotSupported)
}

// IsCompletionStreamNotSupportedError 判断是否是流式传输不支持错误
func IsCompletionStreamNotSupportedError(err error) (is bool) {
	return errors.Is(err, ErrCompletionStreamNotSupported)
}

// IsTooManyEmptyStreamMessagesError 判断是否是流式传输发送了太多空消息错误
func IsTooManyEmptyStreamMessagesError(err error) (is bool) {
	return errors.Is(err, ErrTooManyEmptyStreamMessages)
}

// IsStreamReturnIntervalTimeoutError 判断是否是流式传输返回间隔超时错误
func IsStreamReturnIntervalTimeoutError(err error) (is bool) {
	return errors.Is(err, ErrStreamReturnIntervalTimeout)
}

// IsCanceledError 判断是否是取消错误
func IsCanceledError(err error) (is bool) {
	return errors.Is(err, context.Canceled)
}

// IsDeadlineExceededError 判断是否是截止时间错误
func IsDeadlineExceededError(err error) (is bool) {
	return errors.Is(err, context.DeadlineExceeded)
}

// IsNetError 判断是否是网络错误
func IsNetError(err error) (is bool) {
	var netErr net.Error
	return errors.As(err, &netErr)
}

// ClientError 客户端错误
type ClientError struct {
	RequestID string // 请求ID
	Err       error  // 原始错误
}

// Error 错误信息
func (e *ClientError) Error() (errStr string) {
	return fmt.Sprintf("request_id: %s, error: %v", e.RequestID, e.Err)
}

// RequestID 获取请求ID
func RequestID(err error) (requestId string) {
	if err == nil {
		return ""
	}

	var clientErr *ClientError
	if errors.As(err, &clientErr) {
		return clientErr.RequestID
	}

	return ""
}

// Unwrap 解包错误
func Unwrap(err error) (originalError error) {
	if err == nil {
		return nil
	}
	// 解包 ClientError
	var clientErr *ClientError
	if errors.As(err, &clientErr) {
		if clientErr.Err != nil {
			return clientErr.Err
		}
		return err // 如果内部错误为 nil，返回 ClientError 本身
	}
	// 解包 RequestError
	var requestError *RequestError
	if errors.As(err, &requestError) {
		if requestError.Err != nil {
			return requestError.Err
		}
		return err // 如果内部错误为 nil，返回 RequestError 本身
	}
	// 其他类型的错误
	unwrapped := errors.Unwrap(err)
	if unwrapped == nil {
		return err // 已经是最底层错误，返回原错误
	}
	return unwrapped
}

// Cause 错误根因
func Cause(err error) (causeError error) {
	return doCause(err)
}

// doCause 递归获取错误根因
func doCause(err error) (causeError error) {
	if err == nil {
		return nil
	}
	// 解包错误
	unwrapped := Unwrap(err)
	if unwrapped == nil {
		return err // 已经到达最底层错误，返回当前错误
	}
	// 防止无限递归：如果解包后的错误与原错误相同，直接返回
	if unwrapped == err {
		return err
	}
	return doCause(unwrapped)
}
