/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-12-09 15:57:50
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-09 16:40:09
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtklog

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// ILogger 日志接口
type ILogger interface {
	Debug(ctx context.Context, args ...any)                 // 调试日志
	Debugf(ctx context.Context, format string, args ...any) // 调试日志
	Info(ctx context.Context, args ...any)                  // 信息日志
	Infof(ctx context.Context, format string, args ...any)  // 信息日志
	Warn(ctx context.Context, args ...any)                  // 警告日志
	Warnf(ctx context.Context, format string, args ...any)  // 警告日志
	Error(ctx context.Context, args ...any)                 // 错误日志
	Errorf(ctx context.Context, format string, args ...any) // 错误日志
	Trace(ctx context.Context, args ...any)                 // 跟踪日志
	Tracef(ctx context.Context, format string, args ...any) // 跟踪日志
	Fatal(ctx context.Context, args ...any)                 // 致命日志
	Fatalf(ctx context.Context, format string, args ...any) // 致命日志
	Panic(ctx context.Context, args ...any)                 // 恐慌日志
	Panicf(ctx context.Context, format string, args ...any) // 恐慌日志
}

// Level 日志级别
type Level uint32

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

// DefaultLogger 默认日志实现
type DefaultLogger struct {
	logger *log.Logger
	level  Level
}

// NewDefaultLogger 创建默认日志器
func NewDefaultLogger(level Level) (l *DefaultLogger) {
	return &DefaultLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		level:  level,
	}
}

// Debug 调试日志
func (l *DefaultLogger) Debug(ctx context.Context, args ...any) {
	if DebugLevel > l.level {
		return
	}
	caller := fileInfo(2) // 跳过本函数和调用者
	l.logger.Printf("[DEBUG] [%s] %s", caller, fmt.Sprint(args...))
}

// Debugf 调试日志
func (l *DefaultLogger) Debugf(ctx context.Context, format string, args ...any) {
	if DebugLevel > l.level {
		return
	}
	caller := fileInfo(2) // 跳过本函数和调用者
	l.logger.Printf("[DEBUG] [%s] %s", caller, fmt.Sprintf(format, args...))
}

// Info 信息日志
func (l *DefaultLogger) Info(ctx context.Context, args ...any) {
	if InfoLevel > l.level {
		return
	}
	caller := fileInfo(2) // 跳过本函数和调用者
	l.logger.Printf("[INFO] [%s] %s", caller, fmt.Sprint(args...))
}

// Infof 信息日志
func (l *DefaultLogger) Infof(ctx context.Context, format string, args ...any) {
	if InfoLevel > l.level {
		return
	}
	caller := fileInfo(2) // 跳过本函数和调用者
	l.logger.Printf("[INFO] [%s] %s", caller, fmt.Sprintf(format, args...))
}

// Warn 警告日志
func (l *DefaultLogger) Warn(ctx context.Context, args ...any) {
	if WarnLevel > l.level {
		return
	}
	caller := fileInfo(2) // 跳过本函数和调用者
	l.logger.Printf("[WARN] [%s] %s", caller, fmt.Sprint(args...))
}

// Warnf 警告日志
func (l *DefaultLogger) Warnf(ctx context.Context, format string, args ...any) {
	if WarnLevel > l.level {
		return
	}
	caller := fileInfo(2) // 跳过本函数和调用者
	l.logger.Printf("[WARN] [%s] %s", caller, fmt.Sprintf(format, args...))
}

// Error 错误日志
func (l *DefaultLogger) Error(ctx context.Context, args ...any) {
	if ErrorLevel > l.level {
		return
	}
	caller := fileInfo(2) // 跳过本函数和调用者
	l.logger.Printf("[ERROR] [%s] %s", caller, fmt.Sprint(args...))
}

// Errorf 错误日志
func (l *DefaultLogger) Errorf(ctx context.Context, format string, args ...any) {
	if ErrorLevel > l.level {
		return
	}
	caller := fileInfo(2) // 跳过本函数和调用者
	l.logger.Printf("[ERROR] [%s] %s", caller, fmt.Sprintf(format, args...))
}

// Trace 跟踪日志
func (l *DefaultLogger) Trace(ctx context.Context, args ...any) {
	if TraceLevel > l.level {
		return
	}
	caller := fileInfo(2) // 跳过本函数和调用者
	l.logger.Printf("[TRACE] [%s] %s", caller, fmt.Sprint(args...))
}

// Tracef 跟踪日志
func (l *DefaultLogger) Tracef(ctx context.Context, format string, args ...any) {
	if TraceLevel > l.level {
		return
	}
	caller := fileInfo(2) // 跳过本函数和调用者
	l.logger.Printf("[TRACE] [%s] %s", caller, fmt.Sprintf(format, args...))
}

// Fatal 致命日志
func (l *DefaultLogger) Fatal(ctx context.Context, args ...any) {
	if FatalLevel > l.level {
		return
	}
	caller := fileInfo(2) // 跳过本函数和调用者
	l.logger.Printf("[FATAL] [%s] %s", caller, fmt.Sprint(args...))
	os.Exit(1)
}

// Fatalf 致命日志
func (l *DefaultLogger) Fatalf(ctx context.Context, format string, args ...any) {
	if FatalLevel > l.level {
		return
	}
	caller := fileInfo(2) // 跳过本函数和调用者
	l.logger.Printf("[FATAL] [%s] %s", caller, fmt.Sprintf(format, args...))
	os.Exit(1)
}

// Panic 恐慌日志
func (l *DefaultLogger) Panic(ctx context.Context, args ...any) {
	if PanicLevel > l.level {
		return
	}
	caller := fileInfo(2) // 跳过本函数和调用者
	l.logger.Printf("[PANIC] [%s] %s", caller, fmt.Sprint(args...))
	panic(fmt.Sprint(args...))
}

// Panicf 恐慌日志
func (l *DefaultLogger) Panicf(ctx context.Context, format string, args ...any) {
	if PanicLevel > l.level {
		return
	}
	caller := fileInfo(2) // 跳过本函数和调用者
	l.logger.Printf("[PANIC] [%s] %s", caller, fmt.Sprintf(format, args...))
	panic(fmt.Sprint(args...))
}

// fileInfo 获取调用者的文件名和行号
func fileInfo(skip int) (caller string) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		caller = "<???>:1"
		return
	}
	caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	return
}
