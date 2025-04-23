/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 14:26:01
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-15 15:17:58
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
)

var (
	headerData  = []byte("data: ")
	errorPrefix = []byte(`data: {"error":`)
)

// Streamable 可流式传输的类型
type Streamable interface {
}

// StreamReader 流读取器
type StreamReader[T Streamable] struct {
	emptyMessagesLimit uint
	isFinished         bool
	reader             *bufio.Reader
	response           *http.Response
	errAccumulator     ErrorAccumulator
	unmarshaler        Unmarshaler
	HttpHeader
}

// Recv 接收数据
func (stream *StreamReader[T]) Recv() (response T, err error) {
	var rawLine []byte
	if rawLine, err = stream.RecvRaw(); err != nil {
		return
	}

	if err = stream.unmarshaler.Unmarshal(rawLine, &response); err != nil {
		return
	}
	return
}

// RecvRaw 接收原始数据
func (stream *StreamReader[T]) RecvRaw() (b []byte, err error) {
	if stream.isFinished {
		return nil, io.EOF
	}

	return stream.processLines()
}

// processLines 处理行数据
func (stream *StreamReader[T]) processLines() (b []byte, err error) {
	var (
		emptyMessagesCount uint
		hasErrorPrefix     bool
	)

	for {
		rawLine, readErr := stream.reader.ReadBytes('\n')
		if readErr != nil || hasErrorPrefix {
			if respErr := stream.unmarshalError(); respErr != nil {
				return nil, fmt.Errorf("error, %v", respErr)
			}
			return nil, readErr
		}

		noSpaceLine := bytes.TrimSpace(rawLine)
		if bytes.HasPrefix(noSpaceLine, errorPrefix) {
			hasErrorPrefix = true
		}
		if !bytes.HasPrefix(noSpaceLine, headerData) || hasErrorPrefix {
			if hasErrorPrefix {
				noSpaceLine = bytes.TrimPrefix(noSpaceLine, headerData)
			}
			if writeErr := stream.errAccumulator.Write(noSpaceLine); writeErr != nil {
				return nil, writeErr
			}
			emptyMessagesCount++
			if emptyMessagesCount > stream.emptyMessagesLimit {
				return nil, ErrTooManyEmptyStreamMessages
			}
			continue
		}

		noPrefixLine := bytes.TrimPrefix(noSpaceLine, headerData)
		if string(noPrefixLine) == "[DONE]" {
			stream.isFinished = true
			return nil, io.EOF
		}

		return noPrefixLine, nil
	}
}

// unmarshalError 解析错误响应数据
func (stream *StreamReader[T]) unmarshalError() (errResp map[string]any) {
	var errBytes []byte
	if errBytes = stream.errAccumulator.Bytes(); len(errBytes) == 0 {
		return
	}

	if err := stream.unmarshaler.Unmarshal(errBytes, &errResp); err != nil {
		errResp = nil
		return
	}
	return
}

// Close 关闭流
func (stream *StreamReader[T]) Close() (err error) {
	return stream.response.Body.Close()
}
