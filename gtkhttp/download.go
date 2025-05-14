/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-23 19:46:05
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-13 16:15:20
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import (
	"fmt"
	"github.com/liusuxian/go-toolkit/internal/utils"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

// DownloadFile 从指定的 URL 下载文件到本地路径
func DownloadFile(url, dirPath string) (filePath string, err error) {
	// 发起 HTTP GET 请求
	var resp *http.Response
	if resp, err = http.Get(url); err != nil {
		err = fmt.Errorf("failed to http get: %s, error: %w", url, err)
		return
	}
	defer resp.Body.Close()
	// 检查 HTTP 响应状态码
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to http get: %s, status: %s", url, resp.Status)
		return
	}
	// 创建目录
	if err = utils.MakeDirAll(dirPath); err != nil {
		err = fmt.Errorf("failed to create directory: %s, error: %w", dirPath, err)
		return
	}
	// 构建文件完整路径
	fileName := ExtractFileNameFromURL(url)
	filePath = fmt.Sprintf("%s/%s", dirPath, utils.GenRandomFilename(fileName))
	// 创建文件
	var outFile *os.File
	if outFile, err = os.Create(filePath); err != nil {
		err = fmt.Errorf("failed to create file: %s, error: %w", filePath, err)
		return
	}
	defer outFile.Close()
	// 写入文件
	if _, err = io.Copy(outFile, resp.Body); err != nil {
		err = fmt.Errorf("failed to write file: %s, error: %w", filePath, err)
		return
	}
	return
}

// ExtractFileNameFromURL 从 URL 中提取文件名称
func ExtractFileNameFromURL(rawURL string) (fileName string) {
	// 解析 URL
	var (
		parsedURL *url.URL
		err       error
	)
	if parsedURL, err = url.Parse(rawURL); err != nil {
		// 如果 URL 无法解析，返回空字符串
		return
	}
	// 获取路径的最后一部分作为图片名称
	fileName = path.Base(parsedURL.Path)
	// 如果结果中包含特殊字符，进行解码
	var decodedName string
	if decodedName, err = url.QueryUnescape(fileName); err == nil {
		fileName = decodedName
	}
	// 移除可能的查询参数或锚点
	if idx := strings.IndexAny(fileName, "?#"); idx != -1 {
		fileName = fileName[:idx]
	}
	return
}
