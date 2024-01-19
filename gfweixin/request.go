/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-19 22:39:31
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-20 00:36:19
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gfweixin

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gurl"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/liusuxian/go-toolkit/tools"
	"net/url"
)

// Get 通用Get
func Get(ctx context.Context, url string, params url.Values) (resMap map[string]any, err error) {
	// 发起请求
	var result *gclient.Response
	result, err = g.Client().Get(ctx, url, gurl.BuildQuery(params))
	if err != nil {
		return
	}
	defer result.Close()
	// 处理结果
	resMap = gconv.Map(result.ReadAll())
	errCode := gconv.Int(resMap["errcode"])
	if errCode != 0 {
		err = gerror.NewCode(gcode.New(errCode, gconv.String(resMap["errmsg"]), ""))
		return
	}
	return
}

// Post 通用Post
func Post(ctx context.Context, url string, params url.Values, body map[string]any) (resMap map[string]any, err error) {
	// 发起请求
	if len(params) > 0 {
		url = fmt.Sprintf("%v?%v", url, gurl.BuildQuery(params))
	}
	var jsonBody string
	if jsonBody, err = tools.JsonMarshal(body); err != nil {
		return
	}
	var result *gclient.Response
	result, err = g.Client().Post(ctx, url, jsonBody)
	if err != nil {
		return
	}
	defer result.Close()
	// 处理结果
	resMap = gconv.Map(result.ReadAll())
	errCode := gconv.Int(resMap["errcode"])
	if errCode != 0 {
		err = gerror.NewCode(gcode.New(errCode, gconv.String(resMap["errmsg"]), ""))
		return
	}
	return
}
