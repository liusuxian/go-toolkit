/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-04 11:56:13
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-24 16:27:33
 * @Description: 重试中间件
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"time"
)

// RetryStrategy 重试策略
type RetryStrategy string

const (
	RetryStrategyFixed       RetryStrategy = "fixed"       // 固定间隔
	RetryStrategyLinear      RetryStrategy = "linear"      // 线性递增
	RetryStrategyExponential RetryStrategy = "exponential" // 指数退避
	RetryStrategyJitter      RetryStrategy = "jitter"      // 带抖动的指数退避
)

// RetryCondition 重试条件函数
type RetryCondition func(attempt int, err error) (ok bool)

// RetryCallback 重试失败回调函数
type RetryCallback func(ctx context.Context, requestInfo *RequestInfo)

// RetryMiddlewareConfig 重试中间件配置
type RetryMiddlewareConfig struct {
	MaxAttempts   int            // 最大重试次数（-1表示无限重试，0表示不重试只执行一次，>0表示重试指定次数），默认不重试
	Strategy      RetryStrategy  // 重试策略，可选值: fixed 固定间隔，linear 线性递增，exponential 指数退避，jitter 带抖动的指数退避，默认 exponential
	BaseDelay     time.Duration  // 基础延迟时间，默认 1s
	MaxDelay      time.Duration  // 最大延迟时间，默认 10s
	Multiplier    float64        // 重试间隔倍数（用于指数退避），默认 2.0
	JitterPercent float64        // 抖动百分比（用于抖动策略，范围0-1，如0.1表示±10%），默认 0.1
	Condition     RetryCondition // 重试条件
	OnRetry       RetryCallback  // 重试失败回调函数（同步执行会阻塞重试流程，建议仅用于轻量级操作如日志记录、监控上报等，耗时操作请在回调内使用goroutine异步处理）
}

// RetryMiddleware 重试中间件
type RetryMiddleware struct {
	config RetryMiddlewareConfig
}

// NewRetryMiddleware 创建重试中间件
func NewRetryMiddleware(config RetryMiddlewareConfig) (retry *RetryMiddleware) {
	// 设置最大重试次数（-1表示无限重试，0表示不重试只执行一次，>0表示重试指定次数），默认不重试
	// 只有小于-1的值才会被设置为3
	if config.MaxAttempts < -1 {
		config.MaxAttempts = 3
	}
	// 设置重试策略，可选值: fixed 固定间隔，linear 线性递增，exponential 指数退避，jitter 带抖动的指数退避，默认 exponential
	if config.Strategy == "" {
		config.Strategy = RetryStrategyExponential
	}
	// 设置基础延迟时间，默认 1s
	if config.BaseDelay <= 0 {
		config.BaseDelay = 1 * time.Second
	}
	// 设置最大延迟时间，默认 10s
	if config.MaxDelay <= 0 {
		config.MaxDelay = 10 * time.Second
	}
	// 确保最大延迟不小于基础延迟
	if config.MaxDelay < config.BaseDelay {
		config.MaxDelay = config.BaseDelay * 10 // 设置为基础延迟的10倍
	}
	// 设置重试间隔倍数（用于指数退避），默认 2.0
	if config.Multiplier <= 0 {
		config.Multiplier = 2.0
	}
	// 设置抖动百分比（用于抖动策略，范围0-1，如0.1表示±10%），默认 0.1
	if config.JitterPercent <= 0 || config.JitterPercent > 1 {
		config.JitterPercent = 0.1 // 默认±10%抖动
	}
	// 设置重试条件
	if config.Condition == nil {
		config.Condition = DefaultRetryCondition
	}
	return &RetryMiddleware{
		config: config,
	}
}

