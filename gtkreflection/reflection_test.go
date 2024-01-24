/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-11-27 20:32:59
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-20 00:37:03
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkreflection_test

import (
	"github.com/liusuxian/go-toolkit/gtkreflection"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestOriginValueAndKind(t *testing.T) {
	assert := assert.New(t)
	var s1 = "s"
	out1 := gtkreflection.OriginValueAndKind(s1)
	assert.Equal(out1.InputKind, reflect.String)
	assert.Equal(out1.OriginKind, reflect.String)

	var s2 = "s"
	out2 := gtkreflection.OriginValueAndKind(&s2)
	assert.Equal(out2.InputKind, reflect.Ptr)
	assert.Equal(out2.OriginKind, reflect.String)

	var s3 []int
	out3 := gtkreflection.OriginValueAndKind(s3)
	assert.Equal(out3.InputKind, reflect.Slice)
	assert.Equal(out3.OriginKind, reflect.Slice)

	var s4 []int
	out4 := gtkreflection.OriginValueAndKind(&s4)
	assert.Equal(out4.InputKind, reflect.Ptr)
	assert.Equal(out4.OriginKind, reflect.Slice)
}

func TestOriginTypeAndKind(t *testing.T) {
	assert := assert.New(t)
	var s1 = "s"
	out1 := gtkreflection.OriginTypeAndKind(s1)
	assert.Equal(out1.InputKind, reflect.String)
	assert.Equal(out1.OriginKind, reflect.String)

	var s2 = "s"
	out2 := gtkreflection.OriginTypeAndKind(&s2)
	assert.Equal(out2.InputKind, reflect.Ptr)
	assert.Equal(out2.OriginKind, reflect.String)

	var s3 []int
	out3 := gtkreflection.OriginTypeAndKind(s3)
	assert.Equal(out3.InputKind, reflect.Slice)
	assert.Equal(out3.OriginKind, reflect.Slice)

	var s4 []int
	out4 := gtkreflection.OriginTypeAndKind(&s4)
	assert.Equal(out4.InputKind, reflect.Ptr)
	assert.Equal(out4.OriginKind, reflect.Slice)
}
