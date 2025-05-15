/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-23 10:31:17
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-15 17:56:30
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkoss

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

// AuthorizationDownload 授权给第三方下载
func (s *AliyunOSS) AuthorizationDownload(objectKey string, expiredInSec int64, opts ...Option) (fileUrl string, err error) {
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
	// 授权访问
	return bucket.SignURL(objectKey, oss.HTTPGet, expiredInSec, s.ossOptions...)
}
