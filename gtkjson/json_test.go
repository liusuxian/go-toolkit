/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-22 23:27:03
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-23 00:17:33
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkjson_test

import (
	"github.com/liusuxian/go-toolkit/gtkjson"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestStruct struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

// TestJsonMarshal 测试 JsonMarshal 函数
func TestJsonMarshal(t *testing.T) {
	assert := assert.New(t)
	// 测试普通 map
	t.Run("Normal Map", func(t *testing.T) {
		jsonStr, err := gtkjson.JsonMarshal(map[string]any{
			"a": 1,
			"b": 1.2,
			"c": map[string]any{
				"d": 3,
			},
		})
		assert.NoError(err)
		assert.Equal("{\"a\":1,\"b\":1.2,\"c\":{\"d\":3}}", jsonStr)
	})
	// 测试 HTML 字符不转义
	t.Run("HTML Characters No Escape", func(t *testing.T) {
		jsonStr, err := gtkjson.JsonMarshal(map[string]string{
			"html": "<script>alert('test')</script>",
			"url":  "https://example.com?a=1&b=2",
		})
		assert.NoError(err)
		assert.Contains(jsonStr, "<script>")
		assert.Contains(jsonStr, "&")
		assert.NotContains(jsonStr, "\\u003c")
		assert.NotContains(jsonStr, "\\u003e")
		assert.NotContains(jsonStr, "\\u0026")
	})
	// 测试 nil 值
	t.Run("Nil Value", func(t *testing.T) {
		jsonStr, err := gtkjson.JsonMarshal(nil)
		assert.NoError(err)
		assert.Equal("", jsonStr)
	})
	// 测试结构体
	t.Run("Struct", func(t *testing.T) {
		data := TestStruct{
			Name:  "张三",
			Age:   25,
			Email: "zhangsan@example.com",
		}
		jsonStr, err := gtkjson.JsonMarshal(data)
		assert.NoError(err)
		assert.Contains(jsonStr, "张三")
		assert.Contains(jsonStr, "zhangsan@example.com")
	})
	// 测试数组
	t.Run("Array", func(t *testing.T) {
		jsonStr, err := gtkjson.JsonMarshal([]int{1, 2, 3, 4, 5})
		assert.NoError(err)
		assert.Equal("[1,2,3,4,5]", jsonStr)
	})
	// 测试字符串
	t.Run("String", func(t *testing.T) {
		jsonStr, err := gtkjson.JsonMarshal("hello world")
		assert.NoError(err)
		assert.Equal("\"hello world\"", jsonStr)
	})
}

// TestMustJsonMarshal 测试 MustJsonMarshal 函数
func TestMustJsonMarshal(t *testing.T) {
	assert := assert.New(t)
	// 测试普通 map
	t.Run("Normal Map", func(t *testing.T) {
		jsonStr := gtkjson.MustJsonMarshal(map[string]any{
			"a": 1,
			"b": 1.2,
			"c": map[string]any{
				"d": 3,
			},
		})
		assert.Equal("{\"a\":1,\"b\":1.2,\"c\":{\"d\":3}}", jsonStr)
	})
	// 测试 HTML 字符不转义
	t.Run("HTML Characters No Escape", func(t *testing.T) {
		jsonStr := gtkjson.MustJsonMarshal(map[string]string{
			"content": "<div>测试内容</div>",
		})
		assert.Contains(jsonStr, "<div>")
		assert.Contains(jsonStr, "</div>")
	})
	// 测试 nil 值
	t.Run("Nil Value", func(t *testing.T) {
		jsonStr := gtkjson.MustJsonMarshal(nil)
		assert.Equal("", jsonStr)
	})
}

// TestString 测试 String 函数
func TestString(t *testing.T) {
	assert := assert.New(t)
	// 测试普通 map
	t.Run("Normal Map", func(t *testing.T) {
		str, err := gtkjson.String(map[string]any{
			"name": "李四",
			"age":  30,
		})
		assert.NoError(err)
		assert.Contains(str, "李四")
		assert.Contains(str, "30")
	})
	// 测试 HTML 字符会转义
	t.Run("HTML Characters Escaped", func(t *testing.T) {
		str, err := gtkjson.String(map[string]string{
			"html": "<script>alert('xss')</script>",
			"url":  "https://example.com?a=1&b=2",
		})
		assert.NoError(err)
		assert.Contains(str, "\\u003c")
		assert.Contains(str, "\\u003e")
		assert.Contains(str, "\\u0026")
		assert.NotContains(str, "<script>")
	})
	// 测试 nil 值
	t.Run("Nil Value", func(t *testing.T) {
		str, err := gtkjson.String(nil)
		assert.NoError(err)
		assert.Equal("", str)
	})
	// 测试结构体
	t.Run("Struct", func(t *testing.T) {
		data := TestStruct{
			Name:  "王五",
			Age:   28,
			Email: "wangwu@example.com",
		}
		str, err := gtkjson.String(data)
		assert.NoError(err)
		assert.Contains(str, "王五")
		assert.Contains(str, "28")
		assert.Contains(str, "wangwu@example.com")
	})
}

// TestMustString 测试 MustString 函数
func TestMustString(t *testing.T) {
	assert := assert.New(t)
	// 测试普通 map
	t.Run("Normal Map", func(t *testing.T) {
		str := gtkjson.MustString(map[string]any{
			"x": 100,
			"y": 200,
		})
		assert.Contains(str, "100")
		assert.Contains(str, "200")
	})
	// 测试 nil 值
	t.Run("Nil Value", func(t *testing.T) {
		str := gtkjson.MustString(nil)
		assert.Equal("", str)
	})
	// 测试数组
	t.Run("Array", func(t *testing.T) {
		str := gtkjson.MustString([]string{"apple", "banana", "orange"})
		assert.Contains(str, "apple")
		assert.Contains(str, "banana")
		assert.Contains(str, "orange")
	})
}

// TestBytes 测试 Bytes 函数
func TestBytes(t *testing.T) {
	assert := assert.New(t)
	// 测试普通 map
	t.Run("Normal Map", func(t *testing.T) {
		b, err := gtkjson.Bytes(map[string]any{
			"key": "value",
		})
		assert.NoError(err)
		assert.NotEmpty(b)
		assert.Contains(string(b), "key")
		assert.Contains(string(b), "value")
	})
	// 测试 HTML 字符会转义
	t.Run("HTML Characters Escaped", func(t *testing.T) {
		b, err := gtkjson.Bytes(map[string]string{
			"content": "<div>内容</div>",
		})
		assert.NoError(err)
		assert.Contains(string(b), "\\u003c")
		assert.Contains(string(b), "\\u003e")
	})
	// 测试 nil 值
	t.Run("Nil Value", func(t *testing.T) {
		b, err := gtkjson.Bytes(nil)
		assert.NoError(err)
		assert.Empty(b)
	})
	// 测试结构体
	t.Run("Struct", func(t *testing.T) {
		data := TestStruct{
			Name:  "赵六",
			Age:   35,
			Email: "zhaoliu@example.com",
		}
		b, err := gtkjson.Bytes(data)
		assert.NoError(err)
		assert.NotEmpty(b)
		assert.Contains(string(b), "赵六")
	})
	// 测试布尔值
	t.Run("Boolean", func(t *testing.T) {
		b, err := gtkjson.Bytes(true)
		assert.NoError(err)
		assert.Equal("true", string(b))
	})
}

// TestMustBytes 测试 MustBytes 函数
func TestMustBytes(t *testing.T) {
	assert := assert.New(t)
	// 测试普通 map
	t.Run("Normal Map", func(t *testing.T) {
		b := gtkjson.MustBytes(map[string]int{
			"count": 42,
		})
		assert.NotEmpty(b)
		assert.Contains(string(b), "42")
	})
	// 测试 nil 值
	t.Run("Nil Value", func(t *testing.T) {
		b := gtkjson.MustBytes(nil)
		assert.Empty(b)
	})
	// 测试数字
	t.Run("Number", func(t *testing.T) {
		b := gtkjson.MustBytes(3.14159)
		assert.NotEmpty(b)
		assert.Contains(string(b), "3.14159")
	})
	// 测试复杂嵌套结构
	t.Run("Nested Structure", func(t *testing.T) {
		data := map[string]any{
			"users": []TestStruct{
				{Name: "用户1", Age: 20, Email: "user1@test.com"},
				{Name: "用户2", Age: 25, Email: "user2@test.com"},
			},
			"total": 2,
		}
		b := gtkjson.MustBytes(data)
		assert.NotEmpty(b)
		assert.Contains(string(b), "用户1")
		assert.Contains(string(b), "用户2")
		assert.Contains(string(b), "total")
	})
}
