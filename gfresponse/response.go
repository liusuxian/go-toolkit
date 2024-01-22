/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-19 21:04:44
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-22 19:41:32
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

// RespJsonCtx 标准返回结果数据
func RespJsonCtx(ctx context.Context, rCode gcode.Code, data ...any) {
	RespJson(g.RequestFromCtx(ctx), rCode, data...)
}

// RespJson 标准返回结果数据
func RespJson(req *ghttp.Request, rCode gcode.Code, data ...any) {
	var rData any
	if len(data) > 0 {
		rData = data[0]
	}

	resData := ghttp.DefaultHandlerResponse{
		Code:    rCode.Code(),
		Message: rCode.Message(),
		Data:    rData,
	}
	req.Response.WriteJson(resData)
}

// RespJsonCtxExit 标准返回结果数据并退出
func RespJsonCtxExit(ctx context.Context, rCode gcode.Code, data ...any) {
	RespJsonCtx(ctx, rCode, data...)
	g.RequestFromCtx(ctx).Exit()
}

// RespJsonExit 标准返回结果数据并退出
func RespJsonExit(req *ghttp.Request, rCode gcode.Code, data ...any) {
	RespJson(req, rCode, data...)
	req.Exit()
}

// SuccCtx 成功
func SuccCtx(ctx context.Context, data ...any) {
	RespJsonCtx(ctx, gcode.CodeOK, data...)
}

// Succ 成功
func Succ(req *ghttp.Request, data ...any) {
	RespJson(req, gcode.CodeOK, data...)
}

// SuccCtxExit 成功并退出
func SuccCtxExit(ctx context.Context, data ...any) {
	RespJsonCtxExit(ctx, gcode.CodeOK, data...)
}

// SuccExit 成功并退出
func SuccExit(req *ghttp.Request, data ...any) {
	RespJsonExit(req, gcode.CodeOK, data...)
}

// Redirect 重定向
func Redirect(req *ghttp.Request, link string) {
	req.Response.Header().Set("Location", link)
	req.Response.WriteHeader(302)
}
