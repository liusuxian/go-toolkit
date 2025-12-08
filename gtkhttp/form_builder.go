/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-07 18:27:52
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-08 22:56:52
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

// 转义引号
var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

// FormBuilderHandler 构建表单请求体处理函数
type FormBuilderHandler func(builder FormBuilder) (err error)

// FormBuilder 表单构建器接口
type FormBuilder interface {
	CreateFormFile(fieldname string, file *os.File) (err error)                      //	向表单中添加文件
	CreateFormFileReader(fieldname string, r io.Reader, filename string) (err error) // 向表单中添加文件流
	WriteField(fieldname, value string) (err error)                                  // 向表单中添加字段
	Close() (err error)                                                              // 关闭表单构建器
	FormDataContentType() (mime string)                                              // 获取表单数据的 MIME 类型
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

// CreateFormFile 向表单中添加文件
func (fb *DefaultFormBuilder) CreateFormFile(fieldname string, file *os.File) (err error) {
	return fb.createFormFile(fieldname, file, file.Name())
}

// CreateFormFileReader 向表单中添加文件流
func (fb *DefaultFormBuilder) CreateFormFileReader(fieldname string, r io.Reader, filename string) (err error) {
	if filename == "" {
		if f, ok := r.(interface{ Name() string }); ok {
			filename = f.Name()
		}
	}
	var contentType string
	if f, ok := r.(interface{ ContentType() string }); ok {
		contentType = f.ContentType()
	}

	h := make(textproto.MIMEHeader)
	h.Set(
		"Content-Disposition",
		fmt.Sprintf(
			`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname),
			escapeQuotes(filepath.Base(filename)),
		),
	)
	if contentType != "" {
		h.Set("Content-Type", contentType)
	}

	var fieldWriter io.Writer
	if fieldWriter, err = fb.writer.CreatePart(h); err != nil {
		return
	}

	if _, err = io.Copy(fieldWriter, r); err != nil {
		return
	}
	return
}

// createFormFile 向表单中添加文件
func (fb *DefaultFormBuilder) createFormFile(fieldname string, r io.Reader, filename string) (err error) {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	var fieldWriter io.Writer
	if fieldWriter, err = fb.writer.CreateFormFile(fieldname, filename); err != nil {
		return
	}

	if _, err = io.Copy(fieldWriter, r); err != nil {
		return
	}
	return
}

// WriteField 向表单中添加字段
func (fb *DefaultFormBuilder) WriteField(fieldname, value string) (err error) {
	if fieldname == "" {
		return fmt.Errorf("fieldname cannot be empty")
	}
	return fb.writer.WriteField(fieldname, value)
}

// Close 关闭表单构建器
func (fb *DefaultFormBuilder) Close() (err error) {
	return fb.writer.Close()
}

// FormDataContentType 获取表单数据的 MIME 类型
func (fb *DefaultFormBuilder) FormDataContentType() (mime string) {
	return fb.writer.FormDataContentType()
}

// escapeQuotes 转义引号
func escapeQuotes(s string) (str string) {
	return quoteEscaper.Replace(s)
}
