/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-19 21:15:17
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-19 21:29:48
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gfmiddleware

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-toolkit/gfresponse"
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
	if req.Response.BufferLength() > 0 {
		return
	}

	var (
		msg   string = "OK"
		err          = req.GetError()
		res          = req.GetHandlerResponse()
		rCode        = gerror.Code(err)
	)
	if err != nil {
		if rCode == gcode.CodeNil {
			rCode = gcode.CodeInternalError
		}
		msg = err.Error()
	} else {
		if req.Response.Status > 0 && req.Response.Status != http.StatusOK {
			msg = http.StatusText(req.Response.Status)
			switch req.Response.Status {
			case http.StatusNotFound:
				rCode = gcode.CodeNotFound
			case http.StatusForbidden:
				rCode = gcode.CodeNotAuthorized
			default:
				rCode = gcode.CodeUnknown
			}

			err = gerror.NewCode(rCode, msg)
			req.SetError(err)
		} else {
			rCode = gcode.CodeOK
		}
	}

	req.Response.WriteJson(gfresponse.Resp{
		Code:   rCode.Code(),
		Msg:    msg,
		Detail: rCode.Detail(),
		Data:   res,
	})
}
