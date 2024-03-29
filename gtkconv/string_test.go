/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-04-15 13:42:18
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-22 22:46:20
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkconv_test

import (
	"github.com/liusuxian/go-toolkit/gtkconv"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type AAA struct {
	A int     `json:"a" dc:"a"`
	B float64 `json:"b" dc:"b"`
	C string  `json:"c" dc:"c"`
	D []any   `json:"d" dc:"d"`
}

func TestToStringE(t *testing.T) {
	assert := assert.New(t)
	actualObj, err := gtkconv.ToStringE(nil) // nil
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal("", actualObj)
	}
	actualObj, err = gtkconv.ToStringE(int(1)) // int
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal("1", actualObj)
	}
	actualObj, err = gtkconv.ToStringE(float64(1.56)) // float64
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal("1.56", actualObj)
	}
	actualObj, err = gtkconv.ToStringE(true) // bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal("true", actualObj)
	}
	actualObj, err = gtkconv.ToStringE(false) // bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal("false", actualObj)
	}
	actualObj, err = gtkconv.ToStringE([]byte("1.23")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal("1.23", actualObj)
	}
	actualObj, err = gtkconv.ToStringE(time.Date(2023, 4, 15, 0, 0, 0, 0, time.UTC)) // time.Time
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal("2023-04-15 00:00:00 +0000 UTC", actualObj)
	}
	actualObj, err = gtkconv.ToStringE([]any{1, 1.2, "hello", true}) // []any
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal("[1,1.2,\"hello\",true]", actualObj)
	}
	actualObj, err = gtkconv.ToStringE(map[string]any{"a": 1, "b": 1}) // map[string]any
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal("{\"a\":1,\"b\":1}", actualObj)
	}
	actualObj, err = gtkconv.ToStringE(&AAA{A: 1, B: 1.2, C: "hello", D: []any{1, 1.2, "hello", true}}) // struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal("{\"a\":1,\"b\":1.2,\"c\":\"hello\",\"d\":[1,1.2,\"hello\",true]}", actualObj)
	}
}
