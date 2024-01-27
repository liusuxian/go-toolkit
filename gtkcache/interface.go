/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-27 20:46:12
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-28 03:02:22
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkcache

import (
	"context"
	"github.com/liusuxian/go-toolkit/gtkredis"
	"time"
)

type Func func(ctx context.Context) (val any, err error)

// ICache 缓存接口
type ICache interface {
	// Get 获取缓存
	Get(ctx context.Context, key string) (val any, err error)
	// GetReset 获取缓存，并在缓存命中时重置过期时间
	GetReset(ctx context.Context, key string, timeout time.Duration) (val any, err error)
	// GetMap 批量获取缓存（原子操作）
	GetMap(ctx context.Context, keys ...string) (data map[string]any, err error)
	// Set 设置缓存
	Set(ctx context.Context, key string, val any, timeout ...time.Duration) (err error)
	// SetMap 批量设置缓存（原子操作）
	SetMap(ctx context.Context, data map[string]any, timeout ...time.Duration) (err error)
	// CustomCache 自定义缓存
	CustomCache(ctx context.Context, f Func) (val any, err error)
	// IsExist 缓存是否存在
	IsExist(ctx context.Context, key string) (isExist bool)
	// Delete 删除缓存
	Delete(ctx context.Context, keys ...string) (err error)
	// GetExpire 获取缓存过期时间
	GetExpire(ctx context.Context, key string) (timeout time.Duration, err error)
	// Close 关闭缓存服务
	Close(ctx context.Context) (err error)
}

// IRedisCache Redis 缓存接口
type IRedisCache interface {
	ICache
	// Client Redis 客户端
	Client() (client *gtkredis.RedisClient)
}
