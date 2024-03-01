/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-24 20:51:23
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-29 16:39:26
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package oss

import (
	"context"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/liusuxian/go-toolkit/gtkconv"
	"time"
)

// GetAccessUrl 文件访问Url
func (s *AliyunOSS) GetAccessUrl(ctx context.Context, filePath string, cacheKey string) (fileUrl string, err error) {
	// 从缓存获取url
	if s.Cache != nil {
		var val any
		if val, err = s.Cache.Get(ctx, cacheKey); err != nil {
			return
		}
		if fileUrl, err = gtkconv.ToStringE(val); err != nil {
			return
		}
		if fileUrl != "" {
			return
		}
	}
	// 连接OSS
	var client *oss.Client
	if client, err = oss.New(s.EndpointAccess, s.AccessKeyID, s.AccessKeySecret); err != nil {
		return
	}
	// 获取存储空间
	var bucket *oss.Bucket
	if bucket, err = client.Bucket(s.Bucket); err != nil {
		return
	}
	// 授权访问
	if fileUrl, err = bucket.SignURL(filePath, oss.HTTPGet, 3600); err != nil {
		return
	}
	// 添加缓存
	if s.Cache != nil {
		err = s.Cache.Set(ctx, cacheKey, fileUrl, time.Hour)
	}
	return
}
