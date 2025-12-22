/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-22 23:27:03
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-23 00:11:51
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkjson

import (
	"bytes"
	"encoding/json"
	"strings"
)

// JsonMarshal 将 any 转换为 json 字符串，不对 HTML 字符进行转义
func JsonMarshal(v any) (jsonStr string, err error) {
	if v == nil {
		return
	}
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	if err = encoder.Encode(v); err != nil {
		return
	}
	// 去掉 Encode 自动添加的换行符
	jsonStr = strings.TrimRight(buffer.String(), "\n")
	return
}

// MustJsonMarshal 将 any 转换为 json 字符串，不对 HTML 字符进行转义
func MustJsonMarshal(v any) (jsonStr string) {
	jsonStr, _ = JsonMarshal(v)
	return
}

// String 将 any 转换为 json 字符串（会转义 HTML 字符）
func String(v any) (str string, err error) {
	if v == nil {
		return
	}
	var b []byte
	if b, err = json.Marshal(v); err != nil {
		return
	}
	str = string(b)
	return
}

// MustString 将 any 转换为 json 字符串（会转义 HTML 字符）
func MustString(v any) (str string) {
	str, _ = String(v)
	return
}

// Bytes 将 any 转换为 json 字节切片（会转义 HTML 字符）
func Bytes(v any) (b []byte, err error) {
	if v == nil {
		return
	}
	b, err = json.Marshal(v)
	return
}

// MustBytes 将 any 转换为 json 字节切片（会转义 HTML 字符）
func MustBytes(v any) (b []byte) {
	b, _ = Bytes(v)
	return
}
