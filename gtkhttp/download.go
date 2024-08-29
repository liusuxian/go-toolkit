/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-08-29 17:05:37
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-08-29 17:45:33
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import (
	"fmt"
	"github.com/liusuxian/go-toolkit/gtkfile"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
)

// DownloadFile 从指定的 URL 下载文件到本地路径
func DownloadFile(url, dirPath string) (filePath string, err error) {
	// 发起 HTTP GET 请求
	var resp *http.Response
	if resp, err = http.Get(url); err != nil {
		err = errors.Wrapf(err, "Failed To Http Get: %s", url)
		return
	}
	defer resp.Body.Close()
	// 检查 HTTP 响应状态码
	if resp.StatusCode != http.StatusOK {
		err = errors.Errorf("Failed To Http Get: %s %s", url, resp.Status)
		return
	}
	// 创建目录
	if err = gtkfile.MakeDirAll(dirPath); err != nil {
		err = errors.Wrapf(err, "Failed To Create Directory: %s", dirPath)
		return
	}
	// 构建文件完整路径
	filePath = fmt.Sprintf("%s/%s", dirPath, gtkfile.GenRandomFileName(url))
	// 创建文件
	var outFile *os.File
	if outFile, err = os.Create(filePath); err != nil {
		err = errors.Wrapf(err, "Failed To Create File: %s", filePath)
		return
	}
	defer outFile.Close()
	// 写入文件
	if _, err = io.Copy(outFile, resp.Body); err != nil {
		err = errors.Wrapf(err, "Failed To Write File: %s", filePath)
		return
	}
	return
}
