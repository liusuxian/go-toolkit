/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-04-15 13:39:59
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-24 16:42:51
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

func TestToUint64E(t *testing.T) {
	assert := assert.New(t)
	actualObj, err := gtkconv.ToUint64E(nil) // nil
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint64(0), actualObj)
	}
	actualObj, err = gtkconv.ToUint64E(int(1)) // int
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint64(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint64E(float64(1.56)) // float64
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint64(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint64E(true) // bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint64(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint64E(false) // bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint64(0), actualObj)
	}
	actualObj, err = gtkconv.ToUint64E([]byte("1.23")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint64(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint64E([]byte("1.0")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint64(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint64E([]byte("1.")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint64(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint64E([]byte("1")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint64(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint64E([]byte("a")) // []byte
	errLog(t, err)
	if assert.Error(err) {
		assert.Equal(uint64(0), actualObj)
	}
	actualObj, err = gtkconv.ToUint64E("1.23") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint64(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint64E("1.0") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint64(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint64E("1.") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint64(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint64E("1") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint64(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint64E("-1") // string
	errLog(t, err)
	if assert.Error(err) {
		assert.Equal(uint64(0), actualObj)
	}
	actualObj, err = gtkconv.ToUint64E("b") // string
	errLog(t, err)
	if assert.Error(err) {
		assert.Equal(uint64(0), actualObj)
	}
}

func TestToUint32E(t *testing.T) {
	assert := assert.New(t)
	actualObj, err := gtkconv.ToUint32E(nil) // nil
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint32(0), actualObj)
	}
	actualObj, err = gtkconv.ToUint32E(int(1)) // int
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint32(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint32E(float64(1.56)) // float64
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint32(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint32E(true) // bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint32(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint32E(false) // bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint32(0), actualObj)
	}
	actualObj, err = gtkconv.ToUint32E([]byte("1.23")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint32(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint32E([]byte("1.0")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint32(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint32E([]byte("1.")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint32(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint32E([]byte("1")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint32(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint32E([]byte("a")) // []byte
	errLog(t, err)
	if assert.Error(err) {
		assert.Equal(uint32(0), actualObj)
	}
	actualObj, err = gtkconv.ToUint32E("1.23") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint32(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint32E("1.0") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint32(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint32E("1.") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint32(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint32E("1") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint32(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint32E("-1") // string
	errLog(t, err)
	if assert.Error(err) {
		assert.Equal(uint32(0), actualObj)
	}
	actualObj, err = gtkconv.ToUint32E("b") // string
	errLog(t, err)
	if assert.Error(err) {
		assert.Equal(uint32(0), actualObj)
	}
}

func TestToUint16E(t *testing.T) {
	assert := assert.New(t)
	actualObj, err := gtkconv.ToUint16E(nil) // nil
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint16(0), actualObj)
	}
	actualObj, err = gtkconv.ToUint16E(int(1)) // int
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint16(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint16E(float64(1.56)) // float64
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint16(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint16E(true) // bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint16(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint16E(false) // bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint16(0), actualObj)
	}
	actualObj, err = gtkconv.ToUint16E([]byte("1.23")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint16(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint16E([]byte("1.0")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint16(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint16E([]byte("1.")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint16(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint16E([]byte("1")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint16(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint16E([]byte("a")) // []byte
	errLog(t, err)
	if assert.Error(err) {
		assert.Equal(uint16(0), actualObj)
	}
	actualObj, err = gtkconv.ToUint16E("1.23") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint16(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint16E("1.0") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint16(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint16E("1.") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint16(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint16E("1") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint16(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint16E("-1") // string
	errLog(t, err)
	if assert.Error(err) {
		assert.Equal(uint16(0), actualObj)
	}
	actualObj, err = gtkconv.ToUint16E("b") // string
	errLog(t, err)
	if assert.Error(err) {
		assert.Equal(uint16(0), actualObj)
	}
}

func TestToUint8E(t *testing.T) {
	assert := assert.New(t)
	actualObj, err := gtkconv.ToUint8E(nil) // nil
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint8(0), actualObj)
	}
	actualObj, err = gtkconv.ToUint8E(int(1)) // int
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint8(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint8E(float64(1.56)) // float64
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint8(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint8E(true) // bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint8(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint8E(false) // bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint8(0), actualObj)
	}
	actualObj, err = gtkconv.ToUint8E([]byte("1.23")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint8(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint8E([]byte("1.0")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint8(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint8E([]byte("1.")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint8(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint8E([]byte("1")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint8(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint8E([]byte("a")) // []byte
	errLog(t, err)
	if assert.Error(err) {
		assert.Equal(uint8(0), actualObj)
	}
	actualObj, err = gtkconv.ToUint8E("1.23") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint8(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint8E("1.0") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint8(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint8E("1.") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint8(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint8E("1") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint8(1), actualObj)
	}
	actualObj, err = gtkconv.ToUint8E("-1") // string
	errLog(t, err)
	if assert.Error(err) {
		assert.Equal(uint8(0), actualObj)
	}
	actualObj, err = gtkconv.ToUint8E("b") // string
	errLog(t, err)
	if assert.Error(err) {
		assert.Equal(uint8(0), actualObj)
	}
}

func TestToUintE(t *testing.T) {
	assert := assert.New(t)
	actualObj, err := gtkconv.ToUintE(nil) // nil
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint(0), actualObj)
	}
	actualObj, err = gtkconv.ToUintE(int(1)) // int
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint(1), actualObj)
	}
	actualObj, err = gtkconv.ToUintE(float64(1.56)) // float64
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint(1), actualObj)
	}
	actualObj, err = gtkconv.ToUintE(true) // bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint(1), actualObj)
	}
	actualObj, err = gtkconv.ToUintE(false) // bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint(0), actualObj)
	}
	actualObj, err = gtkconv.ToUintE([]byte("1.23")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint(1), actualObj)
	}
	actualObj, err = gtkconv.ToUintE([]byte("1.0")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint(1), actualObj)
	}
	actualObj, err = gtkconv.ToUintE([]byte("1.")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint(1), actualObj)
	}
	actualObj, err = gtkconv.ToUintE([]byte("1")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint(1), actualObj)
	}
	actualObj, err = gtkconv.ToUintE([]byte("a")) // []byte
	errLog(t, err)
	if assert.Error(err) {
		assert.Equal(uint(0), actualObj)
	}
	actualObj, err = gtkconv.ToUintE("1.23") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint(1), actualObj)
	}
	actualObj, err = gtkconv.ToUintE("1.0") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint(1), actualObj)
	}
	actualObj, err = gtkconv.ToUintE("1.") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint(1), actualObj)
	}
	actualObj, err = gtkconv.ToUintE("1") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(uint(1), actualObj)
	}
	actualObj, err = gtkconv.ToUintE("-1") // string
	errLog(t, err)
	if assert.Error(err) {
		assert.Equal(uint(0), actualObj)
	}
	actualObj, err = gtkconv.ToUintE("b") // string
	errLog(t, err)
	if assert.Error(err) {
		assert.Equal(uint(0), actualObj)
	}
}
