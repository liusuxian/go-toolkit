/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-07 18:29:02
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-23 18:51:06
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import "encoding/json"

// Marshaller 序列化接口
type Marshaller interface {
	Marshal(val any) (b []byte, err error) // 序列化
}

// JSONMarshaller JSON 序列化
type JSONMarshaller struct{}

// Marshal 序列化
func (jm *JSONMarshaller) Marshal(val any) (b []byte, err error) {
	return json.Marshal(val)
}
