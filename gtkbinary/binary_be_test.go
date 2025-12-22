/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-03-01 21:25:11
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-20 22:21:47
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkbinary_test

import (
	"encoding/binary"
	"github.com/liusuxian/go-toolkit/gtkbinary"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBeEncodeAndBeDecode(t *testing.T) {
	assert := assert.New(t)
	for k, v := range testData {
		ve := gtkbinary.BeEncode(v)
		ve1 := gtkbinary.BeEncodeByLength(len(ve), v)

		switch v.(type) {
		case int:
			assert.Equal(gtkbinary.BeDecodeToInt(ve), v)
			assert.Equal(gtkbinary.BeDecodeToInt(ve1), v)
		case int8:
			assert.Equal(gtkbinary.BeDecodeToInt8(ve), v)
			assert.Equal(gtkbinary.BeDecodeToInt8(ve1), v)
		case int16:
			assert.Equal(gtkbinary.BeDecodeToInt16(ve), v)
			assert.Equal(gtkbinary.BeDecodeToInt16(ve1), v)
		case int32:
			assert.Equal(gtkbinary.BeDecodeToInt32(ve), v)
			assert.Equal(gtkbinary.BeDecodeToInt32(ve1), v)
		case int64:
			assert.Equal(gtkbinary.BeDecodeToInt64(ve), v)
			assert.Equal(gtkbinary.BeDecodeToInt64(ve1), v)
		case uint:
			assert.Equal(gtkbinary.BeDecodeToUint(ve), v)
			assert.Equal(gtkbinary.BeDecodeToUint(ve1), v)
		case uint8:
			assert.Equal(gtkbinary.BeDecodeToUint8(ve), v)
			assert.Equal(gtkbinary.BeDecodeToUint8(ve1), v)
		case uint16:
			assert.Equal(gtkbinary.BeDecodeToUint16(ve1), v)
			assert.Equal(gtkbinary.BeDecodeToUint16(ve), v)
		case uint32:
			assert.Equal(gtkbinary.BeDecodeToUint32(ve1), v)
			assert.Equal(gtkbinary.BeDecodeToUint32(ve), v)
		case uint64:
			assert.Equal(gtkbinary.BeDecodeToUint64(ve), v)
			assert.Equal(gtkbinary.BeDecodeToUint64(ve1), v)
		case bool:
			assert.Equal(gtkbinary.BeDecodeToBool(ve), v)
			assert.Equal(gtkbinary.BeDecodeToBool(ve1), v)
		case string:
			assert.Equal(gtkbinary.BeDecodeToString(ve), v)
			assert.Equal(gtkbinary.BeDecodeToString(ve1), v)
		case float32:
			assert.Equal(gtkbinary.BeDecodeToFloat32(ve), v)
			assert.Equal(gtkbinary.BeDecodeToFloat32(ve1), v)
		case float64:
			assert.Equal(gtkbinary.BeDecodeToFloat64(ve), v)
			assert.Equal(gtkbinary.BeDecodeToFloat64(ve1), v)
		default:
			if v == nil {
				continue
			}
			res := make([]byte, len(ve))
			err := gtkbinary.BeDecode(ve, res)
			if err != nil {
				t.Errorf("test data: %s, %v, error:%v", k, v, err)
			}
			assert.Equal(res, v)
		}
	}
}

func TestBeEncodeStruct(t *testing.T) {
	assert := assert.New(t)
	user := User{"wenzi1", 999, "www.baidu.com"}
	ve := gtkbinary.BeEncode(user)
	s := gtkbinary.BeDecodeToString(ve)
	assert.Equal("{wenzi1 999 www.baidu.com}", s)
}

func TestBeEncodeInt(t *testing.T) {
	assert := assert.New(t)
	expected := make([]byte, 4)
	binary.BigEndian.PutUint32(expected, uint32(123456))
	actual := gtkbinary.BeEncode(123456)
	assert.Equal(expected, actual)
}
