/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-04-10 00:16:21
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-22 23:00:40
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkfile_test

import (
	"github.com/liusuxian/go-toolkit/gtkfile"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPathExists(t *testing.T) {
	assert := assert.New(t)
	assert.True(gtkfile.PathExists("."))
	assert.False(gtkfile.PathExists("config/config.yaml"))
}

func TestExtName(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("yaml", gtkfile.ExtName("config/config.yaml"))
	assert.NotEqual("test", gtkfile.ExtName("config/test"))
	assert.Equal("", gtkfile.ExtName("config/test"))
}

func TestGetContents(t *testing.T) {
	assert := assert.New(t)
	assert.NotEmpty(gtkfile.GetContents("file_test.go"))
}

func TestGetBytes(t *testing.T) {
	assert := assert.New(t)
	assert.NotEmpty(gtkfile.GetBytes("file_test.go"))
}

func TestName(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("file", gtkfile.Name("/var/www/file.js"))
	assert.Equal("file", gtkfile.Name("file.js"))
}
