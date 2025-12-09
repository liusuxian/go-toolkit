/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-12-09 15:57:50
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-09 16:41:01
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtklog_test

import (
	"context"
	"github.com/liusuxian/go-toolkit/gtklog"
	"testing"
)

func TestDefaultLogger(t *testing.T) {
	var (
		ctx = context.Background()
		log = gtklog.NewDefaultLogger(gtklog.TraceLevel)
	)
	log.Debug(ctx, "hello world")
	log.Debugf(ctx, "hello world %s", "world")
	log.Info(ctx, "hello world")
	log.Infof(ctx, "hello world %s", "world")
	log.Warn(ctx, "hello world")
	log.Warnf(ctx, "hello world %s", "world")
	log.Error(ctx, "hello world")
	log.Errorf(ctx, "hello world %s", "world")
	log.Trace(ctx, "hello world")
	log.Tracef(ctx, "hello world %s", "world")
}