// Process 处理请求
func (m *RetryMiddleware) Process(ctx context.Context, request any, next MWHandler) (response any, err error) {
	// 从上下文中获取请求信息
	requestInfo := GetRequestInfo(ctx)
	requestInfo.MaxAttempts = m.config.MaxAttempts

	attempt := 0
	for {
		requestInfo.Attempt = attempt
		// 执行请求
		response, err = next(ctx, request)
		// 如果成功或者不需要重试，直接返回
		if err == nil || !m.config.Condition(attempt, err) {
			return
		}
		// 如果重试回调不为空，则执行重试回调
		if m.config.OnRetry != nil && attempt > 0 {
			// 深度拷贝 RequestInfo，避免回调函数修改原始数据
			m.config.OnRetry(ctx, deepCopyRequestInfo(requestInfo))
		}
		// MaxAttempts = 0: 只执行一次，不重试
		if m.config.MaxAttempts == 0 {
			break
		}
		// MaxAttempts > 0: 已达到最大重试次数
		if m.config.MaxAttempts > 0 && attempt >= m.config.MaxAttempts {
			break
		}
		// MaxAttempts = -1: 无限重试，继续执行
		// 计算延迟时间
		delay := m.calculateDelay(attempt + 1)
		// 等待延迟时间
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
			// 继续下次重试
		}
		attempt++
	}
	// 返回最后一次的错误
	return nil, fmt.Errorf("after %d attempts, last error: %w", requestInfo.Attempt, err)
}

// Name 返回中间件名称
func (m *RetryMiddleware) Name() (name string) {
	return "retry"
}

// Priority 返回中间件优先级
func (m *RetryMiddleware) Priority() (priority int) {
	return 20 // 重试中间件优先级较高，在监控之后执行
}

// calculateDelay 计算延迟时间
func (m *RetryMiddleware) calculateDelay(attempt int) (delay time.Duration) {
	switch m.config.Strategy {
	case RetryStrategyFixed:
		delay = m.config.BaseDelay
	case RetryStrategyLinear:
		if attempt > 0 && m.config.BaseDelay > 0 {
			// 防止整数溢出，先检查是否会超过最大延迟
			maxAttempts := int64(m.config.MaxDelay / m.config.BaseDelay)
			if int64(attempt) > maxAttempts {
				delay = m.config.MaxDelay
			} else {
				delay = min(time.Duration(attempt)*m.config.BaseDelay, m.config.MaxDelay)
			}
		} else {
			delay = m.config.BaseDelay
		}
	case RetryStrategyExponential:
		delay = m.calculateExponentialDelay(attempt)
	case RetryStrategyJitter:
		exponentialDelay := m.calculateExponentialDelay(attempt)
		// 添加双向随机抖动（基于配置的百分比）
		jitterRange := float64(exponentialDelay) * m.config.JitterPercent
		jitter := time.Duration((rand.Float64() - 0.5) * 2 * jitterRange)
		// 确保抖动后的延迟在合理范围内
		delay = max(exponentialDelay+jitter, exponentialDelay/2)
	default:
		delay = m.config.BaseDelay
	}
	// 确保不超过最大延迟
	delay = min(delay, m.config.MaxDelay)
	// 确保延迟不为负数或过小
	delay = max(delay, m.config.BaseDelay/2)
	return
}

// calculateExponentialDelay 计算指数退避延迟
func (m *RetryMiddleware) calculateExponentialDelay(attempt int) (delay time.Duration) {
	// 如果基础延迟、最大延迟小于等于0或重试间隔倍数小于等于1，返回基础延迟
	if m.config.BaseDelay <= 0 || m.config.MaxDelay <= 0 || m.config.Multiplier <= 1 {
		return m.config.BaseDelay
	}
	// 计算最大尝试次数
	maxExponent := math.Log(float64(m.config.MaxDelay/m.config.BaseDelay)) / math.Log(m.config.Multiplier)
	// 如果最大尝试次数为无穷大或NaN，或者尝试次数大于最大尝试次数，返回最大延迟
	if math.IsInf(maxExponent, 0) || math.IsNaN(maxExponent) || float64(attempt-1) > maxExponent {
		return m.config.MaxDelay
	}
	// 计算延迟时间
	delayFloat := float64(m.config.BaseDelay) * math.Pow(m.config.Multiplier, float64(attempt-1))
	// 如果延迟时间为无穷大或NaN，返回最大延迟
	if math.IsInf(delayFloat, 0) || math.IsNaN(delayFloat) {
		return m.config.MaxDelay
	}
	// 返回延迟时间
	return time.Duration(delayFloat)
}

