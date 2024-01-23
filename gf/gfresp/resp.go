/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-19 21:04:44
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-23 17:47:10
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
	"net/http"
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

// Resp 响应数据返回
func (resp Response) Resp(req *ghttp.Request) {
	req.Response.WriteJson(resp)
}

// RespExit 响应数据返回并退出
func (resp Response) RespExit(req *ghttp.Request) {
	req.Response.WriteJson(resp)
	req.Exit()
}

// Succ 成功
func Succ(data any) (resp Response) {
	return Response{gcode.CodeOK.Code(), gcode.CodeOK.Message(), data}
}

// Fail 失败
func Fail(code int, msg string, data ...any) (resp Response) {
	var rData any
	if len(data) > 0 {
		rData = data[0]
	}
	return Response{code, msg, rData}
}

// Unauthorized 认证失败
func Unauthorized(msg string, data any) (resp Response) {
	return Response{http.StatusUnauthorized, msg, data}
}

// RespFail 返回失败
func RespFail(req *ghttp.Request, rCode gcode.Code, data ...any) {
	Fail(rCode.Code(), rCode.Message(), data...).Resp(req)
}

// RespFailExit 返回失败并退出
func RespFailExit(req *ghttp.Request, rCode gcode.Code, data ...any) {
	Fail(rCode.Code(), rCode.Message(), data...).RespExit(req)
}

// RespSucc 返回成功
func RespSucc(req *ghttp.Request, data any) {
	Succ(data).Resp(req)
}

// RespSuccExit 返回成功并退出
func RespSuccExit(req *ghttp.Request, data any) {
	Succ(data).RespExit(req)
}

// Redirect 重定向
func Redirect(req *ghttp.Request, link string) {
	req.Response.Header().Set("Location", link)
	req.Response.WriteHeader(302)
}
