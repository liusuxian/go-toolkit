/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-08-29 17:14:55
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-08-29 17:15:17
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

// UploadFileConfig 上传文件配置
type UploadFileConfig struct {
	AllowTypeList []string `json:"allowTypeList" dc:"允许上传的文件类型"`      // 允许上传的文件类型
	MaxSize       int      `json:"maxSize" dc:"单个文件最大上传大小(MB)，默认1MB"` // 单个文件最大上传大小(MB)，默认1MB
	MaxCount      int      `json:"maxCount" dc:"单次上传文件的最大数量，默认10"`    // 单次上传文件的最大数量，默认10
}

// InitUploadFileConfig 初始化上传文件配置
func InitUploadFileConfig(config *UploadFileConfig) {
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
