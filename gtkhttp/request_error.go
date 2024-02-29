/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-26 23:46:02
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-26 23:57:33
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import "fmt"

// RequestError 请求错误
type RequestError struct {
	HTTPStatusCode int
	Err            error
}

// Error
func (re *RequestError) Error() (text string) {
	return fmt.Sprintf("error, status code: %d, message: %v", re.HTTPStatusCode, re.Err)
}

// Unwrap
func (re *RequestError) Unwrap() (err error) {
	return re.Err
}
