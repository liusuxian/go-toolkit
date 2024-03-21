/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-03-18 20:48:59
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-03-21 19:42:09
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtklog_test

import (
	"context"
	"github.com/liusuxian/go-toolkit/gtklog"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewWithOption1(t *testing.T) {
	var (
		ctx     = context.Background()
		assert  = assert.New(t)
		testErr = errors.New("Test Error")
		log     *gtklog.Logger
		err     error
	)
	log, err = gtklog.NewWithOption()
	assert.NoError(err)
	assert.NotNil(log)
	err = log.SetLevel(gtklog.TraceLevel)
	assert.NoError(err)
	err = log.Write(gtklog.DebugLevel, []byte("hello world"))
	assert.NoError(err)
	t.Log("level:", log.GetLevel())
	t.Log("ConfigStr:", log.GetConfigStr())
	log.Trace(ctx, "I am Trace")
	log.Debug(ctx, "I am Debug")
	log.Info(ctx, "I am Info")
	log.Warn(ctx, "I am Warn")
	log.Error(ctx, "I am Error")
	log.Errorf(ctx, "I am Error: %+v", testErr)
}

func TestNewWithOption2(t *testing.T) {
	var (
		ctx     = context.Background()
		assert  = assert.New(t)
		testErr = errors.New("Test Error")
		log     *gtklog.Logger
		err     error
	)
	log, err = gtklog.NewWithOption(func(c *gtklog.Config) {
		c.LogType = "json"
	})
	assert.NoError(err)
	assert.NotNil(log)
	err = log.Write(gtklog.DebugLevel, []byte("hello world"))
	assert.NoError(err)
	t.Log("level:", log.GetLevel())
	t.Log("ConfigStr:", log.GetConfigStr())
	log.Trace(ctx, "I am Trace")
	log.Debug(ctx, "I am Debug")
	log.Info(ctx, "I am Info")
	log.Warn(ctx, "I am Warn")
	log.Error(ctx, "I am Error")
	log.Errorf(ctx, "I am Error: %+v", testErr)
}

// 定义上下文键
const (
	RequestIdKey gtklog.ContextKey = "RequestId"
	AKey         gtklog.ContextKey = "A"
	BKey         gtklog.ContextKey = "B"
	CKey         gtklog.ContextKey = "C"
)

func TestNewWithOption3(t *testing.T) {
	var (
		ctx     = context.Background()
		assert  = assert.New(t)
		testErr = errors.New("Test Error")
		log     *gtklog.Logger
		err     error
	)
	log, err = gtklog.NewWithOption(func(c *gtklog.Config) {
		c.CtxKeys = []gtklog.ContextKey{
			RequestIdKey,
			AKey,
			BKey,
			CKey,
		}
	})
	assert.NoError(err)
	assert.NotNil(log)
	err = log.Write(gtklog.DebugLevel, []byte("hello world"))
	assert.NoError(err)
	t.Log("level:", log.GetLevel())
	t.Log("ConfigStr:", log.GetConfigStr())
	ctx = context.WithValue(ctx, RequestIdKey, "100000")
	ctx = context.WithValue(ctx, AKey, "111")
	ctx = context.WithValue(ctx, BKey, "222")
	ctx = context.WithValue(ctx, CKey, "333")
	log.Trace(ctx, "I am Trace")
	log.Debug(ctx, "I am Debug")
	log.Info(ctx, "I am Info")
	log.Warn(ctx, "I am Warn")
	log.Error(ctx, "I am Error")
	log.Errorf(ctx, "I am Error: %+v", testErr)
}

func TestNewWithOption4(t *testing.T) {
	var (
		ctx     = context.Background()
		assert  = assert.New(t)
		testErr = errors.New("Test Error")
		log     *gtklog.Logger
		err     error
	)
	log, err = gtklog.NewWithOption(func(c *gtklog.Config) {
		c.LogType = "json"
		c.CtxKeys = []gtklog.ContextKey{
			RequestIdKey,
			AKey,
			BKey,
			CKey,
		}
	})
	assert.NoError(err)
	assert.NotNil(log)
	err = log.Write(gtklog.DebugLevel, []byte("hello world"))
	assert.NoError(err)
	t.Log("level:", log.GetLevel())
	t.Log("ConfigStr:", log.GetConfigStr())
	ctx = context.WithValue(ctx, RequestIdKey, "100000")
	ctx = context.WithValue(ctx, AKey, "111")
	ctx = context.WithValue(ctx, BKey, "222")
	ctx = context.WithValue(ctx, CKey, "333")
	log.Trace(ctx, "I am Trace")
	log.Debug(ctx, "I am Debug")
	log.Info(ctx, "I am Info")
	log.Warn(ctx, "I am Warn")
	log.Error(ctx, "I am Error")
	log.Errorf(ctx, "I am Error: %+v", testErr)
}

func TestNewWithOption5(t *testing.T) {
	var (
		ctx     = context.Background()
		assert  = assert.New(t)
		testErr = errors.New("Test Error")
		log     *gtklog.Logger
		err     error
	)
	log, err = gtklog.NewWithOption(func(c *gtklog.Config) {
		c.LogLevelFileName = map[gtklog.Level]string{
			gtklog.TraceLevel: "access.log",
			gtklog.DebugLevel: "access.log",
			gtklog.InfoLevel:  "access.log",
			gtklog.WarnLevel:  "access.log",
			gtklog.ErrorLevel: "error.log",
			gtklog.FatalLevel: "error.log",
			gtklog.PanicLevel: "error.log",
		}
	})
	assert.NoError(err)
	assert.NotNil(log)
	err = log.Write(gtklog.DebugLevel, []byte("hello world"))
	assert.NoError(err)
	t.Log("level:", log.GetLevel())
	t.Log("ConfigStr:", log.GetConfigStr())
	log.Trace(ctx, "I am Trace")
	log.Debug(ctx, "I am Debug")
	log.Info(ctx, "I am Info")
	log.Warn(ctx, "I am Warn")
	log.Error(ctx, "I am Error")
	log.Errorf(ctx, "I am Error: %+v", testErr)
}

func TestLogger(t *testing.T) {
	var (
		ctx     = context.Background()
		assert  = assert.New(t)
		testErr = errors.New("Test Error")
		err     error
	)
	err = gtklog.SetLevel(gtklog.TraceLevel)
	assert.NoError(err)
	err = gtklog.Write(gtklog.DebugLevel, []byte("hello world"))
	assert.NoError(err)
	t.Log("level:", gtklog.GetLevel())
	t.Log("ConfigStr:", gtklog.GetConfigStr())
	gtklog.Trace(ctx, "I am Trace")
	gtklog.Debug(ctx, "I am Debug")
	gtklog.Info(ctx, "I am Info")
	gtklog.Warn(ctx, "I am Warn")
	gtklog.Error(ctx, "I am Error")
	gtklog.Errorf(ctx, "I am Error: %+v", testErr)
}
