/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-05-09 21:15:02
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-05-09 21:21:21
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkvalguard

// GetValue 返回指针指向的值或默认值
func GetValue[T any](ptr *T, defaultValues ...T) (val T) {
	if ptr != nil {
		return *ptr
	}
	if len(defaultValues) > 0 {
		return defaultValues[0]
	}
	var zeroValue T // 类型 T 的默认零值
	return zeroValue
}
