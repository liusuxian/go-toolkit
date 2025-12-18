/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-12-16 23:11:19
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-18 19:14:56
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkcache

import (
	"context"
	"runtime"
	"sync"
	"time"
)

// Item 缓存项
type Item struct {
	Object     any   // 缓存的值
	Expiration int64 // 过期时间（Unix纳秒时间戳），0 表示永不过期
}

// isExpired 检查是否过期
func (item *Item) isExpired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

// keyAndValue 键值对
type keyAndValue struct {
	key   string
	value any
}

// MemoryCache 内存缓存
type MemoryCache struct {
	*memoryCache
}

// memoryCache 内存缓存
type memoryCache struct {
	items     map[string]*Item
	mu        sync.RWMutex
	onEvicted func(key string, value any) // 删除回调函数
	janitor   *janitor                    // 清理器
}

// NewMemoryCache 创建内存缓存
func NewMemoryCache(cleanupInterval ...time.Duration) *MemoryCache {
	items := make(map[string]*Item)
	return newMemoryCacheWithJanitor(items, cleanupInterval...)
}

// NewMemoryCacheFrom 从已有缓存创建内存缓存
func NewMemoryCacheFrom(items map[string]*Item, cleanupInterval ...time.Duration) *MemoryCache {
	return newMemoryCacheWithJanitor(items, cleanupInterval...)
}

// Get 获取缓存
//
//	当`timeout > 0`且缓存命中时，设置/重置`key`的过期时间
func (mc *memoryCache) Get(ctx context.Context, key string, timeout ...time.Duration) (val any, err error) {
	// 获取缓存并刷新过期时间
	if expiration := getExpiration(timeout...); expiration > 0 {
		mc.mu.Lock()
		defer mc.mu.Unlock()

		item, found := mc.items[key]
		if !found || item.isExpired() {
			return nil, nil
		}

		item.Expiration = expiration
		return item.Object, nil
	}
	// 获取缓存
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	item, found := mc.items[key]
	if !found || item.isExpired() {
		return nil, nil
	}
	return item.Object, nil
}

// GetMap 批量获取缓存
//
//	当`timeout > 0`且所有缓存都命中时，设置/重置所有`key`的过期时间，所有`key`过期时间相同
func (mc *memoryCache) GetMap(ctx context.Context, keys []string, timeout ...time.Duration) (data map[string]any, err error) {
	if len(keys) == 0 {
		return
	}

	return
}

// GetOrSet 检索并返回`key`的值，或者当`key`不存在时，则使用`newVal`设置`key`的值
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
func (mc *memoryCache) GetOrSet(ctx context.Context, key string, newVal any, timeout ...time.Duration) (val any, err error) {
	expiration := getExpiration(timeout...)
	mc.mu.Lock()
	defer mc.mu.Unlock()

	item, found := mc.items[key]
	if found && !item.isExpired() {
		if expiration > 0 {
			item.Expiration = expiration
		}
		return item.Object, nil
	}

	mc.items[key] = &Item{
		Object:     newVal,
		Expiration: expiration,
	}
	return newVal, nil
}

// GetOrSetFunc 检索并返回`key`的值，或者当`key`不存在时，则使用函数`f`的结果设置`key`的值
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
//	当`force = true`时，可防止缓存穿透（即使`f`返回`nil`也会缓存）
//	注意：高并发时函数f会被多次执行，可能导致缓存击穿，如需防止请使用 GetOrSetFuncLock
func (mc *memoryCache) GetOrSetFunc(ctx context.Context, key string, f Func, force bool, timeout ...time.Duration) (val any, err error) {
	return
}

// GetOrSetFuncLock 检索并返回`key`的值，或者当`key`不存在时，则使用函数`f`的结果设置`key`的值
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
//	当`force = true`时，可防止缓存穿透（即使`f`返回`nil`也会缓存）
//	注意：使用`singleflight`机制确保相同`key`的函数`f`只执行一次，其他并发请求等待并共享第一个请求的执行结果，有效防止缓存击穿
func (mc *memoryCache) GetOrSetFuncLock(ctx context.Context, key string, f Func, force bool, timeout ...time.Duration) (val any, err error) {
	return
}

// CustomGetOrSetFunc 从缓存中获取指定键`keys`的值，如果缓存未命中，则使用函数`f`的结果设置`keys`的值
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
//	当`force = true`时，可防止缓存穿透（即使`f`返回`nil`也会缓存）
//	注意：高并发时函数f会被多次执行，可能导致缓存击穿，如需防止请使用 CustomGetOrSetFuncLock
func (mc *memoryCache) CustomGetOrSetFunc(ctx context.Context, keys []string, args []any, cc ICustomCache, f Func, force bool, timeout ...time.Duration) (val any, err error) {
	return
}

