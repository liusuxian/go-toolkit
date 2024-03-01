/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-03-01 20:35:01
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-03-01 20:39:57
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkbinary

// Encode
func Encode(values ...any) (bs []byte) {
	return LeEncode(values...)
}

// EncodeByLength
func EncodeByLength(length int, values ...any) (bs []byte) {
	return LeEncodeByLength(length, values...)
}

// Decode
func Decode(b []byte, values ...any) (err error) {
	return LeDecode(b, values...)
}

// EncodeString
func EncodeString(s string) (bs []byte) {
	return LeEncodeString(s)
}

// DecodeToString
func DecodeToString(b []byte) (s string) {
	return LeDecodeToString(b)
}

// EncodeBool
func EncodeBool(b bool) (bs []byte) {
	return LeEncodeBool(b)
}

// EncodeInt
func EncodeInt(i int) (bs []byte) {
	return LeEncodeInt(i)
}

// EncodeUint
func EncodeUint(i uint) (bs []byte) {
	return LeEncodeUint(i)
}

// EncodeInt8
func EncodeInt8(i int8) (bs []byte) {
	return LeEncodeInt8(i)
}

// EncodeUint8
func EncodeUint8(i uint8) (bs []byte) {
	return LeEncodeUint8(i)
}

// EncodeInt16
func EncodeInt16(i int16) (bs []byte) {
	return LeEncodeInt16(i)
}

// EncodeUint16
func EncodeUint16(i uint16) (bs []byte) {
	return LeEncodeUint16(i)
}

// EncodeInt32
func EncodeInt32(i int32) (bs []byte) {
	return LeEncodeInt32(i)
}

// EncodeUint32
func EncodeUint32(i uint32) (bs []byte) {
	return LeEncodeUint32(i)
}

// EncodeInt64
func EncodeInt64(i int64) (bs []byte) {
	return LeEncodeInt64(i)
}

// EncodeUint64
func EncodeUint64(i uint64) (bs []byte) {
	return LeEncodeUint64(i)
}

// EncodeFloat32
func EncodeFloat32(f float32) (bs []byte) {
	return LeEncodeFloat32(f)
}

// EncodeFloat64
func EncodeFloat64(f float64) (bs []byte) {
	return LeEncodeFloat64(f)
}

// DecodeToInt
func DecodeToInt(b []byte) (val int) {
	return LeDecodeToInt(b)
}

// DecodeToUint
func DecodeToUint(b []byte) (val uint) {
	return LeDecodeToUint(b)
}

// DecodeToBool
func DecodeToBool(b []byte) (ok bool) {
	return LeDecodeToBool(b)
}

// DecodeToInt8
func DecodeToInt8(b []byte) (val int8) {
	return LeDecodeToInt8(b)
}

// DecodeToUint8
func DecodeToUint8(b []byte) (val uint8) {
	return LeDecodeToUint8(b)
}

// DecodeToInt16
func DecodeToInt16(b []byte) (val int16) {
	return LeDecodeToInt16(b)
}

// DecodeToUint16
func DecodeToUint16(b []byte) (val uint16) {
	return LeDecodeToUint16(b)
}

// DecodeToInt32
func DecodeToInt32(b []byte) (val int32) {
	return LeDecodeToInt32(b)
}

// DecodeToUint32
func DecodeToUint32(b []byte) (val uint32) {
	return LeDecodeToUint32(b)
}

// DecodeToInt64
func DecodeToInt64(b []byte) (val int64) {
	return LeDecodeToInt64(b)
}

// DecodeToUint64
func DecodeToUint64(b []byte) (val uint64) {
	return LeDecodeToUint64(b)
}

// DecodeToFloat32
func DecodeToFloat32(b []byte) (val float32) {
	return LeDecodeToFloat32(b)
}

// DecodeToFloat64
func DecodeToFloat64(b []byte) (val float64) {
	return LeDecodeToFloat64(b)
}
