/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-03-01 16:52:17
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-22 23:49:40
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkconv_test

import (
	"github.com/liusuxian/go-toolkit/gtkconv"
	"github.com/stretchr/testify/assert"
	"testing"
)

type User struct {
	Name string
	Age  int
	Url  string
}

func TestToByteE(t *testing.T) {
	assert := assert.New(t)
	b, err := gtkconv.ToByteE(1)
	assert.NoError(err)
	assert.Equal(uint8(0x1), b)
}

func TestToRuneE(t *testing.T) {
	assert := assert.New(t)
	b, err := gtkconv.ToRuneE('我')
	assert.NoError(err)
	assert.Equal(int32(25105), b)
}

func TestToRunesE(t *testing.T) {
	assert := assert.New(t)
	b, err := gtkconv.ToRunesE('我')
	assert.NoError(err)
	assert.Equal([]int32{50, 53, 49, 48, 53}, b)
}

func TestToBytesE(t *testing.T) {
	assert := assert.New(t)
	b, err := gtkconv.ToBytesE(User{"wenzi1", 999, "www.baidu.com"})
	assert.NoError(err)
	assert.Equal("{\"Name\":\"wenzi1\",\"Age\":999,\"Url\":\"www.baidu.com\"}", string(b))
}
