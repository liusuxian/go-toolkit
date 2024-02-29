/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-26 21:02:16
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-26 21:21:29
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import (
	"bytes"
	"github.com/pkg/errors"
	"io"
)

// ErrorAccumulator 错误收集器接口
type ErrorAccumulator interface {
	Write(p []byte) (err error) // Write
	Bytes() (errBytes []byte)   // Bytes
}

// errorBuffer 错误`Buffer`接口
type errorBuffer interface {
	io.Writer
	Len() int
	Bytes() []byte
}

// DefaultErrorAccumulator 默认错误收集器
type DefaultErrorAccumulator struct {
	Buffer errorBuffer
}

// NewErrorAccumulator 新建默认错误收集器
func NewErrorAccumulator() (e ErrorAccumulator) {
	return &DefaultErrorAccumulator{
		Buffer: &bytes.Buffer{},
	}
}

// Write
func (e *DefaultErrorAccumulator) Write(p []byte) (err error) {
	if _, err = e.Buffer.Write(p); err != nil {
		err = errors.Errorf("error accumulator write error, %v", err)
		return
	}
	return
}

// Bytes
func (e *DefaultErrorAccumulator) Bytes() (errBytes []byte) {
	if e.Buffer.Len() == 0 {
		return
	}
	errBytes = e.Buffer.Bytes()
	return
}
