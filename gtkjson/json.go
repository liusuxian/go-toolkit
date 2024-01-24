/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-22 23:27:03
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-22 23:29:35
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkjson

import (
	"bytes"
	"encoding/json"
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

// String
func String(v any) (str string, err error) {
	if v == nil {
		return
	}
	var b []byte
	b, err = json.Marshal(v)
	str = string(b)
	return
}

// MustString
func MustString(v any) (str string) {
	str, _ = String(v)
	return
}

// Bytes
func Bytes(v any) (b []byte, err error) {
	b, err = json.Marshal(v)
	return
}

// MustBytes
func MustBytes(v any) (b []byte) {
	b, _ = Bytes(v)
	return
}
