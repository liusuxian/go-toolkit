/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-08-29 17:14:55
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-25 10:33:31
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

// UploadFileConfig 上传文件配置
type UploadFileConfig struct {
	AllowTypeList []string `json:"allowTypeList"` // 允许上传的文件类型
	MaxSize       int      `json:"maxSize"`       // 单个文件最大上传大小(MB)，默认1MB
	MaxCount      int      `json:"maxCount"`      // 单次上传文件的最大数量，默认10
}

// UploadFileService 上传文件服务
type UploadFileService struct {
	config           UploadFileConfig                           // 配置
	uploadFileNameFn func(filename string) (newFilename string) // 上传文件时生成文件名的函数
}

// Option 选项
type Option func(s *UploadFileService)

// WithUploadFileNameFn 设置上传文件时生成文件名的函数
func WithUploadFileNameFn(fn func(filename string) (newFilename string)) (opt Option) {
	return func(s *UploadFileService) {
		s.uploadFileNameFn = fn
	}
}

// NewUploadFileService 创建上传文件服务
func NewUploadFileService(config UploadFileConfig, opts ...Option) (s *UploadFileService) {
	s = &UploadFileService{config: config}
	// 设置配置默认值
	if len(s.config.AllowTypeList) == 0 {
		s.config.AllowTypeList = []string{
			"jpg", "jpeg", "png", "gif",
			"doc", "docx", "xls", "xlsx", "ppt", "pptx", "pdf",
		}
	}
	if s.config.MaxSize == 0 {
		s.config.MaxSize = 1
	}
	if s.config.MaxCount == 0 {
		s.config.MaxCount = 10
	}
	// 设置选项
	for _, opt := range opts {
		opt(s)
	}
	// 设置默认上传文件时生成文件名的函数
	if s.uploadFileNameFn == nil {
		s.uploadFileNameFn = func(filename string) (newFilename string) {
			return filename
		}
	}
	return
}
