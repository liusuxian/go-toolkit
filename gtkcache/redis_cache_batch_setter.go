/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2026-01-15 15:07:27
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2026-01-16 14:25:49
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

// redisBatchSetter Redis 批量设置构建器
type redisBatchSetter struct {
	rc             *RedisCache
	items          []batchSetItem
	defaultTimeout *time.Duration
}

// Add 添加一个 key-value 对到批量设置队列
//
//	当`timeout > 0`时，设置该`key`的过期时间（会覆盖默认过期时间）
//	当`timeout`未指定或`timeout <= 0`时，该`key`使用`SetDefaultTimeout`设置的默认过期时间
//	返回构建器自身，支持链式调用
func (bs *redisBatchSetter) Add(ctx context.Context, key string, val any, timeout ...time.Duration) (batchSetter IBatchSetter) {
	item := batchSetItem{key: key, val: val}
	if len(timeout) > 0 && timeout[0] > 0 {
		item.timeout = &timeout[0]
	}
	bs.items = append(bs.items, item)
	return bs
}

// SetDefaultTimeout 设置默认过期时间（对所有未单独设置过期时间的 key 生效）
//
//	当`timeout > 0`时，所有未单独指定过期时间的`key`将使用此默认过期时间
//	当`timeout <= 0`时，所有未单独指定过期时间的`key`将保持原有的过期时间
//	返回构建器自身，支持链式调用
func (bs *redisBatchSetter) SetDefaultTimeout(ctx context.Context, timeout time.Duration) (batchSetter IBatchSetter) {
	if timeout > 0 {
		bs.defaultTimeout = &timeout
	} else {
		bs.defaultTimeout = nil
	}
	return bs
}

// Execute 执行批量设置操作
func (bs *redisBatchSetter) Execute(ctx context.Context) (err error) {
	if len(bs.items) == 0 {
		err = fmt.Errorf("no items to execute")
		return
	}

	var (
		keys = make([]string, 0, len(bs.items))
		args = make([]any, 0, len(bs.items)*2)
	)
	for _, item := range bs.items {
		keys = append(keys, item.key)
		args = append(args, item.val)
		// 确定过期时间
		timeout := item.timeout
		if timeout == nil {
			timeout = bs.defaultTimeout
		}
		if timeout != nil && *timeout > 0 {
			args = append(args, timeout.Milliseconds())
		} else {
			args = append(args, 0) // 保持原有的过期时间
		}
	}
	// 执行批量设置操作
	_, err = bs.rc.client.EvalSha(ctx, "BATCH_SET_EX", keys, args...)
	return
}
