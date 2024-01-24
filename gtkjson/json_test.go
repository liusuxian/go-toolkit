/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-22 23:27:03
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-24 19:14:42
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkjson_test

import (
	"github.com/liusuxian/go-toolkit/gtkjson"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJsonMarshal(t *testing.T) {
	var (
		err     error
		jsonStr string
		assert  = assert.New(t)
	)
	jsonStr, err = gtkjson.JsonMarshal(map[string]any{
		"a": 1,
		"b": 1.2,
		"c": map[string]any{
			"d": 3,
		},
	})
	if assert.NoError(err) {
		assert.Equal("{\"a\":1,\"b\":1.2,\"c\":{\"d\":3}}\n", jsonStr)
	}
}

func TestMustJsonMarshal(t *testing.T) {
	var (
		jsonStr string
		assert  = assert.New(t)
	)
	jsonStr = gtkjson.MustJsonMarshal(map[string]any{
		"a": 1,
		"b": 1.2,
		"c": map[string]any{
			"d": 3,
		},
	})
	assert.Equal("{\"a\":1,\"b\":1.2,\"c\":{\"d\":3}}\n", jsonStr)
}
