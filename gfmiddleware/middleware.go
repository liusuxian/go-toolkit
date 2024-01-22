/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-19 21:15:17
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-22 18:56:57
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gfmiddleware

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/liusuxian/go-toolkit/gfresponse"
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
		detail   string
		err      = req.GetError()
		respData = req.GetHandlerResponse()
		rCode    = gerror.Code(err)
	)
	if err != nil {
		if rCode == gcode.CodeNil {
			rCode = gcode.CodeInternalError
		}
		detail = err.Error()
		rCode = gcode.WithCode(rCode, detail)
	} else {
		if req.Response.Status > 0 && req.Response.Status != http.StatusOK {
			detail = http.StatusText(req.Response.Status)
			switch req.Response.Status {
			case http.StatusNotFound:
				rCode = gcode.CodeNotFound
			case http.StatusForbidden:
				rCode = gcode.CodeNotAuthorized
			default:
				rCode = gcode.CodeUnknown
			}
			rCode = gcode.WithCode(rCode, detail)
			// It creates error as it can be retrieved by other middlewares.
			err = gerror.NewCode(rCode)
			req.SetError(err)
		} else {
			rCode = gcode.CodeOK
		}
	}

	req.Response.WriteJson(gfresponse.Resp{
		Code:   rCode.Code(),
		Msg:    rCode.Message(),
		Detail: rCode.Detail(),
		Data:   respData,
	})
}
