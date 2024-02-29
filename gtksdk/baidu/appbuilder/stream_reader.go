/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-27 10:20:04
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-27 14:34:22
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package appbuilder

import (
	"bufio"
	"bytes"
	"github.com/liusuxian/go-toolkit/gtkhttp"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

var (
	headerData = []byte("data: ")
)

// streamable 可流式传输的类型
type streamable interface {
	IntegratedResponseResult
}

// streamReader 流读取器
type streamReader[T streamable] struct {
	emptyMessagesLimit uint
	isFinished         bool
	reader             *bufio.Reader
	response           *http.Response
	unmarshaler        gtkhttp.Unmarshaler
}

// Recv 接收数据
func (stream *streamReader[T]) Recv() (response T, err error) {
	if stream.isFinished {
		err = io.EOF
		return
	}

	response, err = stream.processLines()
	return
}

// processLines 处理行数据
func (stream *streamReader[T]) processLines() (response T, err error) {
	response = T{}
	var emptyMessagesCount uint

	for {
		var rawLine []byte
		if rawLine, err = stream.reader.ReadBytes('\n'); err != nil {
			return
		}

		noSpaceLine := bytes.TrimSpace(rawLine)
		if !bytes.HasPrefix(noSpaceLine, headerData) {
			emptyMessagesCount++
			if emptyMessagesCount > stream.emptyMessagesLimit {
				err = errors.New("stream has sent too many empty messages")
				return
			}
			continue
		}

		noPrefixLine := bytes.TrimPrefix(noSpaceLine, headerData)
		if len(noPrefixLine) == 0 {
			stream.isFinished = true
			err = io.EOF
			return
		}

		if err = stream.unmarshaler.Unmarshal(noPrefixLine, &response); err != nil {
			return
		}
		return
	}
}

// Close 关闭流
func (stream *streamReader[T]) Close() {
	stream.response.Body.Close()
}
