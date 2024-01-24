/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-04-16 02:24:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-22 22:44:54
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

func TestToStringSliceE(t *testing.T) {
	assert := assert.New(t)
	actualObj, err := gtkconv.ToStringSliceE([]int{0, 1, 0}) // []int
	errLog(t, err)
	if assert.NoError(err) {
		assert.ElementsMatch([]string{"0", "1", "0"}, actualObj)
	}
	actualObj, err = gtkconv.ToStringSliceE([][]byte{[]byte("1"), []byte("0")}) // [][]byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.ElementsMatch([]string{"1", "0"}, actualObj)
	}
	actualObj, err = gtkconv.ToStringSliceE([]map[string]any{{"a1": 1, "b1": 2}, {"a2": 3, "b2": 4}}) // []map[string]any
	errLog(t, err)
	if assert.NoError(err) {
		assert.ElementsMatch([]string{"{\"a1\":1,\"b1\":2}", "{\"a2\":3,\"b2\":4}"}, actualObj)
	}
	actualObj, err = gtkconv.ToStringSliceE([]map[string]int{{"a1": 1, "b1": 2}, {"a2": 3, "b2": 4}}) // []map[string]int
	errLog(t, err)
	if assert.NoError(err) {
		assert.ElementsMatch([]string{"{\"a1\":1,\"b1\":2}", "{\"a2\":3,\"b2\":4}"}, actualObj)
	}
	actualObj, err = gtkconv.ToStringSliceE([]map[string]bool{{"a1": true, "b1": false}, {"a2": true, "b2": false}}) // []map[string]bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.ElementsMatch([]string{"{\"a1\":true,\"b1\":false}", "{\"a2\":true,\"b2\":false}"}, actualObj)
	}
	actualObj, err = gtkconv.ToStringSliceE([][]int{{1, 2}, {3, 4}}) // [][]int
	errLog(t, err)
	if assert.NoError(err) {
		assert.ElementsMatch([]string{"[1,2]", "[3,4]"}, actualObj)
	}
	actualObj, err = gtkconv.ToStringSliceE([]byte("[1, 1.2, true, \"hello\", \"world\"]")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.ElementsMatch([]string{"1", "1.2", "true", "hello", "world"}, actualObj)
	}
	actualObj, err = gtkconv.ToStringSliceE("[1, 1.2, true, \"hello\", \"world\"]") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.ElementsMatch([]string{"1", "1.2", "true", "hello", "world"}, actualObj)
	}
}
