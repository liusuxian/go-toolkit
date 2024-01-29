/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-29 16:15:07
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-29 17:52:06
 * @Description: 适配 github.com/silenceper/wechat/v2 库的缓存
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkcache

import (
	"context"
	"github.com/liusuxian/go-toolkit/gtkconv"
	"github.com/liusuxian/go-toolkit/gtkredis"
	"time"
)

// WechatCache 微信缓存
type WechatCache struct {
	ctx    context.Context
	client *gtkredis.RedisClient // redis 客户端
}

// NewWechatCache 创建微信缓存
func NewWechatCache(ctx context.Context, opts ...gtkredis.ClientConfigOption) (wc *WechatCache) {
	wc = &WechatCache{
		ctx:    ctx,
		client: gtkredis.NewClient(ctx, opts...),
	}
	return
}

// SetCtx 设置 ctx 参数
func (wc *WechatCache) SetCtx(ctx context.Context) {
	wc.ctx = ctx
}

// Get 获取缓存
func (wc *WechatCache) Get(key string) (val any) {
	return wc.GetContext(wc.ctx, key)
}

// GetContext 获取缓存
func (wc *WechatCache) GetContext(ctx context.Context, key string) (val any) {
	var err error
	if val, err = wc.client.Do(ctx, "GET", key); err != nil {
		return nil
	}
	return
}

// Set 设置缓存
func (wc *WechatCache) Set(key string, val any, timeout time.Duration) (err error) {
	return wc.SetContext(wc.ctx, key, val, timeout)
}

// SetContext 设置缓存
func (wc *WechatCache) SetContext(ctx context.Context, key string, val any, timeout time.Duration) (err error) {
	if int64(timeout.Seconds()) <= 0 {
		_, err = wc.client.Do(ctx, "SET", key, val)
	} else {
		_, err = wc.client.Do(ctx, "SETEX", key, int64(timeout.Seconds()), val)
	}
	return
}

// IsExist 缓存是否存在
func (wc *WechatCache) IsExist(key string) (isExist bool) {
	return wc.IsExistContext(wc.ctx, key)
}

// IsExistContext 判断key是否存在
func (wc *WechatCache) IsExistContext(ctx context.Context, key string) (isExist bool) {
	var (
		val any
		err error
	)
	if val, err = wc.client.Do(ctx, "EXISTS", key); err != nil {
		return
	}
	isExist = gtkconv.ToBool(val)
	return
}

// Delete 删除
func (wc *WechatCache) Delete(key string) (err error) {
	return wc.DeleteContext(wc.ctx, key)
}

// DeleteContext 删除缓存
func (wc *WechatCache) DeleteContext(ctx context.Context, key string) (err error) {
	_, err = wc.client.Do(ctx, "DEL", key)
	return
}
