/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-03-18 20:48:59
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-03-21 19:39:25
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtklog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/liusuxian/go-toolkit/gtkconf"
	"github.com/liusuxian/go-toolkit/gtkfile"
	"github.com/liusuxian/go-toolkit/gtkjson"
	"github.com/liusuxian/go-toolkit/gtkstr"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// 日志级别类型
type Level string

// 日志级别
const (
	PanicLevel Level = "panic"
	FatalLevel Level = "fatal"
	ErrorLevel Level = "error"
	WarnLevel  Level = "warning"
	InfoLevel  Level = "info"
	DebugLevel Level = "debug"
	TraceLevel Level = "trace"
)

type ContextKey string

const (
	CtxValFieldKey   = "ctx"
	onlyWriteMessage = "only_write_msg"

	fieldKeyMsg   = "msg"
	fieldKeyLevel = "level"
	fieldKeyTime  = "time"
	fieldKeyError = "error"
	fieldKeyFunc  = "func"
	fieldKeyFile  = "file"
)

// 要写入日志的数据字段
type Fields map[string]any

// Config 日志配置
type Config struct {
	LogPath            string           `json:"logPath" dc:"日志文件路径，默认 logs"`                                            // 日志文件路径，默认 logs
	LogType            string           `json:"logType" dc:"日志类型，json|text，默认 text"`                                    // 日志类型，json|text，默认 text
	LogLevel           string           `json:"logLevel" dc:"日志级别，panic、fatal、error、warning、info、debug、trace，默认 debug"` // 日志级别，panic、fatal、error、warning、info、debug、trace，默认 debug
	CtxKeys            []ContextKey     `json:"ctxKeys" dc:"自定义 Context 上下文变量名称，自动打印 Context 的变量到日志中，默认为空"`             // 自定义 Context 上下文变量名称，自动打印 Context 的变量到日志中，默认为空
	LogLevelFileName   map[Level]string `json:"logLevelFileName" dc:"日志级别所对应的日志文件名称，默认 gtklog.log"`                     // 日志级别所对应的日志文件名称，默认 gtklog.log
	FileNameDateFormat string           `json:"fileNameDateFormat" dc:"文件名的日期格式，默认 %Y-%m-%d"`                           // 文件名的日期格式，默认 %Y-%m-%d
	TimestampFormat    string           `json:"timestampFormat" dc:"日志中日期时间格式，默认 2006-01-02 15:04:05.000"`              // 日志中日期时间格式，默认 2006-01-02 15:04:05.000
	FileInfoField      string           `json:"fileInfoField" dc:"文件名和行号字段名，默认 caller"`                                 // 文件名和行号字段名，默认 caller
	JSONPrettyPrint    bool             `json:"jsonPrettyPrint" dc:"json日志是否美化输出，默认 false"`                             // json日志是否美化输出，默认 false
	JSONDataKey        string           `json:"jsonDataKey" dc:"json日志条目中，数据字段都会作为该字段的嵌入字段，默认为空"`                       // json日志条目中，数据字段都会作为该字段的嵌入字段，默认为空
	MaxAge             time.Duration    `json:"maxAge" dc:"保留旧日志文件的最长时间，默认 7天"`                                         // 保留旧日志文件的最长时间，默认 7天
	RotationTime       time.Duration    `json:"rotationTime" dc:"日志轮转的时间间隔，默认 24小时"`                                    // 日志轮转的时间间隔，默认 24小时
	RotationSize       int64            `json:"rotationSize" dc:"日志文件达到指定大小时进行轮转，默认 1024*1024*1024*5"`                  // 日志文件达到指定大小时进行轮转，默认 1024*1024*1024*5
	Stdout             bool             `json:"stdout" dc:"是否输出到控制台，默认 false"`                                          // 是否输出到控制台，默认 false
}

// ConfigOption 日志配置选项
type ConfigOption func(c *Config)

// Logger 日志结构
type Logger struct {
	logrus *logrus.Logger // 日志对象
	config *Config        // 日志配置
}

// gLogger 实例
var gLogger *Logger

func init() {
	// 读取配置
	var (
		config *Config
		err    error
	)
	if err = gtkconf.StructKey("logger", &config); err != nil {
		panic(errors.Wrapf(err, "Get Logger Config Error"))
	}
	// 新建日志
	if gLogger, err = NewWithConfig(config); err != nil {
		panic(errors.Wrapf(err, "New Logger Error"))
	}
}