// DefaultRetryCondition 默认重试条件
func DefaultRetryCondition(attempt int, err error) (ok bool) {
	// 网络错误通常可以重试
	if IsNetError(err) {
		return true
	}
	// HTTP状态码错误
	if isRetryableHTTPError(err) {
		return true
	}
	return false
}

// isRetryableHTTPError 判断是否为可重试的HTTP错误
func isRetryableHTTPError(err error) (ok bool) {
	// 检查是否为API错误
	var apiError *APIError
	if errors.As(err, &apiError) {
		return slices.Contains(retryableHTTPStatusCodes, apiError.HTTPStatusCode)
	}
	// 检查是否为请求错误
	var requestError *RequestError
	if errors.As(err, &requestError) {
		return slices.Contains(retryableHTTPStatusCodes, requestError.HTTPStatusCode)
	}
	return false
}

// 可重试的HTTP状态码列表
var retryableHTTPStatusCodes = []int{
	// 4xx 客户端错误中可重试的状态码
	407, // Proxy Authentication Required - 代理认证需要，可能是临时的
	408, // Request Timeout - 请求超时，网络问题导致
	409, // Conflict - 资源冲突，可能是并发导致的临时问题
	421, // Misdirected Request - 请求被错误定向，可能是负载均衡问题
	423, // Locked - 资源被锁定，可能是临时锁定
	424, // Failed Dependency - 依赖操作失败，依赖可能临时不可用
	425, // Too Early - 请求过早发送，服务器要求稍后重试
	429, // Too Many Requests - 请求频率过高，需要限流重试
	449, // Retry With - 微软扩展，明确要求客户端重试
	// 5xx 服务器错误中可重试的状态码
	500, // Internal Server Error - 服务器内部错误，可能是临时故障
	502, // Bad Gateway - 网关错误，上游服务器问题
	503, // Service Unavailable - 服务不可用，通常是临时维护或过载
	504, // Gateway Timeout - 网关超时，上游服务器响应太慢
	507, // Insufficient Storage - 存储空间不足，可能是临时问题
	508, // Loop Detected - 检测到无限循环，可能是配置临时问题
	509, // Bandwidth Limit Exceeded - 带宽限制超出，可以稍后重试
	510, // Not Extended - 服务器需要扩展请求，可能是临时配置问题
	511, // Network Authentication Required - 网络认证需要，可能是临时的
	// Cloudflare 扩展状态码 (网络基础设施相关，通常是临时问题)
	520, // Web Server Returned an Unknown Error - 源服务器返回未知错误
	521, // Web Server Is Down - 源服务器宕机
	522, // Connection Timed Out - 连接超时
	523, // Origin Is Unreachable - 源服务器不可达
	524, // A Timeout Occurred - 发生超时
	525, // SSL Handshake Failed - SSL握手失败
	526, // Invalid SSL Certificate - SSL证书无效
	527, // Railgun Error - Railgun连接错误
	// 其他CDN/代理服务商扩展状态码
	530, // Site Frozen - 站点被冻结，可能是临时的
	598, // Network Read Timeout Error - 网络读取超时 (非标准)
	599, // Network Connect Timeout Error - 网络连接超时 (非标准)
}

// RetryConditions 预定义的重试条件
var RetryConditions = struct {
	// Never 永不重试
	Never RetryCondition
	// Always 总是重试（除了成功的情况）
	Always RetryCondition
	// NetworkOnly 仅网络错误重试
	NetworkOnly RetryCondition
	// HTTPOnly 仅HTTP错误重试
	HTTPOnly RetryCondition
}{
	Never: func(attempt int, err error) (ok bool) {
		return false
	},
	Always: func(attempt int, err error) (ok bool) {
		return true
	},
	NetworkOnly: func(attempt int, err error) (ok bool) {
		return IsNetError(err)
	},
	HTTPOnly: func(attempt int, err error) (ok bool) {
		return isRetryableHTTPError(err)
	},
}
