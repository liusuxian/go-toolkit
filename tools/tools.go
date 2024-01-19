/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-19 22:57:39
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-19 23:13:01
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package tools

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
