/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-12-09 13:44:59
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-09 15:42:00
 * @Description: 带中间件的客户端
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"
)

// IRequestIDGenerator 唯一请求ID生成器接口
type IRequestIDGenerator interface {
	RequestID() (requestId string, err error) // 生成唯一请求ID
}

// MWClient 带中间件的客户端
type MWClient struct {
	requestIDGenerator IRequestIDGenerator // 唯一请求ID生成器
	middlewareChain    *Chain              // 中间件链
}

// MWClientOption 带中间件的客户端选项
type MWClientOption func(c *MWClientConfig)

// MWClientConfig 客户端配置
type MWClientConfig struct {
	requestIDGenerator IRequestIDGenerator // 唯一请求ID生成器
	middlewares        []Middleware
}

// WithRequestIDGenerator 设置唯一请求ID生成器
func WithRequestIDGenerator(generator IRequestIDGenerator) (opt MWClientOption) {
	return func(c *MWClientConfig) {
		c.requestIDGenerator = generator
	}
}

// WithMiddleware 添加中间件
func WithMiddleware(m Middleware) (opt MWClientOption) {
	return func(c *MWClientConfig) {
		c.middlewares = append(c.middlewares, m)
	}
}

// WithLogging 添加日志中间件
func WithLogging(config LoggingMiddlewareConfig) (opt MWClientOption) {
	return func(c *MWClientConfig) {
		c.middlewares = append(c.middlewares, NewLoggingMiddleware(config))
	}
}

// WithRetry 添加重试中间件
func WithRetry(config RetryMiddlewareConfig) (opt MWClientOption) {
	return func(c *MWClientConfig) {
		c.middlewares = append(c.middlewares, NewRetryMiddleware(config))
	}
}

// WithMetrics 添加监控中间件
func WithMetrics(config MetricsMiddlewareConfig) (opt MWClientOption) {
	return func(c *MWClientConfig) {
		c.middlewares = append(c.middlewares, NewMetricsMiddleware(config))
	}
}

// NewMWClient 创建一个带中间件的客户端
func NewMWClient(opts ...MWClientOption) (client *MWClient) {
	// 处理选项
	cfg := &MWClientConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	// 按优先级排序中间件
	sort.Slice(cfg.middlewares, func(i, j int) bool {
		return cfg.middlewares[i].Priority() < cfg.middlewares[j].Priority()
	})
	// 创建中间件链
	middlewareChain := NewChain(cfg.middlewares...)
	// 创建带中间件的客户端
	client = &MWClient{
		requestIDGenerator: cfg.requestIDGenerator,
		middlewareChain:    middlewareChain,
	}
	return
}

// HandlerRequest 处理请求
func (c *MWClient) HandlerRequest(
	ctx context.Context,
	user string, // 代表你的终端用户的唯一标识符
	method string, // 请求方法名称
	request any, // 请求参数
	handler func(ctx context.Context, req any) (resp any, err error), // 处理函数
) (resp any, err error) {
	if c.requestIDGenerator == nil {
		err = errors.New("requestIDGenerator is not set")
		return
	}
	// 生成唯一请求ID
	var requestId string
	if requestId, err = c.requestIDGenerator.RequestID(); err != nil {
		err = &MWClientError{RequestID: requestId, Err: err}
		return
	}
	// 设置请求信息到上下文
	ctx = SetRequestInfo(ctx, &RequestInfo{
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
		err = &MWClientError{RequestID: requestId, Err: err}
		return
	}
	return
}

// GetMetrics 获取指标数据（如果启用了监控中间件）
func (c *MWClient) GetMetrics() (metrics map[string]any) {
	for _, mw := range c.middlewareChain.GetMiddlewares() {
		if metricsMiddleware, ok := mw.(*MetricsMiddleware); ok {
			return metricsMiddleware.GetMetrics()
		}
	}
	return
}

// MWClientError 客户端错误
type MWClientError struct {
	RequestID string // 请求ID
	Err       error  // 原始错误
}

// Error 错误信息
func (e *MWClientError) Error() (errStr string) {
	return fmt.Sprintf("request_id: %s, error: %v", e.RequestID, e.Err)
}

// RequestID 获取请求ID
func RequestID(err error) (requestId string) {
	if err == nil {
		return ""
	}

	var clientErr *MWClientError
	if errors.As(err, &clientErr) {
		return clientErr.RequestID
	}

	return ""
}

// Unwrap 解包错误
func Unwrap(err error) (originalError error) {
	if err == nil {
		return nil
	}
	// 解包 MWClientError
	var clientErr *MWClientError
	if errors.As(err, &clientErr) {
		if clientErr.Err != nil {
			return clientErr.Err
		}
		return err // 如果内部错误为 nil，返回 MWClientError 本身
	}
	// 解包 RequestError
	var requestError *RequestError
	if errors.As(err, &requestError) {
		if requestError.Err != nil {
			return requestError.Err
		}
		return err // 如果内部错误为 nil，返回 RequestError 本身
	}
	// 其他类型的错误
	unwrapped := errors.Unwrap(err)
	if unwrapped == nil {
		return err // 已经是最底层错误，返回原错误
	}
	return unwrapped
}

// Cause 错误根因
func Cause(err error) (causeError error) {
	return doCause(err)
}

// doCause 递归获取错误根因
func doCause(err error) (causeError error) {
	if err == nil {
		return nil
	}
	// 解包错误
	unwrapped := Unwrap(err)
	if unwrapped == nil {
		return err // 已经到达最底层错误，返回当前错误
	}
	// 防止无限递归：如果解包后的错误与原错误相同，直接返回
	if unwrapped == err {
		return err
	}
	return doCause(unwrapped)
}
