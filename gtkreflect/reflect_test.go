/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-11-27 20:32:59
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-29 16:17:00
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkreflect_test

import (
	"github.com/liusuxian/go-toolkit/gtkreflect"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestOriginValueAndKind(t *testing.T) {
	assert := assert.New(t)
	var s1 = "s"
	out1 := gtkreflect.OriginValueAndKind(s1)
	assert.Equal(out1.InputKind, reflect.String)
	assert.Equal(out1.OriginKind, reflect.String)

	var s2 = "s"
	out2 := gtkreflect.OriginValueAndKind(&s2)
	assert.Equal(out2.InputKind, reflect.Ptr)
	assert.Equal(out2.OriginKind, reflect.String)

	var s3 []int
	out3 := gtkreflect.OriginValueAndKind(s3)
	assert.Equal(out3.InputKind, reflect.Slice)
	assert.Equal(out3.OriginKind, reflect.Slice)

	var s4 []int
	out4 := gtkreflect.OriginValueAndKind(&s4)
	assert.Equal(out4.InputKind, reflect.Ptr)
	assert.Equal(out4.OriginKind, reflect.Slice)
}

func TestOriginTypeAndKind(t *testing.T) {
	assert := assert.New(t)
	var s1 = "s"
	out1 := gtkreflect.OriginTypeAndKind(s1)
	assert.Equal(out1.InputKind, reflect.String)
	assert.Equal(out1.OriginKind, reflect.String)

	var s2 = "s"
	out2 := gtkreflect.OriginTypeAndKind(&s2)
	assert.Equal(out2.InputKind, reflect.Ptr)
	assert.Equal(out2.OriginKind, reflect.String)

	var s3 []int
	out3 := gtkreflect.OriginTypeAndKind(s3)
	assert.Equal(out3.InputKind, reflect.Slice)
	assert.Equal(out3.OriginKind, reflect.Slice)

	var s4 []int
	out4 := gtkreflect.OriginTypeAndKind(&s4)
	assert.Equal(out4.InputKind, reflect.Ptr)
	assert.Equal(out4.OriginKind, reflect.Slice)
}

type TestNil struct {
	A int
}

func TestIsNil(t *testing.T) {
	assert := assert.New(t)
	var a interface{} = nil
	var b interface{}
	var c *int = nil
	var d int = 10
	var e string = "hello"
	assert.True(gtkreflect.IsNil(a))
	assert.True(gtkreflect.IsNil(b))
	assert.True(gtkreflect.IsNil(c))
	assert.False(gtkreflect.IsNil(d))
	assert.False(gtkreflect.IsNil(e))
	assert.True(gtkreflect.IsNil(nil))
	var aa any = nil
	var bb any
	assert.True(gtkreflect.IsNil(aa))
	assert.True(gtkreflect.IsNil(bb))
	var cc TestNil
	var dd *TestNil
	assert.False(gtkreflect.IsNil(cc))
	assert.True(gtkreflect.IsNil(dd))
}
