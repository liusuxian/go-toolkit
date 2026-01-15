/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-12-20 00:15:44
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2026-01-15 17:07:15
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkcache

import (
	"encoding/json"
	"time"
)

// batchGetItem 批量获取项
type batchGetItem struct {
	key     string
	timeout *time.Duration
}

// batchSetItem 批量设置项
type batchSetItem struct {
	key     string
	val     any
	timeout *time.Duration
}

// singleflightValue singleflight 的值
type singleflightValue struct {
	val       any
	fromCache bool
}

// getExpiration 获取过期时间戳（Unix纳秒时间戳），0 表示永不过期
func getExpiration(timeout ...time.Duration) int64 {
	if len(timeout) > 0 && timeout[0] > 0 {
		return time.Now().Add(timeout[0]).UnixNano()
	}
	return 0
}

// generateSingleflightKey 生成 singleflight 的唯一 key
func generateSingleflightKey(keys []string, args []any) (uniqueKey string, err error) {
	var (
		data = map[string]any{
			"keys": keys,
			"args": args,
		}
		jsonBytes []byte
	)
	if jsonBytes, err = json.Marshal(data); err != nil {
		return
	}
	uniqueKey = string(jsonBytes)
	return
}
