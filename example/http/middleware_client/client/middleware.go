/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-02 04:49:05
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-09 00:59:49
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package client

import "github.com/liusuxian/go-toolkit/gtkhttp"

// WithMiddleware 添加中间件
func WithMiddleware(m gtkhttp.Middleware) (opt MWClientOption) {
	return func(c *clientOption) {
		c.middlewares = append(c.middlewares, m)
	}
}

// WithLogging 添加日志中间件
func WithLogging(config gtkhttp.LoggingMiddlewareConfig) (opt MWClientOption) {
	return func(c *clientOption) {
		c.middlewares = append(c.middlewares, gtkhttp.NewLoggingMiddleware(config))
	}
}

// WithMetrics 添加监控中间件
func WithMetrics(config gtkhttp.MetricsMiddlewareConfig) (opt MWClientOption) {
	return func(c *clientOption) {
		c.middlewares = append(c.middlewares, gtkhttp.NewMetricsMiddleware(config))
	}
}

// WithRetry 添加重试中间件
func WithRetry(config gtkhttp.RetryMiddlewareConfig) (opt MWClientOption) {
	return func(c *clientOption) {
		c.middlewares = append(c.middlewares, gtkhttp.NewRetryMiddleware(config))
	}
}

// WithDefaultMiddlewares 添加默认中间件（日志、监控、重试）
func WithDefaultMiddlewares() (opt MWClientOption) {
	return func(c *clientOption) {
		c.middlewares = append(c.middlewares,
			gtkhttp.NewMetricsMiddleware(gtkhttp.DefaultMetricsConfig()),
			gtkhttp.NewRetryMiddleware(gtkhttp.DefaultRetryConfig()),
			gtkhttp.NewLoggingMiddleware(gtkhttp.DefaultLoggingConfig()),
		)
	}
}

// GetMetrics 获取指标数据（如果启用了监控中间件）
func (c *MWClient) GetMetrics() (metrics map[string]any) {
	for _, mw := range c.middlewareChain.GetMiddlewares() {
		if metricsMiddleware, ok := mw.(*gtkhttp.MetricsMiddleware); ok {
			return metricsMiddleware.GetMetrics()
		}
	}
	return
}
