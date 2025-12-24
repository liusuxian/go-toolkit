/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-13 12:57:31
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-24 20:11:43
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package utils

import (
	"fmt"
	"github.com/google/uuid"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// PathExists 判断文件或者目录是否存在
func PathExists(path string) (isExist bool) {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// ExtName 获取文件扩展名
func ExtName(path string) (extName string) {
	return strings.TrimLeft(filepath.Ext(path), ".")
}

// GetContents 获取文件的内容
func GetContents(path string) (str string) {
	return string(GetBytes(path))
}

// GetBytes 获取文件的内容
func GetBytes(path string) (buf []byte) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return data
}

// Name 返回路径的最后一个元素，但不包括文件扩展名
func Name(path string) (str string) {
	base := filepath.Base(path)
	if i := strings.LastIndexByte(base, '.'); i != -1 {
		return base[:i]
	}
	return base
}

// GenRandomFilename 生成随机文件名
func GenRandomFilename(filename string) (newFilename string) {
	return uuid.New().String() + filepath.Ext(filename)
}

// MakeDirAll 创建给定路径的所有目录，包括任何必要的父目录
func MakeDirAll(path string) (err error) {
	if !PathExists(path) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return fmt.Errorf("create <%s> error: %w", path, err)
		}
	}
	return
}

// GetFileStat 获取文件的状态信息
func GetFileStat(name string) (fileInfo fs.FileInfo, err error) {
	if fileInfo, err = os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			err = fmt.Errorf("file <%s> not exist", name)
			return
		}
		err = fmt.Errorf("get file <%s> stat error: %w", name, err)
		return
	}
	return
}