// 实现 logrus.Formatter 接口
type textFormatter struct {
	config *Config // 日志配置
}

// Format
func (tf *textFormatter) Format(entry *logrus.Entry) (bs []byte, err error) {
	// 只打印 Message 数据
	var buf bytes.Buffer
	if _, ok := entry.Data[onlyWriteMessage]; ok {
		buf.WriteString(entry.Message)
		buf.WriteString("\n")
		bs = buf.Bytes()
		return
	}
	// 打印全部数据
	buf.WriteString(entry.Time.Format(tf.config.TimestampFormat))
	buf.WriteString(" ")
	buf.WriteString(fmt.Sprintf("[%s]", strings.ToUpper(entry.Level.String())))
	buf.WriteString(" ")
	// 处理自定义 Context 上下文变量
	if val, ok := entry.Data[CtxValFieldKey]; ok {
		ctxValMap := val.(map[ContextKey]any)
		ctxVals := make([]string, 0, len(ctxValMap))
		for _, k := range tf.config.CtxKeys {
			ctxVals = append(ctxVals, fmt.Sprintf("%s:%v", k, ctxValMap[k]))
		}
		buf.WriteString("{")
		buf.WriteString(strings.Join(ctxVals, ", "))
		buf.WriteString("} ")
	}
	// 处理文件名和行号字段
	if val, ok := entry.Data[tf.config.FileInfoField]; ok {
		buf.WriteString(fmt.Sprintf("%s", val))
		buf.WriteString(" ")
	}
	// 处理其他数据字段
	for k, v := range entry.Data {
		if k == CtxValFieldKey || k == tf.config.FileInfoField {
			continue
		}
		buf.WriteString(fmt.Sprintf("%s", v))
		buf.WriteString(" ")
	}
	// 处理 Message 字段
	buf.WriteString(entry.Message)
	buf.WriteString("\n")
	bs = buf.Bytes()
	return
}

// 实现 logrus.Formatter 接口
type jsonFormatter struct {
	config *Config // 日志配置
	// 是否禁用输出中的 HTML 转义
	disableHTMLEscape bool
}

// Format
func (jf *jsonFormatter) Format(entry *logrus.Entry) (bs []byte, err error) {
	// 只打印 Message 数据
	if _, ok := entry.Data[onlyWriteMessage]; ok {
		var buf bytes.Buffer
		buf.WriteString(entry.Message)
		buf.WriteString("\n")
		bs = buf.Bytes()
		return
	}
	// 打印全部数据
	var data Fields
	if jf.config.JSONDataKey != "" {
		data = make(Fields, 5)
		data[jf.config.JSONDataKey] = make(Fields)
	} else {
		data = make(Fields, len(entry.Data)+5)
	}
	data[fieldKeyTime] = entry.Time.Format(jf.config.TimestampFormat)
	data[fieldKeyLevel] = strings.ToUpper(entry.Level.String())
	// 处理自定义 Context 上下文变量
	if val, ok := entry.Data[CtxValFieldKey]; ok {
		if jf.config.JSONDataKey != "" {
			data[jf.config.JSONDataKey].(Fields)[CtxValFieldKey] = val.(map[ContextKey]any)
		} else {
			data[CtxValFieldKey] = val.(map[ContextKey]any)
		}
	}
	// 处理文件名和行号字段
	if val, ok := entry.Data[jf.config.FileInfoField]; ok {
		if jf.config.JSONDataKey != "" {
			data[jf.config.JSONDataKey].(Fields)[jf.config.FileInfoField] = fmt.Sprintf("%s", val)
		} else {
			data[jf.config.FileInfoField] = fmt.Sprintf("%s", val)
		}
	}
	// 处理其他数据字段
	for k, v := range entry.Data {
		if k == CtxValFieldKey || k == jf.config.FileInfoField || k == fieldKeyFile || k == fieldKeyFunc {
			continue
		}

		var vv any
		switch v := v.(type) {
		case error:
			vv = v.Error()
		default:
			vv = v
		}
		if jf.config.JSONDataKey != "" {
			data[jf.config.JSONDataKey].(Fields)[k] = vv
		} else {
			data[k] = vv
		}
	}
	// 处理 Message 字段
	data[fieldKeyMsg] = entry.Message
	// json日志是否美化输出，默认 false
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(!jf.disableHTMLEscape)
	if jf.config.JSONPrettyPrint {
		encoder.SetIndent("", "  ")
	}
	if err = encoder.Encode(data); err != nil {
		err = errors.Errorf("failed to marshal fields to JSON, %v", err)
		return
	}
	bs = buf.Bytes()
	return
}

