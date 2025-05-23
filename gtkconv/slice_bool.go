/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-04-16 02:20:09
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-13 14:46:08
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkconv

import (
	"encoding/json"
	"github.com/liusuxian/go-toolkit/internal/utils"
	"reflect"
)

// ToBoolSliceE  将 any 转换为 []bool 类型
func ToBoolSliceE(i any) (iv []bool, err error) {
	if i == nil {
		return []bool{}, nil
	}

	switch val := i.(type) {
	case []bool:
		return val, nil
	case []any:
		iv = make([]bool, len(val))
		for k, v := range val {
			bl, err := ToBoolE(v)
			if err != nil {
				return []bool{}, convertError(i, "[]bool")
			}
			iv[k] = bl
		}
		return
	case []int64:
		iv = make([]bool, len(val))
		for k, v := range val {
			bl, err := ToBoolE(v)
			if err != nil {
				return []bool{}, convertError(i, "[]bool")
			}
			iv[k] = bl
		}
		return
	case []int32:
		iv = make([]bool, len(val))
		for k, v := range val {
			bl, err := ToBoolE(v)
			if err != nil {
				return []bool{}, convertError(i, "[]bool")
			}
			iv[k] = bl
		}
		return
	case []int16:
		iv = make([]bool, len(val))
		for k, v := range val {
			bl, err := ToBoolE(v)
			if err != nil {
				return []bool{}, convertError(i, "[]bool")
			}
			iv[k] = bl
		}
		return
	case []int8:
		iv = make([]bool, len(val))
		for k, v := range val {
			bl, err := ToBoolE(v)
			if err != nil {
				return []bool{}, convertError(i, "[]bool")
			}
			iv[k] = bl
		}
		return
	case []int:
		iv = make([]bool, len(val))
		for k, v := range val {
			bl, err := ToBoolE(v)
			if err != nil {
				return []bool{}, convertError(i, "[]bool")
			}
			iv[k] = bl
		}
		return
	case []uint64:
		iv = make([]bool, len(val))
		for k, v := range val {
			bl, err := ToBoolE(v)
			if err != nil {
				return []bool{}, convertError(i, "[]bool")
			}
			iv[k] = bl
		}
		return
	case []uint32:
		iv = make([]bool, len(val))
		for k, v := range val {
			bl, err := ToBoolE(v)
			if err != nil {
				return []bool{}, convertError(i, "[]bool")
			}
			iv[k] = bl
		}
		return
	case []uint16:
		iv = make([]bool, len(val))
		for k, v := range val {
			bl, err := ToBoolE(v)
			if err != nil {
				return []bool{}, convertError(i, "[]bool")
			}
			iv[k] = bl
		}
		return
	case []uint8:
		// 检查给定的 i 是否为 JSON 格式的字符串值，并使用 json.UnmarshalUseNumber 进行转换
		if json.Valid(val) {
			anyV := make([]any, len(val))
			if e := json.Unmarshal(val, &anyV); e == nil {
				iv = make([]bool, len(anyV))
				for k, v := range anyV {
					bl, err := ToBoolE(v)
					if err != nil {
						return []bool{}, convertError(i, "[]bool")
					}
					iv[k] = bl
				}
				return
			}
		}
		iv = make([]bool, len(val))
		for k, v := range val {
			bl, err := ToBoolE(v)
			if err != nil {
				return []bool{}, convertError(i, "[]bool")
			}
			iv[k] = bl
		}
		return
	case []uint:
		iv = make([]bool, len(val))
		for k, v := range val {
			bl, err := ToBoolE(v)
			if err != nil {
				return []bool{}, convertError(i, "[]bool")
			}
			iv[k] = bl
		}
		return
	case []string:
		iv = make([]bool, len(val))
		for k, v := range val {
			bl, err := ToBoolE(v)
			if err != nil {
				return []bool{}, convertError(i, "[]bool")
			}
			iv[k] = bl
		}
		return
	case [][]byte:
		iv = make([]bool, len(val))
		for k, v := range val {
			bl, err := ToBoolE(v)
			if err != nil {
				return []bool{}, convertError(i, "[]bool")
			}
			iv[k] = bl
		}
		return
	case string:
		// 检查给定的 i 是否为 JSON 格式的字符串值，并使用 json.UnmarshalUseNumber 进行转换
		anyBytes := []byte(val)
		if json.Valid(anyBytes) {
			anyV := make([]any, len(val))
			if e := json.Unmarshal(anyBytes, &anyV); e == nil {
				iv = make([]bool, len(anyV))
				for k, v := range anyV {
					bl, err := ToBoolE(v)
					if err != nil {
						return []bool{}, convertError(i, "[]bool")
					}
					iv[k] = bl
				}
				return
			}
		}
		return []bool{}, convertError(i, "[]bool")
	default:
		// 使用反射进行转换
		originValueAndKind := utils.OriginValueAndKind(i)
		originKind := originValueAndKind.OriginKind
		if originKind == reflect.Slice || originKind == reflect.Array {
			length := originValueAndKind.OriginValue.Len()
			iv = make([]bool, length)
			for j := 0; j < length; j++ {
				bl, err := ToBoolE(originValueAndKind.OriginValue.Index(j).Interface())
				if err != nil {
					return []bool{}, convertError(i, "[]bool")
				}
				iv[j] = bl
			}
			return
		}

		return []bool{}, convertError(i, "[]bool")
	}
}
