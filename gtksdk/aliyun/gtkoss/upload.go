/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-22 23:33:32
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-24 18:12:52
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkoss

import (
	"bytes"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/liusuxian/go-toolkit/internal/utils"
	"io/fs"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"slices"
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
func (u *UploadFileInfo) GetErr() (err error) {
	return u.err
}

// UploadFromHttp 从 HTTP 请求中上传文件
func (s *AliyunOSS) UploadFromHttp(r *http.Request, dirPath string, opts ...Option) (fileInfo *UploadFileInfo) {
	if r.Method != "POST" {
		fileInfo = &UploadFileInfo{err: fmt.Errorf("unsupported method")}
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
			err:      fmt.Errorf("unsupported file type"),
			FileName: fileHeader.Filename,
			FileSize: fileHeader.Size,
			FileType: fileHeader.Header.Get("Content-type"),
		}
		return
	}
	// 检查上传文件大小是否合法
	if !s.checkSize(fileHeader.Size) {
		fileInfo = &UploadFileInfo{
			err:      fmt.Errorf("unsupported file size"),
			FileName: fileHeader.Filename,
			FileSize: fileHeader.Size,
			FileType: fileHeader.Header.Get("Content-type"),
		}
		return
	}
	// 执行上传
	var filePath string
	if filePath, err = s.doUploadFromHttp(file, fileHeader, dirPath, opts...); err != nil {
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

// UploadFromBytes 从字节切片中上传文件
func (s *AliyunOSS) UploadFromBytes(data []byte, fileName string, dirPath string, opts ...Option) (fileInfo *UploadFileInfo) {
	// 检查数据是否为空
	fileSize := int64(len(data))
	if fileSize == 0 {
		fileInfo = &UploadFileInfo{
			err:      fmt.Errorf("data is empty"),
			FileName: fileName,
			FileType: utils.ExtName(fileName),
		}
		return
	}
	// 判断上传文件类型是否合法
	if !s.checkFileType(fileName) {
		fileInfo = &UploadFileInfo{
			err:      fmt.Errorf("unsupported file type"),
			FileName: fileName,
			FileSize: fileSize,
			FileType: utils.ExtName(fileName),
		}
		return
	}
	// 检查上传文件大小是否合法
	if !s.checkSize(fileSize) {
		fileInfo = &UploadFileInfo{
			err:      fmt.Errorf("unsupported file size"),
			FileName: fileName,
			FileSize: fileSize,
			FileType: utils.ExtName(fileName),
		}
		return
	}
	// 执行上传
	var (
		filePath string
		err      error
	)
	if filePath, err = s.doUploadFromBytes(data, fileName, dirPath, opts...); err != nil {
		fileInfo = &UploadFileInfo{
			err:      err,
			FileName: fileName,
			FileSize: fileSize,
			FileType: utils.ExtName(fileName),
		}
		return
	}
	// 返回
	fileInfo = &UploadFileInfo{
		FileName: fileName,
		FileSize: fileSize,
		FilePath: filePath,
		FileType: utils.ExtName(fileName),
	}
	return
}

// UploadFromFile 通过文件名（包含文件路径）上传
func (s *AliyunOSS) UploadFromFile(dirPath, fileName string, opts ...Option) (fileInfo *UploadFileInfo) {
	// 获取文件的状态信息
	var (
		fileStat fs.FileInfo
		err      error
	)
	if fileStat, err = utils.GetFileStat(fileName); err != nil {
		fileInfo = &UploadFileInfo{err: err}
		return
	}
	// 判断上传文件类型是否合法
	if !s.checkFileType(fileStat.Name()) {
		fileInfo = &UploadFileInfo{
			err:      fmt.Errorf("unsupported file type"),
			FileName: fileName,
			FileSize: fileStat.Size(),
			FileType: utils.ExtName(fileStat.Name()),
		}
		return
	}
	// 检查上传文件大小是否合法
	if !s.checkSize(fileStat.Size()) {
		fileInfo = &UploadFileInfo{
			err:      fmt.Errorf("unsupported file size"),
			FileName: fileName,
			FileSize: fileStat.Size(),
			FileType: utils.ExtName(fileStat.Name()),
		}
		return
	}
	// 执行上传
	var filePath string
	if filePath, err = s.doUploadFromFile(fileName, dirPath, fileStat, opts...); err != nil {
		fileInfo = &UploadFileInfo{err: err}
		return
	}
	// 返回
	fileInfo = &UploadFileInfo{
		FileName: fileName,
		FileSize: fileStat.Size(),
		FilePath: filePath,
		FileType: utils.ExtName(fileStat.Name()),
	}
	return
}

// BatchUploadFromHttp 从 HTTP 请求中批量上传文件
func (s *AliyunOSS) BatchUploadFromHttp(r *http.Request, dirPath string, opts ...Option) (fileInfos []*UploadFileInfo) {
	if r.Method != "POST" {
		fileInfos = []*UploadFileInfo{{err: fmt.Errorf("unsupported method")}}
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
					err:      fmt.Errorf("unsupported file type"),
					FileName: fileHeader.Filename,
					FileSize: fileHeader.Size,
					FileType: fileHeader.Header.Get("Content-type"),
				}
				return
			}
			// 检查上传文件大小是否合法
			if !s.checkSize(fileHeader.Size) {
				fileInfos[idx] = &UploadFileInfo{
					err:      fmt.Errorf("unsupported file size"),
					FileName: fileHeader.Filename,
					FileSize: fileHeader.Size,
					FileType: fileHeader.Header.Get("Content-type"),
				}
				return
			}
			// 执行上传
			var filePath string
			if filePath, err = s.doUploadFromHttp(file, fileHeader, dirPath, opts...); err != nil {
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

// BatchUploadFromFile 通过文件名（包含文件路径）批量上传
func (s *AliyunOSS) BatchUploadFromFile(dirPath string, fileNameList []string, opts ...Option) (fileInfos []*UploadFileInfo) {
	// 检查文件数量
	if len(fileNameList) == 0 {
		fileInfos = []*UploadFileInfo{{err: http.ErrMissingFile}}
		return
	}
	if len(fileNameList) > s.config.MaxCount {
		fileInfos = []*UploadFileInfo{{err: fmt.Errorf("only allowed to upload a maximum of %v files at a time", s.config.MaxCount)}}
		return
	}
	// 并发处理所有文件
	fileInfos = make([]*UploadFileInfo, len(fileNameList))
	var wg sync.WaitGroup
	for k, v := range fileNameList {
		wg.Add(1)
		go func(idx int, fileName string) {
			defer wg.Done()
			// 获取文件的状态信息
			var (
				fileStat fs.FileInfo
				err      error
			)
			if fileStat, err = utils.GetFileStat(fileName); err != nil {
				fileInfos[idx] = &UploadFileInfo{err: err}
				return
			}
			// 判断上传文件类型是否合法
			if !s.checkFileType(fileStat.Name()) {
				fileInfos[idx] = &UploadFileInfo{
					err:      fmt.Errorf("unsupported file type"),
					FileName: fileName,
					FileSize: fileStat.Size(),
					FileType: utils.ExtName(fileStat.Name()),
				}
				return
			}
			// 检查上传文件大小是否合法
			if !s.checkSize(fileStat.Size()) {
				fileInfos[idx] = &UploadFileInfo{
					err:      fmt.Errorf("unsupported file size"),
					FileName: fileName,
					FileSize: fileStat.Size(),
					FileType: utils.ExtName(fileStat.Name()),
				}
				return
			}
			// 执行上传
			var filePath string
			if filePath, err = s.doUploadFromFile(fileName, dirPath, fileStat, opts...); err != nil {
				fileInfos[idx] = &UploadFileInfo{err: err}
				return
			}
			// 返回
			fileInfos[idx] = &UploadFileInfo{
				FileName: fileName,
				FileSize: fileStat.Size(),
				FilePath: filePath,
				FileType: utils.ExtName(fileStat.Name()),
			}
		}(k, v)
	}
	wg.Wait()
	return
}

// checkFileType 判断上传文件类型是否合法
func (s *AliyunOSS) checkFileType(fileName string) (ok bool) {
	return slices.Contains(s.config.AllowTypeList, utils.ExtName(fileName))
}

// checkSize 检查上传文件大小是否合法
func (s *AliyunOSS) checkSize(fileSize int64) (ok bool) {
	return int64(s.config.MaxSize*1024*1024) >= fileSize
}

// doUploadFromHttp 执行上传
func (s *AliyunOSS) doUploadFromHttp(file multipart.File, fileHeader *multipart.FileHeader, dirPath string, opts ...Option) (filePath string, err error) {
	// 构建文件完整路径
	filePath = filepath.Join(dirPath, s.uploadFileNameFn(fileHeader.Filename))
	// 获取存储空间
	var (
		client *oss.Client
		bucket *oss.Bucket
	)
	if client, bucket, err = s.getBucket(opts...); err != nil {
		return
	}
	// 关闭空闲连接
	defer s.closeIdleConnections(client)
	// 上传文件
	if err = bucket.PutObject(filePath, file, s.ossOptions...); err != nil {
		return
	}
	return
}

// doUploadFromBytes 执行上传
func (s *AliyunOSS) doUploadFromBytes(data []byte, fileName, dirPath string, opts ...Option) (filePath string, err error) {
	// 构建文件完整路径
	filePath = filepath.Join(dirPath, s.uploadFileNameFn(fileName))
	// 获取存储空间
	var (
		client *oss.Client
		bucket *oss.Bucket
	)
	if client, bucket, err = s.getBucket(opts...); err != nil {
		return
	}
	// 关闭空闲连接
	defer s.closeIdleConnections(client)
	// 上传文件
	if err = bucket.PutObject(filePath, bytes.NewReader(data), s.ossOptions...); err != nil {
		return
	}
	return
}

// doUploadFromFile 执行上传
func (s *AliyunOSS) doUploadFromFile(fileName, dirPath string, fileStat fs.FileInfo, opts ...Option) (filePath string, err error) {
	// 构建文件完整路径
	filePath = filepath.Join(dirPath, s.uploadFileNameFn(fileStat.Name()))
	// 获取存储空间
	var (
		client *oss.Client
		bucket *oss.Bucket
	)
	if client, bucket, err = s.getBucket(opts...); err != nil {
		return
	}
	// 关闭空闲连接
	defer s.closeIdleConnections(client)
	// 上传文件
	if err = bucket.PutObjectFromFile(filePath, fileName, s.ossOptions...); err != nil {
		return
	}
	return
}