// NewWithOption 新建日志
func NewWithOption(opts ...ConfigOption) (logger *Logger, err error) {
	if logger, err = newLogger(nil, opts...); err != nil {
		return
	}
	var writerMap lfshook.WriterMap
	if writerMap, err = getWriterMap(logger.config); err != nil {
		return
	}
	fileHook := lfshook.NewHook(writerMap, logger.logrus.Formatter)
	logger.logrus.Hooks.Add(fileHook)
	return
}

// NewWithConfig 新建日志
func NewWithConfig(config *Config) (logger *Logger, err error) {
	if logger, err = newLogger(config); err != nil {
		return
	}
	var writerMap lfshook.WriterMap
	if writerMap, err = getWriterMap(logger.config); err != nil {
		return
	}
	fileHook := lfshook.NewHook(writerMap, logger.logrus.Formatter)
	logger.logrus.Hooks.Add(fileHook)
	return
}

// Trace
func (l *Logger) Trace(ctx context.Context, args ...any) {
	l.withFields(ctx, Fields{}).Trace(args...)
}

// Tracef
func (l *Logger) Tracef(ctx context.Context, format string, args ...any) {
	l.withFields(ctx, Fields{}).Tracef(format, args...)
}

// Debug
func (l *Logger) Debug(ctx context.Context, args ...any) {
	l.withFields(ctx, Fields{}).Debug(args...)
}

// Debugf
func (l *Logger) Debugf(ctx context.Context, format string, args ...any) {
	l.withFields(ctx, Fields{}).Debugf(format, args...)
}

// Info
func (l *Logger) Info(ctx context.Context, args ...any) {
	l.withFields(ctx, Fields{}).Info(args...)
}

// Infof
func (l *Logger) Infof(ctx context.Context, format string, args ...any) {
	l.withFields(ctx, Fields{}).Infof(format, args...)
}

// Warn
func (l *Logger) Warn(ctx context.Context, args ...any) {
	l.withFields(ctx, Fields{}).Warn(args...)
}

// Warnf
func (l *Logger) Warnf(ctx context.Context, format string, args ...any) {
	l.withFields(ctx, Fields{}).Warnf(format, args...)
}

// Error
func (l *Logger) Error(ctx context.Context, args ...any) {
	l.withFields(ctx, Fields{}).Error(args...)
}

// Errorf
func (l *Logger) Errorf(ctx context.Context, format string, args ...any) {
	l.withFields(ctx, Fields{}).Errorf(format, args...)
}

// Fatal
func (l *Logger) Fatal(ctx context.Context, args ...any) {
	l.withFields(ctx, Fields{}).Fatal(args...)
}

// Fatalf
func (l *Logger) Fatalf(ctx context.Context, format string, args ...any) {
	l.withFields(ctx, Fields{}).Fatalf(format, args...)
}

// Panic
func (l *Logger) Panic(ctx context.Context, args ...any) {
	l.withFields(ctx, Fields{}).Panic(args...)
}

// Panicf
func (l *Logger) Panicf(ctx context.Context, format string, args ...any) {
	l.withFields(ctx, Fields{}).Panicf(format, args...)
}

// Write 使用日志输出对象写入数据
func (l *Logger) Write(level Level, data []byte) (err error) {
	var lv logrus.Level
	if lv, err = logrus.ParseLevel(string(level)); err != nil {
		return
	}
	w := l.withFields(context.TODO(), Fields{onlyWriteMessage: 1}).WriterLevel(lv)
	defer w.Close()
	_, err = w.Write(data)
	return
}

// GetLevel
func (l *Logger) GetLevel() (level Level) {
	level = Level(l.logrus.GetLevel().String())
	return
}

// SetLevel
func (l *Logger) SetLevel(level Level) (err error) {
	var lv logrus.Level
	if lv, err = logrus.ParseLevel(string(level)); err != nil {
		return
	}
	l.logrus.SetLevel(lv)
	return
}

// GetConfigStr
func (l *Logger) GetConfigStr() (str string) {
	return gtkjson.MustString(l.config)
}

// Trace
func Trace(ctx context.Context, args ...any) {
	gLogger.withFields(ctx, Fields{}).Trace(args...)
}

