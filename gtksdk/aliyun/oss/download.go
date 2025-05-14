/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-23 10:31:17
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-13 13:12:35
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package oss

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

// AuthorizationDownload 授权给第三方下载
func (s *AliyunOSS) AuthorizationDownload(objectKey string, expiredInSec int64, options ...oss.Option) (fileUrl string, err error) {
	// 授权访问
	return s.bucket.SignURL(objectKey, oss.HTTPGet, expiredInSec, options...)
}
