/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-07 02:44:27
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-29 15:55:18
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package utils

import (
	"encoding/json"
	"reflect"
)

// DoRedisArgs 处理`redis`命令参数
func DoRedisArgs(sidx int, args ...any) (err error) {
	for k, v := range args {
		if k > (sidx - 1) {
			reflectInfo := OriginTypeAndKind(v)
			switch reflectInfo.OriginKind {
			case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
				// 忽略切片类型为 []byte 的情况
				if _, ok := v.([]byte); !ok {
					if args[k], err = json.Marshal(v); err != nil {
						return
					}
				}
			}
		}
	}
	return
}
