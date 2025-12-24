/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-03-01 20:35:57
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-24 18:48:38
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkbinary

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"html/template"
	"math"
)

// LeEncode
func LeEncode(vals ...any) (bs []byte) {
	buf := new(bytes.Buffer)
	for i := range vals {
		if vals[i] == nil {
			return buf.Bytes()
		}
		switch val := vals[i].(type) {
		case int:
			buf.Write(LeEncodeInt(val))
		case int8:
			buf.Write(LeEncodeInt8(val))
		case int16:
			buf.Write(LeEncodeInt16(val))
		case int32:
			buf.Write(LeEncodeInt32(val))
		case int64:
			buf.Write(LeEncodeInt64(val))
		case uint:
			buf.Write(LeEncodeUint(val))
		case uint8:
			buf.Write(LeEncodeUint8(val))
		case uint16:
			buf.Write(LeEncodeUint16(val))
		case uint32:
			buf.Write(LeEncodeUint32(val))
		case uint64:
			buf.Write(LeEncodeUint64(val))
		case bool:
			buf.Write(LeEncodeBool(val))
		case string:
			buf.WriteString(val)
		case []byte:
			buf.Write(val)
		case float32:
			buf.Write(LeEncodeFloat32(val))
		case float64:
			buf.Write(LeEncodeFloat64(val))
		case json.Number:
			buf.WriteString(val.String())
		case template.HTML:
			buf.WriteString(string(val))
		case template.URL:
			buf.WriteString(string(val))
		case template.JS:
			buf.WriteString(string(val))
		case template.CSS:
			buf.WriteString(string(val))
		case template.HTMLAttr:
			buf.WriteString(string(val))
		case fmt.Stringer:
			buf.WriteString(val.String())
		case error:
			buf.WriteString(val.Error())
		default:
			if err := binary.Write(buf, binary.LittleEndian, val); err != nil {
				if jsonData, e := json.Marshal(val); e == nil {
					buf.Write(jsonData)
				} else {
					fmt.Fprint(buf, val)
				}
			}
		}
	}
	return buf.Bytes()
}

// LeEncodeByLength
func LeEncodeByLength(length int, vals ...any) (bs []byte) {
	b := LeEncode(vals...)
	if len(b) < length {
		b = append(b, make([]byte, length-len(b))...)
	} else if len(b) > length {
		b = b[0:length]
	}
	return b
}

// LeDecode
func LeDecode(b []byte, vals ...any) (err error) {
	var buf = bytes.NewBuffer(b)
	for i := range vals {
		if err = binary.Read(buf, binary.LittleEndian, vals[i]); err != nil {
			err = fmt.Errorf("binary.read failed: %w", err)
			return
		}
	}
	return
}

// LeEncodeString
func LeEncodeString(s string) (bs []byte) {
	return []byte(s)
}

// LeDecodeToString
func LeDecodeToString(b []byte) (s string) {
	return string(b)
}

// LeEncodeBool
func LeEncodeBool(b bool) (bs []byte) {
	if b {
		return []byte{1}
	} else {
		return []byte{0}
	}
}

// LeEncodeInt
func LeEncodeInt(i int) (bs []byte) {
	if i <= math.MaxInt8 {
		return EncodeInt8(int8(i))
	} else if i <= math.MaxInt16 {
		return EncodeInt16(int16(i))
	} else if i <= math.MaxInt32 {
		return EncodeInt32(int32(i))
	} else {
		return EncodeInt64(int64(i))
	}
}

// LeEncodeUint
func LeEncodeUint(i uint) (bs []byte) {
	if i <= math.MaxUint8 {
		return EncodeUint8(uint8(i))
	} else if i <= math.MaxUint16 {
		return EncodeUint16(uint16(i))
	} else if i <= math.MaxUint32 {
		return EncodeUint32(uint32(i))
	} else {
		return EncodeUint64(uint64(i))
	}
}

// LeEncodeInt8
func LeEncodeInt8(i int8) (bs []byte) {
	return []byte{byte(i)}
}

// LeEncodeUint8
func LeEncodeUint8(i uint8) (bs []byte) {
	return []byte{i}
}

