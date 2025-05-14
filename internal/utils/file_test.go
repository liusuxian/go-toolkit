/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-13 12:59:36
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-13 16:11:23
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package utils_test

import (
	"github.com/liusuxian/go-toolkit/internal/utils"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"testing"
)

func TestPathExists(t *testing.T) {
	assert := assert.New(t)
	assert.True(utils.PathExists("."))
	assert.False(utils.PathExists("config/config.yaml"))
}

func TestExtName(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("yaml", utils.ExtName("config/config.yaml"))
	assert.NotEqual("test", utils.ExtName("config/test"))
	assert.Equal("", utils.ExtName("config/test"))
}

func TestGetContents(t *testing.T) {
	assert := assert.New(t)
	assert.NotEmpty(utils.GetContents("file_test.go"))
}

func TestGetBytes(t *testing.T) {
	assert := assert.New(t)
	assert.NotEmpty(utils.GetBytes("file_test.go"))
}

func TestName(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("file", utils.Name("/var/www/file.js"))
	assert.Equal("file", utils.Name("file.js"))
}

func TestGenRandomFileName(t *testing.T) {
	t.Log(utils.GenRandomFilename("example.png"))
	t.Log(utils.GenRandomFilename("example.png"))
	t.Log(utils.GenRandomFilename("example.png"))
	t.Log(utils.GenRandomFilename("example.png"))
	t.Log(utils.GenRandomFilename("example.png"))
}

func TestGetFileStat(t *testing.T) {
	var (
		assert   = assert.New(t)
		fileInfo fs.FileInfo
		err      error
	)
	fileInfo, err = utils.GetFileStat("example.png")
	assert.Error(err)
	assert.Nil(fileInfo)

	fileInfo, err = utils.GetFileStat("file.go")
	assert.NoError(err)
	assert.NotNil(fileInfo)
	t.Log("fileName:", fileInfo.Name())
	t.Log("fileSize:", fileInfo.Size())
}