// Tracef
func Tracef(ctx context.Context, format string, args ...any) {
	gLogger.withFields(ctx, Fields{}).Tracef(format, args...)
}

// Debug
func Debug(ctx context.Context, args ...any) {
	gLogger.withFields(ctx, Fields{}).Debug(args...)
}

// Debugf
func Debugf(ctx context.Context, format string, args ...any) {
	gLogger.withFields(ctx, Fields{}).Debugf(format, args...)
}

// Info
func Info(ctx context.Context, args ...any) {
	gLogger.withFields(ctx, Fields{}).Info(args...)
}

// Infof
func Infof(ctx context.Context, format string, args ...any) {
	gLogger.withFields(ctx, Fields{}).Infof(format, args...)
}

// Warn
func Warn(ctx context.Context, args ...any) {
	gLogger.withFields(ctx, Fields{}).Warn(args...)
}

// Warnf
func Warnf(ctx context.Context, format string, args ...any) {
	gLogger.withFields(ctx, Fields{}).Warnf(format, args...)
}

// Error
func Error(ctx context.Context, args ...any) {
	gLogger.withFields(ctx, Fields{}).Error(args...)
}

// Errorf
func Errorf(ctx context.Context, format string, args ...any) {
	gLogger.withFields(ctx, Fields{}).Errorf(format, args...)
}

// Fatal
func Fatal(ctx context.Context, args ...any) {
	gLogger.withFields(ctx, Fields{}).Fatal(args...)
}

// Fatalf
func Fatalf(ctx context.Context, format string, args ...any) {
	gLogger.withFields(ctx, Fields{}).Fatalf(format, args...)
}

// Panic
func Panic(ctx context.Context, args ...any) {
	gLogger.withFields(ctx, Fields{}).Panic(args...)
}

// Panicf
func Panicf(ctx context.Context, format string, args ...any) {
	gLogger.withFields(ctx, Fields{}).Panicf(format, args...)
}

// Write 使用日志输出对象写入数据
func Write(level Level, data []byte) (err error) {
	var lv logrus.Level
	if lv, err = logrus.ParseLevel(string(level)); err != nil {
		return
	}
	w := gLogger.withFields(context.TODO(), Fields{onlyWriteMessage: 1}).WriterLevel(lv)
	defer w.Close()
	_, err = w.Write(data)
	return
}

// GetLevel
func GetLevel() (level Level) {
	level = Level(gLogger.logrus.GetLevel().String())
	return
}

// SetLevel
func SetLevel(level Level) (err error) {
	var lv logrus.Level
	if lv, err = logrus.ParseLevel(string(level)); err != nil {
		return
	}
	gLogger.logrus.SetLevel(lv)
	return
}

// GetConfigStr
func GetConfigStr() (str string) {
	return gtkjson.MustString(gLogger.config)
}

// withFields
func (l *Logger) withFields(ctx context.Context, fields Fields) (entry *logrus.Entry) {
	ctxValField := map[ContextKey]any{}
	for _, k := range l.config.CtxKeys {
		ctxValField[k] = ctx.Value(k)
	}
	if len(ctxValField) > 0 {
		fields[CtxValFieldKey] = ctxValField
	}
	fields[l.config.FileInfoField] = fileInfo(3)
	entry = l.logrus.WithFields(logrus.Fields(fields))
	return
}

// fileInfo
func fileInfo(skip int) (caller string) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		caller = "<???>:1"
		return
	}
	caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	return
}

