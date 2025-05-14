/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-04-10 14:24:41
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-13 14:12:56
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkstr_test

import (
	"github.com/liusuxian/go-toolkit/gtkstr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrimAll(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("alog", gtkstr.TrimAll("a.log", "."))
	assert.Equal("alog", gtkstr.TrimAll(" a . log ", "."))
	assert.Equal("ablog", gtkstr.TrimAll("a.b.log", "."))
	assert.Equal("ablog", gtkstr.TrimAll(" a . b . log", "."))
	assert.Equal("", gtkstr.TrimAll("", "."))
	assert.Equal("", gtkstr.TrimAll(" ", "."))
	assert.Equal("", gtkstr.TrimAll("   ", "."))
}

func TestSplit(t *testing.T) {
	assert := assert.New(t)
	assert.ElementsMatch([]string{"a", "log"}, gtkstr.Split("a.log", "."))
	assert.ElementsMatch([]string{"a", "log"}, gtkstr.Split(" a . log ", "."))
	assert.ElementsMatch([]string{"a", "b", "log"}, gtkstr.Split("a.b.log", "."))
	assert.ElementsMatch([]string{"a", "b", "log"}, gtkstr.Split(" a . b . log", "."))
	assert.ElementsMatch([]string{}, gtkstr.Split("", "."))
	assert.ElementsMatch([]string{}, gtkstr.Split(" ", "."))
	assert.ElementsMatch([]string{}, gtkstr.Split("   ", "."))
}
