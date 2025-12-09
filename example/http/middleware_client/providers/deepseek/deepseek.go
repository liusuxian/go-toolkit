/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-10 13:57:27
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-09 15:21:19
 * @Description: DeepSeek服务提供商实现
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
	"strings"
)

const (
	baseURL       = "https://api.deepseek.com"
	apiListModels = "/models"
)

// DeepseekProvider DeepSeek提供商
type DeepseekProvider struct {
	lb *gtkhttp.LoadBalancer // 负载均衡器
}

// NewDeepseekProvider 创建DeepSeek提供商
func NewDeepseekProvider() (provider *DeepseekProvider, err error) {
	if err = godotenv.Load(".env"); err != nil {
		return
	}
	provider = &DeepseekProvider{
		lb: gtkhttp.NewLoadBalancer(strings.Split(gtkenv.Get("DEEPSEEK_API_KEYS"), ",")),
	}
	return
}

// ListModels 列出模型
func (s *DeepseekProvider) ListModels(ctx context.Context, opts ...gtkhttp.HTTPClientOption) (response models.ListModelsResponse, err error) {
	err = common.ExecuteRequest(ctx, &common.ExecuteRequestContext{
		Method:   http.MethodGet,
		BaseURL:  baseURL,
		ApiPath:  apiListModels,
		Opts:     opts,
		LB:       s.lb,
		Response: &response,
	})
	return
}
