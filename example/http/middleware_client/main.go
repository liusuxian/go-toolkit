/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-12-09 00:57:20
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-11 11:32:36
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package main

import (
	"context"
	"fmt"
	"github.com/liusuxian/go-toolkit/example/http/middleware_client/models"
	"github.com/liusuxian/go-toolkit/example/http/middleware_client/providers/deepseek"
	"github.com/liusuxian/go-toolkit/gtkflake"
	"github.com/liusuxian/go-toolkit/gtkhttp"
	"github.com/liusuxian/go-toolkit/gtkjson"
	"time"
)

func isError(err error) {
	if err != nil {
		originalErr := gtkhttp.Unwrap(err)
		fmt.Println("originalErr =", originalErr)
		fmt.Println("Cause Error =", gtkhttp.Cause(err))
		switch {
		case gtkhttp.IsMethodNotSupportedError(originalErr):
			fmt.Println("IsMethodNotSupportedError =", true)
		case gtkhttp.IsCompletionStreamNotSupportedError(originalErr):
			fmt.Println("IsCompletionStreamNotSupportedError =", true)
		case gtkhttp.IsTooManyEmptyStreamMessagesError(originalErr):
			fmt.Println("IsTooManyEmptyStreamMessagesError =", true)
		case gtkhttp.IsStreamReturnIntervalTimeoutError(originalErr):
			fmt.Println("IsStreamReturnIntervalTimeoutError =", true)
		case gtkhttp.IsCanceledError(originalErr):
			fmt.Println("IsCanceledError =", true)
		case gtkhttp.IsDeadlineExceededError(originalErr):
			fmt.Println("IsDeadlineExceededError =", true)
		case gtkhttp.IsNetError(originalErr):
			fmt.Println("IsNetError =", true)
		default:
			fmt.Println("unknown error =", err)
		}
	}
}

// ListModels 列出模型
func ListModels(ctx context.Context, client *gtkhttp.MWClient, request models.ListModelsRequest, opts ...gtkhttp.HTTPClientOption) (response models.ListModelsResponse, err error) {
	// 定义处理函数
	handler := func(ctx context.Context, req any) (resp any, err error) {
		// 列出模型
		var (
			deepseekProvider *deepseek.DeepseekProvider
			e                error
		)
		if deepseekProvider, e = deepseek.NewDeepseekProvider(); e != nil {
			return nil, e
		}
		return deepseekProvider.ListModels(ctx, opts...)
	}
	// 处理请求
	var resp any
	if resp, err = client.HandlerRequest(ctx, request.User, "ListModels", nil, handler); err != nil {
		return
	}
	// 返回结果
	response = resp.(models.ListModelsResponse)
	return
}

func main() {
	// 创建一个分布式唯一ID生成器
	var (
		flake *gtkflake.Flake
		err   error
	)
	if flake, err = gtkflake.New(gtkflake.Settings{}); err != nil {
		fmt.Printf("new flake error = %v\n", err)
		return
	}
	client := gtkhttp.NewMWClient(
		gtkhttp.WithLogging(gtkhttp.LoggingMiddlewareConfig{
			LogRequest: true,
			LogError:   true,
		}),
		gtkhttp.WithRetry(gtkhttp.RetryMiddlewareConfig{
			MaxAttempts:   3,
			Strategy:      gtkhttp.RetryStrategyExponential,
			BaseDelay:     1 * time.Second,
			MaxDelay:      10 * time.Second,
			Multiplier:    2.0,
			JitterPercent: 0.1,
			Condition: func(attempt int, err error) (ok bool) {
				return true
			},
			OnRetry: func(ctx context.Context, requestInfo *gtkhttp.RequestInfo) {
				fmt.Printf("onRetry = %s\n", gtkjson.MustString(requestInfo))
			},
		}),
		gtkhttp.WithMetrics(gtkhttp.MetricsMiddlewareConfig{}),
		gtkhttp.WithRequestIDGenerator(flake),
	)
	defer func() {
		metrics := client.GetMetrics()
		fmt.Printf("metrics = %s\n", gtkjson.MustString(metrics))
	}()
	// 列出模型
	var (
		ctx      = context.Background()
		response models.ListModelsResponse
	)
	if response, err = ListModels(ctx, client, models.ListModelsRequest{
		User: "test",
	}, gtkhttp.WithTimeout(time.Minute*2)); err != nil {
		isError(err)
		fmt.Printf("listModels error = %v, request_id = %s\n", err, gtkhttp.RequestID(err))
		return
	}
	fmt.Printf("listModels response = %s, request_id = %s\n", gtkjson.MustString(response), response.RequestID())
}
