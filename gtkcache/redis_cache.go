/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-27 20:53:08
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-29 16:30:44
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkcache

import (
	"context"
	"github.com/liusuxian/go-toolkit/gtkconv"
	"github.com/liusuxian/go-toolkit/gtkredis"
	"github.com/pkg/errors"
	"time"
)

// RedisCache Redis 缓存
type RedisCache struct {
	ctx    context.Context
	client *gtkredis.RedisClient // redis 客户端
}

// 内置 lua 脚本
var internalScriptMap = map[string]string{
	"GetReset": `
	local val = redis.call('GET', KEYS[1])
	if not val then
		return nil
	end
	redis.call('EXPIRE', KEYS[1], ARGV[1])
	return val
	`,
	"GetMapReset": `
	local vals = redis.call('MGET', unpack(KEYS))
	local allKeysExist = true
	for i, val in ipairs(vals) do
		if not val then
			allKeysExist = false
			break
    end
	end
	if allKeysExist then
    for i = 1, #KEYS do
			redis.call('EXPIRE', KEYS[i], ARGV[1])
    end
	end
	return vals
	`,
	"SetMap": `
	local expireTime = ARGV[#ARGV]
	for i = 1, #KEYS do
		redis.call('SETEX', KEYS[i], expireTime, ARGV[i])
	end
	return 'OK'
	`,
}

// NewRedisCache 创建 RedisCache
func NewRedisCache(ctx context.Context, opts ...gtkredis.ClientConfigOption) (rc *RedisCache) {
	rc = &RedisCache{
		ctx:    ctx,
		client: gtkredis.NewClient(ctx, opts...),
	}
	for k, v := range internalScriptMap {
		if err := rc.client.ScriptLoad(ctx, k, v); err != nil {
			panic(err)
		}
	}
	return
}

// Client Redis 客户端
func (rc *RedisCache) Client() (client *gtkredis.RedisClient) {
	return rc.client
}

// Get 获取缓存
func (rc *RedisCache) Get(ctx context.Context, key string) (val any, err error) {
	val, err = rc.client.Do(ctx, "GET", key)
	return
}

// GetReset 获取缓存，并在缓存命中时重置过期时间（原子操作）
func (rc *RedisCache) GetReset(ctx context.Context, key string, timeout time.Duration) (val any, err error) {
	if int64(timeout.Seconds()) <= 0 {
		err = errors.New("The `timeout` parameter of type `time.Duration` must be greater than 0.")
		return
	}
	val, err = rc.client.EvalSha(ctx, "GetReset", []string{key}, int64(timeout.Seconds()))
	return
}

// GetMap 批量获取缓存（原子操作）
func (rc *RedisCache) GetMap(ctx context.Context, keys ...string) (data map[string]any, err error) {
	data = make(map[string]any)
	if len(keys) == 0 {
		return
	}
	args := make([]any, 0, len(keys))
	for _, v := range keys {
		args = append(args, v)
	}
	var result any
	if result, err = rc.client.Do(ctx, "MGET", args...); err != nil {
		return
	}
	resultList := gtkconv.ToSlice(result)
	for k, v := range keys {
		data[v] = resultList[k]
	}
	return
}

// GetMapReset 批量获取缓存，并在所有缓存都命中时重置所有缓存的过期时间（原子操作）
func (rc *RedisCache) GetMapReset(ctx context.Context, timeout time.Duration, keys ...string) (data map[string]any, err error) {
	data = make(map[string]any)
	if len(keys) == 0 {
		return
	}
	if int64(timeout.Seconds()) <= 0 {
		err = errors.New("The `timeout` parameter of type `time.Duration` must be greater than 0.")
		return
	}
	var (
		keyList = make([]string, 0, len(keys))
		result  any
	)
	keyList = append(keyList, keys...)
	if result, err = rc.client.EvalSha(ctx, "GetMapReset", keyList, int64(timeout.Seconds())); err != nil {
		return
	}
	resultList := gtkconv.ToSlice(result)
	for k, v := range keys {
		data[v] = resultList[k]
	}
	return
}

// Set 设置缓存
func (rc *RedisCache) Set(ctx context.Context, key string, val any, timeout ...time.Duration) (err error) {
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
		_, err = rc.client.Do(ctx, "SET", key, val)
	} else {
		_, err = rc.client.Do(ctx, "SETEX", key, int64(timeout[0].Seconds()), val)
	}
	return
}

// SetMap 批量设置缓存，且所有 key 的过期时间相同（原子操作）
func (rc *RedisCache) SetMap(ctx context.Context, data map[string]any, timeout ...time.Duration) (err error) {
	if len(data) == 0 {
		return
	}
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
		args := make([]any, 0, len(data)*2)
		for k, v := range data {
			args = append(args, k, v)
		}
		_, err = rc.client.Do(ctx, "MSET", args...)
	} else {
		keys := make([]string, 0, len(data))
		args := make([]any, 0, len(data)+1)
		for k, v := range data {
			keys = append(keys, k)
			args = append(args, v)
		}
		args = append(args, int64(timeout[0].Seconds()))
		_, err = rc.client.EvalSha(ctx, "SetMap", keys, args...)
	}
	return
}

// CustomCache 自定义缓存
func (rc *RedisCache) CustomCache(ctx context.Context, f Func) (val any, err error) {
	val, err = f(ctx)
	return
}

// IsExist 缓存是否存在
func (rc *RedisCache) IsExist(ctx context.Context, key string) (isExist bool, err error) {
	var val any
	if val, err = rc.client.Do(ctx, "EXISTS", key); err != nil {
		return
	}
	isExist = gtkconv.ToBool(val)
	return
}

// Delete 删除缓存
func (rc *RedisCache) Delete(ctx context.Context, keys ...string) (err error) {
	if len(keys) == 0 {
		return
	}
	args := make([]any, 0, len(keys))
	for _, v := range keys {
		args = append(args, v)
	}
	_, err = rc.client.Do(ctx, "DEL", args...)
	return
}

// GetExpire 获取缓存过期时间
func (rc *RedisCache) GetExpire(ctx context.Context, key string) (timeout time.Duration, err error) {
	var val any
	if val, err = rc.client.Do(ctx, "TTL", key); err != nil {
		return
	}
	timeout = time.Second * time.Duration(gtkconv.ToInt64(val))
	return
}

// Close 关闭缓存服务
func (rc *RedisCache) Close(ctx context.Context) (err error) {
	err = rc.client.Close()
	return
}
