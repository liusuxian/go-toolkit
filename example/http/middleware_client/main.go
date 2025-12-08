/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-12-09 00:57:20
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-09 01:19:29
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package main

import (
	"context"
	"fmt"
	"github.com/liusuxian/go-toolkit/example/http/middleware_client/client"
	"github.com/liusuxian/go-toolkit/example/http/middleware_client/models"
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

func main() {
	client, err := client.NewClient(client.WithDefaultMiddlewares())
	if err != nil {
		fmt.Printf("NewSDKClient() error = %v\n", err)
		return
	}
	defer func() {
		metrics := client.GetMetrics()
		fmt.Printf("metrics = %s\n", gtkjson.MustString(metrics))
	}()

	ctx := context.Background()
	// 列出模型
	response1, err := client.ListModels(ctx, models.ListModelsRequest{
		User: "test",
	}, gtkhttp.WithTimeout(time.Minute*2))
	isError(err)
	if err != nil {
		fmt.Printf("listModels error = %v, request_id = %s\n", err, gtkhttp.RequestID(err))
		return
	}
	fmt.Printf("listModels response = %s, request_id = %s\n", gtkjson.MustString(response1), response1.RequestID())
}
