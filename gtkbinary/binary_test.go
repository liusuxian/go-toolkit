/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-03-01 21:09:12
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-03-01 21:19:22
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkbinary_test

import (
	"github.com/liusuxian/go-toolkit/gtkbinary"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

type User struct {
	Name string
	Age  int
	Url  string
}

var testData = map[string]interface{}{
	"int":         int(123),
	"int8":        int8(-99),
	"int8.max":    math.MaxInt8,
	"int16":       int16(123),
	"int16.max":   math.MaxInt16,
	"int32":       int32(-199),
	"int32.max":   math.MaxInt32,
	"int64":       int64(123),
	"uint":        uint(123),
	"uint8":       uint8(123),
	"uint8.max":   math.MaxUint8,
	"uint16":      uint16(9999),
	"uint16.max":  math.MaxUint16,
	"uint32":      uint32(123),
	"uint64":      uint64(123),
	"bool.true":   true,
	"bool.false":  false,
	"string":      "hehe haha",
	"byte":        []byte("hehe haha"),
	"float32":     float32(123.456),
	"float32.max": math.MaxFloat32,
	"float64":     float64(123.456),
}

var testBitData = []int{0, 99, 122, 129, 222, 999, 22322}

func TestEncodeAndDecode(t *testing.T) {
	assert := assert.New(t)
	for k, v := range testData {
		ve := gtkbinary.Encode(v)
		ve1 := gtkbinary.EncodeByLength(len(ve), v)
		switch v.(type) {
		case int:
			assert.Equal(gtkbinary.DecodeToInt(ve), v)
			assert.Equal(gtkbinary.DecodeToInt(ve1), v)
		case int8:
			assert.Equal(gtkbinary.DecodeToInt8(ve), v)
			assert.Equal(gtkbinary.DecodeToInt8(ve1), v)
		case int16:
			assert.Equal(gtkbinary.DecodeToInt16(ve), v)
			assert.Equal(gtkbinary.DecodeToInt16(ve1), v)
		case int32:
			assert.Equal(gtkbinary.DecodeToInt32(ve), v)
			assert.Equal(gtkbinary.DecodeToInt32(ve1), v)
		case int64:
			assert.Equal(gtkbinary.DecodeToInt64(ve), v)
			assert.Equal(gtkbinary.DecodeToInt64(ve1), v)
		case uint:
			assert.Equal(gtkbinary.DecodeToUint(ve), v)
			assert.Equal(gtkbinary.DecodeToUint(ve1), v)
		case uint8:
			assert.Equal(gtkbinary.DecodeToUint8(ve), v)
			assert.Equal(gtkbinary.DecodeToUint8(ve1), v)
		case uint16:
			assert.Equal(gtkbinary.DecodeToUint16(ve1), v)
			assert.Equal(gtkbinary.DecodeToUint16(ve), v)
		case uint32:
			assert.Equal(gtkbinary.DecodeToUint32(ve1), v)
			assert.Equal(gtkbinary.DecodeToUint32(ve), v)
		case uint64:
			assert.Equal(gtkbinary.DecodeToUint64(ve), v)
			assert.Equal(gtkbinary.DecodeToUint64(ve1), v)
		case bool:
			assert.Equal(gtkbinary.DecodeToBool(ve), v)
			assert.Equal(gtkbinary.DecodeToBool(ve1), v)
		case string:
			assert.Equal(gtkbinary.DecodeToString(ve), v)
			assert.Equal(gtkbinary.DecodeToString(ve1), v)
		case float32:
			assert.Equal(gtkbinary.DecodeToFloat32(ve), v)
			assert.Equal(gtkbinary.DecodeToFloat32(ve1), v)
		case float64:
			assert.Equal(gtkbinary.DecodeToFloat64(ve), v)
			assert.Equal(gtkbinary.DecodeToFloat64(ve1), v)
		default:
			if v == nil {
				continue
			}
			res := make([]byte, len(ve))
			err := gtkbinary.Decode(ve, res)
			if err != nil {
				t.Errorf("test data: %s, %v, error:%v", k, v, err)
			}
			assert.Equal(res, v)
		}
	}
}

func TestEncodeStruct(t *testing.T) {
	assert := assert.New(t)
	user := User{"wenzi1", 999, "www.baidu.com"}
	ve := gtkbinary.Encode(user)
	s := gtkbinary.DecodeToString(ve)
	assert.Equal("{wenzi1 999 www.baidu.com}", s)
}

func TestBits(t *testing.T) {
	assert := assert.New(t)
	for i := range testBitData {
			bits := make([]gtkbinary.Bit, 0)
			res := gtkbinary.EncodeBits(bits, testBitData[i], 64)
			assert.Equal(gtkbinary.DecodeBits(res), testBitData[i])
			assert.Equal(gtkbinary.DecodeBitsToUint(res), uint(testBitData[i]))
			assert.Equal(gtkbinary.DecodeBytesToBits(gtkbinary.EncodeBitsToBytes(res)), res)
		}
}
