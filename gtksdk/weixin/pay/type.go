/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-13 11:21:19
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-13 11:21:21
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package pay

import "time"

// Time 复制 time.Time 对象，并返回复制体的指针
func Time(t time.Time) *time.Time {
	return &t
}

// String 复制 string 对象，并返回复制体的指针
func String(s string) *string {
	return &s
}

// Bool 复制 bool 对象，并返回复制体的指针
func Bool(b bool) *bool {
	return &b
}

// Float64 复制 float64 对象，并返回复制体的指针
func Float64(f float64) *float64 {
	return &f
}

// Float32 复制 float32 对象，并返回复制体的指针
func Float32(f float32) *float32 {
	return &f
}

// Int64 复制 int64 对象，并返回复制体的指针
func Int64(i int64) *int64 {
	return &i
}

// Int32 复制 int64 对象，并返回复制体的指针
func Int32(i int32) *int32 {
	return &i
}
