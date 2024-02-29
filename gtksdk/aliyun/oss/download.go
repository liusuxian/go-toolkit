/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-23 10:31:17
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-24 20:36:02
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package oss

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

// AuthorizationDownload 授权给第三方下载
func (s *AliyunOSS) AuthorizationDownload(filePath string, expiredInSec int64) (fileUrl string, err error) {
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
	// 授权访问
	fileUrl, err = bucket.SignURL(filePath, oss.HTTPGet, expiredInSec)
	return
}
