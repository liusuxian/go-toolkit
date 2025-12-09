/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-30 15:14:39
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-09 16:19:36
 * @Description: 日志中间件
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/liusuxian/go-toolkit/gtklog"
	"strings"
	"time"
)

const (
	maxSanitizeDepth = 10 // 最大脱敏递归深度
)

// LoggingMiddlewareConfig 日志中间件配置
type LoggingMiddlewareConfig struct {
	Logger          gtklog.ILogger // 日志器
	LogRequest      bool           // 是否记录请求
	LogResponse     bool           // 是否记录响应
	LogError        bool           // 是否记录错误
	SkipSuccessLog  bool           // 是否跳过成功请求的日志
	SensitiveFields []string       // 敏感字段，会被脱敏
}

// LoggingMiddleware 日志中间件
type LoggingMiddleware struct {
	config LoggingMiddlewareConfig
}

// NewLoggingMiddleware 创建日志中间件
func NewLoggingMiddleware(config LoggingMiddlewareConfig) (lm *LoggingMiddleware) {
	if config.Logger == nil {
		config.Logger = gtklog.NewDefaultLogger(gtklog.InfoLevel)
		config.LogRequest = true
		config.LogError = true
	}

	return &LoggingMiddleware{
		config: config,
	}
}

// Process 处理请求
func (m *LoggingMiddleware) Process(ctx context.Context, request any, next MWHandler) (response any, err error) {
	// 从上下文中获取请求信息
	requestInfo := GetRequestInfo(ctx)
	// 记录请求开始日志
	m.logRequestStart(ctx, request, requestInfo)
	// 执行下一个处理器
	processingStartTime := time.Now()
	response, err = next(ctx, request)
	// 更新请求信息
	requestInfo.EndTime = time.Now()
	requestInfo.TotalDurationMs = requestInfo.EndTime.Sub(requestInfo.StartTime).Milliseconds()
	requestInfo.IsSuccess = err == nil
	requestInfo.Error = err
	// 记录请求结束日志
	m.logRequestEnd(ctx, processingStartTime, response, err, requestInfo)
	return
}

// Name 返回中间件名称
func (m *LoggingMiddleware) Name() (name string) {
	return "logging"
}

// Priority 返回中间件优先级
func (m *LoggingMiddleware) Priority() (priority int) {
	return 100 // 日志中间件优先级较低，在其他中间件执行后记录
}

// logRequestStart 记录请求开始日志
func (m *LoggingMiddleware) logRequestStart(ctx context.Context, request any, requestInfo *RequestInfo) {
	// 是否记录请求
	if m.config.LogRequest {
		// 创建一个别名结构体
		type Alias RequestInfo
		startTemp := struct {
			EndTime         *time.Time `json:"end_time,omitempty"`
			TotalDurationMs int64      `json:"total_duration_ms,omitempty"`
			IsSuccess       *bool      `json:"is_success,omitempty"`
			Error           string     `json:"error,omitempty"`
			Request         any        `json:"request"`
			Alias
		}{
			EndTime:         nil,
			TotalDurationMs: 0,
			IsSuccess:       nil,
			Error:           "",
			Request:         nil,
			Alias:           Alias(*requestInfo),
		}
		// 脱敏处理请求数据
		if reqData := m.sanitizeData(request); reqData != nil {
			startTemp.Request = reqData
		}
		m.config.Logger.Infof(ctx, "request started: %s", toMustString(startTemp))
	}
}

// logRequestEnd 记录请求结束日志
func (m *LoggingMiddleware) logRequestEnd(ctx context.Context, processingStartTime time.Time, response any, err error, requestInfo *RequestInfo) {
	if err != nil {
		// 是否记录错误
		if m.config.LogError {
			// 创建一个别名结构体
			type Alias RequestInfo
			endTemp := struct {
				DurationMs      int64  `json:"duration_ms,omitempty"`
				TotalDurationMs int64  `json:"total_duration_ms,omitempty"`
				Error           string `json:"error,omitempty"`
				Alias
			}{
				DurationMs:      requestInfo.EndTime.Sub(processingStartTime).Milliseconds(),
				TotalDurationMs: requestInfo.TotalDurationMs,
				Error:           requestInfo.Error.Error(),
				Alias:           Alias(*requestInfo),
			}
			m.config.Logger.Errorf(ctx, "request failed: %s", toMustString(endTemp))
		}
	} else {
		// 是否跳过成功请求的日志
		if !m.config.SkipSuccessLog {
			// 创建一个别名结构体
			type Alias RequestInfo
			endTemp := struct {
				DurationMs      int64  `json:"duration_ms,omitempty"`
				TotalDurationMs int64  `json:"total_duration_ms,omitempty"`
				Error           string `json:"error,omitempty"`
				Response        *any   `json:"response,omitempty"`
				Alias
			}{
				DurationMs:      requestInfo.EndTime.Sub(processingStartTime).Milliseconds(),
				TotalDurationMs: requestInfo.TotalDurationMs,
				Error:           "",
				Response:        nil,
				Alias:           Alias(*requestInfo),
			}
			// 是否记录响应
			if m.config.LogResponse {
				// 脱敏处理响应数据
				if respData := m.sanitizeData(response); respData != nil {
					endTemp.Response = &respData
				}
			}
			m.config.Logger.Infof(ctx, "request completed: %s", toMustString(endTemp))
		}
	}
}

// sanitizeData 脱敏数据
func (m *LoggingMiddleware) sanitizeData(data any) (newData any) {
	if data == nil {
		return nil
	}
	// 将数据转换为map进行处理
	var (
		jsonData []byte
		err      error
	)
	if jsonData, err = json.Marshal(data); err != nil {
		return fmt.Sprintf("failed to marshal data: %v", err)
	}
	var result any
	if err = json.Unmarshal(jsonData, &result); err != nil {
		return string(jsonData) // 如果无法解析，直接返回字符串
	}
	// 递归脱敏，从深度0开始
	return m.sanitizeValue(result, 0)
}

// sanitizeValue 递归脱敏值，添加深度限制
func (m *LoggingMiddleware) sanitizeValue(value any, depth int) (newValue any) {
	// 防止无限递归
	if depth > maxSanitizeDepth {
		return "<max_depth_reached>"
	}

	switch v := value.(type) {
	case map[string]any:
		result := make(map[string]any)
		for key, val := range v {
			// 检查是否为敏感字段
			isSensitive := false
			for _, field := range m.config.SensitiveFields {
				if strings.Contains(strings.ToLower(key), strings.ToLower(field)) {
					isSensitive = true
					break
				}
			}
			if isSensitive {
				result[key] = "***"
			} else {
				result[key] = m.sanitizeValue(val, depth+1) // 递归处理，深度+1
			}
		}
		return result
	case []any:
		result := make([]any, len(v))
		for i, val := range v {
			result[i] = m.sanitizeValue(val, depth+1) // 递归处理数组元素，深度+1
		}
		return result
	default:
		return v // 基本类型直接返回
	}
}
