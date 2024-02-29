/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-26 21:32:04
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-26 23:26:10
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import (
	"encoding/json"
	"io"
	"net/http"
)

// IsFailureStatusCode 是否失败状态码
func IsFailureStatusCode(resp *http.Response) (ok bool) {
	return resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest
}

// DecodeString 解码字符串
func DecodeString(body io.Reader, output *string) (err error) {
	var b []byte
	if b, err = io.ReadAll(body); err != nil {
		return
	}
	*output = string(b)
	return
}

// DecodeResponse 解码响应数据
func DecodeResponse(body io.Reader, v any) (err error) {
	if v == nil {
		return
	}
	if result, ok := v.(*string); ok {
		return DecodeString(body, result)
	}
	return json.NewDecoder(body).Decode(v)
}