// LeEncodeInt16
func LeEncodeInt16(i int16) (bs []byte) {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(i))
	return b
}

// LeEncodeUint16
func LeEncodeUint16(i uint16) (bs []byte) {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, i)
	return b
}

// LeEncodeInt32
func LeEncodeInt32(i int32) (bs []byte) {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(i))
	return b
}

// LeEncodeUint32
func LeEncodeUint32(i uint32) (bs []byte) {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, i)
	return b
}

// LeEncodeInt64
func LeEncodeInt64(i int64) (bs []byte) {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return b
}

// LeEncodeUint64
func LeEncodeUint64(i uint64) (bs []byte) {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, i)
	return b
}

// LeEncodeFloat32
func LeEncodeFloat32(f float32) (bs []byte) {
	bits := math.Float32bits(f)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, bits)
	return b
}

// LeEncodeFloat64
func LeEncodeFloat64(f float64) (bs []byte) {
	bits := math.Float64bits(f)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, bits)
	return b
}

// LeDecodeToInt
func LeDecodeToInt(b []byte) (val int) {
	if len(b) < 2 {
		return int(LeDecodeToUint8(b))
	} else if len(b) < 3 {
		return int(LeDecodeToUint16(b))
	} else if len(b) < 5 {
		return int(LeDecodeToUint32(b))
	} else {
		return int(LeDecodeToUint64(b))
	}
}

// LeDecodeToUint
func LeDecodeToUint(b []byte) (val uint) {
	if len(b) < 2 {
		return uint(LeDecodeToUint8(b))
	} else if len(b) < 3 {
		return uint(LeDecodeToUint16(b))
	} else if len(b) < 5 {
		return uint(LeDecodeToUint32(b))
	} else {
		return uint(LeDecodeToUint64(b))
	}
}

// LeDecodeToBool
func LeDecodeToBool(b []byte) (ok bool) {
	if len(b) == 0 {
		return false
	}
	if bytes.Equal(b, make([]byte, len(b))) {
		return false
	}
	return true
}

// LeDecodeToInt8
func LeDecodeToInt8(b []byte) (val int8) {
	if len(b) == 0 {
		panic(`empty slice given`)
	}
	return int8(b[0])
}

// LeDecodeToUint8
func LeDecodeToUint8(b []byte) (val uint8) {
	if len(b) == 0 {
		panic(`empty slice given`)
	}
	return b[0]
}

// LeDecodeToInt16
func LeDecodeToInt16(b []byte) (val int16) {
	return int16(binary.LittleEndian.Uint16(LeFillUpSize(b, 2)))
}

// LeDecodeToUint16
func LeDecodeToUint16(b []byte) (val uint16) {
	return binary.LittleEndian.Uint16(LeFillUpSize(b, 2))
}

// LeDecodeToInt32
func LeDecodeToInt32(b []byte) (val int32) {
	return int32(binary.LittleEndian.Uint32(LeFillUpSize(b, 4)))
}

// LeDecodeToUint32
func LeDecodeToUint32(b []byte) (val uint32) {
	return binary.LittleEndian.Uint32(LeFillUpSize(b, 4))
}

// LeDecodeToInt64
func LeDecodeToInt64(b []byte) (val int64) {
	return int64(binary.LittleEndian.Uint64(LeFillUpSize(b, 8)))
}

// LeDecodeToUint64
func LeDecodeToUint64(b []byte) (val uint64) {
	return binary.LittleEndian.Uint64(LeFillUpSize(b, 8))
}

// LeDecodeToFloat32
func LeDecodeToFloat32(b []byte) (val float32) {
	return math.Float32frombits(binary.LittleEndian.Uint32(LeFillUpSize(b, 4)))
}

// LeDecodeToFloat64
func LeDecodeToFloat64(b []byte) (val float64) {
	return math.Float64frombits(binary.LittleEndian.Uint64(LeFillUpSize(b, 8)))
}

// LeFillUpSize
func LeFillUpSize(b []byte, l int) (bs []byte) {
	if len(b) >= l {
		return b[:l]
	}
	c := make([]byte, l)
	copy(c, b)
	return c
}
