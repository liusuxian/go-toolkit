/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-08-29 17:06:47
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-15 17:50:14
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkoss

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

// EndpointType 节点类型
type EndpointType string

const (
	EndpointTypeAccelerate EndpointType = "accelerate"
	EndpointTypeInternal   EndpointType = "internal"
	EndpointTypeAccess     EndpointType = "access"
)

// AliyunOSS 阿里云OSS服务
type AliyunOSS struct {
	config           OSSConfig                                  // 配置
	uploadFileNameFn func(filename string) (newFilename string) // 上传文件时生成文件名的函数
	cache            Cache                                      // 缓存器
	ossClientOptions []oss.ClientOption                         // 阿里云OSS客户端选项
	ossOptions       []oss.Option                               // 阿里云OSS操作选项
	endpointType     EndpointType                               // 节点类型
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

// WithOSSClientOptions 设置阿里云OSS客户端选项
func WithOSSClientOptions(ossClientOptions ...oss.ClientOption) (opt Option) {
	return func(s *AliyunOSS) {
		s.ossClientOptions = ossClientOptions
	}
}

// WithOSSOptions 设置阿里云OSS操作选项
func WithOSSOptions(ossOptions ...oss.Option) (opt Option) {
	return func(s *AliyunOSS) {
		s.ossOptions = ossOptions
	}
}

// WithEndpointType 设置节点类型
func WithEndpointType(endpointType EndpointType) (opt Option) {
	return func(s *AliyunOSS) {
		s.endpointType = endpointType
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
	// 设置默认节点类型
	if s.endpointType == "" {
		s.endpointType = EndpointTypeAccelerate
	}
	return
}

// getBucket 获取存储空间
func (s *AliyunOSS) getBucket(opts ...Option) (client *oss.Client, bucket *oss.Bucket, err error) {
	// 设置选项
	for _, opt := range opts {
		opt(s)
	}
	// 根据节点类型获取节点
	var endpoint string
	switch s.endpointType {
	case EndpointTypeAccelerate:
		endpoint = s.config.EndpointAccelerate
	case EndpointTypeInternal:
		endpoint = s.config.EndpointInternal
	case EndpointTypeAccess:
		endpoint = s.config.EndpointAccess
	default:
		endpoint = s.config.EndpointAccelerate
	}
	// 连接OSS
	if client, err = oss.New(endpoint, s.config.AccessKeyID, s.config.AccessKeySecret, s.ossClientOptions...); err != nil {
		return
	}
	// 获取存储空间
	if bucket, err = client.Bucket(s.config.Bucket); err != nil {
		return
	}
	return
}

// closeIdleConnections 关闭空闲连接
func (s *AliyunOSS) closeIdleConnections(client *oss.Client) {
	if client.HTTPClient != nil {
		client.HTTPClient.CloseIdleConnections()
	}
}
