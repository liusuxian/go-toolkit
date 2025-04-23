/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-07-15 17:56:08
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-23 19:46:35
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import (
	"fmt"
	"github.com/liusuxian/go-toolkit/gtkarr"
	"github.com/liusuxian/go-toolkit/gtkfile"
	"github.com/pkg/errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
)

// UploadFileInfo 上传的文件信息
type UploadFileInfo struct {
	err      error  // 上传失败时返回的错误信息
	FileName string `json:"fileName" dc:"文件名"`  // 文件名
	FileSize int64  `json:"fileSize" dc:"文件大小"` // 文件大小
	FilePath string `json:"filePath" dc:"文件路径"` // 文件路径
	FileType string `json:"fileType" dc:"文件类型"` // 文件类型
}

// GetErr 获取上传失败时返回的错误信息
func (ufi *UploadFileInfo) GetErr() (err error) {
	return ufi.err
}

// Upload 上传
func (ufc *UploadFileConfig) Upload(r *http.Request, dirPath string) (fileInfo *UploadFileInfo) {
	if r.Method != "POST" {
		fileInfo = &UploadFileInfo{err: errors.New("Unsupported Method")}
		return
	}
	var (
		file       multipart.File
		fileHeader *multipart.FileHeader
		err        error
	)
	if file, fileHeader, err = r.FormFile("file"); err != nil {
		fileInfo = &UploadFileInfo{err: err}
		return
	}
	defer file.Close()

	// 判断上传文件类型是否合法
	if !ufc.checkFileType(fileHeader.Filename) {
		fileInfo = &UploadFileInfo{
			err:      errors.New("Unsupported File Type"),
			FileName: fileHeader.Filename,
			FileSize: fileHeader.Size,
			FileType: fileHeader.Header.Get("Content-type"),
		}
		return
	}
	// 检查上传文件大小是否合法
	if !ufc.checkSize(fileHeader.Size) {
		fileInfo = &UploadFileInfo{
			err:      errors.New("Unsupported File Size"),
			FileName: fileHeader.Filename,
			FileSize: fileHeader.Size,
			FileType: fileHeader.Header.Get("Content-type"),
		}
		return
	}
	// 执行上传
	var filePath string
	if filePath, err = ufc.doUpload(file, fileHeader, dirPath); err != nil {
		fileInfo = &UploadFileInfo{err: err}
		return
	}
	// 返回
	fileInfo = &UploadFileInfo{
		FileName: fileHeader.Filename,
		FileSize: fileHeader.Size,
		FilePath: filePath,
		FileType: fileHeader.Header.Get("Content-type"),
	}
	return
}

// BatchUpload 批量上传
func (ufc *UploadFileConfig) BatchUpload(r *http.Request, dirPath string) (fileInfos []*UploadFileInfo) {
	if r.Method != "POST" {
		fileInfos = []*UploadFileInfo{{err: errors.New("Unsupported Method")}}
		return
	}
	var err error
	if err = r.ParseMultipartForm(defaultMaxMemory); err != nil {
		fileInfos = []*UploadFileInfo{{err: err}}
		return
	}
	// 检查文件数量
	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		fileInfos = []*UploadFileInfo{{err: http.ErrMissingFile}}
		return
	}
	if len(files) > ufc.MaxCount {
		fileInfos = []*UploadFileInfo{{err: errors.Errorf("Only Allowed To Upload A Maximum Of %v Files At A Time", ufc.MaxCount)}}
		return
	}
	// 并发处理所有文件
	fileInfos = make([]*UploadFileInfo, len(files))
	var wg sync.WaitGroup
	for k, v := range files {
		wg.Add(1)
		go func(idx int, fileHeader *multipart.FileHeader) {
			defer wg.Done()
			// 打开文件
			var (
				file multipart.File
				err  error
			)
			if file, err = fileHeader.Open(); err != nil {
				fileInfos[idx] = &UploadFileInfo{err: err}
				return
			}
			defer file.Close()

			// 判断上传文件类型是否合法
			if !ufc.checkFileType(fileHeader.Filename) {
				fileInfos[idx] = &UploadFileInfo{
					err:      errors.New("Unsupported File Type"),
					FileName: fileHeader.Filename,
					FileSize: fileHeader.Size,
					FileType: fileHeader.Header.Get("Content-type"),
				}
				return
			}
			// 检查上传文件大小是否合法
			if !ufc.checkSize(fileHeader.Size) {
				fileInfos[idx] = &UploadFileInfo{
					err:      errors.New("Unsupported File Size"),
					FileName: fileHeader.Filename,
					FileSize: fileHeader.Size,
					FileType: fileHeader.Header.Get("Content-type"),
				}
				return
			}
			// 执行上传
			var filePath string
			if filePath, err = ufc.doUpload(file, fileHeader, dirPath); err != nil {
				fileInfos[idx] = &UploadFileInfo{err: err}
				return
			}
			// 返回
			fileInfos[idx] = &UploadFileInfo{
				FileName: fileHeader.Filename,
				FileSize: fileHeader.Size,
				FilePath: filePath,
				FileType: fileHeader.Header.Get("Content-type"),
			}
		}(k, v)
	}
	wg.Wait()
	return
}

// checkFileType 判断上传文件类型是否合法
func (ufc *UploadFileConfig) checkFileType(fileName string) (ok bool) {
	return gtkarr.ContainsStr(ufc.AllowTypeList, gtkfile.ExtName(fileName))
}

// checkSize 检查上传文件大小是否合法
func (ufc *UploadFileConfig) checkSize(fileSize int64) (ok bool) {
	return int64(ufc.MaxSize*1024*1024) >= fileSize
}

// doUpload 执行上传
func (ufc *UploadFileConfig) doUpload(file multipart.File, fileHeader *multipart.FileHeader, dirPath string) (filePath string, err error) {
	// 创建目录
	if err = gtkfile.MakeDirAll(dirPath); err != nil {
		err = errors.Wrapf(err, "Failed To Create Directory: %s", dirPath)
		return
	}
	// 构建文件完整路径
	filePath = fmt.Sprintf("%s/%s", dirPath, gtkfile.GenRandomFileName(fileHeader.Filename))
	// 创建文件
	var outFile *os.File
	if outFile, err = os.Create(filePath); err != nil {
		err = errors.Wrapf(err, "Failed To Create File: %s", filePath)
		return
	}
	defer outFile.Close()
	// 写入文件
	if _, err = io.Copy(outFile, file); err != nil {
		err = errors.Wrapf(err, "Failed To Write File: %s", filePath)
		return
	}
	return
}
