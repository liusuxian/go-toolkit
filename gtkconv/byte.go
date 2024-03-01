/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-03-01 16:52:17
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-03-01 22:32:07
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkconv

import (
	"encoding/json"
	"github.com/liusuxian/go-toolkit/gtkbinary"
	"github.com/liusuxian/go-toolkit/gtkreflect"
	"reflect"
)

// ToByteE 将 any 转换为 byte 类型
func ToByteE(i any) (iv byte, err error) {
	if val, ok := i.(byte); ok {
		return val, nil
	}
	return ToUint8E(i)
}

// ToRuneE 将 any 转换为 rune 类型
func ToRuneE(i any) (iv rune, err error) {
	if val, ok := i.(rune); ok {
		return val, nil
	}
	return ToInt32E(i)
}

// ToRunesE 将 any 转换为 []rune 类型
func ToRunesE(i any) (iv []rune, err error) {
	if val, ok := i.([]rune); ok {
		return val, nil
	}
	if v, e := ToStringE(i); e != nil {
		return []rune{}, nil
	} else {
		return []rune(v), nil
	}
}

// ToBytesE 将 any 转换为 []byte 类型
func ToBytesE(i any) (iv []byte, err error) {
	if i == nil {
		return []byte{}, nil
	}

	switch val := i.(type) {
	case string:
		return []byte(val), nil
	case []byte:
		return val, nil
	default:
		originValueAndKind := gtkreflect.OriginValueAndKind(i)
		switch originValueAndKind.OriginKind {
		case reflect.Map, reflect.Struct:
			return json.Marshal(i)
		case reflect.Array, reflect.Slice:
			length := originValueAndKind.OriginValue.Len()
			iv = make([]byte, length)
			for j := 0; j < length; j++ {
				var intv uint8
				if intv, err = ToUint8E(originValueAndKind.OriginValue.Index(j).Interface()); err != nil {
					return []byte{}, convertError(i, "[]byte")
				}
				iv[j] = intv
			}
			return iv, nil
		}
		return gtkbinary.Encode(i), nil
	}
}
