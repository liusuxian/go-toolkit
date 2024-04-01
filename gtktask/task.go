/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-04-01 13:15:12
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-04-01 23:15:48
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtktask

import (
	"context"
	"github.com/pkg/errors"
	"sync"
	"time"
)

// SPolling 轮询对象
type SPolling struct {
	list  []bool // 数据切片，false:当前下标不可用 true:当前下标可用
	index uint   // 当前下标位置
	lock  sync.Mutex
}

// RetryFunc 重试函数的类型
type RetryFunc func(ctx context.Context) (err error)

// NewPolling 新建轮询
func NewPolling(list []bool, startIndex ...uint) (s *SPolling, err error) {
	if len(list) == 0 {
		err = errors.Errorf("list must not be empty")
		return
	}
	s = &SPolling{
		list: make([]bool, len(list)),
	}
	copy(s.list, list)
	if len(startIndex) > 0 {
		s.index = startIndex[0]
	}
	return
}

// SetIsAvailable 设置下标是否可用
func (s *SPolling) SetIsAvailable(index uint, isAvailable bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	total := uint(len(s.list))
	if total > 0 && index <= total-1 {
		s.list[index] = isAvailable
	}
}

// Polling 轮询
func (s *SPolling) Polling() (index uint, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	total := uint(len(s.list))
	hasAvailable := false
	for i := s.index; i < total; i++ {
		if s.list[i] {
			index = i
			s.index = i + 1
			hasAvailable = true
			break
		}
	}
	// 如果切片中没有可用元素，直接返回错误
	if !hasAvailable {
		err = errors.New("no available index found")
		return
	}
	// 重置 index 为切片的起始位置
	if s.index >= total {
		s.index = 0
	}
	return
}

// Retry 重试
//
//	f: 要执行的函数
//	maxRetries: 最大重试次数（包含首次尝试）
//	delay: 默认重试之间的延迟时间。当配置了`delayList`时，该参数将失效
//	increaseDelay: 是否让延迟时间随着重试次数增加而线性增加。当配置了`delayList`时，该参数将失效
//	delayList: 自定义延迟列表
func Retry(ctx context.Context, f RetryFunc, maxRetries uint, delay time.Duration, increaseDelay bool, delayList ...time.Duration) (err error) {
	for i := uint(0); i < maxRetries; i++ {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			return
		default:
		}

		if err = f(ctx); err != nil {
			// 检查是否已经是最后一次尝试
			if i == maxRetries-1 {
				return
			}

			if len(delayList) > 0 && uint(len(delayList)) > i {
				// 使用自定义延迟列表
				time.Sleep(delayList[i])
			} else if len(delayList) > 0 {
				// 延迟列表长度不够时返回错误
				err = errors.Errorf("not enough delay values provided, required: %d", maxRetries-1)
				return
			} else {
				if increaseDelay {
					// 重试延迟随重试次数线性增加
					time.Sleep(delay * time.Duration(i+1))
				} else {
					// 每次重试的延迟时间保持不变
					time.Sleep(delay)
				}
			}
			continue
		}
		return
	}
	return
}

// GetGoroutinesAndTasks 根据任务总数计算所需协程数和每个协程处理的任务数量
//
//	total: 任务总数
//	expected: 每个协程期望处理的任务数量
func GetGoroutinesAndTasks(total, expected uint) (goroutineNum uint, tasks []uint) {
	// 如果总任务数小于或等于每个协程预期的任务数，则只需要一个协程来处理所有任务
	if total <= expected {
		return 1, []uint{total}
	}
	// 根据总任务数除以每个协程预期的任务数来计算需要的协程数
	// 这里添加`expected - 1`是为了在除法中实现向上取整，确保即使有余数也能分配足够的协程
	goroutineNum = (total + expected - 1) / expected
	// 初始化一个切片来保存每个协程将要处理的任务数
	tasks = make([]uint, goroutineNum)
	// 尽可能均匀地将任务分配给每个协程
	for k := range tasks {
		tasks[k] = total / goroutineNum
	}
	// 如果有余数，则将剩余的任务逐一分配给部分协程，以保证所有任务都能被处理
	remainder := total % goroutineNum
	for i := uint(0); i < remainder; i++ {
		tasks[i]++
	}
	return
}
