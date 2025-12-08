/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-13 13:14:03
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-08 15:02:16
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package utils_test

import (
	"github.com/liusuxian/go-toolkit/internal/utils"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestOriginValueAndKind(t *testing.T) {
	assert := assert.New(t)
	var s1 = "s"
	out1 := utils.OriginValueAndKind(s1)
	assert.Equal(out1.InputKind, reflect.String)
	assert.Equal(out1.OriginKind, reflect.String)

	var s2 = "s"
	out2 := utils.OriginValueAndKind(&s2)
	assert.Equal(out2.InputKind, reflect.Ptr)
	assert.Equal(out2.OriginKind, reflect.String)

	var s3 []int
	out3 := utils.OriginValueAndKind(s3)
	assert.Equal(out3.InputKind, reflect.Slice)
	assert.Equal(out3.OriginKind, reflect.Slice)

	var s4 []int
	out4 := utils.OriginValueAndKind(&s4)
	assert.Equal(out4.InputKind, reflect.Ptr)
	assert.Equal(out4.OriginKind, reflect.Slice)
}

func TestOriginTypeAndKind(t *testing.T) {
	assert := assert.New(t)
	var s1 = "s"
	out1 := utils.OriginTypeAndKind(s1)
	assert.Equal(out1.InputKind, reflect.String)
	assert.Equal(out1.OriginKind, reflect.String)

	var s2 = "s"
	out2 := utils.OriginTypeAndKind(&s2)
	assert.Equal(out2.InputKind, reflect.Ptr)
	assert.Equal(out2.OriginKind, reflect.String)

	var s3 []int
	out3 := utils.OriginTypeAndKind(s3)
	assert.Equal(out3.InputKind, reflect.Slice)
	assert.Equal(out3.OriginKind, reflect.Slice)

	var s4 []int
	out4 := utils.OriginTypeAndKind(&s4)
	assert.Equal(out4.InputKind, reflect.Ptr)
	assert.Equal(out4.OriginKind, reflect.Slice)
}

type TestNil struct {
	A int
}

func TestIsNil(t *testing.T) {
	assert := assert.New(t)
	var a any = nil
	var b any
	var c *int = nil
	var d int = 10
	var e string = "hello"
	assert.True(utils.IsNil(a))
	assert.True(utils.IsNil(b))
	assert.True(utils.IsNil(c))
	assert.False(utils.IsNil(d))
	assert.False(utils.IsNil(e))
	assert.True(utils.IsNil(nil))
	var aa any = nil
	var bb any
	assert.True(utils.IsNil(aa))
	assert.True(utils.IsNil(bb))
	var cc TestNil
	var dd *TestNil
	assert.False(utils.IsNil(cc))
	assert.True(utils.IsNil(dd))
}
