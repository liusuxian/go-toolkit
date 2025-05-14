/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-02-19 21:04:58
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-13 16:11:42
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkfile

import (
	"github.com/liusuxian/go-toolkit/internal/utils"
	"io/fs"
)

// PathExists 判断文件或者目录是否存在
func PathExists(path string) (isExist bool) {
	return utils.PathExists(path)
}

// ExtName 获取文件扩展名
func ExtName(path string) (extName string) {
	return utils.ExtName(path)
}

// GetContents 获取文件的内容
func GetContents(path string) (str string) {
	return utils.GetContents(path)
}

// GetBytes 获取文件的内容
func GetBytes(path string) (buf []byte) {
	return utils.GetBytes(path)
}

// Name 返回路径的最后一个元素，但不包括文件扩展名
func Name(path string) (str string) {
	return utils.Name(path)
}

// GenRandomFilename 生成随机文件名
func GenRandomFilename(filename string) (newFilename string) {
	return utils.GenRandomFilename(filename)
}

// MakeDirAll 创建给定路径的所有目录，包括任何必要的父目录
func MakeDirAll(path string) (err error) {
	return utils.MakeDirAll(path)
}

// GetFileStat 获取文件的状态信息
func GetFileStat(name string) (fileInfo fs.FileInfo, err error) {
	return utils.GetFileStat(name)
}
