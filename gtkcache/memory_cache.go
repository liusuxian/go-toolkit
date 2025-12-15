/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-12-15 12:50:11
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-15 21:35:07
 * @Description: ICache 接口的内存缓存实现（基于 github.com/patrickmn/go-cache）
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkcache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/liusuxian/go-toolkit/internal/utils"
	gocache "github.com/patrickmn/go-cache"
	"golang.org/x/sync/singleflight"
	"io"
	"time"
)

// MemoryCache 内存缓存
type MemoryCache struct {
	cache  *gocache.Cache
	group  singleflight.Group // 用于防止缓存击穿，确保相同 key 的函数只执行一次
	closed bool
}

// NewMemoryCache 创建内存缓存
func NewMemoryCache(ctx context.Context, cleanupInterval ...time.Duration) (mc *MemoryCache, err error) {
	interval := time.Duration(0) // 默认懒惰删除
	if len(cleanupInterval) > 0 {
		interval = cleanupInterval[0]
	}

	mc = &MemoryCache{
		cache:  gocache.New(gocache.NoExpiration, interval),
		closed: false,
	}
	return
}

// Get 获取缓存
//
//	当`timeout > 0`且缓存命中时，设置/重置`key`的过期时间
func (mc *MemoryCache) Get(ctx context.Context, key string, timeout ...time.Duration) (val any, err error) {
	var isExist bool
	if val, isExist = mc.cache.Get(key); !isExist {
		return
	}
	// 重置过期时间
	if len(timeout) > 0 && timeout[0] > 0 {
		mc.cache.Replace(key, val, timeout[0])
	}
	return
}

// GetMap 批量获取缓存
//
//	当`timeout > 0`且所有缓存都命中时，设置/重置所有`key`的过期时间，所有`key`过期时间相同
func (mc *MemoryCache) GetMap(ctx context.Context, keys []string, timeout ...time.Duration) (data map[string]any, err error) {
	if len(keys) == 0 {
		return
	}

	data = make(map[string]any)
	allHit := true
	for _, key := range keys {
		var (
			val     any
			isExist bool
		)
		if val, isExist = mc.cache.Get(key); !isExist {
			data[key] = nil
			allHit = false
			continue
		}
		data[key] = val
	}
	// 只有所有`key`都命中时才重置过期时间
	if allHit && len(timeout) > 0 && timeout[0] > 0 {
		for key, val := range data {
			mc.cache.Replace(key, val, timeout[0])
		}
	}
	return
}

// GetOrSet 检索并返回`key`的值，或者当`key`不存在时，则使用`newVal`设置`key`的值
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
func (mc *MemoryCache) GetOrSet(ctx context.Context, key string, newVal any, timeout ...time.Duration) (val any, err error) {
	var isExist bool
	if val, isExist = mc.cache.Get(key); isExist {
		// 重置过期时间
		if len(timeout) > 0 && timeout[0] > 0 {
			mc.cache.Replace(key, val, timeout[0])
		}
		return
	}
	// 设置新值
	mc.cache.Set(key, newVal, mc.getTimeout(timeout...))
	val = newVal
	return
}

// GetOrSetFunc 检索并返回`key`的值，或者当`key`不存在时，则使用函数`f`的结果设置`key`的值
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
//	当`force = true`时，可防止缓存穿透
func (mc *MemoryCache) GetOrSetFunc(ctx context.Context, key string, f Func, force bool, timeout ...time.Duration) (val any, err error) {
	var isExist bool
	if val, isExist = mc.cache.Get(key); isExist {
		// 重置过期时间
		if len(timeout) > 0 && timeout[0] > 0 {
			mc.cache.Replace(key, val, timeout[0])
		}
		return
	}
	// 执行函数获取新值
	var newVal any
	if newVal, err = f(ctx); err != nil {
		return
	}
	if utils.IsNil(newVal) && !force {
		return
	}
	// 此处不判断`newVal == nil`是因为防止缓存穿透
	mc.cache.Set(key, newVal, mc.getTimeout(timeout...))
	val = newVal
	return
}

// GetOrSetFuncLock 检索并返回`key`的值，或者当`key`不存在时，则使用函数`f`的结果设置`key`的值，函数`f`是在写入互斥锁内执行，以确保并发安全
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
//	当`force = true`时，可防止缓存穿透
func (mc *MemoryCache) GetOrSetFuncLock(ctx context.Context, key string, f Func, force bool, timeout ...time.Duration) (val any, err error) {
	var isExist bool
	if val, isExist = mc.cache.Get(key); isExist {
		// 重置过期时间
		if len(timeout) > 0 && timeout[0] > 0 {
			mc.cache.Replace(key, val, timeout[0])
		}
		return
	}
	// 使用 singleflight 确保同一个 key 的函数只执行一次
	val, err, _ = mc.group.Do(key, func() (v any, e error) {
		return mc.GetOrSetFunc(ctx, key, f, force, timeout...)
	})
	return
}

