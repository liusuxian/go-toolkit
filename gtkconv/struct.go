/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-04-18 17:32:08
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-13 13:50:18
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkconv

import (
	"encoding/json"
	"fmt"
)

// ToStructE 将 any 转换为 struct/[]struct 类型
func ToStructE(input, output any, opts ...DecoderConfigOption) (err error) {
	if input == nil {
		return nil
	}
	if output == nil {
		return fmt.Errorf("object output cannot be nil")
	}

	switch val := input.(type) {
	case []byte:
		// 如果它是 JSON 字符串，自动反序列化它
		if json.Valid(val) {
			if e := json.Unmarshal(val, output); e == nil {
				return
			}
		}
	case string:
		// 如果它是 JSON 字符串，自动反序列化它
		anyBytes := []byte(val)
		if json.Valid(anyBytes) {
			if e := json.Unmarshal(anyBytes, output); e == nil {
				return
			}
		}
	}

	return decode(input, defaultDecoderConfig(output, opts...))
}
