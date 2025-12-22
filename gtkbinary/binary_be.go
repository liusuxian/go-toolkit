/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-03-01 20:51:11
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-20 22:17:36
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkbinary

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

// BeEncode
func BeEncode(vals ...any) (bs []byte) {
	buf := new(bytes.Buffer)
	for i := range vals {
		if vals[i] == nil {
			return buf.Bytes()
		}

		switch val := vals[i].(type) {
		case int:
			buf.Write(BeEncodeInt(val))
		case int8:
			buf.Write(BeEncodeInt8(val))
		case int16:
			buf.Write(BeEncodeInt16(val))
		case int32:
			buf.Write(BeEncodeInt32(val))
		case int64:
			buf.Write(BeEncodeInt64(val))
		case uint:
			buf.Write(BeEncodeUint(val))
		case uint8:
			buf.Write(BeEncodeUint8(val))
		case uint16:
			buf.Write(BeEncodeUint16(val))
		case uint32:
			buf.Write(BeEncodeUint32(val))
		case uint64:
			buf.Write(BeEncodeUint64(val))
		case bool:
			buf.Write(BeEncodeBool(val))
		case string:
			buf.Write(BeEncodeString(val))
		case []byte:
			buf.Write(val)
		case float32:
			buf.Write(BeEncodeFloat32(val))
		case float64:
			buf.Write(BeEncodeFloat64(val))
		default:
			if err := binary.Write(buf, binary.BigEndian, val); err != nil {
				buf.Write(BeEncodeString(fmt.Sprintf("%v", val)))
			}
		}
	}
	return buf.Bytes()
}

// BeEncodeByLength
func BeEncodeByLength(length int, vals ...any) (bs []byte) {
	b := BeEncode(vals...)
	if len(b) < length {
		b = append(b, make([]byte, length-len(b))...)
	} else if len(b) > length {
		b = b[0:length]
	}
	return b
}

// BeDecode
func BeDecode(b []byte, vals ...any) (err error) {
	var buf = bytes.NewBuffer(b)
	for i := range vals {
		if err = binary.Read(buf, binary.BigEndian, vals[i]); err != nil {
			err = fmt.Errorf("binary.read failed: %w", err)
			return
		}
	}
	return
}

// BeEncodeString
func BeEncodeString(s string) (bs []byte) {
	return []byte(s)
}

// BeDecodeToString
func BeDecodeToString(b []byte) (s string) {
	return string(b)
}

// BeEncodeBool
func BeEncodeBool(b bool) (bs []byte) {
	if b {
		return []byte{1}
	} else {
		return []byte{0}
	}
}

// BeEncodeInt
func BeEncodeInt(i int) (bs []byte) {
	if i <= math.MaxInt8 {
		return BeEncodeInt8(int8(i))
	} else if i <= math.MaxInt16 {
		return BeEncodeInt16(int16(i))
	} else if i <= math.MaxInt32 {
		return BeEncodeInt32(int32(i))
	} else {
		return BeEncodeInt64(int64(i))
	}
}

// BeEncodeUint
func BeEncodeUint(i uint) (bs []byte) {
	if i <= math.MaxUint8 {
		return BeEncodeUint8(uint8(i))
	} else if i <= math.MaxUint16 {
		return BeEncodeUint16(uint16(i))
	} else if i <= math.MaxUint32 {
		return BeEncodeUint32(uint32(i))
	} else {
		return BeEncodeUint64(uint64(i))
	}
}

// BeEncodeInt8
func BeEncodeInt8(i int8) (bs []byte) {
	return []byte{byte(i)}
}

// BeEncodeUint8
func BeEncodeUint8(i uint8) (bs []byte) {
	return []byte{i}
}

// BeEncodeInt16
func BeEncodeInt16(i int16) (bs []byte) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(i))
	return b
}

// BeEncodeUint16
func BeEncodeUint16(i uint16) (bs []byte) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, i)
	return b
}

