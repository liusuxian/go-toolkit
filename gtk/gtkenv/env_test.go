/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-05-06 13:50:03
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-22 22:53:51
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkenv_test

import (
	"github.com/liusuxian/go-toolkit/gtk/gtkenv"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAll(t *testing.T) {
	assert := assert.New(t)
	assert.NotEqual([]string{}, gtkenv.All())
}

func TestMap(t *testing.T) {
	assert := assert.New(t)
	assert.NotEqual(map[string]string{}, gtkenv.Map())
}

func TestGet(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("", gtkenv.Get("a123"))
	assert.Equal("321", gtkenv.Get("a123", "321"))
}

func TestSet(t *testing.T) {
	assert := assert.New(t)
	err := gtkenv.Set("a123", "321")
	if assert.NoError(err) {
		assert.Equal("321", gtkenv.Get("a123"))
		gtkenv.Remove("a123")
	}
}

func TestSetMap(t *testing.T) {
	assert := assert.New(t)
	err := gtkenv.SetMap(map[string]string{
		"a123": "321",
		"b123": "321",
		"c123": "321",
	})
	if assert.NoError(err) {
		assert.Equal("321", gtkenv.Get("a123"))
		assert.Equal("321", gtkenv.Get("b123"))
		assert.Equal("321", gtkenv.Get("c123"))
		gtkenv.Remove("a123", "b123", "c123")
	}
}

func TestContains(t *testing.T) {
	assert := assert.New(t)
	assert.False(gtkenv.Contains("a123"))
	assert.False(gtkenv.Contains("b123"))
	assert.False(gtkenv.Contains("c123"))
	assert.True(gtkenv.Contains("GOROOT"))
}

func TestRemove(t *testing.T) {
	assert := assert.New(t)
	err := gtkenv.Remove("a123", "b123", "c123")
	if assert.NoError(err) {
		assert.False(gtkenv.Contains("a123"))
		assert.False(gtkenv.Contains("b123"))
		assert.False(gtkenv.Contains("c123"))
	}
}

func TestMapFromEnv(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(map[string]string{"a123": "321", "b123": "321"}, gtkenv.MapFromEnv([]string{"a123=321", "b123=321"}))
}

func TestMapToEnv(t *testing.T) {
	assert := assert.New(t)
	assert.ElementsMatch([]string{"a123=321", "b123=321"}, gtkenv.MapToEnv(map[string]string{"a123": "321", "b123": "321"}))
}

func TestFilter(t *testing.T) {
	assert := assert.New(t)
	assert.ElementsMatch([]string{"a123=321", "b123=321"}, gtkenv.Filter([]string{"a123=321", "b123=321", "a123=321"}))
}
