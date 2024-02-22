/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-06 18:30:38
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-22 18:33:07
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkurl_test

import (
	"github.com/liusuxian/go-toolkit/gtkurl"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsUrlEncoded(t *testing.T) {
	var (
		assert = assert.New(t)
		err    error
		ok     bool
	)
	ok, err = gtkurl.IsUrlEncoded("@4F8CjvWTBMxkaCH2Y41sFvL91WSGOvuAMpFxqA2gKlcabfH13nLnfmIujgn96zUFQ6LAdFdzlZ+by3kUJg3TMA==")
	assert.NoError(err)
	assert.False(ok)
	ok, err = gtkurl.IsUrlEncoded("%404F8CjvWTBMxkaCH2Y41sFvL91WSGOvuAMpFxqA2gKlcabfH13nLnfmIujgn96zUFQ6LAdFdzlZ%2Bby3kUJg3TMA%3D%3D")
	assert.NoError(err)
	assert.True(ok)
	ok, err = gtkurl.IsUrlEncoded("")
	assert.NoError(err)
	assert.False(ok)
	ok, err = gtkurl.IsUrlEncoded("undefined")
	assert.NoError(err)
	assert.True(ok)
	ok, err = gtkurl.IsUrlEncoded("%ZZ")
	assert.Error(err)
	assert.False(ok)
	ok, err = gtkurl.IsUrlEncoded("Hello%")
	assert.Error(err)
	assert.False(ok)
	ok, err = gtkurl.IsUrlEncoded("Hello%2")
	assert.Error(err)
	assert.False(ok)
}

func TestQueryDecode(t *testing.T) {
	var (
		assert = assert.New(t)
		err    error
		str    string
	)
	str, err = gtkurl.QueryDecode("@4F8CjvWTBMxkaCH2Y41sFvL91WSGOvuAMpFxqA2gKlcabfH13nLnfmIujgn96zUFQ6LAdFdzlZ+by3kUJg3TMA==")
	assert.NoError(err)
	assert.Equal("@4F8CjvWTBMxkaCH2Y41sFvL91WSGOvuAMpFxqA2gKlcabfH13nLnfmIujgn96zUFQ6LAdFdzlZ+by3kUJg3TMA==", str)
	str, err = gtkurl.QueryDecode("%404F8CjvWTBMxkaCH2Y41sFvL91WSGOvuAMpFxqA2gKlcabfH13nLnfmIujgn96zUFQ6LAdFdzlZ%2Bby3kUJg3TMA%3D%3D")
	assert.NoError(err)
	assert.Equal("@4F8CjvWTBMxkaCH2Y41sFvL91WSGOvuAMpFxqA2gKlcabfH13nLnfmIujgn96zUFQ6LAdFdzlZ+by3kUJg3TMA==", str)
	str, err = gtkurl.QueryDecode("")
	assert.NoError(err)
	assert.Equal("", str)
	str, err = gtkurl.QueryDecode("undefined")
	assert.NoError(err)
	assert.Equal("undefined", str)
	str, err = gtkurl.QueryDecode("%ZZ")
	assert.Error(err)
	assert.Equal("", str)
	str, err = gtkurl.QueryDecode("Hello%")
	assert.Error(err)
	assert.Equal("", str)
	str, err = gtkurl.QueryDecode("Hello%2")
	assert.Error(err)
	assert.Equal("", str)
}

func TestIsUrl(t *testing.T) {
	var (
		assert = assert.New(t)
		ok     bool
	)
	ok = gtkurl.IsUrl("http://example.com")
	assert.True(ok)
	ok = gtkurl.IsUrl("https://example.com")
	assert.True(ok)
	ok = gtkurl.IsUrl("ftp://example.com")
	assert.False(ok)
	ok = gtkurl.IsUrl("example.com")
	assert.False(ok)
	ok = gtkurl.IsUrl("hello world")
	assert.False(ok)
}
