/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-05-05 15:36:47
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-22 22:39:11
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkconv_test

import (
	"github.com/liusuxian/go-toolkit/gtk/gtkconv"
	"github.com/stretchr/testify/assert"
	"testing"
)

type FFFFFFF struct {
	A any
	B any
	C any
}

type GGGGGGG struct {
	A any `json:"a" dc:"a"`
	B any `json:"b" dc:"b"`
	C any `json:"c" dc:"c"`
}

func TestToStringMapStringE(t *testing.T) {
	assert := assert.New(t)
	actualObj, err := gtkconv.ToStringMapStringE(map[any]any{"a": "1", "b": 2.6, "c": true}) // map[any]any
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]string{"a": "1", "b": "2.6", "c": "true"}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapStringE([]byte(`{"a": "1.6", "b": 2.7, "c": true}`)) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]string{"a": "1.6", "b": "2.7", "c": "true"}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapStringE(`{"a": "1.6", "b": 2.7, "c": true}`) // json
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]string{"a": "1.6", "b": "2.7", "c": "true"}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapStringE(map[string]string{"a": "1.6", "b": "2.7", "c": "3.1"}) // map[string]string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]string{"a": "1.6", "b": "2.7", "c": "3.1"}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapStringE(FFFFFFF{A: 1.6, B: false, C: "2.7"}) // struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]string{"A": "1.6", "B": "false", "C": "2.7"}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapStringE(&FFFFFFF{A: 1.6, B: false, C: "2.7"}) // *struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]string{"A": "1.6", "B": "false", "C": "2.7"}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapStringE(GGGGGGG{A: 1.6, B: false, C: "2.7"}) // struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]string{"a": "1.6", "b": "false", "c": "2.7"}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapStringE(&GGGGGGG{A: 1.6, B: false, C: "2.7"}) // *struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]string{"a": "1.6", "b": "false", "c": "2.7"}, actualObj)
	}
}

func TestToStringMapStringSliceE(t *testing.T) {
	assert := assert.New(t)
	actualObj, err := gtkconv.ToStringMapStringSliceE(map[any]any{"a": []any{1, 1.2, true}}) // map[any]any
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string][]string{"a": {"1", "1.2", "true"}}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapStringSliceE(map[any][]any{"a": {1, 1.2, true}}) // map[any][]any
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string][]string{"a": {"1", "1.2", "true"}}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapStringSliceE(map[any]string{"a": "1.6", "b": "[\"1\", \"2.1\", \"true\"]"}) // map[any]string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string][]string{"a": {"1.6"}, "b": {"1", "2.1", "true"}}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapStringSliceE(map[string]any{"a": "1.6", "b": "[\"1\", \"2.1\", \"true\"]"}) // map[string]any
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string][]string{"a": {"1.6"}, "b": {"1", "2.1", "true"}}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapStringSliceE(map[string]string{"a": "1.6", "b": "[\"1\", \"2.1\", \"true\"]"}) // map[string]string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string][]string{"a": {"1.6"}, "b": {"1", "2.1", "true"}}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapStringSliceE([]byte(`{"a": [1.6, 2.7, true]}`)) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string][]string{"a": {"1.6", "2.7", "true"}}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapStringSliceE(`{"a": [1.6, 2.7, true]}`) // json
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string][]string{"a": {"1.6", "2.7", "true"}}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapStringSliceE(FFFFFFF{A: "1.6", B: "false", C: "2.7"}) // struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string][]string{"A": {"1.6"}, "B": {"false"}, "C": {"2.7"}}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapStringSliceE(&FFFFFFF{A: "1.6", B: "false", C: "2.7"}) // *struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string][]string{"A": {"1.6"}, "B": {"false"}, "C": {"2.7"}}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapStringSliceE(GGGGGGG{A: "1.6", B: "false", C: "2.7"}) // struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string][]string{"a": {"1.6"}, "b": {"false"}, "c": {"2.7"}}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapStringSliceE(&GGGGGGG{A: "1.6", B: "false", C: "2.7"}) // *struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string][]string{"a": {"1.6"}, "b": {"false"}, "c": {"2.7"}}, actualObj)
	}
}
