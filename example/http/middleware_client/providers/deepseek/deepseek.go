/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-10 13:57:27
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-09 01:33:04
 * @Description: DeepSeek服务提供商实现，采用单例模式，在包导入时自动注册到提供商工厂
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package deepseek

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/liusuxian/go-toolkit/example/http/middleware_client/common"
	"github.com/liusuxian/go-toolkit/example/http/middleware_client/models"
	"github.com/liusuxian/go-toolkit/gtkenv"
	"github.com/liusuxian/go-toolkit/gtkhttp"
	"net/http"
)

const (
	baseURL   = "https://api.deepseek.com"
	apiModels = "/models"
)

// ListModels 列出模型
func ListModels(ctx context.Context, opts ...gtkhttp.HTTPClientOption) (response models.ListModelsResponse, err error) {
	if err = godotenv.Load(".env"); err != nil {
		return
	}
	err = common.ExecuteRequest(ctx, &common.ExecuteRequestContext{
		Method:   http.MethodGet,
		BaseURL:  baseURL,
		ApiPath:  apiModels,
		ApiKey:   gtkenv.Get("TEST_DEEPSEEK_API_KEY"),
		Opts:     opts,
		Response: &response,
	})
	return
}
