/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2026-01-26 18:56:00
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2026-01-27 10:54:23
 * @Description:
 *
 * Copyright (c) 2026 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkredis_test

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redsync/redsync/v4"
	"github.com/liusuxian/go-toolkit/gtkredis"
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 测试基本的加锁和解锁
func TestRedsyncBasicLockUnlock(t *testing.T) {
	var (
		ctx    = context.Background()
		r      = miniredis.RunT(t)
		assert = assert.New(t)
	)
	client, err := gtkredis.NewClient(ctx, &gtkredis.ClientConfig{
		Addr:     r.Addr(),
		Username: "default",
		Password: "",
		DB:       1,
	})
	assert.NoError(err)
	defer client.Close()
	// 创建锁
	mutex := client.NewMutex("test:redsync:basic", redsync.WithExpiry(5*time.Second))
	// 第一次加锁应该成功
	err = mutex.Lock()
	assert.NoError(err)
	// 再次加锁应该失败（锁已被持有）
	mutex2 := client.NewMutex("test:redsync:basic")
	err = mutex2.Lock()
	assert.Error(err)
	// 解锁
	ok, err := mutex.Unlock()
	assert.True(ok)
	assert.NoError(err)
	// 解锁后再次加锁应该成功
	err = mutex2.Lock()
	assert.NoError(err)
	ok, err = mutex2.Unlock()
	assert.True(ok)
	assert.NoError(err)
}

// 测试锁续期
func TestRedsyncExtend(t *testing.T) {
	var (
		ctx    = context.Background()
		r      = miniredis.RunT(t)
		assert = assert.New(t)
	)
	client, err := gtkredis.NewClient(ctx, &gtkredis.ClientConfig{
		Addr:     r.Addr(),
		Username: "default",
		Password: "",
		DB:       1,
	})
	assert.NoError(err)
	defer client.Close()

	mutex := client.NewMutex("test:redsync:extend", redsync.WithExpiry(2*time.Second))
	// 加锁
	err = mutex.Lock()
	assert.NoError(err)
	// 获取初始 TTL
	ttl1 := mutex.Until()
	assert.InDelta(int64(2*time.Second), int64(time.Until(ttl1)), float64(100*time.Millisecond))
	// 等待 1 秒
	time.Sleep(1 * time.Second)
	// 续期
	ok, err := mutex.Extend()
	assert.True(ok)
	assert.NoError(err)
	// 获取续期后的 TTL
	ttl2 := mutex.Until()
	assert.InDelta(int64(2*time.Second), int64(time.Until(ttl2)), float64(100*time.Millisecond))
	// ttl2 应该在 ttl1 之后（过期时间被延长了）
	assert.True(ttl2.After(ttl1), "TTL should increase after extend, got %v -> %v", ttl1, ttl2)
	ok, err = mutex.Unlock()
	assert.True(ok)
	assert.NoError(err)
}

// 测试并发场景下的锁
func TestRedsyncConcurrency(t *testing.T) {
	var (
		ctx    = context.Background()
		r      = miniredis.RunT(t)
		assert = assert.New(t)
	)
	client, err := gtkredis.NewClient(ctx, &gtkredis.ClientConfig{
		Addr:     r.Addr(),
		Username: "default",
		Password: "",
		DB:       1,
	})
	assert.NoError(err)
	defer client.Close()

	var (
		lockName   = "test:redsync:concurrent"
		goroutines = 10
		counter    = int32(0)
	)

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := range goroutines {
		go func(id int) {
			defer wg.Done()

			mutex := client.NewMutex(
				lockName,
				redsync.WithExpiry(5*time.Second),
				redsync.WithTries(32),
				redsync.WithRetryDelay(100*time.Millisecond),
			)
			// 阻塞式获取锁
			err = mutex.Lock()
			assert.NoError(err)
			t.Logf("Goroutine %d acquired lock\n", id)
			// 模拟临界区操作
			current := atomic.LoadInt32(&counter)
			time.Sleep(50 * time.Millisecond) // 模拟耗时操作
			atomic.StoreInt32(&counter, current+1)
			t.Logf("Goroutine %d releasing lock (counter: %d)\n", id, counter)
			// 释放锁
			ok, err := mutex.Unlock()
			assert.True(ok)
			assert.NoError(err)
		}(i)
	}
	wg.Wait()
	// 验证计数器
	assert.False(counter != int32(goroutines))
	t.Logf("Final counter: %d (expected: %d)\n", counter, goroutines)
}

// 测试使用 context 取消
func TestRedsyncContext(t *testing.T) {
	var (
		ctx    = context.Background()
		r      = miniredis.RunT(t)
		assert = assert.New(t)
	)
	client, err := gtkredis.NewClient(ctx, &gtkredis.ClientConfig{
		Addr:     r.Addr(),
		Username: "default",
		Password: "",
		DB:       1,
	})
	assert.NoError(err)
	defer client.Close()

	mutex1 := client.NewMutex("test:redsync:context", redsync.WithExpiry(10*time.Second))
	// 首先持有锁
	err = mutex1.Lock()
	assert.NoError(err)
	// 尝试用另一个 mutex 获取同一把锁，但会因为 context 超时而失败
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	mutex2 := client.NewMutex("test:redsync:context", redsync.WithExpiry(10*time.Second))
	err = mutex2.LockContext(ctx)
	assert.Error(err)
	// 释放第一个锁
	ok, err := mutex1.Unlock()
	assert.True(ok)
	assert.NoError(err)
}

// 测试 Redsync 实例的多个锁
func TestRedsyncMultipleLocks(t *testing.T) {
	var (
		ctx    = context.Background()
		r      = miniredis.RunT(t)
		assert = assert.New(t)
	)
	client, err := gtkredis.NewClient(ctx, &gtkredis.ClientConfig{
		Addr:     r.Addr(),
		Username: "default",
		Password: "",
		DB:       1,
	})
	assert.NoError(err)
	defer client.Close()
	// 创建 Redsync 实例
	rs := client.NewRedsync()
	// 创建多个锁
	mutex1 := rs.NewMutex("test:redsync:lock1")
	mutex2 := rs.NewMutex("test:redsync:lock2")
	mutex3 := rs.NewMutex("test:redsync:lock3")
	// 同时获取多个锁
	err = mutex1.Lock()
	assert.NoError(err)
	err = mutex2.Lock()
	assert.NoError(err)
	err = mutex3.Lock()
	assert.NoError(err)
	// 释放所有锁
	ok, err := mutex1.Unlock()
	assert.True(ok)
	assert.NoError(err)
	ok, err = mutex2.Unlock()
	assert.True(ok)
	assert.NoError(err)
	ok, err = mutex3.Unlock()
	assert.True(ok)
	assert.NoError(err)
}