// BeEncodeInt32
func BeEncodeInt32(i int32) (bs []byte) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(i))
	return b
}

// BeEncodeUint32
func BeEncodeUint32(i uint32) (bs []byte) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, i)
	return b
}

// BeEncodeInt64
func BeEncodeInt64(i int64) (bs []byte) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}

// BeEncodeUint64
func BeEncodeUint64(i uint64) (bs []byte) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return b
}

// BeEncodeFloat32
func BeEncodeFloat32(f float32) (bs []byte) {
	bits := math.Float32bits(f)
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, bits)
	return b
}

// BeEncodeFloat64
func BeEncodeFloat64(f float64) (bs []byte) {
	bits := math.Float64bits(f)
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, bits)
	return b
}

// BeDecodeToInt
func BeDecodeToInt(b []byte) (val int) {
	if len(b) < 2 {
		return int(BeDecodeToUint8(b))
	} else if len(b) < 3 {
		return int(BeDecodeToUint16(b))
	} else if len(b) < 5 {
		return int(BeDecodeToUint32(b))
	} else {
		return int(BeDecodeToUint64(b))
	}
}

// BeDecodeToUint
func BeDecodeToUint(b []byte) (val uint) {
	if len(b) < 2 {
		return uint(BeDecodeToUint8(b))
	} else if len(b) < 3 {
		return uint(BeDecodeToUint16(b))
	} else if len(b) < 5 {
		return uint(BeDecodeToUint32(b))
	} else {
		return uint(BeDecodeToUint64(b))
	}
}

// BeDecodeToBool
func BeDecodeToBool(b []byte) (ok bool) {
	if len(b) == 0 {
		return false
	}
	if bytes.Equal(b, make([]byte, len(b))) {
		return false
	}
	return true
}

// BeDecodeToInt8
func BeDecodeToInt8(b []byte) (val int8) {
	if len(b) == 0 {
		panic(`empty slice given`)
	}
	return int8(b[0])
}

// BeDecodeToUint8
func BeDecodeToUint8(b []byte) (val uint8) {
	if len(b) == 0 {
		panic(`empty slice given`)
	}
	return b[0]
}

// BeDecodeToInt16
func BeDecodeToInt16(b []byte) (val int16) {
	return int16(binary.BigEndian.Uint16(BeFillUpSize(b, 2)))
}

// BeDecodeToUint16
func BeDecodeToUint16(b []byte) (val uint16) {
	return binary.BigEndian.Uint16(BeFillUpSize(b, 2))
}

// BeDecodeToInt32
func BeDecodeToInt32(b []byte) (val int32) {
	return int32(binary.BigEndian.Uint32(BeFillUpSize(b, 4)))
}

// BeDecodeToUint32
func BeDecodeToUint32(b []byte) (val uint32) {
	return binary.BigEndian.Uint32(BeFillUpSize(b, 4))
}

// BeDecodeToInt64
func BeDecodeToInt64(b []byte) (val int64) {
	return int64(binary.BigEndian.Uint64(BeFillUpSize(b, 8)))
}

// BeDecodeToUint64
func BeDecodeToUint64(b []byte) (val uint64) {
	return binary.BigEndian.Uint64(BeFillUpSize(b, 8))
}

// BeDecodeToFloat32
func BeDecodeToFloat32(b []byte) (val float32) {
	return math.Float32frombits(binary.BigEndian.Uint32(BeFillUpSize(b, 4)))
}

// BeDecodeToFloat64
func BeDecodeToFloat64(b []byte) (val float64) {
	return math.Float64frombits(binary.BigEndian.Uint64(BeFillUpSize(b, 8)))
}

// BeFillUpSize
func BeFillUpSize(b []byte, l int) (bs []byte) {
	if len(b) >= l {
		return b[:l]
	}
	c := make([]byte, l)
	copy(c[l-len(b):], b)
	return c
}
