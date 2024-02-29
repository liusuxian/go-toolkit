/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-26 20:45:35
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-26 20:52:53
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import "encoding/json"

// Unmarshaler 反序列化接口
type Unmarshaler interface {
	Unmarshal(data []byte, v any) (err error) // 反序列化
}

// JSONUnmarshaler `JSON`反序列化
type JSONUnmarshaler struct{}

// Unmarshal 反序列化
func (jm *JSONUnmarshaler) Unmarshal(data []byte, v any) (err error) {
	return json.Unmarshal(data, v)
}
