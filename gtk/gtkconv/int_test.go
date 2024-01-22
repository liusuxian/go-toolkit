/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-04-15 13:38:55
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-22 22:33:28
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

func TestToInt64E(t *testing.T) {
	assert := assert.New(t)
	actualObj, err := gtkconv.ToInt64E(nil) // nil
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int64(0), actualObj)
	}
	actualObj, err = gtkconv.ToInt64E(int(1)) // int
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int64(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt64E(float64(1.56)) // float64
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int64(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt64E(true) // bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int64(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt64E(false) // bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int64(0), actualObj)
	}
	actualObj, err = gtkconv.ToInt64E([]byte("1.23")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int64(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt64E([]byte("1.0")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int64(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt64E([]byte("1.")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int64(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt64E([]byte("1")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int64(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt64E([]byte("a")) // []byte
	errLog(t, err)
	if assert.Error(err) {
		assert.Equal(int64(0), actualObj)
	}
	actualObj, err = gtkconv.ToInt64E("1.23") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int64(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt64E("1.0") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int64(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt64E("1.") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int64(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt64E("1") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int64(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt64E("b") // string
	errLog(t, err)
	if assert.Error(err) {
		assert.Equal(int64(0), actualObj)
	}
}

func TestToInt32E(t *testing.T) {
	assert := assert.New(t)
	actualObj, err := gtkconv.ToInt32E(nil) // nil
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int32(0), actualObj)
	}
	actualObj, err = gtkconv.ToInt32E(int(1)) // int
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int32(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt32E(float64(1.56)) // float64
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int32(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt32E(true) // bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int32(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt32E(false) // bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int32(0), actualObj)
	}
	actualObj, err = gtkconv.ToInt32E([]byte("1.23")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int32(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt32E([]byte("1.0")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int32(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt32E([]byte("1.")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int32(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt32E([]byte("1")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int32(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt32E([]byte("a")) // []byte
	errLog(t, err)
	if assert.Error(err) {
		assert.Equal(int32(0), actualObj)
	}
	actualObj, err = gtkconv.ToInt32E("1.23") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int32(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt32E("1.0") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int32(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt32E("1.") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int32(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt32E("1") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int32(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt32E("b") // string
	errLog(t, err)
	if assert.Error(err) {
		assert.Equal(int32(0), actualObj)
	}
}

func TestToInt16E(t *testing.T) {
	assert := assert.New(t)
	actualObj, err := gtkconv.ToInt16E(nil) // nil
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int16(0), actualObj)
	}
	actualObj, err = gtkconv.ToInt16E(int(1)) // int
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int16(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt16E(float64(1.56)) // float64
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int16(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt16E(true) // bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int16(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt16E(false) // bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int16(0), actualObj)
	}
	actualObj, err = gtkconv.ToInt16E([]byte("1.23")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int16(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt16E([]byte("1.0")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int16(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt16E([]byte("1.")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int16(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt16E([]byte("1")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int16(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt16E([]byte("a")) // []byte
	errLog(t, err)
	if assert.Error(err) {
		assert.Equal(int16(0), actualObj)
	}
	actualObj, err = gtkconv.ToInt16E("1.23") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int16(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt16E("1.0") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int16(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt16E("1.") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int16(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt16E("1") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int16(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt16E("b") // string
	errLog(t, err)
	if assert.Error(err) {
		assert.Equal(int16(0), actualObj)
	}
}

func TestToInt8E(t *testing.T) {
	assert := assert.New(t)
	actualObj, err := gtkconv.ToInt8E(nil) // nil
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int8(0), actualObj)
	}
	actualObj, err = gtkconv.ToInt8E(int(1)) // int
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int8(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt8E(float64(1.56)) // float64
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int8(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt8E(true) // bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int8(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt8E(false) // bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int8(0), actualObj)
	}
	actualObj, err = gtkconv.ToInt8E([]byte("1.23")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int8(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt8E([]byte("1.0")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int8(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt8E([]byte("1.")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int8(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt8E([]byte("1")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int8(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt8E([]byte("a")) // []byte
	errLog(t, err)
	if assert.Error(err) {
		assert.Equal(int8(0), actualObj)
	}
	actualObj, err = gtkconv.ToInt8E("1.23") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int8(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt8E("1.0") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int8(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt8E("1.") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int8(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt8E("1") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int8(1), actualObj)
	}
	actualObj, err = gtkconv.ToInt8E("b") // string
	errLog(t, err)
	if assert.Error(err) {
		assert.Equal(int8(0), actualObj)
	}
}

func TestToIntE(t *testing.T) {
	assert := assert.New(t)
	actualObj, err := gtkconv.ToIntE(nil) // nil
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int(0), actualObj)
	}
	actualObj, err = gtkconv.ToIntE(int(1)) // int
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int(1), actualObj)
	}
	actualObj, err = gtkconv.ToIntE(float64(1.56)) // float64
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int(1), actualObj)
	}
	actualObj, err = gtkconv.ToIntE(true) // bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int(1), actualObj)
	}
	actualObj, err = gtkconv.ToIntE(false) // bool
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int(0), actualObj)
	}
	actualObj, err = gtkconv.ToIntE([]byte("1.23")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int(1), actualObj)
	}
	actualObj, err = gtkconv.ToIntE([]byte("1.0")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int(1), actualObj)
	}
	actualObj, err = gtkconv.ToIntE([]byte("1.")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int(1), actualObj)
	}
	actualObj, err = gtkconv.ToIntE([]byte("1")) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int(1), actualObj)
	}
	actualObj, err = gtkconv.ToIntE([]byte("a")) // []byte
	errLog(t, err)
	if assert.Error(err) {
		assert.Equal(int(0), actualObj)
	}
	actualObj, err = gtkconv.ToIntE("1.23") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int(1), actualObj)
	}
	actualObj, err = gtkconv.ToIntE("1.0") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int(1), actualObj)
	}
	actualObj, err = gtkconv.ToIntE("1.") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int(1), actualObj)
	}
	actualObj, err = gtkconv.ToIntE("1") // string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(int(1), actualObj)
	}
	actualObj, err = gtkconv.ToIntE("b") // []byte
	errLog(t, err)
	if assert.Error(err) {
		assert.Equal(int(0), actualObj)
	}
}