// CustomGetOrSetFuncLock 从缓存中获取指定键`keys`的值，如果缓存未命中，则使用函数`f`的结果设置`keys`的值
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
//	当`force = true`时，可防止缓存穿透（即使`f`返回`nil`也会缓存）
//	注意：使用`singleflight`机制确保相同`key`的函数`f`只执行一次，其他并发请求等待并共享第一个请求的执行结果，有效防止缓存击穿
func (mc *memoryCache) CustomGetOrSetFuncLock(ctx context.Context, keys []string, args []any, cc ICustomCache, f Func, force bool, timeout ...time.Duration) (val any, err error) {
	return
}

// Set 设置缓存
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
func (mc *memoryCache) Set(ctx context.Context, key string, val any, timeout ...time.Duration) (err error) {
	expiration := getExpiration(timeout...)
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.items[key] = &Item{
		Object:     val,
		Expiration: expiration,
	}
	return
}

// SetMap 批量设置缓存，所有`key`的过期时间相同
//
//	当`timeout > 0`时，设置/重置所有`key`的过期时间，所有`key`过期时间相同
func (mc *memoryCache) SetMap(ctx context.Context, data map[string]any, timeout ...time.Duration) (err error) {
	if len(data) == 0 {
		return
	}

	return
}

// SetIfNotExist 当`key`不存在时，则使用`val`设置`key`的值，返回是否设置成功
//
//	当`timeout > 0`且`key`设置成功时，设置`key`的过期时间
func (mc *memoryCache) SetIfNotExist(ctx context.Context, key string, val any, timeout ...time.Duration) (ok bool, err error) {
	return
}

// SetIfNotExistFunc 当`key`不存在时，则使用函数`f`的结果设置`key`的值，返回是否设置成功
//
//	当`timeout > 0`且`key`设置成功时，设置`key`的过期时间
//	当`force = true`时，可防止缓存穿透（即使`f`返回`nil`也会缓存）
//	注意：高并发时函数f会被多次执行，可能导致缓存击穿，如需防止请使用 SetIfNotExistFuncLock
func (mc *memoryCache) SetIfNotExistFunc(ctx context.Context, key string, f Func, force bool, timeout ...time.Duration) (ok bool, err error) {
	return
}

// SetIfNotExistFuncLock 当`key`不存在时，则使用函数`f`的结果设置`key`的值，返回是否设置成功，函数`f`是在写入互斥锁内执行，以确保并发安全
//
//	当`timeout > 0`且`key`设置成功时，设置`key`的过期时间
//	当`force = true`时，可防止缓存穿透（即使`f`返回`nil`也会缓存）
func (mc *memoryCache) SetIfNotExistFuncLock(ctx context.Context, key string, f Func, force bool, timeout ...time.Duration) (ok bool, err error) {
	return
}

// Update 当`key`存在时，则使用`val`更新`key`的值，返回`key`的旧值
//
//	当`timeout > 0`且`key`更新成功时，更新`key`的过期时间
func (mc *memoryCache) Update(ctx context.Context, key string, val any, timeout ...time.Duration) (oldVal any, isExist bool, err error) {
	expiration := getExpiration(timeout...)
	mc.mu.Lock()
	defer mc.mu.Unlock()

	item, found := mc.items[key]
	if !found || item.isExpired() {
		return nil, false, nil
	}

	ov := item.Object
	item.Object = val
	if expiration > 0 {
		item.Expiration = expiration
	}
	return ov, true, nil
}

// UpdateExpire 当`key`存在时，则更新`key`的过期时间，返回`key`的旧的过期时间值
//
//	当`key`不存在时，则返回-1
//	当`key`存在但没有设置过期时间时，则返回0
//	当`key`存在且设置了过期时间时，则返回过期时间
//	当`timeout > 0`且`key`存在时，更新`key`的过期时间
func (mc *memoryCache) UpdateExpire(ctx context.Context, key string, timeout time.Duration) (oldTimeout time.Duration, err error) {
	expiration := getExpiration(timeout)
	mc.mu.Lock()
	defer mc.mu.Unlock()

	item, found := mc.items[key]
	if !found || item.isExpired() {
		return -1, nil
	}

	if item.Expiration == 0 {
		if expiration > 0 {
			item.Expiration = expiration
		}
		return 0, nil
	}

	oldExpiration := item.Expiration
	if expiration > 0 {
		item.Expiration = expiration
	}
	return time.Duration(oldExpiration - time.Now().UnixNano()), nil
}

