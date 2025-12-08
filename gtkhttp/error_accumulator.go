/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-07 18:25:49
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-08 23:00:25
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import (
	"bytes"
	"fmt"
	"io"
)

// ErrorAccumulator 错误收集器接口
type ErrorAccumulator interface {
	Write(p []byte) (err error) // Write
	Bytes() (errBytes []byte)   // Bytes
}

// errorBuffer 错误 Buffer 接口
type errorBuffer interface {
	io.Writer
	Len() (n int)
	Bytes() (b []byte)
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
		return fmt.Errorf("error accumulator write error, %w", err)
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
