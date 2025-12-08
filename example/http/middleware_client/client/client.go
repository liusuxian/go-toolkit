/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-12-08 23:17:18
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-09 01:01:22
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/liusuxian/go-toolkit/example/http/middleware_client/models"
	"github.com/liusuxian/go-toolkit/example/http/middleware_client/providers/deepseek"
	"github.com/liusuxian/go-toolkit/gtkflake"
	"github.com/liusuxian/go-toolkit/gtkhttp"
	"sort"
	"time"
)

var (
	ErrFailedToCreateFlakeInstance = errors.New("failed to create flake instance") // 创建分布式唯一ID生成器失败
)

// MWClient 带中间件的客户端
type MWClient struct {
	flakeInstance   *gtkflake.Flake // 分布式唯一ID生成器
	middlewareChain *gtkhttp.Chain  // 中间件链
}

// MWClientOption 带中间件的客户端选项
type MWClientOption func(c *clientOption)

// clientOption 客户端选项
type clientOption struct {
	middlewares []gtkhttp.Middleware
}

// NewClient 创建一个带中间件的客户端
func NewClient(opts ...MWClientOption) (client *MWClient, err error) {
	// 创建一个分布式唯一ID生成器
	var flakeInstance *gtkflake.Flake
	if flakeInstance, err = gtkflake.New(gtkflake.Settings{}); err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), ErrFailedToCreateFlakeInstance)
		return
	}
	// 处理选项
	cliOpt := &clientOption{}
	for _, opt := range opts {
		opt(cliOpt)
	}
	// 按优先级排序中间件
	sort.Slice(cliOpt.middlewares, func(i, j int) bool {
		return cliOpt.middlewares[i].Priority() < cliOpt.middlewares[j].Priority()
	})
	// 创建中间件链
	middlewareChain := gtkhttp.NewChain(cliOpt.middlewares...)
	// 创建带中间件的客户端
	client = &MWClient{
		flakeInstance:   flakeInstance,
		middlewareChain: middlewareChain,
	}
	return
}

// ListModels 列出模型
func (c *MWClient) ListModels(ctx context.Context, request models.ListModelsRequest, opts ...gtkhttp.HTTPClientOption) (response models.ListModelsResponse, err error) {
	// 定义处理函数
	handler := func(ctx context.Context, req any) (resp any, err error) {
		// 列出模型
		return deepseek.ListModels(ctx, opts...)
	}
	// 处理请求
	var resp any
	if resp, err = c.handlerRequest(ctx, request.User, "ListModels", nil, handler); err != nil {
		return
	}
	// 返回结果
	response = resp.(models.ListModelsResponse)
	return
}

// handlerRequest 处理请求
func (c *MWClient) handlerRequest(
	ctx context.Context,
	user string,
	method string,
	request any,
	handler func(ctx context.Context, req any) (resp any, err error),
) (resp any, err error) {
	// 生成唯一请求ID
	var requestId string
	if requestId, err = c.flakeInstance.RequestID(); err != nil {
		err = &gtkhttp.ClientError{RequestID: requestId, Err: err}
		return
	}
	// 设置请求信息到上下文
	ctx = gtkhttp.SetRequestInfo(ctx, &gtkhttp.RequestInfo{
		Method:    method,
		StartTime: time.Now(),
		RequestID: requestId,
		User:      user,
	})
	// 定义最终处理函数
	finalHandler := func(ctx context.Context, req any) (resp any, err error) {
		// 执行具体的处理逻辑
		return handler(ctx, req)
	}
	// 执行中间件链
	if resp, err = c.middlewareChain.Execute(ctx, request, finalHandler); err != nil {
		err = &gtkhttp.ClientError{RequestID: requestId, Err: err}
		return
	}
	return
}
