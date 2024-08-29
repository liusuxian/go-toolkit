/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-08-29 17:06:47
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-08-29 17:07:53
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package oss

import "github.com/liusuxian/go-toolkit/gtkcache"

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

// AliyunOSS 阿里云 OSS 信息
type AliyunOSS struct {
	Bucket             string          `json:"bucket" dc:"bucket名称"`                 // bucket名称
	EndpointAccelerate string          `json:"endpointAccelerate" dc:"传输加速节点"`       // 传输加速节点
	EndpointInternal   string          `json:"endpointInternal" dc:"内网访问节点"`         // 内网访问节点
	EndpointAccess     string          `json:"endpointAccess" dc:"外网访问节点"`           // 外网访问节点
	AccessKeyID        string          `json:"accessKeyID" dc:"accessKeyID"`         // accessKeyID
	AccessKeySecret    string          `json:"accessKeySecret" dc:"accessKeySecret"` // accessKeySecret
	AllowTypeList      []string        `json:"allowTypeList" dc:"允许上传的文件类型"`         // 允许上传的文件类型
	MaxSize            int             `json:"maxSize" dc:"单个文件最大上传大小(MB)，默认1MB"`    // 单个文件最大上传大小(MB)，默认1MB
	MaxCount           int             `json:"maxCount" dc:"单次上传文件的最大数量，默认10"`       // 单次上传文件的最大数量，默认10
	Cache              gtkcache.ICache `json:"cache" dc:"缓存器"`                       // 缓存器
}

// InitAliyunOSS 初始化阿里云`OSS`信息
func InitAliyunOSS(config *AliyunOSS) {
	if len(config.AllowTypeList) == 0 {
		config.AllowTypeList = []string{
			"jpg", "jpeg", "png", "gif",
			"doc", "docx", "xls", "xlsx", "ppt", "pptx", "pdf",
		}
	}
	if config.MaxSize == 0 {
		config.MaxSize = 1
	}
	if config.MaxCount == 0 {
		config.MaxCount = 10
	}
}
