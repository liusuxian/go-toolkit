/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-22 23:33:32
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-08-29 16:58:24
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package oss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/liusuxian/go-toolkit/gtkarr"
	"github.com/liusuxian/go-toolkit/gtkfile"
	"github.com/pkg/errors"
	"io/fs"
	"mime/multipart"
	"net/http"
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

// Upload 上传
func (s *AliyunOSS) Upload(r *http.Request, dirPath string) (fileInfo *UploadFileInfo) {
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
	if !s.checkFileType(fileHeader.Filename) {
		fileInfo = &UploadFileInfo{
			err:      errors.New("Unsupported File Type"),
			FileName: fileHeader.Filename,
			FileSize: fileHeader.Size,
			FileType: fileHeader.Header.Get("Content-type"),
		}
		return
	}
	// 检查上传文件大小是否合法
	if !s.checkSize(fileHeader.Size) {
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
	if filePath, err = s.doUpload(file, fileHeader, dirPath); err != nil {
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

// UploadFromFile 通过文件名（包含文件路径）上传
func (s *AliyunOSS) UploadFromFile(dirPath, fileName string) (fileInfo *UploadFileInfo) {
	// 获取文件的状态信息
	var (
		fileStat fs.FileInfo
		err      error
	)
	if fileStat, err = gtkfile.GetFileStat(fileName); err != nil {
		fileInfo = &UploadFileInfo{err: err}
		return
	}
	// 判断上传文件类型是否合法
	if !s.checkFileType(fileStat.Name()) {
		fileInfo = &UploadFileInfo{
			err:      errors.New("Unsupported File Type"),
			FileName: fileName,
			FileSize: fileStat.Size(),
			FileType: gtkfile.ExtName(fileStat.Name()),
		}
		return
	}
	// 检查上传文件大小是否合法
	if !s.checkSize(fileStat.Size()) {
		fileInfo = &UploadFileInfo{
			err:      errors.New("Unsupported File Size"),
			FileName: fileName,
			FileSize: fileStat.Size(),
			FileType: gtkfile.ExtName(fileStat.Name()),
		}
		return
	}
	// 执行上传
	var filePath string
	if filePath, err = s.doUploadFromFile(fileName, dirPath, fileStat); err != nil {
		fileInfo = &UploadFileInfo{err: err}
		return
	}
	// 返回
	fileInfo = &UploadFileInfo{
		FileName: fileName,
		FileSize: fileStat.Size(),
		FilePath: filePath,
		FileType: gtkfile.ExtName(fileStat.Name()),
	}
	return
}

// BatchUpload 批量上传
func (s *AliyunOSS) BatchUpload(r *http.Request, dirPath string) (fileInfos []*UploadFileInfo) {
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
	if len(files) > s.MaxCount {
		fileInfos = []*UploadFileInfo{{err: errors.Errorf("Only Allowed To Upload A Maximum Of %v Files At A Time", s.MaxCount)}}
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
					err:      errors.New("Unsupported File Type"),
					FileName: fileHeader.Filename,
					FileSize: fileHeader.Size,
					FileType: fileHeader.Header.Get("Content-type"),
				}
				return
			}
			// 检查上传文件大小是否合法
			if !s.checkSize(fileHeader.Size) {
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
			if filePath, err = s.doUpload(file, fileHeader, dirPath); err != nil {
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
func (s *AliyunOSS) BatchUploadFromFile(dirPath string, fileNameList ...string) (fileInfos []*UploadFileInfo) {
	// 检查文件数量
	if len(fileNameList) == 0 {
		fileInfos = []*UploadFileInfo{{err: http.ErrMissingFile}}
		return
	}
	if len(fileNameList) > s.MaxCount {
		fileInfos = []*UploadFileInfo{{err: errors.Errorf("Only Allowed To Upload A Maximum Of %v Files At A Time", s.MaxCount)}}
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
			if fileStat, err = gtkfile.GetFileStat(fileName); err != nil {
				fileInfos[idx] = &UploadFileInfo{err: err}
				return
			}
			// 判断上传文件类型是否合法
			if !s.checkFileType(fileStat.Name()) {
				fileInfos[idx] = &UploadFileInfo{
					err:      errors.New("Unsupported File Type"),
					FileName: fileName,
					FileSize: fileStat.Size(),
					FileType: gtkfile.ExtName(fileStat.Name()),
				}
				return
			}
			// 检查上传文件大小是否合法
			if !s.checkSize(fileStat.Size()) {
				fileInfos[idx] = &UploadFileInfo{
					err:      errors.New("Unsupported File Size"),
					FileName: fileName,
					FileSize: fileStat.Size(),
					FileType: gtkfile.ExtName(fileStat.Name()),
				}
				return
			}
			// 执行上传
			var filePath string
			if filePath, err = s.doUploadFromFile(fileName, dirPath, fileStat); err != nil {
				fileInfos[idx] = &UploadFileInfo{err: err}
				return
			}
			// 返回
			fileInfos[idx] = &UploadFileInfo{
				FileName: fileName,
				FileSize: fileStat.Size(),
				FilePath: filePath,
				FileType: gtkfile.ExtName(fileStat.Name()),
			}
		}(k, v)
	}
	wg.Wait()
	return
}

// checkFileType 判断上传文件类型是否合法
func (s *AliyunOSS) checkFileType(fileName string) (ok bool) {
	return gtkarr.ContainsStr(s.AllowTypeList, gtkfile.ExtName(fileName))
}

// checkSize 检查上传文件大小是否合法
func (s *AliyunOSS) checkSize(fileSize int64) (ok bool) {
	return int64(s.MaxSize*1024*1024) >= fileSize
}

// doUpload 执行上传
func (s *AliyunOSS) doUpload(file multipart.File, fileHeader *multipart.FileHeader, dirPath string) (filePath string, err error) {
	// 连接OSS
	var client *oss.Client
	if client, err = oss.New(s.EndpointAccelerate, s.AccessKeyID, s.AccessKeySecret); err != nil {
		return
	}
	// 获取存储空间
	var bucket *oss.Bucket
	if bucket, err = client.Bucket(s.Bucket); err != nil {
		return
	}
	// 构建文件完整路径
	filePath = fmt.Sprintf("%s/%s", dirPath, gtkfile.GenRandomFileName(fileHeader.Filename))
	// 上传文件
	if err = bucket.PutObject(filePath, file); err != nil {
		return
	}
	return
}

// doUploadFromFile 执行上传
func (s *AliyunOSS) doUploadFromFile(fileName, dirPath string, fileStat fs.FileInfo) (filePath string, err error) {
	// 连接OSS
	var client *oss.Client
	if client, err = oss.New(s.EndpointAccelerate, s.AccessKeyID, s.AccessKeySecret); err != nil {
		return
	}
	// 获取存储空间
	var bucket *oss.Bucket
	if bucket, err = client.Bucket(s.Bucket); err != nil {
		return
	}
	// 构建文件完整路径
	filePath = fmt.Sprintf("%s/%s", dirPath, gtkfile.GenRandomFileName(fileStat.Name()))
	// 上传文件
	if err = bucket.PutObjectFromFile(filePath, fileName); err != nil {
		return
	}
	return
}
