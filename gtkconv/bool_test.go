/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-04-15 13:34:47
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-22 22:53:36
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkconv_test

import (
	"github.com/liusuxian/go-toolkit/gtkconv"
	"github.com/stretchr/testify/assert"
	"testing"
)

func errLog(t *testing.T, err error) {
	if err != nil {
		t.Logf("Error: %+v\n", err.Error())
	}
}

func TestToBoolE(t *testing.T) {
	assert := assert.New(t)
	actualObj, err := gtkconv.ToBoolE(nil) // nil
	errLog(t, err)
	if assert.NoError(err) {
		assert.False(actualObj)
	}
	actualObj, err = gtkconv.ToBoolE(true) // bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.True(actualObj)
	}
	actualObj, err = gtkconv.ToBoolE(false) // bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.False(actualObj)
	}
	actualObj, err = gtkconv.ToBoolE(byte('a')) // int8
	errLog(t, err)
	if assert.NoError(err) {
		assert.True(actualObj)
	}
	actualObj, err = gtkconv.ToBoolE(0) // int
	errLog(t, err)
	if assert.NoError(err) {
		assert.False(actualObj)
	}
	actualObj, err = gtkconv.ToBoolE(-1) // int
	errLog(t, err)
	if assert.NoError(err) {
		assert.False(actualObj)
	}
	actualObj, err = gtkconv.ToBoolE(1) // int
	errLog(t, err)
	if assert.NoError(err) {
		assert.True(actualObj)
	}
	actualObj, err = gtkconv.ToBoolE([]byte("true")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.True(actualObj)
	}
	actualObj, err = gtkconv.ToBoolE([]byte("")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.False(actualObj)
	}
	actualObj, err = gtkconv.ToBoolE([]byte{}) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.False(actualObj)
	}
	actualObj, err = gtkconv.ToBoolE("true") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.True(actualObj)
	}
	actualObj, err = gtkconv.ToBoolE("false") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.False(actualObj)
	}
	actualObj, err = gtkconv.ToBoolE("") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.False(actualObj)
	}
	actualObj, err = gtkconv.ToBoolE(" ") // string
	errLog(t, err)
	if assert.Error(err) {
		assert.False(actualObj)
	}
	actualObj, err = gtkconv.ToBoolE("0") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.False(actualObj)
	}
	actualObj, err = gtkconv.ToBoolE("hello") // string
	errLog(t, err)
	if assert.Error(err) {
		assert.False(actualObj)
	}
	actualObj, err = gtkconv.ToBoolE(float64(1.23)) // float64
	errLog(t, err)
	if assert.NoError(err) {
		assert.True(actualObj)
	}
}
