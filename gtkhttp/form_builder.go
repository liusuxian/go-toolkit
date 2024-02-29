/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-26 20:49:54
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-26 20:59:55
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import (
	"io"
	"mime/multipart"
)

// FormBuilder 表单构建器接口
type FormBuilder interface {
	WriteField(fieldName, value string) (err error) // 向表单中添加字段
	Close() (err error)                             // 关闭表单构建器
	FormDataContentType() (mime string)             // 获取表单数据的`MIME`类型
}

// DefaultFormBuilder 默认表单构建器
type DefaultFormBuilder struct {
	writer *multipart.Writer
}

// NewFormBuilder 新建默认表单构建器
func NewFormBuilder(body io.Writer) (fb *DefaultFormBuilder) {
	return &DefaultFormBuilder{
		writer: multipart.NewWriter(body),
	}
}

// WriteField 向表单中添加字段
func (fb *DefaultFormBuilder) WriteField(fieldName, value string) (err error) {
	return fb.writer.WriteField(fieldName, value)
}

// Close 关闭表单构建器
func (fb *DefaultFormBuilder) Close() (err error) {
	return fb.writer.Close()
}

// FormDataContentType 获取表单数据的`MIME`类型
func (fb *DefaultFormBuilder) FormDataContentType() (mime string) {
	return fb.writer.FormDataContentType()
}
