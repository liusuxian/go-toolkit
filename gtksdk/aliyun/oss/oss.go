/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-08-29 17:06:47
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-13 16:14:51
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package oss

import (
	"context"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/liusuxian/go-toolkit/internal/utils"
	"time"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

// OSSConfig 阿里云OSS配置
type OSSConfig struct {
	Bucket             string   `json:"bucket"`             // bucket名称
	EndpointAccelerate string   `json:"endpointAccelerate"` // 传输加速节点
	EndpointInternal   string   `json:"endpointInternal"`   // 内网访问节点
	EndpointAccess     string   `json:"endpointAccess"`     // 外网访问节点
	AccessKeyID        string   `json:"accessKeyID"`        // accessKeyID
	AccessKeySecret    string   `json:"accessKeySecret"`    // accessKeySecret
	AllowTypeList      []string `json:"allowTypeList"`      // 允许上传的文件类型
	MaxSize            int      `json:"maxSize"`            // 单个文件最大上传大小(MB)，默认1MB
	MaxCount           int      `json:"maxCount"`           // 单次上传文件的最大数量，默认10
}

// Cache 缓存
type Cache interface {
	Get(ctx context.Context, key string, timeout ...time.Duration) (val any, err error) // 获取缓存
	Set(ctx context.Context, key string, val any, timeout ...time.Duration) (err error) // 设置缓存
}

// AliyunOSS 阿里云OSS服务
type AliyunOSS struct {
	config           OSSConfig                                  // 配置
	uploadFileNameFn func(filename string) (newFilename string) // 上传文件时生成文件名的函数
	cache            Cache                                      // 缓存器
	client           *oss.Client                                // 客户端
	bucket           *oss.Bucket                                // 存储空间
}

// Option 选项
type Option func(s *AliyunOSS)

// WithUploadFileNameFn 设置上传文件时生成文件名的函数
func WithUploadFileNameFn(fn func(filename string) (newFilename string)) (opt Option) {
	return func(s *AliyunOSS) {
		s.uploadFileNameFn = fn
	}
}

// WithCache 设置缓存器
func WithCache(cache Cache) (opt Option) {
	return func(s *AliyunOSS) {
		s.cache = cache
	}
}

// NewAliyunOSS 新建阿里云OSS服务
func NewAliyunOSS(config OSSConfig, opts ...Option) (s *AliyunOSS, err error) {
	s = &AliyunOSS{config: config}
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
			return utils.GenRandomFilename(filename)
		}
	}
	// 连接OSS
	if s.client, err = oss.New(s.config.EndpointAccelerate, s.config.AccessKeyID, s.config.AccessKeySecret); err != nil {
		return
	}
	// 获取存储空间
	if s.bucket, err = s.client.Bucket(s.config.Bucket); err != nil {
		return
	}
	return
}

// CloseIdleConnections 关闭空闲连接
func (s *AliyunOSS) CloseIdleConnections() {
	s.client.HTTPClient.CloseIdleConnections()
}
