/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-04-01 13:15:12
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-04-01 17:52:18
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtktask_test

import (
	"context"
	"github.com/liusuxian/go-toolkit/gtktask"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestPolling(t *testing.T) {
	var (
		assert       = assert.New(t)
		total   uint = 5
		polling      = gtktask.NewPolling()
	)
	for i := 0; i < int(total); i++ {
		index := polling.Polling(total)
		assert.Equal(uint(i), index)
	}
	for i := 0; i < int(total); i++ {
		index := polling.Polling(total)
		assert.Equal(uint(i), index)
	}
}

func TestRetry(t *testing.T) {
	var (
		ctx    = context.Background()
		assert = assert.New(t)
		err    error
	)
	err = gtktask.Retry(ctx, func(ctx context.Context) (err error) {
		t.Logf("Retry1: %+v\n", time.Now())
		return errors.New("test error1")
	}, 3, time.Second, false)
	assert.Error(err)
	t.Logf("\n")

	err = gtktask.Retry(ctx, func(ctx context.Context) (err error) {
		t.Logf("Retry2: %+v\n", time.Now())
		return errors.New("test error2")
	}, 3, time.Second, true)
	assert.Error(err)
	t.Logf("\n")

	err = gtktask.Retry(ctx, func(ctx context.Context) (err error) {
		t.Logf("Retry3: %+v\n", time.Now())
		return errors.New("test error3")
	}, 3, 0, false, time.Second, time.Second*3)
	assert.Error(err)
	t.Logf("\n")

	// 创建一个可取消的上下文
	var wg sync.WaitGroup
	wg.Add(2)
	cancelCtx, cancelFun := context.WithCancel(ctx)
	go func() {
		defer wg.Done()
		err = gtktask.Retry(cancelCtx, func(ctx context.Context) (err error) {
			t.Logf("Retry4: %+v\n", time.Now())
			return errors.New("test error4")
		}, 3, 0, false, time.Second, time.Second*3)
		t.Logf("err: %v\n", err)
	}()
	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 2)
		cancelFun()
	}()
	wg.Wait()
}

func TestGetGoroutinesAndTasks(t *testing.T) {
	var (
		assert       = assert.New(t)
		goroutineNum uint
		tasks        []uint
	)
	goroutineNum, tasks = gtktask.GetGoroutinesAndTasks(1, 2)
	assert.Equal(uint(1), goroutineNum)
	assert.Equal([]uint{1}, tasks)
	goroutineNum, tasks = gtktask.GetGoroutinesAndTasks(2, 2)
	assert.Equal(uint(1), goroutineNum)
	assert.Equal([]uint{2}, tasks)
	goroutineNum, tasks = gtktask.GetGoroutinesAndTasks(3, 2)
	assert.Equal(uint(2), goroutineNum)
	assert.Equal([]uint{2, 1}, tasks)
	goroutineNum, tasks = gtktask.GetGoroutinesAndTasks(4, 2)
	assert.Equal(uint(2), goroutineNum)
	assert.Equal([]uint{2, 2}, tasks)
	goroutineNum, tasks = gtktask.GetGoroutinesAndTasks(5, 2)
	assert.Equal(uint(3), goroutineNum)
	assert.Equal([]uint{2, 2, 1}, tasks)
	goroutineNum, tasks = gtktask.GetGoroutinesAndTasks(6, 2)
	assert.Equal(uint(3), goroutineNum)
	assert.Equal([]uint{2, 2, 2}, tasks)
	goroutineNum, tasks = gtktask.GetGoroutinesAndTasks(7, 2)
	assert.Equal(uint(4), goroutineNum)
	assert.Equal([]uint{2, 2, 2, 1}, tasks)
	goroutineNum, tasks = gtktask.GetGoroutinesAndTasks(8, 2)
	assert.Equal(uint(4), goroutineNum)
	assert.Equal([]uint{2, 2, 2, 2}, tasks)
	goroutineNum, tasks = gtktask.GetGoroutinesAndTasks(9, 2)
	assert.Equal(uint(5), goroutineNum)
	assert.Equal([]uint{2, 2, 2, 2, 1}, tasks)

	goroutineNum, tasks = gtktask.GetGoroutinesAndTasks(1, 3)
	assert.Equal(uint(1), goroutineNum)
	assert.Equal([]uint{1}, tasks)
	goroutineNum, tasks = gtktask.GetGoroutinesAndTasks(2, 3)
	assert.Equal(uint(1), goroutineNum)
	assert.Equal([]uint{2}, tasks)
	goroutineNum, tasks = gtktask.GetGoroutinesAndTasks(3, 3)
	assert.Equal(uint(1), goroutineNum)
	assert.Equal([]uint{3}, tasks)
	goroutineNum, tasks = gtktask.GetGoroutinesAndTasks(4, 3)
	assert.Equal(uint(2), goroutineNum)
	assert.Equal([]uint{2, 2}, tasks)
	goroutineNum, tasks = gtktask.GetGoroutinesAndTasks(5, 3)
	assert.Equal(uint(2), goroutineNum)
	assert.Equal([]uint{3, 2}, tasks)
	goroutineNum, tasks = gtktask.GetGoroutinesAndTasks(6, 3)
	assert.Equal(uint(2), goroutineNum)
	assert.Equal([]uint{3, 3}, tasks)
	goroutineNum, tasks = gtktask.GetGoroutinesAndTasks(7, 3)
	assert.Equal(uint(3), goroutineNum)
	assert.Equal([]uint{3, 2, 2}, tasks)
	goroutineNum, tasks = gtktask.GetGoroutinesAndTasks(8, 3)
	assert.Equal(uint(3), goroutineNum)
	assert.Equal([]uint{3, 3, 2}, tasks)
	goroutineNum, tasks = gtktask.GetGoroutinesAndTasks(9, 3)
	assert.Equal(uint(3), goroutineNum)
	assert.Equal([]uint{3, 3, 3}, tasks)
}
