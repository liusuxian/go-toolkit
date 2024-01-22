/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-05-04 14:22:21
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2023-05-14 02:16:46
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

type AddressNoTag struct {
	Street string
	City   string
}

type PersonNoTag struct {
	Name    string
	Age     int
	Address AddressNoTag
	sex     int
}

type AddressTag struct {
	Street string `json:"street" struct:"street1" dc:"street"`
	City   string `json:"city" struct:"city1" dc:"city"`
}

type PersonTag struct {
	Name    string      `json:"name" struct:"name1" dc:"name"`
	Age     int         `json:"age" struct:"age1" dc:"age"`
	Address *AddressTag `json:"address" struct:"address1" dc:"address"`
	sex     int         `json:"sex" struct:"sex1" dc:"sex"`
}

func TestToStringMapE(t *testing.T) {
	assert := assert.New(t)
	actualObj, err := gtkconv.ToStringMapE(map[any]any{"a": "hello", "b": []any{"hello", "true"}, "c": map[string]any{"a": "hello", "b": []any{"hello", "true"}}}) // map[any]any
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]any{"a": "hello", "b": []any{"hello", "true"}, "c": map[string]any{"a": "hello", "b": []any{"hello", "true"}}}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapE([]byte(`{"a":"hello","b":["hello","true"],"c":{"a":"hello","b":["hello","true"]}}`)) // []byte
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]any{"a": "hello", "b": []any{"hello", "true"}, "c": map[string]any{"a": "hello", "b": []any{"hello", "true"}}}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapE(`{"a":"hello","b":["hello","true"],"c":{"a":"hello","b":["hello","true"]}}`) // json
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]any{"a": "hello", "b": []any{"hello", "true"}, "c": map[string]any{"a": "hello", "b": []any{"hello", "true"}}}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapE(map[string]string{"a": "hello", "b": "world", "c": "true"}) // map[string]string
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]any{"a": "hello", "b": "world", "c": "true"}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapE(PersonNoTag{Name: "lsx", Age: 18, Address: AddressNoTag{Street: "hz-123", City: "hz"}, sex: 1}) // struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]any{"Name": "lsx", "Age": 18, "Address": map[string]any{"City": "hz", "Street": "hz-123"}}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapE(&PersonNoTag{Name: "lsx", Age: 18, Address: AddressNoTag{Street: "hz-123", City: "hz"}, sex: 1}) // *struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]any{"Name": "lsx", "Age": 18, "Address": map[string]any{"City": "hz", "Street": "hz-123"}}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapE(PersonTag{Name: "lsx", Age: 18, Address: &AddressTag{Street: "hz-123", City: "hz"}, sex: 1}) // struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]any{"name": "lsx", "age": 18, "address": map[string]any{"city": "hz", "street": "hz-123"}}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapE(&PersonTag{Name: "lsx", Age: 18, Address: &AddressTag{Street: "hz-123", City: "hz"}, sex: 1}) // *struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]any{"name": "lsx", "age": 18, "address": map[string]any{"city": "hz", "street": "hz-123"}}, actualObj)
	}
	actualObj, err = gtkconv.ToStringMapE(&PersonTag{Name: "lsx", Age: 18, Address: &AddressTag{Street: "hz-123", City: "hz"}, sex: 1}, func(dc *gtkconv.DecoderConfig) {
		dc.TagName = "struct"
	}) // *struct
	errLog(t, err)
	if assert.NoError(err) {
		assert.Equal(map[string]any{"name1": "lsx", "age1": 18, "address1": map[string]any{"city1": "hz", "street1": "hz-123"}}, actualObj)
	}
}
