/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2026-01-26 18:56:00
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2026-01-26 19:04:33
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
	if err := mutex.Lock(); err != nil {
		t.Fatalf("Failed to acquire lock: %v", err)
	}
	t.Log("Lock acquired")
	// 再次加锁应该失败（锁已被持有）
	mutex2 := client.NewMutex("test:redsync:basic")
	if err := mutex2.Lock(); err == nil {
		t.Fatal("Expected lock to fail, but it succeeded")
	}
	t.Log("Second lock attempt failed as expected")
	// 解锁
	if ok, err := mutex.Unlock(); !ok || err != nil {
		t.Fatalf("Failed to unlock: %v", err)
	}
	t.Log("Lock released")
	// 解锁后再次加锁应该成功
	if err := mutex2.Lock(); err != nil {
		t.Fatalf("Failed to acquire lock after unlock: %v", err)
	}
	t.Log("Lock re-acquired successfully")
	mutex2.Unlock()
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
	if err := mutex.Lock(); err != nil {
		t.Fatalf("Failed to acquire lock: %v", err)
	}
	t.Log("Lock acquired")

	// 获取初始 TTL
	ttl1 := mutex.Until()
	t.Logf("Initial TTL: %v\n", ttl1)

	// 等待 1 秒
	time.Sleep(1 * time.Second)

	// 续期
	if ok, err := mutex.Extend(); !ok || err != nil {
		t.Fatalf("Failed to extend lock: %v", err)
	}
	t.Log("Lock extended")

	// 获取续期后的 TTL
	ttl2 := mutex.Until()
	t.Logf("TTL after extend: %v\n", ttl2)

	// ttl2 应该在 ttl1 之后（过期时间被延长了）
	if !ttl2.After(ttl1) {
		t.Fatalf("TTL should increase after extend, got %v -> %v", ttl1, ttl2)
	}
	mutex.Unlock()
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

	lockName := "test:redsync:concurrent"
	goroutines := 10
	counter := int32(0)

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
			if err := mutex.Lock(); err != nil {
				t.Errorf("Goroutine %d failed to acquire lock: %v", id, err)
				return
			}
			t.Logf("Goroutine %d acquired lock\n", id)
			// 模拟临界区操作
			current := atomic.LoadInt32(&counter)
			time.Sleep(50 * time.Millisecond) // 模拟耗时操作
			atomic.StoreInt32(&counter, current+1)
			t.Logf("Goroutine %d releasing lock (counter: %d)\n", id, counter)
			// 释放锁
			if ok, err := mutex.Unlock(); !ok || err != nil {
				t.Errorf("Goroutine %d failed to unlock: %v", id, err)
			}
		}(i)
	}
	wg.Wait()
	// 验证计数器
	if counter != int32(goroutines) {
		t.Fatalf("Expected counter to be %d, got %d", goroutines, counter)
	}
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
	// 第一个 goroutine 持有锁
	if err := mutex1.Lock(); err != nil {
		t.Fatalf("Failed to acquire lock: %v", err)
	}
	t.Log("Lock acquired by first goroutine")
	// 第二个 goroutine 尝试获取锁，但会因为 context 取消而中止
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	mutex2 := client.NewMutex("test:redsync:context", redsync.WithExpiry(10*time.Second))

	start := time.Now()
	err = mutex2.LockContext(ctx)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("Expected lock to fail due to context cancellation")
	}
	t.Logf("Lock failed as expected after %v: %v\n", elapsed, err)
	// 第一个 goroutine 释放锁
	mutex1.Unlock()
	t.Log("Lock released by first goroutine")
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
	if err := mutex1.Lock(); err != nil {
		t.Fatalf("Failed to acquire lock1: %v", err)
	}
	if err := mutex2.Lock(); err != nil {
		t.Fatalf("Failed to acquire lock2: %v", err)
	}
	if err := mutex3.Lock(); err != nil {
		t.Fatalf("Failed to acquire lock3: %v", err)
	}
	t.Log("All three locks acquired")
	// 释放所有锁
	mutex1.Unlock()
	mutex2.Unlock()
	mutex3.Unlock()
	t.Log("All three locks released")
}
