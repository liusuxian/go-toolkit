/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-17 18:24:31
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-09 14:31:41
 * @Description: 监控中间件
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import (
	"context"
	"errors"
	"maps"
	"net"
	"slices"
	"strings"
	"sync"
	"time"
)

const (
	maxResponseTimeRecords = 1000 // 每个key最多保存1000条响应时间记录
)

// MetricsCollector 指标收集器接口
type MetricsCollector interface {
	// 记录请求开始
	RecordRequestStart(method string)
	// 记录请求完成
	RecordRequestComplete(method string, durationMs int64, success bool)
	// 记录错误
	RecordError(method, errorType string)
	// 记录重试
	RecordRetry(method string, retryCount int)
	// 获取指标数据
	GetMetrics() (metrics map[string]any)
	// 重置指标
	Reset()
}

// DefaultMetricsCollector 默认指标收集器
type DefaultMetricsCollector struct {
	mu sync.RWMutex
	// 请求计数器
	totalRequests   map[string]int64 // 总请求数
	successRequests map[string]int64 // 成功请求数
	failedRequests  map[string]int64 // 失败请求数
	// 响应时间统计
	responseTimes map[string][]int64 // 响应时间列表（毫秒）
	// 错误统计
	errorCounts map[string]int64 // 错误计数
	// 重试统计
	retryCounts map[string]int64 // 重试计数
	// 活跃请求数
	activeRequests map[string]int64 // 当前活跃请求数
	// 时间范围内的统计
	startTime time.Time // 统计开始时间
}

// NewDefaultMetricsCollector 创建默认指标收集器
func NewDefaultMetricsCollector() (metricsCollector *DefaultMetricsCollector) {
	return &DefaultMetricsCollector{
		totalRequests:   make(map[string]int64),
		successRequests: make(map[string]int64),
		failedRequests:  make(map[string]int64),
		responseTimes:   make(map[string][]int64),
		errorCounts:     make(map[string]int64),
		retryCounts:     make(map[string]int64),
		activeRequests:  make(map[string]int64),
		startTime:       time.Now(),
	}
}

// RecordRequestStart 记录请求开始
func (c *DefaultMetricsCollector) RecordRequestStart(method string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.getKey(method)
	c.totalRequests[key]++
	c.activeRequests[key]++
}

// RecordRequestComplete 记录请求完成
func (c *DefaultMetricsCollector) RecordRequestComplete(method string, durationMs int64, success bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.getKey(method)
	if _, exists := c.responseTimes[key]; !exists {
		c.responseTimes[key] = make([]int64, 0, maxResponseTimeRecords)
	}
	// 如果已达到最大记录数，移除最旧的记录（使用更高效的方式）
	if len(c.responseTimes[key]) >= maxResponseTimeRecords {
		// 保留后面80%的记录，丢弃前面20%
		keepFromIndex := maxResponseTimeRecords / 5 // 20%的位置
		c.responseTimes[key] = slices.Delete(c.responseTimes[key], 0, keepFromIndex)
	}
	c.responseTimes[key] = append(c.responseTimes[key], durationMs)
	// 记录成功/失败
	if success {
		c.successRequests[key]++
	} else {
		c.failedRequests[key]++
	}
	// 减少活跃请求数
	if c.activeRequests[key] > 0 {
		c.activeRequests[key]--
	}
}

// RecordError 记录错误
func (c *DefaultMetricsCollector) RecordError(method string, errorType string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.getKey(method) + ":" + errorType
	c.errorCounts[key]++
}

// RecordRetry 记录重试
func (c *DefaultMetricsCollector) RecordRetry(method string, retryCount int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.getKey(method)
	c.retryCounts[key] += int64(retryCount)
}

