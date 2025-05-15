/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-07-13 20:17:18
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-15 17:54:22
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkoss

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

// DeleteObjects 删除多个对象
func (s *AliyunOSS) DeleteObjects(objectKeys []string, opts ...Option) (deletedObjects []string, err error) {
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
	// 删除对象
	var result oss.DeleteObjectsResult
	if result, err = bucket.DeleteObjects(objectKeys); err != nil {
		return
	}
	deletedObjects = result.DeletedObjects
	return
}
