/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-04-10 00:20:56
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-22 22:58:13
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkslice_test

import (
	"github.com/liusuxian/go-toolkit/gtk/gtkslice"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsContains(t *testing.T) {
	assert := assert.New(t)
	assert.False(gtkslice.IsContains([]int{}, 0))
	assert.True(gtkslice.IsContains([]int{0, 1, 2}, 0))
	assert.True(gtkslice.IsContains([]int{0, 1, 2}, 1))
	assert.True(gtkslice.IsContains([]int{0, 1, 2}, 2))
	assert.False(gtkslice.IsContains([]int{0, 1, 2}, 3))
	assert.True(gtkslice.IsContains([]int{10, 9, 9, 1, 6, 6, 5, 5, 4, 4, 4, 3, 3, 2}, 1))
	assert.True(gtkslice.IsContains([]int{10, 9, 9, 1, 6, 6, 5, 5, 4, 4, 4, 3, 3, 2}, 2))
	assert.True(gtkslice.IsContains([]int{10, 9, 9, 1, 6, 6, 5, 5, 4, 4, 4, 3, 3, 2}, 3))
	assert.False(gtkslice.IsContains([]int{10, 9, 9, 1, 6, 6, 5, 5, 4, 4, 4, 3, 3, 2}, 0))

	assert.False(gtkslice.IsContains([]float32{}, 1.0))
	assert.True(gtkslice.IsContains([]float32{1.0, 1.1, 1.2}, 1.0))
	assert.True(gtkslice.IsContains([]float32{1.0, 1.1, 1.2}, 1.1))
	assert.True(gtkslice.IsContains([]float32{1.0, 1.1, 1.2}, 1.2))
	assert.False(gtkslice.IsContains([]float32{1.0, 1.1, 1.2}, 1.3))
	assert.True(gtkslice.IsContains([]float32{10.1, 9.1, 9.1, 1.1, 6.1, 6.1, 5.1, 5.1, 4.1, 4.1, 4.1, 3.1, 3.1, 2.1}, 1.1))
	assert.True(gtkslice.IsContains([]float32{10.1, 9.1, 9.1, 1.1, 6.1, 6.1, 5.1, 5.1, 4.1, 4.1, 4.1, 3.1, 3.1, 2.1}, 2.1))
	assert.True(gtkslice.IsContains([]float32{10.1, 9.1, 9.1, 1.1, 6.1, 6.1, 5.1, 5.1, 4.1, 4.1, 4.1, 3.1, 3.1, 2.1}, 3.1))
	assert.False(gtkslice.IsContains([]float32{10.1, 9.1, 9.1, 1.1, 6.1, 6.1, 5.1, 5.1, 4.1, 4.1, 4.1, 3.1, 3.1, 2.1}, 1.0))

	assert.False(gtkslice.IsContains([]string{"hello", "jack", "hello", "world", "tom", "hay", "tom"}, "lsx"))
	assert.False(gtkslice.IsContains([]string{"hello", "jack", "hello", "world", "tom", "hay", "tom"}, "hell"))
	assert.True(gtkslice.IsContains([]string{"hello", "jack", "hello", "world", "tom", "hay", "tom"}, "hay"))

	assert.False(gtkslice.IsContains([]string{"hello", "我", "是", "中", "国", "人"}, "Hello"))
	assert.False(gtkslice.IsContains([]string{"hello", "我", "是", "中", "国", "人"}, "哈"))
	assert.True(gtkslice.IsContains([]string{"hello", "我", "是", "中", "国", "人"}, "中"))

	assert.False(gtkslice.IsContains([]rune("hello我是中国人"), 'H'))
	assert.False(gtkslice.IsContains([]rune("hello我是中国人"), '哈'))
	assert.True(gtkslice.IsContains([]rune("hello我是中国人"), '中'))
	assert.True(gtkslice.IsContains([]rune("hello我是中国人"), 'h'))

	assert.False(gtkslice.IsContains([]byte{4, 3, 1, 5, 6, 9, 8, 2, 1, 7}, 0))
	assert.True(gtkslice.IsContains([]byte{4, 3, 1, 5, 6, 9, 8, 2, 1, 7}, 2))

	assert.False(gtkslice.IsContains([]byte("Hello World"), 'h'))
	assert.False(gtkslice.IsContains([]byte("Hello World"), 'w'))
	assert.True(gtkslice.IsContains([]byte("Hello World"), 'H'))
	assert.True(gtkslice.IsContains([]byte("Hello World"), 'W'))
}
