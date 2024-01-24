/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-04-10 00:20:56
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-24 19:09:56
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkarr_test

import (
	"github.com/liusuxian/go-toolkit/gtkarr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContains(t *testing.T) {
	assert := assert.New(t)
	assert.False(gtkarr.ContainsInt([]int{}, 0))
	assert.True(gtkarr.ContainsInt([]int{0, 1, 2}, 0))
	assert.True(gtkarr.ContainsInt([]int{0, 1, 2}, 1))
	assert.True(gtkarr.ContainsInt([]int{0, 1, 2}, 2))
	assert.False(gtkarr.ContainsInt([]int{0, 1, 2}, 3))
	assert.True(gtkarr.ContainsInt([]int{10, 9, 9, 1, 6, 6, 5, 5, 4, 4, 4, 3, 3, 2}, 1))
	assert.True(gtkarr.ContainsInt([]int{10, 9, 9, 1, 6, 6, 5, 5, 4, 4, 4, 3, 3, 2}, 2))
	assert.True(gtkarr.ContainsInt([]int{10, 9, 9, 1, 6, 6, 5, 5, 4, 4, 4, 3, 3, 2}, 3))
	assert.False(gtkarr.ContainsInt([]int{10, 9, 9, 1, 6, 6, 5, 5, 4, 4, 4, 3, 3, 2}, 0))

	assert.False(gtkarr.ContainsFloat32([]float32{}, 1.0))
	assert.True(gtkarr.ContainsFloat32([]float32{1.0, 1.1, 1.2}, 1.0))
	assert.True(gtkarr.ContainsFloat32([]float32{1.0, 1.1, 1.2}, 1.1))
	assert.True(gtkarr.ContainsFloat32([]float32{1.0, 1.1, 1.2}, 1.2))
	assert.False(gtkarr.ContainsFloat32([]float32{1.0, 1.1, 1.2}, 1.3))
	assert.True(gtkarr.ContainsFloat32([]float32{10.1, 9.1, 9.1, 1.1, 6.1, 6.1, 5.1, 5.1, 4.1, 4.1, 4.1, 3.1, 3.1, 2.1}, 1.1))
	assert.True(gtkarr.ContainsFloat32([]float32{10.1, 9.1, 9.1, 1.1, 6.1, 6.1, 5.1, 5.1, 4.1, 4.1, 4.1, 3.1, 3.1, 2.1}, 2.1))
	assert.True(gtkarr.ContainsFloat32([]float32{10.1, 9.1, 9.1, 1.1, 6.1, 6.1, 5.1, 5.1, 4.1, 4.1, 4.1, 3.1, 3.1, 2.1}, 3.1))
	assert.False(gtkarr.ContainsFloat32([]float32{10.1, 9.1, 9.1, 1.1, 6.1, 6.1, 5.1, 5.1, 4.1, 4.1, 4.1, 3.1, 3.1, 2.1}, 1.0))

	assert.False(gtkarr.ContainsStr([]string{"hello", "jack", "hello", "world", "tom", "hay", "tom"}, "lsx"))
	assert.False(gtkarr.ContainsStr([]string{"hello", "jack", "hello", "world", "tom", "hay", "tom"}, "hell"))
	assert.True(gtkarr.ContainsStr([]string{"hello", "jack", "hello", "world", "tom", "hay", "tom"}, "hay"))

	assert.False(gtkarr.ContainsStr([]string{"hello", "我", "是", "中", "国", "人"}, "Hello"))
	assert.False(gtkarr.ContainsStr([]string{"hello", "我", "是", "中", "国", "人"}, "哈"))
	assert.True(gtkarr.ContainsStr([]string{"hello", "我", "是", "中", "国", "人"}, "中"))

	assert.False(gtkarr.ContainsRune([]rune("hello我是中国人"), 'H'))
	assert.False(gtkarr.ContainsRune([]rune("hello我是中国人"), '哈'))
	assert.True(gtkarr.ContainsRune([]rune("hello我是中国人"), '中'))
	assert.True(gtkarr.ContainsRune([]rune("hello我是中国人"), 'h'))

	assert.False(gtkarr.ContainsByte([]byte{4, 3, 1, 5, 6, 9, 8, 2, 1, 7}, 0))
	assert.True(gtkarr.ContainsByte([]byte{4, 3, 1, 5, 6, 9, 8, 2, 1, 7}, 2))

	assert.False(gtkarr.ContainsByte([]byte("Hello World"), 'h'))
	assert.False(gtkarr.ContainsByte([]byte("Hello World"), 'w'))
	assert.True(gtkarr.ContainsByte([]byte("Hello World"), 'H'))
	assert.True(gtkarr.ContainsByte([]byte("Hello World"), 'W'))

	assert.True(gtkarr.ContainsAny([]any{1, 1.2, "he"}, 1.2))
	assert.True(gtkarr.ContainsAny([]any{1, 1.2, "he"}, "he"))
}