// GetMetrics 获取指标数据
func (c *DefaultMetricsCollector) GetMetrics() (metrics map[string]any) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	metrics = make(map[string]any)
	// 深拷贝总请求数数据
	totalRequests := make(map[string]int64)
	maps.Copy(totalRequests, c.totalRequests)
	metrics["total_requests"] = totalRequests
	// 深拷贝成功请求数数据
	successRequests := make(map[string]int64)
	maps.Copy(successRequests, c.successRequests)
	metrics["success_requests"] = successRequests
	// 深拷贝失败请求数数据
	failedRequests := make(map[string]int64)
	maps.Copy(failedRequests, c.failedRequests)
	metrics["failed_requests"] = failedRequests
	// 深拷贝错误计数数据
	errorCounts := make(map[string]int64)
	maps.Copy(errorCounts, c.errorCounts)
	metrics["error_counts"] = errorCounts
	// 深拷贝重试计数数据
	retryCounts := make(map[string]int64)
	maps.Copy(retryCounts, c.retryCounts)
	metrics["retry_counts"] = retryCounts
	// 深拷贝活跃请求数数据
	activeRequests := make(map[string]int64)
	maps.Copy(activeRequests, c.activeRequests)
	metrics["active_requests"] = activeRequests
	// 统计开始时间
	metrics["start_time"] = c.startTime
	// 计算成功率和平均响应时间
	successRates := make(map[string]float64)
	avgResponseTimes := make(map[string]float64)
	for key, total := range totalRequests {
		if total > 0 {
			success := successRequests[key]
			successRates[key] = float64(success/total) * 100
		}
		// 计算平均响应时间
		if times, exists := c.responseTimes[key]; exists && len(times) > 0 {
			var sum int64
			for _, t := range times {
				sum += t
			}
			avgResponseTimes[key] = float64(sum) / float64(len(times))
		}
	}
	// 添加成功率、平均响应时间和运行时间
	metrics["success_rates"] = successRates
	metrics["avg_response_times"] = avgResponseTimes
	metrics["uptime_seconds"] = time.Since(c.startTime).Seconds()
	return
}

// Reset 重置指标数据
func (c *DefaultMetricsCollector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.totalRequests = make(map[string]int64)
	c.successRequests = make(map[string]int64)
	c.failedRequests = make(map[string]int64)
	c.responseTimes = make(map[string][]int64)
	c.errorCounts = make(map[string]int64)
	c.retryCounts = make(map[string]int64)
	c.activeRequests = make(map[string]int64)
	c.startTime = time.Now()
}

// getKey 获取指标键值
func (c *DefaultMetricsCollector) getKey(subKeys ...string) (key string) {
	subKeyList := make([]string, 0, len(subKeys))
	for _, subKey := range subKeys {
		if subKey != "" {
			subKeyList = append(subKeyList, subKey)
		}
	}

	if len(subKeyList) == 0 {
		return "unknown"
	}
	return strings.Join(subKeyList, ":")
}

// MetricsMiddlewareConfig 监控中间件配置
type MetricsMiddlewareConfig struct {
	Collector MetricsCollector // 指标收集器
}

// MetricsMiddleware 监控中间件
type MetricsMiddleware struct {
	config MetricsMiddlewareConfig
}

// NewMetricsMiddleware 创建监控中间件
func NewMetricsMiddleware(config MetricsMiddlewareConfig) (mm *MetricsMiddleware) {
	if config.Collector == nil {
		config.Collector = NewDefaultMetricsCollector()
	}

	return &MetricsMiddleware{
		config: config,
	}
}

// Process 处理请求
func (m *MetricsMiddleware) Process(ctx context.Context, request any, next MWHandler) (response any, err error) {
	// 从上下文中获取请求信息
	requestInfo := GetRequestInfo(ctx)
	// 记录请求开始
	m.config.Collector.RecordRequestStart(requestInfo.Method)
	// 执行下一个处理器
	response, err = next(ctx, request)
	// 记录请求完成
	m.config.Collector.RecordRequestComplete(requestInfo.Method, requestInfo.TotalDurationMs, requestInfo.IsSuccess)
	// 记录错误
	if err != nil {
		errorType := m.classifyError(err)
		m.config.Collector.RecordError(requestInfo.Method, errorType)
	}
	// 记录重试次数
	if requestInfo.Attempt > 0 {
		m.config.Collector.RecordRetry(requestInfo.Method, requestInfo.Attempt)
	}
	return
}

// Name 返回中间件名称
func (m *MetricsMiddleware) Name() (name string) {
	return "metrics"
}

// Priority 返回中间件优先级
func (m *MetricsMiddleware) Priority() (priority int) {
	return 10 // 监控中间件优先级较高，尽早执行
}

// GetMetrics 获取指标数据
func (m *MetricsMiddleware) GetMetrics() (metrics map[string]any) {
	return m.config.Collector.GetMetrics()
}

// classifyError 分类错误类型
func (m *MetricsMiddleware) classifyError(err error) (errorType string) {
	if err == nil {
		return "none"
	}
	var netErr net.Error
	if errors.As(err, &netErr) {
		return "net_err"
	}
	// 检查是否为API错误
	var apiError *APIError
	if errors.As(err, &apiError) {
		return "api_error"
	}
	// 检查是否为请求错误
	var requestError *RequestError
	if errors.As(err, &requestError) {
		return "request_error"
	}
	// 其他未知错误
	return "unknown"
}
