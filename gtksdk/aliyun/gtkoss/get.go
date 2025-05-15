/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-13 11:23:37
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-15 17:59:59
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkoss

import (
	"context"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
)

// GetObject 获取文件
func (s *AliyunOSS) GetObject(ctx context.Context, objectKey string, opts ...Option) (b []byte, err error) {
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
	// 获取文件
	var obj io.ReadCloser
	if obj, err = bucket.GetObject(objectKey, s.ossOptions...); err != nil {
		return
	}
	defer obj.Close()
	// 返回
	return io.ReadAll(obj)
}

// GetObjectWithURL 获取文件
func (s *AliyunOSS) GetObjectWithURL(ctx context.Context, signedURL string, opts ...Option) (b []byte, err error) {
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
	// 获取文件
	var obj io.ReadCloser
	if obj, err = bucket.GetObjectWithURL(signedURL, s.ossOptions...); err != nil {
		return
	}
	defer obj.Close()
	// 返回
	return io.ReadAll(obj)
}

// GetObjectToFile 获取文件并保存到本地
func (s *AliyunOSS) GetObjectToFile(ctx context.Context, objectKey, filePath string, opts ...Option) (err error) {
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
	// 获取文件并保存到本地
	return bucket.GetObjectToFile(objectKey, filePath, s.ossOptions...)
}

// GetObjectToFileWithURL 获取文件并保存到本地
func (s *AliyunOSS) GetObjectToFileWithURL(ctx context.Context, signedURL, filePath string, opts ...Option) (err error) {
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
	// 获取文件并保存到本地
	return bucket.GetObjectToFileWithURL(signedURL, filePath, s.ossOptions...)
}
