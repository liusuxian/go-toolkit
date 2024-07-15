/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-02-19 21:04:58
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-07-15 20:22:45
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkfile

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"hash"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

// GenRandomFileName 生成随机文件名，默认使用 MD5 算法生成文件名，可选 SHA256 算法
func GenRandomFileName(originFileName string, isSha256 ...bool) (fileName string) {
	baseName := strconv.FormatInt(time.Now().UnixNano(), 36)
	randomPart, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	input := fmt.Sprintf("%s-%s-%s", originFileName, baseName, randomPart.String())

	var hasher hash.Hash
	if len(isSha256) > 0 && isSha256[0] {
		hasher = sha256.New()
	} else {
		hasher = md5.New()
	}
	hasher.Write([]byte(input))
	hashed := hex.EncodeToString(hasher.Sum(nil))
	ext := filepath.Ext(originFileName)
	fileName = fmt.Sprintf("%s%s", hashed, ext)
	return
}

// MakeDirAll 创建给定路径的所有目录，包括任何必要的父目录
func MakeDirAll(path string) (err error) {
	if !PathExists(path) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return errors.Errorf("create <%s> error: %s", path, err)
		}
	}
	return
}
