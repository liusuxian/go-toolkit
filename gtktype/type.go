/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-14 15:13:27
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-14 15:22:46
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtktype

import "time"

// String 复制 string 对象，并返回复制体的指针
func String(a string) *string {
	return &a
}

// StringValue 获取 string 对象的值
func StringValue(a *string) string {
	if a == nil {
		return ""
	}
	return *a
}

// Int 复制 int 对象，并返回复制体的指针
func Int(a int) *int {
	return &a
}

// IntValue 获取 int 对象的值
func IntValue(a *int) int {
	if a == nil {
		return 0
	}
	return *a
}

// Int8 复制 int8 对象，并返回复制体的指针
func Int8(a int8) *int8 {
	return &a
}

// Int8Value 获取 int8 对象的值
func Int8Value(a *int8) int8 {
	if a == nil {
		return 0
	}
	return *a
}

// Int16 复制 int16 对象，并返回复制体的指针
func Int16(a int16) *int16 {
	return &a
}

// Int16Value 获取 int16 对象的值
func Int16Value(a *int16) int16 {
	if a == nil {
		return 0
	}
	return *a
}

// Int32 复制 int32 对象，并返回复制体的指针
func Int32(a int32) *int32 {
	return &a
}

// Int32Value 获取 int32 对象的值
func Int32Value(a *int32) int32 {
	if a == nil {
		return 0
	}
	return *a
}

// Int64 复制 int64 对象，并返回复制体的指针
func Int64(a int64) *int64 {
	return &a
}

// Int64Value 获取 int64 对象的值
func Int64Value(a *int64) int64 {
	if a == nil {
		return 0
	}
	return *a
}

// Bool 复制 bool 对象，并返回复制体的指针
func Bool(a bool) *bool {
	return &a
}

// BoolValue 获取 bool 对象的值
func BoolValue(a *bool) bool {
	if a == nil {
		return false
	}
	return *a
}

// Uint 复制 uint 对象，并返回复制体的指针
func Uint(a uint) *uint {
	return &a
}

// UintValue 获取 uint 对象的值
func UintValue(a *uint) uint {
	if a == nil {
		return 0
	}
	return *a
}

// Uint8 复制 uint8 对象，并返回复制体的指针
func Uint8(a uint8) *uint8 {
	return &a
}

// Uint8Value 获取 uint8 对象的值
func Uint8Value(a *uint8) uint8 {
	if a == nil {
		return 0
	}
	return *a
}

// Uint16 复制 uint16 对象，并返回复制体的指针
func Uint16(a uint16) *uint16 {
	return &a
}

// Uint16Value 获取 uint16 对象的值
func Uint16Value(a *uint16) uint16 {
	if a == nil {
		return 0
	}
	return *a
}

// Uint32 复制 uint32 对象，并返回复制体的指针
func Uint32(a uint32) *uint32 {
	return &a
}

// Uint32Value 获取 uint32 对象的值
func Uint32Value(a *uint32) uint32 {
	if a == nil {
		return 0
	}
	return *a
}

// Uint64 复制 uint64 对象，并返回复制体的指针
func Uint64(a uint64) *uint64 {
	return &a
}

// Uint64Value 获取 uint64 对象的值
func Uint64Value(a *uint64) uint64 {
	if a == nil {
		return 0
	}
	return *a
}

// Float32 复制 float32 对象，并返回复制体的指针
func Float32(a float32) *float32 {
	return &a
}

// Float32Value 获取 float32 对象的值
func Float32Value(a *float32) float32 {
	if a == nil {
		return 0
	}
	return *a
}

// Float64 复制 float64 对象，并返回复制体的指针
func Float64(a float64) *float64 {
	return &a
}

// Float64Value 获取 float64 对象的值
func Float64Value(a *float64) float64 {
	if a == nil {
		return 0
	}
	return *a
}

