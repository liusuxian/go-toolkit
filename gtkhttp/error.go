/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 14:45:31
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-23 18:50:11
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import (
	"errors"
	"fmt"
)

var (
	ErrTooManyEmptyStreamMessages = errors.New("stream has sent too many empty messages")
)

// RequestError 请求错误
type RequestError struct {
	HTTPStatus     string // HTTP 状态描述
	HTTPStatusCode int    // HTTP 状态码
	Err            error  // 错误信息
	Body           []byte // 响应体
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
