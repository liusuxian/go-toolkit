/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-12-20 01:28:44
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-20 02:37:34
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkcache_test

import (
	"context"
	"fmt"
	"github.com/liusuxian/go-toolkit/gtkcache"
	"github.com/liusuxian/go-toolkit/gtkconv"
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type BBB struct {
	A     int
	B     float64
	C     string
	cache gtkcache.ICache
}

func (b *BBB) Get(ctx context.Context, keys []string, args []any, timeout ...time.Duration) (val any, err error) {
	return
}

func (b *BBB) Add(ctx context.Context, keys []string, args []any, newVal any, timeout ...time.Duration) (val any, err error) {
	// 模拟批量设置两个 key
	if len(keys) >= 2 {
		_ = b.cache.Set(ctx, keys[0], 1000, timeout...)
		_ = b.cache.Set(ctx, keys[1], 2000, timeout...)
		val = []any{1000, 2000}
	}
	return
}

func TestMemoryCacheString(t *testing.T) {
	var (
		ctx    = context.Background()
		assert = assert.New(t)
		cache  gtkcache.ICache
	)
	cache = gtkcache.NewMemoryCache()
	assert.NotNil(cache)

	var (
		val     any
		isExist bool
		ok      bool
		timeout time.Duration
	)
	val, err := cache.Get(ctx, "test_key_1", time.Second*10)
	assert.NoError(err)
	assert.Nil(val)
	isExist, err = cache.IsExist(ctx, "test_key_1")
	assert.NoError(err)
	assert.False(isExist)
	timeout, err = cache.GetExpire(ctx, "test_key_1")
	assert.NoError(err)
	assert.Equal(time.Duration(-1), timeout)

	err = cache.Set(ctx, "test_key_2", nil, time.Second)
	assert.NoError(err)
	val, err = cache.GetOrSet(ctx, "test_key_2", 200, time.Second*2)
	assert.NoError(err)
	// MemoryCache 存储原值，nil 返回 nil
	assert.Nil(val)
	a1 := BBB{}
	gtkconv.ToStruct(val, &a1)
	assert.Equal(BBB{A: 0, B: 0, C: ""}, a1)
	timeout, err = cache.GetExpire(ctx, "test_key_2")
	assert.NoError(err)
	// 允许微秒级误差
	assert.InDelta(int64(time.Second*2), int64(timeout), float64(time.Millisecond))
	val, err = cache.GetOrSet(ctx, "test_key_22", map[string]any{"a": 1}, time.Second)
	assert.NoError(err)
	// MemoryCache 存储原对象，不是 JSON 字符串
	assert.Equal(map[string]any{"a": 1}, val)
	val, err = cache.GetOrSetFunc(ctx, "test_key_3", func(ctx context.Context) (val any, err error) {
		return
	}, true, time.Second)
	assert.NoError(err)
	// MemoryCache 中 nil 返回 nil
	assert.Nil(val)
	err = cache.SetMap(ctx, map[string]any{"a": 1, "b": map[string]any{"b": 100}}, time.Second)
	assert.NoError(err)
	ok, err = cache.SetIfNotExist(ctx, "test_key_3", 100, time.Second)
	assert.NoError(err)
	assert.False(ok)
	ok, err = cache.SetIfNotExist(ctx, "test_key_4", nil, time.Second)
	assert.NoError(err)
	assert.True(ok)
	ok, err = cache.SetIfNotExistFunc(ctx, "test_key_4", func(ctx context.Context) (val any, err error) {
		return
	}, true, time.Second)
	assert.NoError(err)
	assert.False(ok)
	ok, err = cache.SetIfNotExistFunc(ctx, "test_key_5", func(ctx context.Context) (val any, err error) {
		return
	}, true, time.Second)
	assert.NoError(err)
	assert.True(ok)

	val, err = cache.CustomGetOrSetFunc(ctx, []string{"test_key_10", "test_key_11"}, []any{}, &BBB{cache: cache}, func(ctx context.Context) (val any, err error) {
		return map[string]any{
			"test_key_10": 1,
			"test_key_11": 2,
		}, nil
	}, true, time.Second)
	assert.NoError(err)
	assert.Equal([]any{1000, 2000}, val)
}

func TestMemoryCacheString2(t *testing.T) {
	var (
		ctx    = context.Background()
		assert = assert.New(t)
		cache  gtkcache.ICache
	)
	cache = gtkcache.NewMemoryCache()
	assert.NotNil(cache)

	var (
		val     any
		isExist bool
		timeout time.Duration
	)
	err := cache.Set(ctx, "test_key_1", 100, time.Second*10)
	assert.NoError(err)
	val, err = cache.Get(ctx, "test_key_1")
	assert.NoError(err)
	assert.Equal(100, gtkconv.ToInt(val))
	isExist, err = cache.IsExist(ctx, "test_key_1")
	assert.NoError(err)
	assert.True(isExist)
	timeout, err = cache.GetExpire(ctx, "test_key_1")
	assert.NoError(err)
	assert.InDelta(int64(time.Second*10), int64(timeout), float64(time.Millisecond))
	val, isExist, err = cache.Update(ctx, "test_key_1", 200)
	assert.NoError(err)
	assert.True(isExist)
	assert.Equal(100, gtkconv.ToInt(val))
	timeout, err = cache.GetExpire(ctx, "test_key_1")
	assert.NoError(err)
	assert.InDelta(int64(time.Second*10), int64(timeout), float64(time.Millisecond))
	val, isExist, err = cache.Update(ctx, "test_key_1", 300, time.Second*20)
	assert.NoError(err)
	assert.True(isExist)
	assert.Equal(200, gtkconv.ToInt(val))
	timeout, err = cache.GetExpire(ctx, "test_key_1")
	assert.NoError(err)
	assert.InDelta(int64(time.Second*20), int64(timeout), float64(time.Millisecond))
	val, isExist, err = cache.Update(ctx, "test_key_2", 300, time.Second*20)
	assert.NoError(err)
	assert.False(isExist)
	assert.Nil(val)
	timeout, err = cache.GetExpire(ctx, "test_key_2")
	assert.NoError(err)
	assert.Equal(time.Duration(-1), timeout)
	// MemoryCache 的 UpdateExpire 对 timeout=0 不报错，只是不更新
	timeout, err = cache.UpdateExpire(ctx, "test_key_2", 0)
	assert.NoError(err)
	assert.Equal(time.Duration(-1), timeout)
	timeout, err = cache.UpdateExpire(ctx, "test_key_2", time.Second*10)
	assert.NoError(err)
	assert.Equal(time.Duration(-1), timeout)
	// MemoryCache 对存在的 key 调用 UpdateExpire(0) 不报错，保持原过期时间
	timeout, err = cache.UpdateExpire(ctx, "test_key_1", 0)
	assert.NoError(err)
	assert.InDelta(int64(time.Second*20), int64(timeout), float64(time.Millisecond))
	timeout, err = cache.UpdateExpire(ctx, "test_key_1", time.Second*30)
	assert.NoError(err)
	assert.InDelta(int64(time.Second*20), int64(timeout), float64(time.Millisecond))
	timeout, err = cache.UpdateExpire(ctx, "test_key_1", time.Second*40)
	assert.NoError(err)
	assert.InDelta(int64(time.Second*30), int64(timeout), float64(time.Millisecond))
}

func TestMemoryCacheString3(t *testing.T) {
	var (
		ctx    = context.Background()
		assert = assert.New(t)
		cache  gtkcache.ICache
	)
	cache = gtkcache.NewMemoryCache()
	assert.NotNil(cache)

	var (
		isExist bool
		timeout time.Duration
	)
	err := cache.Set(ctx, "test_key_1", 200)
	assert.NoError(err)
	timeout, err = cache.GetExpire(ctx, "test_key_1")
	assert.NoError(err)
	assert.Equal(time.Duration(0), timeout)
	err = cache.Delete(ctx, "test_key_1", "test_key_2")
	assert.NoError(err)
	isExist, err = cache.IsExist(ctx, "test_key_1")
	assert.NoError(err)
	assert.False(isExist)
}

func TestMemoryCacheString4(t *testing.T) {
	var (
		ctx    = context.Background()
		assert = assert.New(t)
		cache  gtkcache.ICache
	)
	cache = gtkcache.NewMemoryCache()
	assert.NotNil(cache)

	var (
		val     any
		timeout time.Duration
		data    map[string]any
	)
	err := cache.Set(ctx, "test_key_1", 100, time.Second*10)
	assert.NoError(err)
	data, err = cache.GetMap(ctx, []string{})
	assert.NoError(err)
	assert.Equal(map[string]any{}, data)
	data, err = cache.GetMap(ctx, []string{"test_key_1", "test_key_2"}, 0)
	assert.NoError(err)
	assert.Equal(map[string]any{"test_key_1": 100, "test_key_2": nil}, data)
	data, err = cache.GetMap(ctx, []string{"test_key_1", "test_key_2"})
	assert.NoError(err)
	assert.Equal(map[string]any{"test_key_1": 100, "test_key_2": nil}, data)
	data, err = cache.GetMap(ctx, []string{"test_key_1", "test_key_2"}, time.Second*20)
	assert.NoError(err)
	assert.Equal(map[string]any{"test_key_1": 100, "test_key_2": nil}, data)
	timeout, err = cache.GetExpire(ctx, "test_key_1")
	assert.NoError(err)
	assert.InDelta(int64(time.Second*10), int64(timeout), float64(time.Millisecond))
	err = cache.SetMap(ctx, map[string]any{"test_key_1": 100, "test_key_2": 200}, time.Second*5)
	assert.NoError(err)
	data, err = cache.GetMap(ctx, []string{"test_key_1", "test_key_2"})
	assert.NoError(err)
	assert.Equal(map[string]any{"test_key_1": 100, "test_key_2": 200}, data)
	timeout, err = cache.GetExpire(ctx, "test_key_1")
	assert.NoError(err)
	assert.InDelta(int64(time.Second*5), int64(timeout), float64(time.Millisecond))
	timeout, err = cache.GetExpire(ctx, "test_key_2")
	assert.NoError(err)
	assert.InDelta(int64(time.Second*5), int64(timeout), float64(time.Millisecond))
	val, err = cache.Get(ctx, "test_key_1", time.Second*10)
	assert.NoError(err)
	assert.Equal(100, gtkconv.ToInt(val))
	timeout, err = cache.GetExpire(ctx, "test_key_1")
	assert.NoError(err)
	assert.InDelta(int64(time.Second*10), int64(timeout), float64(time.Millisecond))
	data, err = cache.GetMap(ctx, []string{"test_key_1", "test_key_2"}, time.Second*60)
	assert.NoError(err)
	assert.Equal(map[string]any{"test_key_1": 100, "test_key_2": 200}, data)
	timeout, err = cache.GetExpire(ctx, "test_key_1")
	assert.NoError(err)
	assert.InDelta(int64(time.Second*60), int64(timeout), float64(time.Millisecond))
	timeout, err = cache.GetExpire(ctx, "test_key_2")
	assert.NoError(err)
	assert.InDelta(int64(time.Second*60), int64(timeout), float64(time.Millisecond))
}

func TestMemoryCacheGetOrSetFunc(t *testing.T) {
	var (
		ctx    = context.Background()
		assert = assert.New(t)
		cache  gtkcache.ICache
	)
	cache = gtkcache.NewMemoryCache()
	assert.NotNil(cache)

	var (
		executeCount int32
		wg           sync.WaitGroup
		concurrency  = 1000
		key          = "test_concurrent_key"
	)
	// 模拟1000个并发请求同时访问同一个不存在的key
	for i := range concurrency {
		index := i
		wg.Go(func() {
			val, err := cache.GetOrSetFunc(ctx, key, func(ctx context.Context) (any, error) {
				atomic.AddInt32(&executeCount, 1)
				time.Sleep(10 * time.Millisecond)
				return fmt.Sprintf("value_%d", index), nil
			}, false, time.Second*10)
			assert.NoError(err)
			assert.NotNil(val)
		})
	}
	wg.Wait()
	finalCount := atomic.LoadInt32(&executeCount)
	assert.Equal(int32(1), finalCount, "函数应该只执行1次，实际执行了%d次", finalCount)
}

func TestMemoryCacheSetIfNotExistFunc(t *testing.T) {
	var (
		ctx    = context.Background()
		assert = assert.New(t)
		cache  gtkcache.ICache
	)
	cache = gtkcache.NewMemoryCache()
	assert.NotNil(cache)

	var (
		executeCount int32
		wg           sync.WaitGroup
		concurrency  = 1000
		key          = "test_setif_key"
		successCount int32
	)
	// 模拟1000个并发请求同时尝试设置同一个key
	for i := range concurrency {
		index := i
		wg.Go(func() {
			ok, err := cache.SetIfNotExistFunc(ctx, key, func(ctx context.Context) (any, error) {
				atomic.AddInt32(&executeCount, 1)
				time.Sleep(10 * time.Millisecond)
				return fmt.Sprintf("value_%d", index), nil
			}, false, time.Second*10)
			assert.NoError(err)
			if ok {
				atomic.AddInt32(&successCount, 1)
			}
		})
	}
	wg.Wait()
	finalCount := atomic.LoadInt32(&executeCount)
	finalSuccess := atomic.LoadInt32(&successCount)
	assert.Equal(int32(1), finalCount, "函数应该只执行1次，实际执行了%d次", finalCount)
	assert.Equal(int32(1), finalSuccess, "应该只有1个请求设置成功，实际成功了%d次", finalSuccess)
}

type SimpleMemoryCustomCache struct {
	cache gtkcache.ICache
}

func (s *SimpleMemoryCustomCache) Get(ctx context.Context, keys []string, args []any, timeout ...time.Duration) (val any, err error) {
	// 简单实现：使用第一个 key
	if len(keys) > 0 {
		if memoryCache, ok := s.cache.(*gtkcache.MemoryCache); ok {
			return memoryCache.Get(ctx, keys[0], timeout...)
		}
		return s.cache.Get(ctx, keys[0], timeout...)
	}
	return
}

func (s *SimpleMemoryCustomCache) Add(ctx context.Context, keys []string, args []any, newVal any, timeout ...time.Duration) (val any, err error) {
	// 简单实现：使用第一个 key 设置值
	if len(keys) > 0 {
		if memoryCache, ok := s.cache.(*gtkcache.MemoryCache); ok {
			_, val = memoryCache.Add(keys[0], newVal, timeout...)
			return
		}
		err = s.cache.Set(ctx, keys[0], newVal, timeout...)
		return newVal, err
	}
	return
}

func TestMemoryCacheCustomGetOrSetFunc(t *testing.T) {
	var (
		ctx    = context.Background()
		assert = assert.New(t)
		cache  gtkcache.ICache
	)
	cache = gtkcache.NewMemoryCache()
	assert.NotNil(cache)

	var (
		executeCount int32
		wg           sync.WaitGroup
		concurrency  = 1000
		keys         = []string{"test_custom_key"}
		args         = []any{"arg1"}
	)
	customCache := &SimpleMemoryCustomCache{cache: cache}
	// 模拟1000个并发请求同时访问同一个自定义缓存key
	for i := range concurrency {
		index := i
		wg.Go(func() {
			val, err := cache.CustomGetOrSetFunc(ctx, keys, args, customCache, func(ctx context.Context) (any, error) {
				atomic.AddInt32(&executeCount, 1)
				time.Sleep(10 * time.Millisecond)
				return fmt.Sprintf("custom_value_%d", index), nil
			}, false, time.Second*10)
			assert.NoError(err)
			assert.NotNil(val)
		})
	}
	wg.Wait()
	finalCount := atomic.LoadInt32(&executeCount)
	assert.Equal(int32(1), finalCount, "函数应该只执行1次，实际执行了%d次", finalCount)
}
