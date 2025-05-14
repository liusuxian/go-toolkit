/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-24 20:51:23
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-13 13:11:50
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package oss

import (
	"context"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"time"
)

// GetAccessUrl 文件访问Url
func (s *AliyunOSS) GetAccessUrl(ctx context.Context, objectKey, cacheKey string, expiredInSec int64, options ...oss.Option) (fileUrl string, err error) {
	// 从缓存获取url
	if s.cache != nil {
		var val any
		if val, err = s.cache.Get(ctx, cacheKey); err != nil {
			return
		}
		// 如果缓存中存在url，则直接返回
		var ok bool
		if fileUrl, ok = val.(string); ok {
			if fileUrl != "" {
				return
			}
		}
	}
	// 授权访问
	if fileUrl, err = s.bucket.SignURL(objectKey, oss.HTTPGet, expiredInSec, options...); err != nil {
		return
	}
	// 添加缓存
	if s.cache != nil {
		err = s.cache.Set(ctx, cacheKey, fileUrl, time.Duration(expiredInSec)*time.Second)
	}
	return
}
