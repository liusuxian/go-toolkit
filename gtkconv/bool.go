/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-04-15 13:26:31
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-23 00:27:08
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkconv

import (
	"encoding/json"
	"strconv"
	"strings"
)

// ToBoolE 将 any 转换为 bool 类型
func ToBoolE(i any) (bl bool, err error) {
	i = indirect(i)

	switch val := i.(type) {
	case nil:
		return false, nil
	case bool:
		return val, nil
	case int64:
		return val > 0, nil
	case int32:
		return val > 0, nil
	case int16:
		return val > 0, nil
	case int8:
		return val > 0, nil
	case int:
		return val > 0, nil
	case uint64:
		return val > 0, nil
	case uint32:
		return val > 0, nil
	case uint16:
		return val > 0, nil
	case uint8:
		return val > 0, nil
	case uint:
		return val > 0, nil
	case float64:
		return val > 0, nil
	case float32:
		return val > 0, nil
	case []byte:
		return ToBoolE(string(val))
	case string:
		if val == "" {
			return false, nil
		}
		if strings.ToUpper(val) == "OK" {
			return true, nil
		}
		iv, err := strconv.ParseBool(val)
		if err == nil {
			return iv, nil
		}
		return false, convertError(i, "bool")
	case json.Number:
		iv, err := ToInt64E(val)
		if err == nil {
			return iv > 0, nil
		}
		return false, convertError(i, "bool")
	default:
		return false, convertError(i, "bool")
	}
}
