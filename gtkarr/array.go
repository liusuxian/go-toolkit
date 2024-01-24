/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-22 22:31:39
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-23 15:57:14
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkarr

// ContainsStr 字符串包含检测
func ContainsStr(s []string, e string) (ok bool) {
	if len(s) == 0 {
		return false
	}
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return
}

// ContainsRune 字符串包含检测
func ContainsRune(s []rune, e rune) (ok bool) {
	if len(s) == 0 {
		return false
	}
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return
}

// ContainsInt 整型包含检测
func ContainsInt(s []int, e int) (ok bool) {
	if len(s) == 0 {
		return false
	}
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return
}

// ContainsInt8 整型包含检测
func ContainsInt8(s []int8, e int8) (ok bool) {
	if len(s) == 0 {
		return false
	}
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return
}

// ContainsInt16 整型包含检测
func ContainsInt16(s []int16, e int16) (ok bool) {
	if len(s) == 0 {
		return false
	}
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return
}

// ContainsInt32 整型包含检测
func ContainsInt32(s []int32, e int32) (ok bool) {
	if len(s) == 0 {
		return false
	}
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return
}

// ContainsInt64 整型包含检测
func ContainsInt64(s []int64, e int64) (ok bool) {
	if len(s) == 0 {
		return false
	}
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return
}

// ContainsUint 整型包含检测
func ContainsUint(s []uint, e uint) (ok bool) {
	if len(s) == 0 {
		return false
	}
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return
}

// ContainsUint8 整型包含检测
func ContainsUint8(s []uint8, e uint8) (ok bool) {
	if len(s) == 0 {
		return false
	}
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return
}

// ContainsUint16 整型包含检测
func ContainsUint16(s []uint16, e uint16) (ok bool) {
	if len(s) == 0 {
		return false
	}
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return
}

// ContainsUint32 整型包含检测
func ContainsUint32(s []uint32, e uint32) (ok bool) {
	if len(s) == 0 {
		return false
	}
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return
}

// ContainsUint64 整型包含检测
func ContainsUint64(s []uint64, e uint64) (ok bool) {
	if len(s) == 0 {
		return false
	}
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return
}

// ContainsFloat32 浮点型包含检测
func ContainsFloat32(s []float32, e float32) (ok bool) {
	if len(s) == 0 {
		return false
	}
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return
}

// ContainsFloat64 浮点型包含检测
func ContainsFloat64(s []float64, e float64) (ok bool) {
	if len(s) == 0 {
		return false
	}
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return
}

// ContainsByte 字节类型包含检测
func ContainsByte(s []byte, e byte) (ok bool) {
	if len(s) == 0 {
		return false
	}
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return
}

// ContainsAny 任意类型包含检测
func ContainsAny(s []any, e any) (ok bool) {
	if len(s) == 0 {
		return false
	}
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return
}
