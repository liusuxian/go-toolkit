/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-19 22:57:39
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-20 00:04:07
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package tools

import (
	"bytes"
	"encoding/json"
	"strings"
)

// JsonMarshal
func JsonMarshal(v any) (jsonStr string, err error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(v)
	jsonStr = buffer.String()
	return
}

// MustJsonMarshal
func MustJsonMarshal(v any) (jsonStr string) {
	jsonStr, _ = JsonMarshal(v)
	return
}

// String 结构体或切片字符串输出
func String(v any) (str string, err error) {
	if v == nil {
		return
	}
	var b []byte
	b, err = json.Marshal(v)
	str = string(b)
	return
}

// MustString 结构体或切片字符串输出
func MustString(v any) (str string) {
	str, _ = String(v)
	return
}

// Bytes 结构体或切片字节流输出
func Bytes(v any) (b []byte, err error) {
	b, err = json.Marshal(v)
	return
}

// MustBytes 结构体或切片字节流输出
func MustBytes(v any) (b []byte) {
	b, _ = Bytes(v)
	return
}

// ContainsStr 字符串数据包含检测
func ContainsStr(s []string, ele string) (ok bool) {
	if len(s) == 0 {
		return
	}
	for _, v := range s {
		if strings.Compare(v, ele) == 0 {
			ok = true
			return
		}
	}
	return
}
