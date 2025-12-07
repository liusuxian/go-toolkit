/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-18 02:57:43
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-07 03:39:57
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkmq

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
)

// LogLevel 日志级别
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// Logger 日志接口
type Logger interface {
	Debug(ctx context.Context, format string, args ...any) // 调试日志
	Info(ctx context.Context, format string, args ...any)  // 信息日志
	Warn(ctx context.Context, format string, args ...any)  // 警告日志
	Error(ctx context.Context, format string, args ...any) // 错误日志
}

// DefaultLogger 默认日志实现
type DefaultLogger struct {
	logger *log.Logger
	level  LogLevel
}

// NewDefaultLogger 创建默认日志器
func NewDefaultLogger(level LogLevel) (l *DefaultLogger) {
	return &DefaultLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		level:  level,
	}
}

// Debug 调试日志
func (l *DefaultLogger) Debug(ctx context.Context, format string, args ...any) {
	if LogLevelDebug < l.level {
		return
	}
	caller := fileInfo(2) // 跳过本函数和调用者
	l.logger.Printf("[DEBUG] [%s] %s", caller, fmt.Sprintf(format, args...))
}

// Info 信息日志
func (l *DefaultLogger) Info(ctx context.Context, format string, args ...any) {
	if LogLevelInfo < l.level {
		return
	}
	caller := fileInfo(2) // 跳过本函数和调用者
	l.logger.Printf("[INFO] [%s] %s", caller, fmt.Sprintf(format, args...))
}

// Warn 警告日志
func (l *DefaultLogger) Warn(ctx context.Context, format string, args ...any) {
	if LogLevelWarn < l.level {
		return
	}
	caller := fileInfo(2) // 跳过本函数和调用者
	l.logger.Printf("[WARN] [%s] %s", caller, fmt.Sprintf(format, args...))
}

// Error 错误日志
func (l *DefaultLogger) Error(ctx context.Context, format string, args ...any) {
	if LogLevelError < l.level {
		return
	}
	caller := fileInfo(2) // 跳过本函数和调用者
	l.logger.Printf("[ERROR] [%s] %s", caller, fmt.Sprintf(format, args...))
}

// fileInfo 获取调用者的文件名和行号
func fileInfo(skip int) (caller string) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		caller = "<???>:1"
		return
	}
	caller = fmt.Sprintf("%s:%d", file, line)
	return
}