// IntSlice 复制 int 对象的切片，并返回复制体的指针
func IntSlice(a []int) []*int {
	if a == nil {
		return nil
	}
	res := make([]*int, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// IntValueSlice 获取 int 对象的切片
func IntValueSlice(a []*int) []int {
	if a == nil {
		return nil
	}
	res := make([]int, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// Int8Slice 复制 int8 对象的切片，并返回复制体的指针
func Int8Slice(a []int8) []*int8 {
	if a == nil {
		return nil
	}
	res := make([]*int8, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// Int8ValueSlice 获取 int8 对象的切片
func Int8ValueSlice(a []*int8) []int8 {
	if a == nil {
		return nil
	}
	res := make([]int8, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// Int16Slice 复制 int16 对象的切片，并返回复制体的指针
func Int16Slice(a []int16) []*int16 {
	if a == nil {
		return nil
	}
	res := make([]*int16, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// Int16ValueSlice 获取 int16 对象的切片
func Int16ValueSlice(a []*int16) []int16 {
	if a == nil {
		return nil
	}
	res := make([]int16, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// Int32Slice 复制 int32 对象的切片，并返回复制体的指针
func Int32Slice(a []int32) []*int32 {
	if a == nil {
		return nil
	}
	res := make([]*int32, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// Int32ValueSlice 获取 int32 对象的切片
func Int32ValueSlice(a []*int32) []int32 {
	if a == nil {
		return nil
	}
	res := make([]int32, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// Int64Slice 复制 int64 对象的切片，并返回复制体的指针
func Int64Slice(a []int64) []*int64 {
	if a == nil {
		return nil
	}
	res := make([]*int64, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// Int64ValueSlice 获取 int64 对象的切片
func Int64ValueSlice(a []*int64) []int64 {
	if a == nil {
		return nil
	}
	res := make([]int64, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// UintSlice 复制 uint 对象的切片，并返回复制体的指针
func UintSlice(a []uint) []*uint {
	if a == nil {
		return nil
	}
	res := make([]*uint, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// UintValueSlice 获取 uint 对象的切片
func UintValueSlice(a []*uint) []uint {
	if a == nil {
		return nil
	}
	res := make([]uint, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// Uint8Slice 复制 uint8 对象的切片，并返回复制体的指针
func Uint8Slice(a []uint8) []*uint8 {
	if a == nil {
		return nil
	}
	res := make([]*uint8, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// Uint8ValueSlice 获取 uint8 对象的切片
func Uint8ValueSlice(a []*uint8) []uint8 {
	if a == nil {
		return nil
	}
	res := make([]uint8, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// Uint16Slice 复制 uint16 对象的切片，并返回复制体的指针
func Uint16Slice(a []uint16) []*uint16 {
	if a == nil {
		return nil
	}
	res := make([]*uint16, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// Uint16ValueSlice 获取 uint16 对象的切片
func Uint16ValueSlice(a []*uint16) []uint16 {
	if a == nil {
		return nil
	}
	res := make([]uint16, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// Uint32Slice 复制 uint32 对象的切片，并返回复制体的指针
func Uint32Slice(a []uint32) []*uint32 {
	if a == nil {
		return nil
	}
	res := make([]*uint32, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// Uint32ValueSlice 获取 uint32 对象的切片
func Uint32ValueSlice(a []*uint32) []uint32 {
	if a == nil {
		return nil
	}
	res := make([]uint32, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// Uint64Slice 复制 uint64 对象的切片，并返回复制体的指针
func Uint64Slice(a []uint64) []*uint64 {
	if a == nil {
		return nil
	}
	res := make([]*uint64, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// Uint64ValueSlice 获取 uint64 对象的切片
func Uint64ValueSlice(a []*uint64) []uint64 {
	if a == nil {
		return nil
	}
	res := make([]uint64, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// Float32Slice 复制 float32 对象的切片，并返回复制体的指针
func Float32Slice(a []float32) []*float32 {
	if a == nil {
		return nil
	}
	res := make([]*float32, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// Float32ValueSlice 获取 float32 对象的切片
func Float32ValueSlice(a []*float32) []float32 {
	if a == nil {
		return nil
	}
	res := make([]float32, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// Float64Slice 复制 float64 对象的切片，并返回复制体的指针
func Float64Slice(a []float64) []*float64 {
	if a == nil {
		return nil
	}
	res := make([]*float64, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// Float64ValueSlice 获取 float64 对象的切片
func Float64ValueSlice(a []*float64) []float64 {
	if a == nil {
		return nil
	}
	res := make([]float64, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// StringSlice 复制 string 对象的切片，并返回复制体的指针
func StringSlice(a []string) []*string {
	if a == nil {
		return nil
	}
	res := make([]*string, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// StringSliceValue 获取 string 对象的切片
func StringSliceValue(a []*string) []string {
	if a == nil {
		return nil
	}
	res := make([]string, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// BoolSlice 复制 bool 对象的切片，并返回复制体的指针
func BoolSlice(a []bool) []*bool {
	if a == nil {
		return nil
	}
	res := make([]*bool, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// BoolSliceValue 获取 bool 对象的切片
func BoolSliceValue(a []*bool) []bool {
	if a == nil {
		return nil
	}
	res := make([]bool, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

// Time 复制 time.Time 对象，并返回复制体的指针
func Time(t time.Time) *time.Time {
	return &t
}

// TimeValue 获取 time.Time 对象的值
func TimeValue(a *time.Time) time.Time {
	if a == nil {
		return time.Time{}
	}
	return *a
}

// TimeSlice 复制 time.Time 对象的切片，并返回复制体的指针
func TimeSlice(a []time.Time) []*time.Time {
	if a == nil {
		return nil
	}
	res := make([]*time.Time, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

// TimeSliceValue 获取 time.Time 对象的切片
func TimeSliceValue(a []*time.Time) []time.Time {
	if a == nil {
		return nil
	}
	res := make([]time.Time, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}
