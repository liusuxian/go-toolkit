/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-04-18 18:18:23
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-22 22:40:10
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

func TestToSliceE(t *testing.T) {
	assert := assert.New(t)
	actualObj, err := gtkconv.ToSliceE("[1, 1.2, true, \"hello\"]") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.ElementsMatch([]any{1.0, 1.2, true, "hello"}, actualObj)
	}
	actualObj, err = gtkconv.ToSliceE([]byte("[1, 1.2, true, \"hello\"]")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.ElementsMatch([]any{1.0, 1.2, true, "hello"}, actualObj)
	}
	actualObj, err = gtkconv.ToSliceE("[\"h\", \"e\", \"l\", \"l\", \"o\"]") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.ElementsMatch([]any{"h", "e", "l", "l", "o"}, actualObj)
	}
	actualObj, err = gtkconv.ToSliceE(1) // int
	errLog(t, err)
	if assert.NoError(err) {
		assert.ElementsMatch([]any{1}, actualObj)
	}
	actualObj, err = gtkconv.ToSliceE("hello") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.ElementsMatch([]any{"hello"}, actualObj)
	}
	actualObj, err = gtkconv.ToSliceE([]map[string]any{{"a1": 1, "b1": 2}, {"a2": 3, "b2": 4}}) // []map[string]any
	errLog(t, err)
	if assert.NoError(err) {
		assert.ElementsMatch([]any{map[string]any{"a1": 1, "b1": 2}, map[string]any{"a2": 3, "b2": 4}}, actualObj)
	}
	actualObj, err = gtkconv.ToSliceE([]map[string]int{{"a1": 1, "b1": 2}, {"a2": 3, "b2": 4}}) // []map[string]int
	errLog(t, err)
	if assert.NoError(err) {
		assert.ElementsMatch([]any{map[string]int{"a1": 1, "b1": 2}, map[string]int{"a2": 3, "b2": 4}}, actualObj)
	}
	actualObj, err = gtkconv.ToSliceE([]map[string]bool{{"a1": true, "b1": false}, {"a2": true, "b2": false}}) // []map[string]bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.ElementsMatch([]any{map[string]bool{"a1": true, "b1": false}, map[string]bool{"a2": true, "b2": false}}, actualObj)
	}
	actualObj, err = gtkconv.ToSliceE([][]int{{1, 2}, {3, 4}}) // [][]int
	errLog(t, err)
	if assert.NoError(err) {
		assert.ElementsMatch([]any{[]int{1, 2}, []int{3, 4}}, actualObj)
	}
}
