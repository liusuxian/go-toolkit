/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-04-16 02:21:26
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-22 22:42:00
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

func TestToBoolSliceE(t *testing.T) {
	assert := assert.New(t)
	actualObj, err := gtkconv.ToBoolSliceE([]int{0, 1, 0}) // []int
	errLog(t, err)
	if assert.NoError(err) {
		assert.ElementsMatch([]bool{false, true, false}, actualObj)
	}
	actualObj, err = gtkconv.ToBoolSliceE([]string{"true", "false"}) // []string
	errLog(t, err)
	if assert.NoError(err) {
		assert.ElementsMatch([]bool{true, false}, actualObj)
	}
	actualObj, err = gtkconv.ToBoolSliceE([][]byte{[]byte("1"), []byte("0")}) // [][]byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.ElementsMatch([]bool{true, false}, actualObj)
	}
	actualObj, err = gtkconv.ToBoolSliceE([]map[string]any{{"a1": 1, "b1": 2}, {"a2": 3, "b2": 4}}) // []map[string]any
	errLog(t, err)
	if assert.Error(err) {
		assert.ElementsMatch([]bool{}, actualObj)
	}
	actualObj, err = gtkconv.ToBoolSliceE([]byte("[1, 0, true, false, \"true\", \"false\"]")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.ElementsMatch([]bool{true, false, true, false, true, false}, actualObj)
	}
	actualObj, err = gtkconv.ToBoolSliceE("[1, 0, true, false, \"true\", \"false\"]") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.ElementsMatch([]bool{true, false, true, false, true, false}, actualObj)
	}
}