// IsExist 缓存是否存在
func (mc *memoryCache) IsExist(ctx context.Context, key string) (isExist bool, err error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	item, found := mc.items[key]
	return found && !item.isExpired(), nil
}

// Size 缓存中的key数量
func (mc *memoryCache) Size(ctx context.Context) (size int, err error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	count := 0
	for _, v := range mc.items {
		if !v.isExpired() {
			count++
		}
	}
	return count, nil
}

// Delete 删除缓存
func (mc *memoryCache) Delete(ctx context.Context, keys ...string) (err error) {
	if len(keys) == 0 {
		return
	}

	var evictedItems []keyAndValue
	mc.mu.Lock()
	for _, k := range keys {
		v, evicted := mc.delete(k)
		if evicted {
			evictedItems = append(evictedItems, keyAndValue{k, v})
		}
	}
	mc.mu.Unlock()
	for _, v := range evictedItems {
		mc.onEvicted(v.key, v.value)
	}
	return
}

// GetExpire 获取缓存`key`的过期时间
//
//	当`key`不存在时，则返回-1
//	当`key`存在但没有设置过期时间时，则返回0
//	当`key`存在且设置了过期时间时，则返回过期时间
func (mc *memoryCache) GetExpire(ctx context.Context, key string) (timeout time.Duration, err error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	item, found := mc.items[key]
	if !found || item.isExpired() {
		return -1, nil
	}

	if item.Expiration == 0 {
		return 0, nil
	}

	return time.Duration(item.Expiration - time.Now().UnixNano()), nil
}

// Close 关闭缓存服务
func (mc *memoryCache) Close(ctx context.Context) (err error) {
	mc.Flush() // 清空缓存
	if mc.janitor != nil {
		mc.janitor.stop <- true // 停止清理器
	}
	return
}

// OnEvicted 设置删除回调函数
func (mc *memoryCache) OnEvicted(f func(key string, value any)) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.onEvicted = f
}

// Flush 清空缓存
func (mc *memoryCache) Flush() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.items = make(map[string]*Item)
}

// DeleteExpired 删除过期缓存
func (mc *memoryCache) DeleteExpired() {
	var (
		evictedItems []keyAndValue
		now          = time.Now().UnixNano()
	)
	mc.mu.Lock()
	for k, v := range mc.items {
		if v.Expiration > 0 && now > v.Expiration {
			ov, evicted := mc.delete(k)
			if evicted {
				evictedItems = append(evictedItems, keyAndValue{k, ov})
			}
		}
	}
	mc.mu.Unlock()
	for _, v := range evictedItems {
		mc.onEvicted(v.key, v.value)
	}
}

// delete 删除缓存
func (mc *memoryCache) delete(key string) (any, bool) {
	if mc.onEvicted != nil {
		if v, found := mc.items[key]; found {
			delete(mc.items, key)
			return v.Object, true
		}
	}
	delete(mc.items, key)
	return nil, false
}

// newMemoryCache 创建内存缓存
func newMemoryCache(items map[string]*Item) *memoryCache {
	return &memoryCache{
		items: items,
	}
}

// newMemoryCacheWithJanitor 创建内存缓存并启动清理器
func newMemoryCacheWithJanitor(items map[string]*Item, cleanupInterval ...time.Duration) *MemoryCache {
	var (
		mc = newMemoryCache(items)
		MC = &MemoryCache{mc}
	)
	if len(cleanupInterval) > 0 && cleanupInterval[0] > 0 {
		runJanitor(mc, cleanupInterval[0])
		runtime.SetFinalizer(MC, stopJanitor)
	}
	return MC
}

// getExpiration 获取过期时间戳（Unix纳秒时间戳），0 表示永不过期
func getExpiration(timeout ...time.Duration) int64 {
	if len(timeout) > 0 && timeout[0] > 0 {
		return time.Now().Add(timeout[0]).UnixNano()
	}
	return 0
}

// janitor 清理器
type janitor struct {
	interval time.Duration
	stop     chan bool
}

// Run 启动清理任务
func (j *janitor) Run(mc *memoryCache) {
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mc.DeleteExpired()
		case <-j.stop:
			return
		}
	}
}

// stopJanitor 停止清理任务
func stopJanitor(mc *MemoryCache) {
	mc.janitor.stop <- true
}

// runJanitor 启动清理器
func runJanitor(mc *memoryCache, interval time.Duration) {
	j := &janitor{
		interval: interval,
		stop:     make(chan bool, 1),
	}
	mc.janitor = j
	go j.Run(mc)
}
