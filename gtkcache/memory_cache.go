/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-12-16 23:11:19
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2026-01-16 13:01:34
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkcache

import (
	"context"
	"github.com/liusuxian/go-toolkit/internal/utils"
	"golang.org/x/sync/singleflight"
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
	group     singleflight.Group          // 用于防止缓存击穿，确保相同 key 的函数只执行一次
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
//	注意：如需为每个`key`设置/重置不同的过期时间，请使用`BatchGet`
func (mc *memoryCache) GetMap(ctx context.Context, keys []string, timeout ...time.Duration) (data map[string]any, err error) {
	dataMap := make(map[string]any)
	if len(keys) == 0 {
		return dataMap, nil
	}

	// 批量获取缓存并刷新过期时间
	if expiration := getExpiration(timeout...); expiration > 0 {
		mc.mu.Lock()
		defer mc.mu.Unlock()

		var (
			allHit = true
			mcMap  = make(map[string]*Item)
		)
		for _, key := range keys {
			item, found := mc.items[key]
			if !found || item.isExpired() {
				dataMap[key] = nil
				allHit = false
				continue
			}
			dataMap[key] = item.Object
			mcMap[key] = item
		}
		if allHit {
			for _, item := range mcMap {
				item.Expiration = expiration
			}
		}
		return dataMap, nil
	}
	// 批量获取缓存
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	for _, key := range keys {
		item, found := mc.items[key]
		if !found || item.isExpired() {
			dataMap[key] = nil
			continue
		}
		dataMap[key] = item.Object
	}
	return dataMap, nil
}

// BatchGet 批量获取缓存
//
//	支持为每个`key`设置/重置不同的过期时间
//	当所有`key`使用相同过期时间时，可以使用更简洁的`GetMap`方法
//	defaultTimeout: 可选参数，设置默认过期时间（对所有未单独设置过期时间的 key 生效）
//	当`timeout > 0`且缓存命中时，所有未单独指定过期时间的`key`将使用此默认过期时间
//	当`timeout <= 0`时，所有未单独指定过期时间的`key`将保持原有的过期时间
func (mc *memoryCache) BatchGet(ctx context.Context, fn func(add func(key string, timeout ...time.Duration)), defaultTimeout ...time.Duration) (values map[string]any, err error) {
	// 声明为接口类型
	var getter IBatchGetter = &memoryBatchGetter{
		mc:    mc,
		items: make([]batchGetItem, 0),
	}
	// 设置默认过期时间（对所有未单独设置过期时间的 key 生效）
	if len(defaultTimeout) > 0 && defaultTimeout[0] > 0 {
		getter = getter.SetDefaultTimeout(ctx, defaultTimeout[0])
	}
	// 添加一个 key 到批量获取队列
	addFunc := func(key string, timeout ...time.Duration) {
		getter = getter.Add(ctx, key, timeout...)
	}
	// 执行用户提供的函数，将数据添加到批量设置队列中
	fn(addFunc)
	// 执行批量设置操作
	return getter.Execute(ctx)
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
//	注意：使用`singleflight`机制确保相同`key`的函数`f`只执行一次，其他并发请求等待并共享第一个请求的执行结果，有效防止缓存击穿
func (mc *memoryCache) GetOrSetFunc(ctx context.Context, key string, f Func, force bool, timeout ...time.Duration) (val any, err error) {
	// 获取缓存
	oldVal, err := mc.Get(ctx, key, timeout...)
	if err != nil {
		return nil, err
	}
	if oldVal != nil {
		return oldVal, nil
	}
	// 使用 singleflight 确保函数只执行一次
	result, err, _ := mc.group.Do(key, func() (any, error) {
		// 获取缓存（double-check）
		cVal, err := mc.Get(ctx, key, timeout...)
		if err != nil {
			return nil, err
		}
		if cVal != nil {
			return singleflightValue{val: cVal, fromCache: true}, nil
		}
		fVal, err := f(ctx)
		if err != nil {
			return nil, err
		}
		return singleflightValue{val: fVal, fromCache: false}, nil
	})
	if err != nil {
		return nil, err
	}
	sfVal := result.(singleflightValue)
	if sfVal.fromCache {
		return sfVal.val, nil
	}
	if utils.IsNil(sfVal.val) && !force {
		return nil, nil
	}
	// 添加缓存
	_, newVal := mc.Add(key, sfVal.val, timeout...)
	return newVal, nil
}

// CustomGetOrSetFunc 从缓存中获取指定键`keys`的值，如果缓存未命中，则使用函数`f`的结果设置`keys`的值
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
//	当`force = true`时，可防止缓存穿透（即使`f`返回`nil`也会缓存）
//	注意：使用`singleflight`机制确保相同`key`的函数`f`只执行一次，其他并发请求等待并共享第一个请求的执行结果，有效防止缓存击穿
func (mc *memoryCache) CustomGetOrSetFunc(ctx context.Context, keys []string, args []any, cc ICustomCache, f Func, force bool, timeout ...time.Duration) (val any, err error) {
	// 获取缓存
	oldVal, err := cc.Get(ctx, keys, args, timeout...)
	if err != nil {
		return nil, err
	}
	if oldVal != nil {
		return oldVal, nil
	}
	// 生成 singleflight 的唯一 key
	sfKey, err := generateSingleflightKey(keys, args)
	if err != nil {
		return nil, err
	}
	// 使用 singleflight 确保函数只执行一次
	result, err, _ := mc.group.Do(sfKey, func() (any, error) {
		// 获取缓存（double-check）
		cVal, err := cc.Get(ctx, keys, args, timeout...)
		if err != nil {
			return nil, err
		}
		if cVal != nil {
			return singleflightValue{val: cVal, fromCache: true}, nil
		}
		fVal, err := f(ctx)
		if err != nil {
			return nil, err
		}
		return singleflightValue{val: fVal, fromCache: false}, nil
	})
	if err != nil {
		return nil, err
	}
	sfVal := result.(singleflightValue)
	if sfVal.fromCache {
		return sfVal.val, nil
	}
	if utils.IsNil(sfVal.val) && !force {
		return nil, nil
	}
	// 添加缓存
	return cc.Add(ctx, keys, args, sfVal.val, timeout...)
}

// Set 设置缓存
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
func (mc *memoryCache) Set(ctx context.Context, key string, val any, timeout ...time.Duration) (err error) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	var expiration int64
	if len(timeout) > 0 && timeout[0] > 0 {
		// 设置新的过期时间
		expiration = time.Now().Add(timeout[0]).UnixNano()
	} else {
		// 保持原有的过期时间（KEEPTTL 行为）
		oldItem, found := mc.items[key]
		if found && !oldItem.isExpired() {
			expiration = oldItem.Expiration
		}
	}

	mc.items[key] = &Item{
		Object:     val,
		Expiration: expiration,
	}
	return nil
}

// SetMap 批量设置缓存，所有`key`的过期时间相同
//
//	当`timeout > 0`时，设置/重置所有`key`的过期时间，所有`key`过期时间相同
//	注意：如需为每个`key`设置不同的过期时间，请使用`BatchSet`
func (mc *memoryCache) SetMap(ctx context.Context, data map[string]any, timeout ...time.Duration) (err error) {
	if len(data) == 0 {
		return nil
	}

	now := time.Now()
	mc.mu.Lock()
	defer mc.mu.Unlock()

	for key, val := range data {
		var expiration int64
		if len(timeout) > 0 && timeout[0] > 0 {
			// 设置新的过期时间
			expiration = now.Add(timeout[0]).UnixNano()
		} else {
			// 保持原有的过期时间（KEEPTTL 行为）
			oldItem, found := mc.items[key]
			if found && !oldItem.isExpired() {
				expiration = oldItem.Expiration
			}
		}

		mc.items[key] = &Item{
			Object:     val,
			Expiration: expiration,
		}
	}
	return nil
}

// BatchSet 批量设置缓存
//
//	支持为每个`key`设置不同的过期时间
//	当所有`key`使用相同过期时间时，可以使用更简洁的`SetMap`方法
//	defaultTimeout: 可选参数，设置默认过期时间（对所有未单独设置过期时间的 key 生效）
//	当`defaultTimeout > 0`时，所有未单独指定过期时间的`key`将使用此默认过期时间
//	当`defaultTimeout <= 0`时，所有未单独指定过期时间的`key`将保持原有的过期时间
func (mc *memoryCache) BatchSet(ctx context.Context, fn func(add func(key string, val any, timeout ...time.Duration)), defaultTimeout ...time.Duration) (err error) {
	// 声明为接口类型
	var setter IBatchSetter = &memoryBatchSetter{
		mc:    mc,
		items: make([]batchSetItem, 0),
	}
	// 设置默认过期时间（对所有未单独设置过期时间的 key 生效）
	if len(defaultTimeout) > 0 && defaultTimeout[0] > 0 {
		setter = setter.SetDefaultTimeout(ctx, defaultTimeout[0])
	}
	// 添加一个 key-value 对到批量设置队列
	addFunc := func(key string, val any, timeout ...time.Duration) {
		setter = setter.Add(ctx, key, val, timeout...)
	}
	// 执行用户提供的函数，将数据添加到批量设置队列中
	fn(addFunc)
	// 执行批量设置操作
	return setter.Execute(ctx)
}

// SetIfNotExist 当`key`不存在时，则使用`val`设置`key`的值，返回是否设置成功
//
//	当`timeout > 0`且`key`设置成功时，设置`key`的过期时间
func (mc *memoryCache) SetIfNotExist(ctx context.Context, key string, val any, timeout ...time.Duration) (ok bool, err error) {
	// 添加缓存
	isSuccess, _ := mc.Add(key, val, timeout...)
	return isSuccess, nil
}

// SetIfNotExistFunc 当`key`不存在时，则使用函数`f`的结果设置`key`的值，返回是否设置成功
//
//	当`timeout > 0`且`key`设置成功时，设置`key`的过期时间
//	当`force = true`时，可防止缓存穿透（即使`f`返回`nil`也会缓存）
//	注意：使用`singleflight`机制确保相同`key`的函数`f`只执行一次，其他并发请求等待并共享第一个请求的执行结果，有效防止缓存击穿
func (mc *memoryCache) SetIfNotExistFunc(ctx context.Context, key string, f Func, force bool, timeout ...time.Duration) (ok bool, err error) {
	// 缓存是否存在
	isExist, err := mc.IsExist(ctx, key)
	if err != nil {
		return false, err
	}
	if isExist {
		return false, nil
	}
	// 使用 singleflight 确保函数只执行一次
	var result any
	if result, err, _ = mc.group.Do(key, func() (v any, e error) {
		// 缓存是否存在（double-check）
		cIsExist, err := mc.IsExist(ctx, key)
		if err != nil {
			return nil, err
		}
		if cIsExist {
			return singleflightValue{val: nil, fromCache: true}, nil
		}
		// 执行函数获取新值
		fVal, err := f(ctx)
		if err != nil {
			return nil, err
		}
		return singleflightValue{val: fVal, fromCache: false}, nil
	}); err != nil {
		return
	}
	sfVal := result.(singleflightValue)
	if sfVal.fromCache {
		return false, nil
	}
	if utils.IsNil(sfVal.val) && !force {
		return false, nil
	}
	// 添加缓存
	isSuccess, _ := mc.Add(key, sfVal.val, timeout...)
	return isSuccess, nil
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
		return nil
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
	return nil
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
	return nil
}

// OnEvicted 设置删除回调函数
func (mc *memoryCache) OnEvicted(f func(key string, value any)) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.onEvicted = f
}

// Add 添加缓存
//
//	如果`key`已存在且未过期，则返回现有值和 false（表示添加失败）
//	如果`key`不存在或已过期，则添加新值并返回该值和 true（表示添加成功）
func (mc *memoryCache) Add(key string, val any, timeout ...time.Duration) (isSuccess bool, result any) {
	expiration := getExpiration(timeout...)
	mc.mu.Lock()
	defer mc.mu.Unlock()

	item, found := mc.items[key]
	if found && !item.isExpired() {
		return false, item.Object
	}

	mc.items[key] = &Item{
		Object:     val,
		Expiration: expiration,
	}
	return true, val
}

// Items 获取所有未过期的缓存项
func (mc *memoryCache) Items() (items map[string]Item) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	items = make(map[string]Item, len(mc.items))
	now := time.Now().UnixNano()
	for k, v := range mc.items {
		if v.Expiration > 0 && now > v.Expiration {
			continue
		}
		items[k] = *v
	}
	return
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
