/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2026-01-15 19:21:24
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2026-01-16 14:26:38
 * @Description:
 *
 * Copyright (c) 2026 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkcache

import (
	"context"
	"fmt"
	"time"
)

// memoryBatchGetter MemoryCache 批量获取构建器
type memoryBatchGetter struct {
	mc             *memoryCache
	items          []batchGetItem
	defaultTimeout *time.Duration
}

// Add 添加一个 key 到批量获取队列
//
//	当`timeout > 0`且该`key`缓存命中时，设置/重置该`key`的过期时间（会覆盖默认过期时间）
//	当`timeout`未指定或`timeout <= 0`时，该`key`使用`SetDefaultTimeout`设置的默认过期时间
//	返回构建器自身，支持链式调用
func (bg *memoryBatchGetter) Add(ctx context.Context, key string, timeout ...time.Duration) (batchGetter IBatchGetter) {
	item := batchGetItem{key: key}
	if len(timeout) > 0 && timeout[0] > 0 {
		item.timeout = &timeout[0]
	}
	bg.items = append(bg.items, item)
	return bg
}

// SetDefaultTimeout 设置默认过期时间（对所有未单独设置过期时间的 key 生效）
//
//	当`timeout > 0`且缓存命中时，所有未单独指定过期时间的`key`将使用此默认过期时间
//	当`timeout <= 0`时，所有未单独指定过期时间的`key`将保持原有的过期时间
//	返回构建器自身，支持链式调用
func (bg *memoryBatchGetter) SetDefaultTimeout(ctx context.Context, timeout time.Duration) (batchGetter IBatchGetter) {
	if timeout > 0 {
		bg.defaultTimeout = &timeout
	} else {
		bg.defaultTimeout = nil
	}
	return bg
}

// Execute 执行批量获取操作
//
//	返回 map[key]value，不存在或已过期的`key`不会出现在结果`map`中
func (bg *memoryBatchGetter) Execute(ctx context.Context) (values map[string]any, err error) {
	if len(bg.items) == 0 {
		err = fmt.Errorf("no items to execute")
		return
	}

	values = make(map[string]any)
	// 智能选择锁类型
	if bg.needsToResetExpiration() {
		bg.mc.mu.Lock()
		defer bg.mc.mu.Unlock()

		now := time.Now()
		for _, item := range bg.items {
			mcItem, found := bg.mc.items[item.key]
			if found && !mcItem.isExpired() {
				values[item.key] = mcItem.Object
				// 确定过期时间
				timeout := item.timeout
				if timeout == nil {
					timeout = bg.defaultTimeout
				}
				if timeout != nil && *timeout > 0 {
					mcItem.Expiration = now.Add(*timeout).UnixNano()
				}
			}
		}
	} else {
		bg.mc.mu.RLock()
		defer bg.mc.mu.RUnlock()

		for _, item := range bg.items {
			mcItem, found := bg.mc.items[item.key]
			if found && !mcItem.isExpired() {
				values[item.key] = mcItem.Object
			}
		}
	}
	return
}

// needsToResetExpiration 检查是否有任意一个 key 需要设置/重置过期时间
func (bg *memoryBatchGetter) needsToResetExpiration() (need bool) {
	// 检查是否有默认过期时间
	if bg.defaultTimeout != nil && *bg.defaultTimeout > 0 {
		return true
	}
	// 检查是否有任意一个 item 指定了过期时间
	for _, item := range bg.items {
		if item.timeout != nil && *item.timeout > 0 {
			return true
		}
	}
	return false
}
