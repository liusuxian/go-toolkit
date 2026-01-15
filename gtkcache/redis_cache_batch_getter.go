/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2026-01-15 18:59:18
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2026-01-15 19:43:44
 * @Description:
 *
 * Copyright (c) 2026 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkcache

import (
	"context"
	"fmt"
	"github.com/liusuxian/go-toolkit/gtkconv"
	"time"
)

// redisBatchGetter Redis 批量获取构建器
type redisBatchGetter struct {
	rc             *RedisCache
	items          []batchGetItem
	defaultTimeout *time.Duration
}

// Add 添加一个 key 到批量获取队列
//
//	当`timeout > 0`且该`key`缓存命中时，设置/重置该`key`的过期时间（会覆盖默认过期时间）
//	当`timeout`未指定或`timeout <= 0`时，该`key`使用`SetDefaultTimeout`设置的默认过期时间
//	返回构建器自身，支持链式调用
func (bg *redisBatchGetter) Add(ctx context.Context, key string, timeout ...time.Duration) (batchGetter IBatchGetter) {
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
func (bg *redisBatchGetter) SetDefaultTimeout(ctx context.Context, timeout time.Duration) (batchGetter IBatchGetter) {
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
//	执行成功后，自动清空构建器中的数据（不建议继续使用该构建器）
//	执行失败时，保留构建器中的数据，可以直接再次调用本方法进行重试
//	建议：为每次批量操作创建新的构建器实例
func (bg *redisBatchGetter) Execute(ctx context.Context) (values map[string]any, err error) {
	if len(bg.items) == 0 {
		err = fmt.Errorf("no items to execute")
		return
	}

	var (
		keys = make([]string, 0, len(bg.items))
		args = make([]any, 0, len(bg.items))
	)
	for _, item := range bg.items {
		keys = append(keys, item.key)
		// 确定过期时间
		timeout := item.timeout
		if timeout == nil {
			timeout = bg.defaultTimeout
		}
		if timeout != nil && *timeout > 0 {
			args = append(args, timeout.Milliseconds())
		} else {
			args = append(args, 0) // 保持原有的过期时间
		}
	}

	var result any
	if result, err = bg.rc.client.EvalSha(ctx, "BATCH_GET_EX", keys, args); err != nil {
		return
	}

	for i := range bg.items {
		bg.items[i] = batchGetItem{}
	}
	bg.items = nil
	bg.defaultTimeout = nil

	values, err = gtkconv.ToStringMapE(result)
	return
}
