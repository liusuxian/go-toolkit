/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-18 02:57:43
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-19 22:03:49
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkkafka

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
	Debugf(ctx context.Context, format string, args ...any) // 调试日志
	Infof(ctx context.Context, format string, args ...any)  // 信息日志
	Errorf(ctx context.Context, format string, args ...any) // 错误日志
}

// defaultLogger 默认日志实现
type defaultLogger struct {
	logger *log.Logger
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

// Debugf 调试日志
func (l *defaultLogger) Debugf(ctx context.Context, format string, args ...any) {
	caller := fileInfo(2) // 跳过本函数和调用者
	l.logger.Printf("[DEBUG] %s %s", caller, fmt.Sprintf(format, args...))
}

// Infof 信息日志
func (l *defaultLogger) Infof(ctx context.Context, format string, args ...any) {
	caller := fileInfo(2) // 跳过本函数和调用者
	l.logger.Printf("[INFO] %s %s", caller, fmt.Sprintf(format, args...))
}

// Errorf 错误日志
func (l *defaultLogger) Errorf(ctx context.Context, format string, args ...any) {
	caller := fileInfo(2) // 跳过本函数和调用者
	l.logger.Printf("[ERROR] %s %s", caller, fmt.Sprintf(format, args...))
}

// newDefaultLogger 新建默认日志
func newDefaultLogger() (logger *defaultLogger) {
	return &defaultLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}
