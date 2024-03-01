/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-03-01 21:00:42
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-03-01 21:04:08
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkbinary

// Bit Binary bit (0 | 1)
type Bit int8

// EncodeBits
func EncodeBits(bits []Bit, i int, l int) (ts []Bit) {
	return EncodeBitsWithUint(bits, uint(i), l)
}

// EncodeBitsWithUint
func EncodeBitsWithUint(bits []Bit, ui uint, l int) (ts []Bit) {
	a := make([]Bit, l)
	for i := l - 1; i >= 0; i-- {
		a[i] = Bit(ui & 1)
		ui >>= 1
	}
	if bits != nil {
		return append(bits, a...)
	}
	return a
}

// EncodeBitsToBytes
func EncodeBitsToBytes(bits []Bit) (bs []byte) {
	if len(bits)%8 != 0 {
		for i := 0; i < len(bits)%8; i++ {
			bits = append(bits, 0)
		}
	}
	b := make([]byte, 0)
	for i := 0; i < len(bits); i += 8 {
		b = append(b, byte(DecodeBitsToUint(bits[i:i+8])))
	}
	return b
}

// DecodeBits
func DecodeBits(bits []Bit) (val int) {
	v := 0
	for _, i := range bits {
		v = v<<1 | int(i)
	}
	return v
}

// DecodeBitsToUint
func DecodeBitsToUint(bits []Bit) (val uint) {
	v := uint(0)
	for _, i := range bits {
		v = v<<1 | uint(i)
	}
	return v
}

// DecodeBytesToBits
func DecodeBytesToBits(bs []byte) (ts []Bit) {
	bits := make([]Bit, 0)
	for _, b := range bs {
		bits = EncodeBitsWithUint(bits, uint(b), 8)
	}
	return bits
}
