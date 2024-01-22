/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-19 21:39:57
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-22 10:55:57
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gfcache

import (
	"context"
	"github.com/gogf/gf/v2/container/gvar"
	"time"
)

// CacheAdapter 缓存适配器
type CacheAdapter interface {
	// Get 获取缓存
	Get(ctx context.Context, key any) (val *gvar.Var, err error)
	// Set 设置缓存
	Set(ctx context.Context, key, value any, duration time.Duration) (err error)
	// Clear 清除缓存
	Clear(ctx context.Context, keys ...any) (err error)
}
