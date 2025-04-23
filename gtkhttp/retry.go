/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-11 11:48:30
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-11 16:57:27
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import (
	"context"
	"time"
)

// RetryFunc 重试函数的类型
type RetryFunc func(ctx context.Context) (err error)

// Retry 重试
//
//	f: 要执行的函数
//	maxRetries: 最大重试次数。当配置了`delayList`时，该参数将失效
//	retryDelay: 默认重试之间的延迟时间。当配置了`delayList`时，该参数将失效
//	increaseDelay: 是否让延迟时间随着重试次数增加而线性增加。当配置了`delayList`时，该参数将失效
//	delayList: 自定义延迟列表
func Retry(ctx context.Context, f RetryFunc, maxRetries uint, retryDelay time.Duration, increaseDelay bool, delayList ...time.Duration) (err error) {
	if len(delayList) > 0 {
		maxRetries = uint(len(delayList))
	}
	for retry := uint(0); retry <= maxRetries; retry++ {
		if retry > 0 {
			// 重试前等待
			select {
			case <-ctx.Done():
				err = ctx.Err()
				return
			default:
				if len(delayList) > 0 {
					// 使用自定义延迟列表
					<-time.After(delayList[retry-1])
				} else {
					if increaseDelay {
						// 重试延迟随重试次数线性增加
						<-time.After(retryDelay * time.Duration(retry))
					} else {
						// 每次重试的延迟时间保持不变
						<-time.After(retryDelay)
					}
				}
			}
		}
		if err = f(ctx); err == nil {
			return
		}
	}
	return
}
