/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-19 21:56:00
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-19 21:56:04
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gflogger_test

import (
	"context"
	"go-toolkit/gflogger"
	"testing"
)

func TestLogger(t *testing.T) {
	ctx := context.Background()
	gflogger.Print(ctx, "hello")
	gflogger.Printf(ctx, "hello %s", "world")
}
