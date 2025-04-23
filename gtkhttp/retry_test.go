/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-11 11:48:30
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-23 18:51:55
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp_test

import (
	"context"
	"errors"
	utils "github.com/liusuxian/go-toolkit/gtkhttp"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	// 测试第一次就成功的情况
	t.Run("SuccessOnFirstAttempt", func(t *testing.T) {
		callCount := 0
		err := utils.Retry(context.Background(), func(ctx context.Context) error {
			callCount++
			return nil // 立即成功
		}, 3, 10*time.Millisecond, false)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if callCount != 1 {
			t.Fatalf("Expected function to be called once, got %d", callCount)
		}
	})
	// 测试需要重试几次才成功的情况
	t.Run("SuccessAfterRetries", func(t *testing.T) {
		callCount := 0
		err := utils.Retry(context.Background(), func(ctx context.Context) error {
			callCount++
			if callCount < 3 {
				return errors.New("temporary error")
			}
			return nil // 第三次调用成功
		}, 5, 10*time.Millisecond, false)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if callCount != 3 {
			t.Fatalf("Expected function to be called 3 times, got %d", callCount)
		}
	})
	// 测试超出最大重试次数仍然失败的情况
	t.Run("FailAfterMaxRetries", func(t *testing.T) {
		callCount := 0
		err := utils.Retry(context.Background(), func(ctx context.Context) error {
			callCount++
			return errors.New("persistent error")
		}, 3, 10*time.Millisecond, false)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}
		if callCount != 4 { // 初始尝试 + 3次重试
			t.Fatalf("Expected function to be called 4 times, got %d", callCount)
		}
	})
	// 测试上下文取消的情况
	t.Run("ContextCancellation", func(t *testing.T) {
		callCount := 0
		ctx, cancel := context.WithCancel(context.Background())
		// 创建一个通道来同步取消时间
		doneChan := make(chan struct{})
		go func() {
			err := utils.Retry(ctx, func(ctx context.Context) error {
				callCount++
				// 第一次调用后取消上下文
				if callCount == 1 {
					cancel()
					// 等待一点时间确保取消生效
					time.Sleep(20 * time.Millisecond)
				}
				return errors.New("error to trigger retry")
			}, 5, 100*time.Millisecond, false)
			if !errors.Is(err, context.Canceled) {
				t.Errorf("Expected context.Canceled error, got %v", err)
			}
			close(doneChan)
		}()
		// 等待测试完成或超时
		select {
		case <-doneChan:
			// 测试正常完成
		case <-time.After(500 * time.Millisecond):
			t.Fatal("Test timed out")
		}
		// 应该只调用一次，因为上下文在第一次重试前就被取消了
		if callCount != 1 {
			t.Fatalf("Expected function to be called once, got %d", callCount)
		}
	})
	// 测试自定义延迟列表
	t.Run("CustomDelayList", func(t *testing.T) {
		callCount := 0
		startTime := time.Now()
		delayList := []time.Duration{50 * time.Millisecond, 100 * time.Millisecond, 150 * time.Millisecond}
		err := utils.Retry(context.Background(), func(ctx context.Context) error {
			callCount++
			if callCount <= 3 {
				return errors.New("temporary error")
			}
			return nil // 第四次调用成功
		}, 2, 10*time.Millisecond, false, delayList...)
		elapsedTime := time.Since(startTime)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		// 检查调用次数
		if callCount != 4 { // 初始尝试 + 根据delayList的3次重试
			t.Fatalf("Expected function to be called 4 times, got %d", callCount)
		}
		// 验证总耗时，应该至少是所有延迟的总和
		// 考虑到函数执行和调度的开销，我们检查时间是否在合理范围内
		minExpectedTime := 300 * time.Millisecond // 50ms + 100ms + 150ms
		if elapsedTime < minExpectedTime {
			t.Fatalf("Expected elapsed time to be at least %v, got %v", minExpectedTime, elapsedTime)
		}
	})
	// 测试线性增加延迟
	t.Run("IncreasingDelay", func(t *testing.T) {
		callCount := 0
		startTime := time.Now()
		err := utils.Retry(context.Background(), func(ctx context.Context) error {
			callCount++
			if callCount <= 3 {
				return errors.New("temporary error")
			}
			return nil // 第四次调用成功
		}, 3, 50*time.Millisecond, true)
		elapsedTime := time.Since(startTime)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		// 检查调用次数
		if callCount != 4 { // 初始尝试 + 3次重试
			t.Fatalf("Expected function to be called 4 times, got %d", callCount)
		}
		// 验证总耗时，应该至少是所有延迟的总和 (0*50ms + 1*50ms + 2*50ms + 3*50ms)
		minExpectedTime := 300 * time.Millisecond // 0 + 50ms + 100ms + 150ms
		if elapsedTime < minExpectedTime {
			t.Fatalf("Expected elapsed time to be at least %v, got %v", minExpectedTime, elapsedTime)
		}
	})
}
