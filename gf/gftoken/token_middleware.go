/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-23 23:38:09
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-24 01:15:15
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gftoken

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/liusuxian/go-toolkit/gf/gflogger"
	"strings"
)

// GroupMiddleware
func (m *Token) GroupMiddleware(ctx context.Context, group *ghttp.RouterGroup) {
	// 初始化配置信息
	m.InitConfig()
	// 设置为Group模式
	m.MiddlewareType = MiddlewareTypeGroup
	gflogger.Info(ctx, "[gftoken][params:"+m.String()+"]start... ")
	// 缓存模式
	if m.CacheMode > CacheModeFile {
		gflogger.Panic(ctx, "[gftoken]CacheMode set error")
	}
	// 初始化文件缓存
	if m.CacheMode == 3 {
		m.initFileCache(ctx)
	}
	group.Middleware(m.authMiddleware)
}

// BindMiddleware
func (m *Token) BindMiddleware(ctx context.Context, s *ghttp.Server) {
	// 初始化配置信息
	m.InitConfig()
	// 设置为Group模式
	m.MiddlewareType = MiddlewareTypeBind
	gflogger.Info(ctx, "[gftoken][params:"+m.String()+"]start... ")
	// 缓存模式
	if m.CacheMode > CacheModeFile {
		gflogger.Panic(ctx, "[gftoken]CacheMode set error")
	}
	// 初始化文件缓存
	if m.CacheMode == 3 {
		m.initFileCache(ctx)
	}
	for _, authPath := range m.AuthPaths {
		tmpPath := authPath
		if !strings.HasSuffix(authPath, "/*") {
			tmpPath += "/*"
		}
		s.BindMiddleware(tmpPath, m.authMiddleware)
	}
}

// GlobalMiddleware
func (m *Token) GlobalMiddleware(ctx context.Context, s *ghttp.Server) {
	// 初始化配置信息
	m.InitConfig()
	// 设置为Group模式
	m.MiddlewareType = MiddlewareTypeGlobal
	gflogger.Info(ctx, "[gftoken][params:"+m.String()+"]start... ")
	// 缓存模式
	if m.CacheMode > CacheModeFile {
		gflogger.Panic(ctx, "[gftoken]CacheMode set error")
	}
	// 初始化文件缓存
	if m.CacheMode == 3 {
		m.initFileCache(ctx)
	}
	s.BindMiddlewareDefault(m.authMiddleware)
}
