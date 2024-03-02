/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-29 16:41:35
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-03-02 18:02:54
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkresp

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/liusuxian/go-toolkit/gtkconv"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

const (
	CodeSuccess int = 0
)

// Response 通用响应数据
type Response struct {
	Code    int    `json:"code"    dc:"错误码(0:成功, 非0:错误)"`   // 错误码(0:成功, 非0:错误)
	Message string `json:"message" dc:"错误消息"`               // 错误消息
	Data    any    `json:"data"    dc:"根据API定义，对特定请求的结果数据"` // 根据`API`定义，对特定请求的结果数据
}

// Success 判断是否成功
func (resp Response) Success() (ok bool) {
	return resp.Code == CodeSuccess
}

// DataString 获取`Data`转字符串
func (resp Response) DataString() (data string) {
	return gtkconv.ToString(resp.Data)
}

// DataInt 获取`Data`转`Int`
func (resp Response) DataInt() (data int) {
	return gtkconv.ToInt(resp.Data)
}

// GetString 获取`Data`值转字符串
func (resp Response) GetString(key string) (data string) {
	return gtkconv.ToString(resp.Get(key))
}

// GetInt 获取`Data`值转`Int`
func (resp Response) GetInt(key string) (data int) {
	return gtkconv.ToInt(resp.Get(key))
}

// Get 获取`Data`值
func (resp Response) Get(key string) (data any) {
	m := gtkconv.ToStringMap(resp.Data)
	if m == nil {
		return nil
	}
	return m[key]
}

// Json
func (resp Response) Json() (b []byte, err error) {
	return json.Marshal(resp)
}

// MustJson
func (resp Response) MustJson() (b []byte) {
	b, _ = json.Marshal(resp)
	return
}

// Succ 成功
func Succ(data any) (resp Response) {
	return Response{Data: data}
}

// Fail 失败
func Fail(code int, msg string, data ...any) (resp Response) {
	var rData any
	if len(data) > 0 {
		rData = data[0]
	}
	return Response{Code: code, Message: msg, Data: rData}
}

// Unauthorized 认证失败
func Unauthorized(msg string, data any) (resp Response) {
	return Response{Code: http.StatusUnauthorized, Message: msg, Data: data}
}

// RespSucc 返回成功
func RespSucc(w http.ResponseWriter, data any) {
	WriteJson(w, Succ(data))
}

// RespFail 返回失败
func RespFail(w http.ResponseWriter, code int, msg string, data ...any) {
	WriteJson(w, Fail(code, msg, data...))
}

// Write
func Write(w http.ResponseWriter, data ...any) {
	if len(data) == 0 {
		return
	}
	for _, v := range data {
		switch val := v.(type) {
		case []byte:
			w.Write(val)
		case string:
			w.Write([]byte(val))
		default:
			w.Write(gtkconv.ToBytes(v))
		}
	}
}

// Writef
func Writef(w http.ResponseWriter, format string, params ...any) {
	Write(w, fmt.Sprintf(format, params...))
}

// Writeln
func Writeln(w http.ResponseWriter, data ...any) {
	if len(data) == 0 {
		Write(w, "\n")
		return
	}
	Write(w, append(data, "\n")...)
}

// Writefln
func Writefln(w http.ResponseWriter, format string, params ...any) {
	Writeln(w, fmt.Sprintf(format, params...))
}

// WriteStatus
func WriteStatus(w http.ResponseWriter, status int, data ...any) {
	w.WriteHeader(status)
	if len(data) > 0 {
		Write(w, data...)
	} else {
		Write(w, http.StatusText(status))
	}
}

// WriteJson
func WriteJson(w http.ResponseWriter, resp Response) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if b, err := resp.Json(); err != nil {
		w.Write(Fail(http.StatusInternalServerError, errors.Wrap(err, "WriteJson Failed").Error()).MustJson())
	} else {
		w.Write(b)
	}
}

// RespSSESucc 返回`SSE`事件成功
func RespSSESucc(w http.ResponseWriter, data any) {
	// 设置`SSE`的响应头信息
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// 发送数据：'data: <jsonData>\n\n'
	fmt.Fprintf(w, "data: %s\n\n", Succ(data).MustJson())
	// 确保即时发送数据
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

// RespSSEFail 返回`SSE`事件失败
func RespSSEFail(w http.ResponseWriter, code int, msg string, data ...any) {
	// 设置`SSE`的响应头信息
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// 发送数据：'data: <jsonData>\n\n'
	fmt.Fprintf(w, "data: %s\n\n", Fail(code, msg, data...).MustJson())
	// 确保即时发送数据
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

// Redirect 重定向
func Redirect(w http.ResponseWriter, link string) {
	w.Header().Set("Location", link)
	w.WriteHeader(302)
}

// WriteSuccMessage 写成功响应消息
func WriteSuccMessage(ws *websocket.Conn, messageType int, data any) (err error) {
	var dataBytes []byte
	if dataBytes, err = gtkconv.ToBytesE(Succ(data)); err != nil {
		return
	}
	return ws.WriteMessage(messageType, dataBytes)
}

// WriteFailMessage 写失败响应消息
func WriteFailMessage(ws *websocket.Conn, messageType, code int, msg string, data ...any) (err error) {
	var dataBytes []byte
	if dataBytes, err = gtkconv.ToBytesE(Fail(code, msg, data...)); err != nil {
		return
	}
	return ws.WriteMessage(messageType, dataBytes)
}

// WriteMessage 写任意消息
func WriteMessage(ws *websocket.Conn, messageType int, data any) (err error) {
	var dataBytes []byte
	if dataBytes, err = gtkconv.ToBytesE(data); err != nil {
		return
	}
	return ws.WriteMessage(messageType, dataBytes)
}

// WriteControl 使用给定的截止时间写入一个控制消息
//
//	允许的消息类型包括 `CloseMessage`，`PingMessage`，`PongMessage`
func WriteControl(ws *websocket.Conn, messageType int, data any, deadline time.Time) (err error) {
	var dataBytes []byte
	if dataBytes, err = gtkconv.ToBytesE(data); err != nil {
		return
	}
	return ws.WriteControl(messageType, dataBytes, deadline)
}
