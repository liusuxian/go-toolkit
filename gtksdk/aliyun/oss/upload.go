/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-22 23:33:32
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-29 18:20:52
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package oss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/liusuxian/go-toolkit/gtkarr"
	"github.com/liusuxian/go-toolkit/gtkcache"
	"github.com/liusuxian/go-toolkit/gtkfile"
	"github.com/pkg/errors"
	"mime/multipart"
	"net/http"
	"sync"
)

// AliyunOSS 阿里云`OSS`信息
type AliyunOSS struct {
	Bucket             string          `json:"bucket" dc:"bucket名称"`                 // bucket名称
	EndpointAccelerate string          `json:"endpointAccelerate" dc:"传输加速节点"`       // 传输加速节点
	EndpointInternal   string          `json:"endpointInternal" dc:"内网访问节点"`         // 内网访问节点
	EndpointAccess     string          `json:"endpointAccess" dc:"外网访问节点"`           // 外网访问节点
	AccessKeyID        string          `json:"accessKeyID" dc:"accessKeyID"`         // accessKeyID
	AccessKeySecret    string          `json:"accessKeySecret" dc:"accessKeySecret"` // accessKeySecret
	AllowTypeList      []string        `json:"allowTypeList" dc:"允许上传的文件类型"`         // 允许上传的文件类型
	MaxSize            int             `json:"maxSize" dc:"单个文件最大上传大小(MB)，默认1MB"`    // 单个文件最大上传大小(MB)，默认1MB
	MaxCount           int             `json:"maxCount" dc:"单次上传文件的最大数量，默认10"`       // 单次上传文件的最大数量，默认10
	Cache              gtkcache.ICache `json:"cache" dc:"缓存器"`                       // 缓存器
}

// UploadFileInfo 上传的文件信息
type UploadFileInfo struct {
	err      error  // 上传失败时返回的错误信息
	FileName string `json:"fileName" dc:"文件名"`  // 文件名
	FileSize int64  `json:"fileSize" dc:"文件大小"` // 文件大小
	FileUrl  string `json:"fileUrl" dc:"文件Url"` // 文件Url
	FileType string `json:"fileType" dc:"文件类型"` // 文件类型
}

// InitAliyunOSS 初始化阿里云`OSS`信息
func InitAliyunOSS(config *AliyunOSS) {
	if len(config.AllowTypeList) == 0 {
		config.AllowTypeList = []string{
			"jpg", "jpeg", "png", "gif",
			"doc", "docx", "xls", "xlsx", "ppt", "pptx", "pdf",
		}
	}
	if config.MaxSize == 0 {
		config.MaxSize = 1
	}
	if config.MaxCount == 0 {
		config.MaxCount = 10
	}
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
	var err error
	if err = r.ParseMultipartForm(10 << 20); err != nil {
		fileInfo = &UploadFileInfo{err: err}
		return
	}
	var (
		file       multipart.File
		fileHeader *multipart.FileHeader
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
	var fileUrl string
	if fileUrl, err = s.doUpload(file, fileHeader, dirPath); err != nil {
		fileInfo = &UploadFileInfo{err: err}
		return
	}
	// 返回
	fileInfo = &UploadFileInfo{
		FileName: fileHeader.Filename,
		FileSize: fileHeader.Size,
		FileUrl:  fileUrl,
		FileType: fileHeader.Header.Get("Content-type"),
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
	if err = r.ParseMultipartForm(10 << 20); err != nil {
		fileInfos = []*UploadFileInfo{{err: err}}
		return
	}
	// 检查文件数量
	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		fileInfos = []*UploadFileInfo{{err: errors.New("No Files Uploaded")}}
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
			var fileUrl string
			if fileUrl, err = s.doUpload(file, fileHeader, dirPath); err != nil {
				fileInfos[idx] = &UploadFileInfo{err: err}
				return
			}
			// 返回
			fileInfos[idx] = &UploadFileInfo{
				FileName: fileHeader.Filename,
				FileSize: fileHeader.Size,
				FileUrl:  fileUrl,
				FileType: fileHeader.Header.Get("Content-type"),
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
func (s *AliyunOSS) doUpload(file multipart.File, fileHeader *multipart.FileHeader, dirPath string) (fileUrl string, err error) {
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
	// 上传文件
	fileUrl = fmt.Sprintf("%s/%s", dirPath, gtkfile.GenRandomFileName(fileHeader.Filename))
	if err = bucket.PutObject(fileUrl, file); err != nil {
		return
	}
	return
}