// CustomGetOrSetFunc 从缓存中获取指定键`keys`的值，如果缓存未命中，则使用函数`f`的结果设置`keys`的值
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
//	当`force = true`时，可防止缓存穿透
func (mc *MemoryCache) CustomGetOrSetFunc(ctx context.Context, keys []string, args []any, cc ICustomCache, f Func, force bool, timeout ...time.Duration) (val any, err error) {
	// 获取缓存
	if val, err = cc.Get(ctx, keys, args, timeout...); err != nil {
		return
	}
	if val != nil {
		return
	}
	// 执行函数获取新值
	var newVal any
	if newVal, err = f(ctx); err != nil {
		return
	}
	if utils.IsNil(newVal) && !force {
		return
	}
	// 此处不判断`newVal == nil`是因为防止缓存穿透
	val, err = cc.Set(ctx, keys, args, newVal, timeout...)
	return
}

// CustomGetOrSetFuncLock 从缓存中获取指定键`keys`的值，如果缓存未命中，则使用函数`f`的结果设置`keys`的值，函数`f`是在写入互斥锁内执行，以确保并发安全
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
//	当`force = true`时，可防止缓存穿透
func (mc *MemoryCache) CustomGetOrSetFuncLock(ctx context.Context, keys []string, args []any, cc ICustomCache, f Func, force bool, timeout ...time.Duration) (val any, err error) {
	// 获取缓存
	if val, err = cc.Get(ctx, keys, args, timeout...); err != nil {
		return
	}
	if val != nil {
		return
	}
	// 生成 singleflight 的唯一 key
	var sfKey string
	if sfKey, err = mc.generateSingleflightKey(keys, args); err != nil {
		return
	}
	// 使用 singleflight 确保同一个 key 的函数只执行一次
	val, err, _ = mc.group.Do(sfKey, func() (v any, e error) {
		return mc.CustomGetOrSetFunc(ctx, keys, args, cc, f, force, timeout...)
	})
	return
}

// Set 设置缓存
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
func (mc *MemoryCache) Set(ctx context.Context, key string, val any, timeout ...time.Duration) (err error) {
	mc.cache.Set(key, val, mc.getTimeout(timeout...))
	return
}

// SetMap 批量设置缓存，所有`key`的过期时间相同
//
//	当`timeout > 0`时，设置/重置所有`key`的过期时间，所有`key`过期时间相同
func (mc *MemoryCache) SetMap(ctx context.Context, data map[string]any, timeout ...time.Duration) (err error) {
	if len(data) == 0 {
		return
	}

	duration := mc.getTimeout(timeout...)
	for key, val := range data {
		mc.cache.Set(key, val, duration)
	}
	return
}

// SetIfNotExist 当`key`不存在时，则使用`val`设置`key`的值，返回是否设置成功
//
//	当`timeout > 0`且`key`设置成功时，设置`key`的过期时间
func (mc *MemoryCache) SetIfNotExist(ctx context.Context, key string, val any, timeout ...time.Duration) (ok bool, err error) {
	if e := mc.cache.Add(key, val, mc.getTimeout(timeout...)); e != nil {
		return
	}
	ok = true
	return
}

// SetIfNotExistFunc 当`key`不存在时，则使用函数`f`的结果设置`key`的值，返回是否设置成功
//
//	当`timeout > 0`且`key`设置成功时，设置`key`的过期时间
//	当`force = true`时，可防止缓存穿透
func (mc *MemoryCache) SetIfNotExistFunc(ctx context.Context, key string, f Func, force bool, timeout ...time.Duration) (ok bool, err error) {
	if _, isExist := mc.cache.Get(key); isExist {
		return
	}
	// 执行函数获取新值
	var val any
	if val, err = f(ctx); err != nil {
		return
	}
	if utils.IsNil(val) && !force {
		return
	}
	// 此处不判断`val == nil`是因为防止缓存穿透
	return mc.SetIfNotExist(ctx, key, val, timeout...)
}

// SetIfNotExistFuncLock 当`key`不存在时，则使用函数`f`的结果设置`key`的值，返回是否设置成功，函数`f`是在写入互斥锁内执行，以确保并发安全
//
//	当`timeout > 0`且`key`设置成功时，设置`key`的过期时间
//	当`force = true`时，可防止缓存穿透
func (mc *MemoryCache) SetIfNotExistFuncLock(ctx context.Context, key string, f Func, force bool, timeout ...time.Duration) (ok bool, err error) {
	if _, isExist := mc.cache.Get(key); isExist {
		return
	}
	// 使用 singleflight 确保函数只执行一次
	var val any
	if val, err, _ = mc.group.Do(key, func() (v any, e error) {
		// double-check：再次检查是否已存在
		if _, isExist := mc.cache.Get(key); isExist {
			return
		}
		// 只执行函数，不设置缓存
		return f(ctx)
	}); err != nil {
		return
	}
	// 再次检查是否已存在（可能在等待 singleflight 期间被设置）
	if _, isExist := mc.cache.Get(key); isExist {
		return
	}
	if utils.IsNil(val) && !force {
		return
	}
	// 此处不判断`val == nil`是因为防止缓存穿透
	return mc.SetIfNotExist(ctx, key, val, timeout...)
}

