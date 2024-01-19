/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-11-27 20:32:59
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-19 21:33:54
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package reflection_test

import (
	"github.com/stretchr/testify/assert"
	"go-toolkit/reflection"
	"reflect"
	"testing"
)

func TestOriginValueAndKind(t *testing.T) {
	assert := assert.New(t)
	var s1 = "s"
	out1 := reflection.OriginValueAndKind(s1)
	assert.Equal(out1.InputKind, reflect.String)
	assert.Equal(out1.OriginKind, reflect.String)

	var s2 = "s"
	out2 := reflection.OriginValueAndKind(&s2)
	assert.Equal(out2.InputKind, reflect.Ptr)
	assert.Equal(out2.OriginKind, reflect.String)

	var s3 []int
	out3 := reflection.OriginValueAndKind(s3)
	assert.Equal(out3.InputKind, reflect.Slice)
	assert.Equal(out3.OriginKind, reflect.Slice)

	var s4 []int
	out4 := reflection.OriginValueAndKind(&s4)
	assert.Equal(out4.InputKind, reflect.Ptr)
	assert.Equal(out4.OriginKind, reflect.Slice)
}

func TestOriginTypeAndKind(t *testing.T) {
	assert := assert.New(t)
	var s1 = "s"
	out1 := reflection.OriginTypeAndKind(s1)
	assert.Equal(out1.InputKind, reflect.String)
	assert.Equal(out1.OriginKind, reflect.String)

	var s2 = "s"
	out2 := reflection.OriginTypeAndKind(&s2)
	assert.Equal(out2.InputKind, reflect.Ptr)
	assert.Equal(out2.OriginKind, reflect.String)

	var s3 []int
	out3 := reflection.OriginTypeAndKind(s3)
	assert.Equal(out3.InputKind, reflect.Slice)
	assert.Equal(out3.OriginKind, reflect.Slice)

	var s4 []int
	out4 := reflection.OriginTypeAndKind(&s4)
	assert.Equal(out4.InputKind, reflect.Ptr)
	assert.Equal(out4.OriginKind, reflect.Slice)
}
