/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-03-01 21:20:27
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-24 19:05:23
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

func TestLeEncodeAndLeDecode(t *testing.T) {
	assert := assert.New(t)
	for k, v := range testData {
		ve := gtkbinary.LeEncode(v)
		ve1 := gtkbinary.LeEncodeByLength(len(ve), v)

		switch v.(type) {
		case int:
			assert.Equal(gtkbinary.LeDecodeToInt(ve), v)
			assert.Equal(gtkbinary.LeDecodeToInt(ve1), v)
		case int8:
			assert.Equal(gtkbinary.LeDecodeToInt8(ve), v)
			assert.Equal(gtkbinary.LeDecodeToInt8(ve1), v)
		case int16:
			assert.Equal(gtkbinary.LeDecodeToInt16(ve), v)
			assert.Equal(gtkbinary.LeDecodeToInt16(ve1), v)
		case int32:
			assert.Equal(gtkbinary.LeDecodeToInt32(ve), v)
			assert.Equal(gtkbinary.LeDecodeToInt32(ve1), v)
		case int64:
			assert.Equal(gtkbinary.LeDecodeToInt64(ve), v)
			assert.Equal(gtkbinary.LeDecodeToInt64(ve1), v)
		case uint:
			assert.Equal(gtkbinary.LeDecodeToUint(ve), v)
			assert.Equal(gtkbinary.LeDecodeToUint(ve1), v)
		case uint8:
			assert.Equal(gtkbinary.LeDecodeToUint8(ve), v)
			assert.Equal(gtkbinary.LeDecodeToUint8(ve1), v)
		case uint16:
			assert.Equal(gtkbinary.LeDecodeToUint16(ve1), v)
			assert.Equal(gtkbinary.LeDecodeToUint16(ve), v)
		case uint32:
			assert.Equal(gtkbinary.LeDecodeToUint32(ve1), v)
			assert.Equal(gtkbinary.LeDecodeToUint32(ve), v)
		case uint64:
			assert.Equal(gtkbinary.LeDecodeToUint64(ve), v)
			assert.Equal(gtkbinary.LeDecodeToUint64(ve1), v)
		case bool:
			assert.Equal(gtkbinary.LeDecodeToBool(ve), v)
			assert.Equal(gtkbinary.LeDecodeToBool(ve1), v)
		case string:
			assert.Equal(gtkbinary.LeDecodeToString(ve), v)
			assert.Equal(gtkbinary.LeDecodeToString(ve1), v)
		case float32:
			assert.Equal(gtkbinary.LeDecodeToFloat32(ve), v)
			assert.Equal(gtkbinary.LeDecodeToFloat32(ve1), v)
		case float64:
			assert.Equal(gtkbinary.LeDecodeToFloat64(ve), v)
			assert.Equal(gtkbinary.LeDecodeToFloat64(ve1), v)
		default:
			if v == nil {
				continue
			}
			res := make([]byte, len(ve))
			err := gtkbinary.LeDecode(ve, res)
			if err != nil {
				t.Errorf("test data: %s, %v, error:%v", k, v, err)
			}
			assert.Equal(res, v)
		}
	}
}

func TestLeEncodeStruct(t *testing.T) {
	assert := assert.New(t)
	user := User{"wenzi1", 999, "www.baidu.com"}
	ve := gtkbinary.LeEncode(user)
	s := gtkbinary.LeDecodeToString(ve)
	assert.Equal("{\"Name\":\"wenzi1\",\"Age\":999,\"Url\":\"www.baidu.com\"}", s)
}

func TestLeEncodeInt(t *testing.T) {
	assert := assert.New(t)
	expected := make([]byte, 4)
	binary.LittleEndian.PutUint32(expected, uint32(123456))
	actual := gtkbinary.LeEncode(123456)
	assert.Equal(expected, actual)
}