// Update 当`key`存在时，则使用`val`更新`key`的值，返回`key`的旧值
//
//	当`timeout > 0`且`key`更新成功时，更新`key`的过期时间
func (mc *MemoryCache) Update(ctx context.Context, key string, val any, timeout ...time.Duration) (oldVal any, isExist bool, err error) {
	if oldVal, isExist = mc.cache.Get(key); !isExist {
		return
	}

	if e := mc.cache.Replace(key, val, mc.getTimeout(timeout...)); e != nil {
		isExist = false
		return
	}
	return
}

// UpdateExpire 当`key`存在时，则更新`key`的过期时间，返回`key`的旧的过期时间值
//
//	当`key`不存在时，则返回-1
//	当`key`存在但没有设置过期时间时，则返回0
//	当`key`存在且设置了过期时间时，则返回过期时间
//	当`timeout > 0`且`key`存在时，更新`key`的过期时间
func (mc *MemoryCache) UpdateExpire(ctx context.Context, key string, timeout time.Duration) (oldTimeout time.Duration, err error) {
	if timeout <= 0 {
		err = errors.New("timeout must be greater than 0")
		return
	}
	// 使用 GetWithExpiration 获取值和过期时间
	var (
		val        any
		expiration time.Time
		isExist    bool
	)
	if val, expiration, isExist = mc.cache.GetWithExpiration(key); !isExist {
		oldTimeout = -1
		return
	}
	// 计算旧的剩余过期时间
	if expiration.IsZero() {
		oldTimeout = 0
		if e := mc.cache.Replace(key, val, timeout); e != nil {
			oldTimeout = -1
			return
		}
		return
	}

	oldTimeout = time.Until(expiration)
	if oldTimeout <= 0 {
		// 已经过期，删除并返回 -1
		mc.cache.Delete(key)
		oldTimeout = -1
		return
	}
	if e := mc.cache.Replace(key, val, timeout); e != nil {
		oldTimeout = -1
		return
	}
	return
}

// IsExist 缓存是否存在
func (mc *MemoryCache) IsExist(ctx context.Context, key string) (isExist bool, err error) {
	_, isExist = mc.cache.Get(key)
	return
}

// Size 缓存中的key数量
func (mc *MemoryCache) Size(ctx context.Context) (size int, err error) {
	return mc.cache.ItemCount(), nil
}

// Delete 删除缓存
func (mc *MemoryCache) Delete(ctx context.Context, keys ...string) (err error) {
	if len(keys) == 0 {
		return
	}

	for _, key := range keys {
		mc.cache.Delete(key)
	}
	return
}

// GetExpire 获取缓存`key`的过期时间
//
//	当`key`不存在时，则返回-1
//	当`key`存在但没有设置过期时间时，则返回0
//	当`key`存在且设置了过期时间时，则返回剩余过期时间
func (mc *MemoryCache) GetExpire(ctx context.Context, key string) (timeout time.Duration, err error) {
	// 使用 GetWithExpiration 获取值和过期时间
	var (
		expiration time.Time
		isExist    bool
	)
	if _, expiration, isExist = mc.cache.GetWithExpiration(key); !isExist {
		timeout = -1
		return
	}
	// 计算旧的剩余过期时间
	if expiration.IsZero() {
		return
	}
	timeout = time.Until(expiration)
	if timeout <= 0 {
		// 已经过期，删除并返回 -1
		mc.cache.Delete(key)
		timeout = -1
		return
	}
	return
}

// Close 关闭缓存服务
func (mc *MemoryCache) Close(ctx context.Context) (err error) {
	if mc.closed {
		return
	}

	mc.cache.Flush()
	mc.closed = true
	return
}

// Save 保存缓存到 writer
func (mc *MemoryCache) Save(w io.Writer) (err error) {
	return mc.cache.Save(w)
}

// Load 从 reader 加载缓存
func (mc *MemoryCache) Load(r io.Reader) (err error) {
	return mc.cache.Load(r)
}

// SaveToFile 保存缓存到文件
func (mc *MemoryCache) SaveToFile(filename string) (err error) {
	return mc.cache.SaveFile(filename)
}

// LoadFromFile 从文件加载缓存
func (mc *MemoryCache) LoadFromFile(filename string) (err error) {
	return mc.cache.LoadFile(filename)
}

// getTimeout 获取过期时间
func (mc *MemoryCache) getTimeout(timeout ...time.Duration) (duration time.Duration) {
	duration = gocache.NoExpiration
	if len(timeout) > 0 && timeout[0] > 0 {
		duration = timeout[0]
	}
	return
}

// generateSingleflightKey 生成 singleflight 的唯一 key
func (mc *MemoryCache) generateSingleflightKey(keys []string, args []any) (uniqueKey string, err error) {
	var (
		data = map[string]any{
			"keys": keys,
			"args": args,
		}
		jsonBytes []byte
	)
	if jsonBytes, err = json.Marshal(data); err != nil {
		return
	}
	uniqueKey = string(jsonBytes)
	return
}
