/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-19 21:39:57
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-19 21:42:24
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package cache

import (
	"context"
	"time"
)

// CacheAdapter 缓存适配器
type CacheAdapter interface {
	GetCache(ctx context.Context, key string) (value any, err error)                         // 获取缓存
	SetCache(ctx context.Context, key string, value any, duration time.Duration) (err error) // 设置缓存
	ClearCache(ctx context.Context, keys ...string) (err error)                              // 清空缓存
}
