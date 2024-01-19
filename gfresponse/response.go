/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-19 21:04:44
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-19 21:12:26
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gfresponse

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Resp 数据返回通用JSON数据结构
type Resp struct {
	Code   int    `json:"code"`   // 错误码(0:成功, 非0:错误码)
	Msg    string `json:"msg"`    // 提示信息
	Detail any    `json:"detail"` // 提示详细信息
	Data   any    `json:"data"`   // 返回数据(业务接口定义具体数据结构)
}

// RespJsonCtx 标准返回结果数据
func RespJsonCtx(ctx context.Context, code int, msg string, detail any, data ...any) {
	RespJson(g.RequestFromCtx(ctx), code, msg, detail, data...)
}

// RespJson 标准返回结果数据
func RespJson(req *ghttp.Request, code int, msg string, detail any, data ...any) {
	var rData any
	if len(data) > 0 {
		rData = data[0]
	}

	resData := Resp{
		Code:   code,
		Msg:    msg,
		Detail: detail,
		Data:   rData,
	}
	req.Response.WriteJson(resData)
}

// RespJsonCtxByGcode 标准返回结果数据
func RespJsonCtxByGcode(ctx context.Context, ggcode gcode.Code, err error, data ...any) {
	if err != nil {
		ggcode = gcode.WithCode(ggcode, err.Error())
	}
	RespJsonCtx(ctx, ggcode.Code(), ggcode.Message(), ggcode.Detail(), data...)
}

// RespJsonByGcode 标准返回结果数据
func RespJsonByGcode(req *ghttp.Request, ggcode gcode.Code, err error, data ...any) {
	if err != nil {
		ggcode = gcode.WithCode(ggcode, err.Error())
	}
	RespJson(req, ggcode.Code(), ggcode.Message(), ggcode.Detail(), data...)
}

// RespJsonCtxExit 标准返回结果数据并退出
func RespJsonCtxExit(ctx context.Context, code int, msg string, detail any, data ...any) {
	RespJsonCtx(ctx, code, msg, detail, data...)
	g.RequestFromCtx(ctx).Exit()
}

// RespJsonExit 标准返回结果数据并退出
func RespJsonExit(req *ghttp.Request, code int, msg string, detail any, data ...any) {
	RespJson(req, code, msg, detail, data...)
	req.Exit()
}

// RespJsonCtxExitByGcode 标准返回结果数据并退出
func RespJsonCtxExitByGcode(ctx context.Context, ggcode gcode.Code, err error, data ...any) {
	if err != nil {
		ggcode = gcode.WithCode(ggcode, err.Error())
	}
	RespJsonCtx(ctx, ggcode.Code(), ggcode.Message(), ggcode.Detail(), data...)
	g.RequestFromCtx(ctx).Exit()
}

// RespJsonExitByGcode 标准返回结果数据并退出
func RespJsonExitByGcode(req *ghttp.Request, ggcode gcode.Code, err error, data ...any) {
	if err != nil {
		ggcode = gcode.WithCode(ggcode, err.Error())
	}
	RespJson(req, ggcode.Code(), ggcode.Message(), ggcode.Detail(), data...)
	req.Exit()
}

// SuccCtx 成功
func SuccCtx(ctx context.Context, data ...any) {
	RespJsonCtx(ctx, gcode.CodeOK.Code(), gcode.CodeOK.Message(), gcode.CodeOK.Detail(), data...)
}

// Succ 成功
func Succ(req *ghttp.Request, data ...any) {
	RespJson(req, gcode.CodeOK.Code(), gcode.CodeOK.Message(), gcode.CodeOK.Detail(), data...)
}

// SuccCtxExit 成功并退出
func SuccCtxExit(ctx context.Context, data ...any) {
	RespJsonCtxExit(ctx, gcode.CodeOK.Code(), gcode.CodeOK.Message(), gcode.CodeOK.Detail(), data...)
}

// SuccExit 成功并退出
func SuccExit(req *ghttp.Request, data ...any) {
	RespJsonExit(req, gcode.CodeOK.Code(), gcode.CodeOK.Message(), gcode.CodeOK.Detail(), data...)
}
