/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-19 21:33:07
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-13 13:23:47
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkreflect

import (
	"github.com/liusuxian/go-toolkit/internal/utils"
	"reflect"
)

// OriginValueAndKindOutput
type OriginValueAndKindOutput struct {
	InputValue  reflect.Value
	InputKind   reflect.Kind
	OriginValue reflect.Value
	OriginKind  reflect.Kind
}

// OriginValueAndKind 检索并返回原始的反射值和数据类型
func OriginValueAndKind(value any) (out OriginValueAndKindOutput) {
	o := utils.OriginValueAndKind(value)
	return OriginValueAndKindOutput{
		InputValue:  o.InputValue,
		InputKind:   o.InputKind,
		OriginValue: o.OriginValue,
		OriginKind:  o.OriginKind,
	}
}

// OriginTypeAndKindOutput
type OriginTypeAndKindOutput struct {
	InputType  reflect.Type
	InputKind  reflect.Kind
	OriginType reflect.Type
	OriginKind reflect.Kind
}

// OriginTypeAndKind 检索并返回原始的反射类型和数据类型
func OriginTypeAndKind(value any) (out OriginTypeAndKindOutput) {
	o := utils.OriginTypeAndKind(value)
	return OriginTypeAndKindOutput{
		InputType:  o.InputType,
		InputKind:  o.InputKind,
		OriginType: o.OriginType,
		OriginKind: o.OriginKind,
	}
}

// ValueToInterface 将 reflect 值转换为其 any 类型
func ValueToInterface(v reflect.Value) (value any, ok bool) {
	return utils.ValueToInterface(v)
}

// IsNil 检查给定的值是否是`nil`
func IsNil(value any) (isNil bool) {
	return utils.IsNil(value)
}
