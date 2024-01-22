/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-05-05 14:19:25
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-22 22:39:36
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

type FFFF struct {
	A any
	B any
	C any
}

type GGGG struct {
	A any `json:"a" dc:"a"`
	B any `json:"b" dc:"b"`
	C any `json:"c" dc:"c"`
}

func TestToStringMapUint64E(t *testing.T) {
	assert := assert.New(t)
	actualObj, err := gtkconv.ToStringMapUint64E(map[any]any{"a": "1", "b": 2.6, "c": true}) // map[any]any
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint64{"a": 1, "b": 2, "c": 1}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint64E([]byte(`{"a": "1.6", "b": 2.7, "c": true}`)) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint64{"a": 1, "b": 2, "c": 1}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint64E(`{"a": "1.6", "b": 2.7, "c": true}`) // json
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint64{"a": 1, "b": 2, "c": 1}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint64E(map[string]string{"a": "1.6", "b": "2.7", "c": "3.1"}) // map[string]string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint64{"a": 1, "b": 2, "c": 3}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint64E(FFFF{A: 1.6, B: false, C: "2.7"}) // struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint64{"A": 1, "B": 0, "C": 2}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint64E(&FFFF{A: 1.6, B: false, C: "2.7"}) // *struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint64{"A": 1, "B": 0, "C": 2}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint64E(GGGG{A: 1.6, B: false, C: "2.7"}) // struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint64{"a": 1, "b": 0, "c": 2}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint64E(&GGGG{A: 1.6, B: false, C: "2.7"}) // *struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint64{"a": 1, "b": 0, "c": 2}, actualObj)
	}
}

func TestToStringMapUint32E(t *testing.T) {
	assert := assert.New(t)
	actualObj, err := gtkconv.ToStringMapUint32E(map[any]any{"a": "1", "b": 2.6, "c": true}) // map[any]any
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint32{"a": 1, "b": 2, "c": 1}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint32E([]byte(`{"a": "1.6", "b": 2.7, "c": true}`)) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint32{"a": 1, "b": 2, "c": 1}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint32E(`{"a": "1.6", "b": 2.7, "c": true}`) // json
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint32{"a": 1, "b": 2, "c": 1}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint32E(map[string]string{"a": "1.6", "b": "2.7", "c": "3.1"}) // map[string]string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint32{"a": 1, "b": 2, "c": 3}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint32E(FFFF{A: 1.6, B: false, C: "2.7"}) // struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint32{"A": 1, "B": 0, "C": 2}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint32E(&FFFF{A: 1.6, B: false, C: "2.7"}) // *struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint32{"A": 1, "B": 0, "C": 2}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint32E(GGGG{A: 1.6, B: false, C: "2.7"}) // struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint32{"a": 1, "b": 0, "c": 2}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint32E(&GGGG{A: 1.6, B: false, C: "2.7"}) // *struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint32{"a": 1, "b": 0, "c": 2}, actualObj)
	}
}

func TestToStringMapUint16E(t *testing.T) {
	assert := assert.New(t)
	actualObj, err := gtkconv.ToStringMapUint16E(map[any]any{"a": "1", "b": 2.6, "c": true}) // map[any]any
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint16{"a": 1, "b": 2, "c": 1}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint16E([]byte(`{"a": "1.6", "b": 2.7, "c": true}`)) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint16{"a": 1, "b": 2, "c": 1}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint16E(`{"a": "1.6", "b": 2.7, "c": true}`) // json
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint16{"a": 1, "b": 2, "c": 1}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint16E(map[string]string{"a": "1.6", "b": "2.7", "c": "3.1"}) // map[string]string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint16{"a": 1, "b": 2, "c": 3}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint16E(FFFF{A: 1.6, B: false, C: "2.7"}) // struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint16{"A": 1, "B": 0, "C": 2}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint16E(&FFFF{A: 1.6, B: false, C: "2.7"}) // *struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint16{"A": 1, "B": 0, "C": 2}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint16E(GGGG{A: 1.6, B: false, C: "2.7"}) // struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint16{"a": 1, "b": 0, "c": 2}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint16E(&GGGG{A: 1.6, B: false, C: "2.7"}) // *struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint16{"a": 1, "b": 0, "c": 2}, actualObj)
	}
}

func TestToStringMapUint8E(t *testing.T) {
	assert := assert.New(t)
	actualObj, err := gtkconv.ToStringMapUint8E(map[any]any{"a": "1", "b": 2.6, "c": true}) // map[any]any
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint8{"a": 1, "b": 2, "c": 1}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint8E([]byte(`{"a": "1.6", "b": 2.7, "c": true}`)) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint8{"a": 1, "b": 2, "c": 1}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint8E(`{"a": "1.6", "b": 2.7, "c": true}`) // json
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint8{"a": 1, "b": 2, "c": 1}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint8E(map[string]string{"a": "1.6", "b": "2.7", "c": "3.1"}) // map[string]string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint8{"a": 1, "b": 2, "c": 3}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint8E(FFFF{A: 1.6, B: false, C: "2.7"}) // struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint8{"A": 1, "B": 0, "C": 2}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint8E(&FFFF{A: 1.6, B: false, C: "2.7"}) // *struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint8{"A": 1, "B": 0, "C": 2}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint8E(GGGG{A: 1.6, B: false, C: "2.7"}) // struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint8{"a": 1, "b": 0, "c": 2}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUint8E(&GGGG{A: 1.6, B: false, C: "2.7"}) // *struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint8{"a": 1, "b": 0, "c": 2}, actualObj)
	}
}

func TestToStringMapUintE(t *testing.T) {
	assert := assert.New(t)
	actualObj, err := gtkconv.ToStringMapUintE(map[any]any{"a": "1", "b": 2.6, "c": true}) // map[any]any
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint{"a": 1, "b": 2, "c": 1}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUintE([]byte(`{"a": "1.6", "b": 2.7, "c": true}`)) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint{"a": 1, "b": 2, "c": 1}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUintE(`{"a": "1.6", "b": 2.7, "c": true}`) // json
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint{"a": 1, "b": 2, "c": 1}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUintE(map[string]string{"a": "1.6", "b": "2.7", "c": "3.1"}) // map[string]string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint{"a": 1, "b": 2, "c": 3}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUintE(FFFF{A: 1.6, B: false, C: "2.7"}) // struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint{"A": 1, "B": 0, "C": 2}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUintE(&FFFF{A: 1.6, B: false, C: "2.7"}) // *struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint{"A": 1, "B": 0, "C": 2}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUintE(GGGG{A: 1.6, B: false, C: "2.7"}) // struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint{"a": 1, "b": 0, "c": 2}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapUintE(&GGGG{A: 1.6, B: false, C: "2.7"}) // *struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]uint{"a": 1, "b": 0, "c": 2}, actualObj)
	}
}
