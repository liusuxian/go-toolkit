/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-07-13 20:17:18
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-07-13 21:09:13
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package oss

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

// DeleteObjects 删除多个对象
func (s *AliyunOSS) DeleteObjects(objectKeys ...string) (err error) {
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
	// 删除多个对象
	_, err = bucket.DeleteObjects(objectKeys)
	return
}
