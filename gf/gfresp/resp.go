/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-19 21:04:44
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-23 16:41:48
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gfresp

import (
	"encoding/json"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/liusuxian/go-toolkit/gf/gflogger"
	"net/http"
)

const (
	FAIL         = -2
	ERROR        = -3
	UNAUTHORIZED = http.StatusUnauthorized
)

// Response 通用响应数据结构
type Response struct {
	Code    int         `json:"code"    dc:"错误码(0:成功, 非0:错误)"`     // 错误码(0:成功, 非0:错误)
	Message string      `json:"message" dc:"错误消息"`                 // 错误消息
	Data    interface{} `json:"data"    dc:"根据 API 定义，对特定请求的结果数据"` // 根据 API 定义，对特定请求的结果数据
}

// Success 判断是否成功
func (resp Response) Success() (ok bool) {
	return resp.Code == gcode.CodeOK.Code()
}

// DataString 获取Data转字符串
func (resp Response) DataString() (data string) {
	return gconv.String(resp.Data)
}

// DataInt 获取Data转Int
func (resp Response) DataInt() (data int) {
	return gconv.Int(resp.Data)
}

// GetString 获取Data值转字符串
func (resp Response) GetString(key string) (data string) {
	return gconv.String(resp.Get(key))
}

// GetInt 获取Data值转Int
func (resp Response) GetInt(key string) (data int) {
	return gconv.Int(resp.Get(key))
}

// Get 获取Data值
func (resp Response) Get(key string) (data *gvar.Var) {
	m := gconv.Map(resp.Data)
	if m == nil {
		return nil
	}
	return gvar.New(m[key])
}

// Json
func (resp Response) Json() (str string) {
	b, _ := json.Marshal(resp)
	return string(b)
}

// Succ 成功
func Succ(data any) (resp Response) {
	return Response{gcode.CodeOK.Code(), gcode.CodeOK.Message(), data}
}

// Fail 失败
func Fail(msg string) (resp Response) {
	return Response{FAIL, msg, ""}
}

// FailData 失败设置Data
func FailData(msg string, data any) (resp Response) {
	return Response{FAIL, msg, data}
}

// Error 错误
func Error(msg string) (resp Response) {
	return Response{ERROR, msg, ""}
}

// ErrorData 错误设置Data
func ErrorData(msg string, data any) (resp Response) {
	return Response{ERROR, msg, data}
}

// Unauthorized 认证失败
func Unauthorized(msg string, data any) (resp Response) {
	return Response{UNAUTHORIZED, msg, data}
}

// Resp 响应数据返回
func Resp(req *ghttp.Request, rCode gcode.Code, err error, data ...any) {
	var rData any
	if len(data) > 0 {
		rData = data[0]
	}

	if err != nil {
		gflogger.Errorf(req.GetCtx(), "Response Error: %+v", err)
	}

	resData := Response{
		Code:    rCode.Code(),
		Message: rCode.Message(),
		Data:    rData,
	}
	req.Response.WriteJson(resData)
}

// RespNoErr 响应数据返回
func RespNoErr(req *ghttp.Request, rCode gcode.Code, data ...any) {
	Resp(req, rCode, nil, data...)
}

// RespExit 响应数据返回并退出
func RespExit(req *ghttp.Request, rCode gcode.Code, err error, data ...any) {
	Resp(req, rCode, err, data...)
	req.Exit()
}

// RespNoErrExit 响应数据返回并退出
func RespNoErrExit(req *ghttp.Request, rCode gcode.Code, data ...any) {
	Resp(req, rCode, nil, data...)
	req.Exit()
}

// RespSucc 成功
func RespSucc(req *ghttp.Request, data ...any) {
	Resp(req, gcode.CodeOK, nil, data...)
}

// RespSuccExit 成功并退出
func RespSuccExit(req *ghttp.Request, data ...any) {
	Resp(req, gcode.CodeOK, nil, data...)
	req.Exit()
}

// Redirect 重定向
func Redirect(req *ghttp.Request, link string) {
	req.Response.Header().Set("Location", link)
	req.Response.WriteHeader(302)
}