// newLogger 新建日志
func newLogger(cfg *Config, opts ...ConfigOption) (logger *Logger, err error) {
	// 新建日志对象
	if cfg == nil {
		cfg = &Config{}
	}
	logger = &Logger{
		logrus: logrus.New(),
		config: cfg,
	}
	// 启用 ReportCaller
	logger.logrus.SetReportCaller(true)
	// 处理日志配置选项
	for _, opt := range opts {
		opt(logger.config)
	}
	// 日志文件路径，默认 logs
	if logger.config.LogPath == "" {
		logger.config.LogPath = "logs"
	}
	// 日志类型，json|text，默认 text
	if logger.config.LogType == "" {
		logger.config.LogType = "text"
	}
	// 日志级别，panic、fatal、error、warning、info、debug、trace，默认 debug
	if logger.config.LogLevel == "" {
		logger.config.LogLevel = strings.ToLower(logrus.DebugLevel.String())
	}
	// 自定义 Context 上下文变量名称，自动打印 Context 的变量到日志中，默认为空
	if logger.config.CtxKeys == nil {
		logger.config.CtxKeys = []ContextKey{}
	}
	// 日志级别所对应的日志文件名称，默认 gtklog.log
	if len(logger.config.LogLevelFileName) == 0 {
		logger.config.LogLevelFileName = map[Level]string{
			PanicLevel: "gtklog.log",
			FatalLevel: "gtklog.log",
			ErrorLevel: "gtklog.log",
			WarnLevel:  "gtklog.log",
			InfoLevel:  "gtklog.log",
			DebugLevel: "gtklog.log",
			TraceLevel: "gtklog.log",
		}
	}
	// 文件名的日期格式，默认 %Y-%m-%d
	if logger.config.FileNameDateFormat == "" {
		logger.config.FileNameDateFormat = "%Y-%m-%d"
	}
	// 日志中日期时间格式，默认 2006-01-02 15:04:05.000
	if logger.config.TimestampFormat == "" {
		logger.config.TimestampFormat = "2006-01-02 15:04:05.000"
	}
	// 文件名和行号字段名，默认 caller
	if logger.config.FileInfoField == "" {
		logger.config.FileInfoField = "caller"
	}
	// 保留旧日志文件的最长时间，默认 7天
	if logger.config.MaxAge == time.Duration(0) {
		logger.config.MaxAge = time.Hour * 24 * 7
	}
	// 日志轮转的时间间隔，默认 24小时
	if logger.config.RotationTime == time.Duration(0) {
		logger.config.RotationTime = time.Hour * 24
	}
	// 日志文件达到指定大小时进行轮转，默认 1024*1024*1024*5
	if logger.config.RotationSize == 0 {
		logger.config.RotationSize = 1024 * 1024 * 1024 * 5
	}
	// 创建日志目录
	if err = makeDirAll(logger.config.LogPath); err != nil {
		return
	}
	// 是否输出到控制台
	if !logger.config.Stdout {
		logger.logrus.SetOutput(io.Discard)
	}
	// 日志级别，panic、fatal、error、warning、info、debug、trace，默认 debug
	if logger.logrus.Level, err = logrus.ParseLevel(logger.config.LogLevel); err != nil {
		return
	}
	// 日志类型
	switch logger.config.LogType {
	case "json":
		logger.logrus.Formatter = &jsonFormatter{
			config: logger.config,
		}
	default:
		logger.logrus.Formatter = &textFormatter{
			config: logger.config,
		}
	}
	return
}

// getFileName 获取日志文件名称
func getFileName(levelFileName, fileNameDateFormat string) (fileName string) {
	fileNameList := gtkstr.Split(levelFileName, ".")
	fileNameListLen := len(fileNameList)
	if fileNameListLen == 1 {
		fileName = fmt.Sprintf("%s-%s.log", fileNameList[0], fileNameDateFormat)
	} else if fileNameListLen >= 2 {
		fileName = fmt.Sprintf("%s-%s.%s", strings.Join(fileNameList[:fileNameListLen-1], "-"), fileNameDateFormat, fileNameList[fileNameListLen-1])
	} else {
		fileName = fmt.Sprintf("gtklog-%s.log", fileNameDateFormat)
	}
	return
}

// getWriterMap
func getWriterMap(config *Config) (writerMap lfshook.WriterMap, err error) {
	writerMap = make(lfshook.WriterMap, len(config.LogLevelFileName))
	for level, levelFileName := range config.LogLevelFileName {
		var lv logrus.Level
		if lv, err = logrus.ParseLevel(string(level)); err != nil {
			return
		}
		if writerMap[lv], err = rotatelogs.New(
			fmt.Sprintf("%s/%s", config.LogPath, getFileName(levelFileName, config.FileNameDateFormat)),
			rotatelogs.WithMaxAge(config.MaxAge),
			rotatelogs.WithRotationTime(config.RotationTime),
			rotatelogs.WithRotationSize(config.RotationSize),
		); err != nil {
			return
		}
	}
	return
}

// makeDirAll 创建日志目录
func makeDirAll(logPath string) (err error) {
	if !gtkfile.PathExists(logPath) {
		if err = os.MkdirAll(logPath, os.ModePerm); err != nil {
			return errors.Errorf("create <%s> error: %s", logPath, err)
		}
	}
	return
}
