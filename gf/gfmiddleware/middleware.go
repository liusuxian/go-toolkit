/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-19 21:15:17
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-22 21:35:59
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gfmiddleware

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"net/http"
)

// CORS 允许接口跨域请求
func CORS(req *ghttp.Request) {
	req.Response.CORSDefault()
	req.Middleware.Next()
}

// HandlerResponse 自定义返回中间件
func HandlerResponse(req *ghttp.Request) {
	req.Middleware.Next()

	// There's custom buffer content, it then exits current handler.
	if req.Response.BufferLength() > 0 {
		return
	}

	var (
		err   = req.GetError()
		res   = req.GetHandlerResponse()
		rCode = gerror.Code(err)
	)
	if err != nil {
		if rCode == gcode.CodeNil {
			rCode = gcode.CodeInternalError
		}
	} else {
		if req.Response.Status > 0 && req.Response.Status != http.StatusOK {
			msg := http.StatusText(req.Response.Status)
			switch req.Response.Status {
			case http.StatusNotFound:
				rCode = gcode.CodeNotFound
			case http.StatusForbidden:
				rCode = gcode.CodeNotAuthorized
			default:
				rCode = gcode.CodeUnknown
			}
			// It creates error as it can be retrieved by other middlewares.
			err = gerror.NewCode(rCode, msg)
			req.SetError(err)
		} else {
			rCode = gcode.CodeOK
		}
	}

	req.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    rCode.Code(),
		Message: rCode.Message(),
		Data:    res,
	})
}
