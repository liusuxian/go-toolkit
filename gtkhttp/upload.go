/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-07-15 17:56:08
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-25 12:20:58
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import (
	"errors"
	"fmt"
	"github.com/liusuxian/go-toolkit/internal/utils"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"sync"
)

var (
	ErrMissingFile         = errors.New("no such file")
	ErrUnsupportedMethod   = errors.New("unsupported method")
	ErrUnsupportedFileType = errors.New("unsupported file type")
	ErrUnsupportedFileSize = errors.New("unsupported file size")
)

// UploadFileInfo 上传的文件信息
type UploadFileInfo struct {
	FileName string `json:"file_name"` // 文件名
	FileSize int64  `json:"file_size"` // 文件大小
	FilePath string `json:"file_path"` // 文件路径
	FileType string `json:"file_type"` // 文件类型
	err      error  `json:"-"`         // 上传失败时返回的错误信息
}

// GetErr 获取上传失败时返回的错误信息
func (ufi *UploadFileInfo) GetErr() (err error) {
	return ufi.err
}

// Upload 上传
func (s *UploadFileService) Upload(r *http.Request, dirPath string) (fileInfo *UploadFileInfo) {
	if r.Method != "POST" {
		fileInfo = &UploadFileInfo{err: ErrUnsupportedMethod}
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
	if !s.checkFileType(fileHeader.Filename) {
		fileInfo = &UploadFileInfo{
			err:      ErrUnsupportedFileType,
			FileName: fileHeader.Filename,
			FileSize: fileHeader.Size,
			FileType: fileHeader.Header.Get("Content-type"),
		}
		return
	}
	// 检查上传文件大小是否合法
	if !s.checkSize(fileHeader.Size) {
		fileInfo = &UploadFileInfo{
			err:      ErrUnsupportedFileSize,
			FileName: fileHeader.Filename,
			FileSize: fileHeader.Size,
			FileType: fileHeader.Header.Get("Content-type"),
		}
		return
	}
	// 执行上传
	var filePath string
	if filePath, err = s.doUpload(file, fileHeader, dirPath); err != nil {
		fileInfo = &UploadFileInfo{
			err:      err,
			FileName: fileHeader.Filename,
			FileSize: fileHeader.Size,
			FileType: fileHeader.Header.Get("Content-type"),
		}
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
func (s *UploadFileService) BatchUpload(r *http.Request, dirPath string) (fileInfos []*UploadFileInfo) {
	if r.Method != "POST" {
		fileInfos = []*UploadFileInfo{{err: ErrUnsupportedMethod}}
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
		fileInfos = []*UploadFileInfo{{err: ErrMissingFile}}
		return
	}
	if len(files) > s.config.MaxCount {
		fileInfos = []*UploadFileInfo{{err: fmt.Errorf("only allowed to upload a maximum of %v files at a time", s.config.MaxCount)}}
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
			if !s.checkFileType(fileHeader.Filename) {
				fileInfos[idx] = &UploadFileInfo{
					err:      ErrUnsupportedFileType,
					FileName: fileHeader.Filename,
					FileSize: fileHeader.Size,
					FileType: fileHeader.Header.Get("Content-type"),
				}
				return
			}
			// 检查上传文件大小是否合法
			if !s.checkSize(fileHeader.Size) {
				fileInfos[idx] = &UploadFileInfo{
					err:      ErrUnsupportedFileSize,
					FileName: fileHeader.Filename,
					FileSize: fileHeader.Size,
					FileType: fileHeader.Header.Get("Content-type"),
				}
				return
			}
			// 执行上传
			var filePath string
			if filePath, err = s.doUpload(file, fileHeader, dirPath); err != nil {
				fileInfos[idx] = &UploadFileInfo{
					err:      err,
					FileName: fileHeader.Filename,
					FileSize: fileHeader.Size,
					FileType: fileHeader.Header.Get("Content-type"),
				}
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
func (s *UploadFileService) checkFileType(fileName string) (ok bool) {
	return slices.Contains(s.config.AllowTypeList, utils.ExtName(fileName))
}

// checkSize 检查上传文件大小是否合法
func (s *UploadFileService) checkSize(fileSize int64) (ok bool) {
	return int64(s.config.MaxSize*1024*1024) >= fileSize
}

// doUpload 执行上传
func (s *UploadFileService) doUpload(file multipart.File, fileHeader *multipart.FileHeader, dirPath string) (filePath string, err error) {
	// 创建目录
	if err = utils.MakeDirAll(dirPath); err != nil {
		err = fmt.Errorf("failed to create directory: %s, error: %w", dirPath, err)
		return
	}
	// 构建文件完整路径
	filePath = filepath.Join(dirPath, s.uploadFileNameFn(fileHeader.Filename))
	// 创建文件
	var outFile *os.File
	if outFile, err = os.Create(filePath); err != nil {
		err = fmt.Errorf("failed to create file: %s, error: %w", filePath, err)
		return
	}
	defer outFile.Close()
	// 写入文件
	if _, err = io.Copy(outFile, file); err != nil {
		err = fmt.Errorf("failed to write file: %s, error: %w", filePath, err)
		return
	}
	return
}
