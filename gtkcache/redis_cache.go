/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-27 20:53:08
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2026-01-15 23:24:00
 * @Description: IRedisCache 接口的实现
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkcache

import (
	"context"
	"errors"
	"github.com/liusuxian/go-toolkit/gtkconv"
	"github.com/liusuxian/go-toolkit/gtkredis"
	"github.com/liusuxian/go-toolkit/internal/utils"
	"golang.org/x/sync/singleflight"
	"time"
)

// RedisCache Redis 缓存
type RedisCache struct {
	ctx    context.Context
	client *gtkredis.RedisClient // redis 客户端
	group  singleflight.Group    // 用于防止缓存击穿，确保相同 key 的函数只执行一次
}

// 内置 lua 脚本
var internalScriptMap = map[string]string{
	"ADD_EX": `
	local val = redis.call('GET', KEYS[1])
	if not val then
		if tonumber(ARGV[2], 10) > 0 then
			redis.call('PSETEX', KEYS[1], ARGV[2], ARGV[1])
		else
			redis.call('SET', KEYS[1], ARGV[1])
		end
		return ARGV[1]
	end
	return val
	`,

	"MGET_EX": `
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
			redis.call('PEXPIRE', KEYS[i], ARGV[1])
		end
	end
	return vals
	`,

	"BATCH_GET_EX": `
	-- KEYS: [key1, key2, key3, ...]
	-- ARGV: [timeout1, timeout2, timeout3, ...]
	local result = {}
	for i = 1, #KEYS do
		local val = redis.call('GET', KEYS[i])
		if val then
			result[KEYS[i]] = val
			local timeout = tonumber(ARGV[i], 10)
			if timeout > 0 then
				redis.call('PEXPIRE', KEYS[i], timeout)
			end
		end
	end
	local jsonStr = cjson.encode(result)
	if jsonStr == "[]" then
		return "{}"
	end
	return jsonStr
	`,

	"GETORSET": `
	local val = redis.call('GET', KEYS[1])
	if not val then
		if tonumber(ARGV[2], 10) > 0 then
			redis.call('PSETEX', KEYS[1], ARGV[2], ARGV[1])
		else
			redis.call('SET', KEYS[1], ARGV[1])
		end
		return ARGV[1]
	end
	if tonumber(ARGV[2], 10) > 0 then
		redis.call('PEXPIRE', KEYS[1], ARGV[2])
	end
	return val
	`,

	"MSET_KEEPTTL": `
	-- KEYS: [key1, key2, key3, ...]
	-- ARGV: [val1, val2, val3, ...]
	for i = 1, #KEYS do
		redis.call('SET', KEYS[i], ARGV[i], 'KEEPTTL')
	end
	return 'OK'
	`,

	"MSET_EX": `
	for i = 1, #KEYS do
		redis.call('PSETEX', KEYS[i], ARGV[#ARGV], ARGV[i])
	end
	return 'OK'
	`,

	"BATCH_SET_EX": `
	-- KEYS: [key1, key2, key3, ...]
	-- ARGV: [val1, timeout1, val2, timeout2, val3, timeout3, ...]
	for i = 1, #KEYS do
		local val = ARGV[(i-1)*2 + 1]
		local timeout = tonumber(ARGV[(i-1)*2 + 2], 10)
		if timeout > 0 then
			redis.call('PSETEX', KEYS[i], timeout, val)
		else
			redis.call('SET', KEYS[i], val, 'KEEPTTL')
		end
	end
	return 'OK'
	`,

	"UPDATE_EX": `
	local pttl = redis.call('PTTL', KEYS[1])
	if pttl == -2 then
		return -1
	end
	redis.call('PEXPIRE', KEYS[1], ARGV[1])
	if pttl == -1 then
		return 0
	end
	return pttl
	`,

	"SADD_EX": `
	local count = redis.call('SADD', KEYS[1], unpack(ARGV, 1, #ARGV - 1))
	redis.call('PEXPIRE', KEYS[1], ARGV[#ARGV])
	return count
	`,

	"SISMEMBER_EX": `
	local isExist = redis.call('EXISTS', KEYS[1])
	if isExist == 0 then
		return 0
	end
	redis.call('PEXPIRE', KEYS[1], ARGV[2])
	return redis.call('SISMEMBER', KEYS[1], ARGV[1])
	`,

	"SMEMBERS_EX": `
	local isExist = redis.call('EXISTS', KEYS[1])
	if isExist == 0 then
		return nil
	end
	redis.call('PEXPIRE', KEYS[1], ARGV[1])
	return redis.call('SMEMBERS', KEYS[1])
	`,

	"SPOP_EX": `
	local members = redis.call('SPOP', KEYS[1], ARGV[1])
	local isExist = redis.call('EXISTS', KEYS[1])
	if isExist == 0 then
		return members
	end
	redis.call('PEXPIRE', KEYS[1], ARGV[2])
	return members
	`,

	"SUNION_EX": `
	local allKeysExist = true
	for i = 1, #KEYS do
		local isExist = redis.call('EXISTS', KEYS[i])
		if isExist == 0 then
			allKeysExist = false
			break
		end
	end
	if allKeysExist then
		for i = 1, #KEYS do
			redis.call('PEXPIRE', KEYS[i], ARGV[1])
		end
	end
	return redis.call('SUNION', unpack(KEYS))
	`,

	"SCARD_EX": `
	local count = redis.call('SCARD', KEYS[1])
	if count == 0 then
		return count
	end
	redis.call('PEXPIRE', KEYS[1], ARGV[1])
	return count
	`,

	"SREM_EX": `
	local count = redis.call('SREM', KEYS[1], unpack(ARGV, 1, #ARGV - 1))
	local isExist = redis.call('EXISTS', KEYS[1])
	if isExist == 0 then
		return count
	end
	redis.call('PEXPIRE', KEYS[1], ARGV[#ARGV])
	return count
	`,

	"ZADD_EX": `
	local count = redis.call('ZADD', KEYS[1], unpack(ARGV, 1, #ARGV - 1))
	redis.call('PEXPIRE', KEYS[1], ARGV[#ARGV])
	return count
	`,

	"RANGE_EX": `
	local isExist = redis.call('EXISTS', KEYS[1])
	if isExist == 0 then
		return nil
	end
	redis.call('PEXPIRE', KEYS[1], ARGV[#ARGV])
	return redis.call(ARGV[1], KEYS[1], unpack(ARGV, 2, #ARGV - 1))
	`,

	"PAGE_EX": `
	local total = redis.call('ZCARD', KEYS[1])
	if total == 0 then
		return cjson.encode({total=total})
	end
	if tonumber(ARGV[5], 10) > 0 then
		redis.call('PEXPIRE', KEYS[1], ARGV[5])
	end
	local start = (ARGV[2] - 1) * ARGV[3]
	local stop = (ARGV[2] * ARGV[3]) - 1
	if start >= total then
		return cjson.encode({total=total})
	end
	if ARGV[4] ~= "" then
		return cjson.encode({total=total, members=redis.call(ARGV[1], KEYS[1], start, stop, ARGV[4])})
	end
	return cjson.encode({total=total, members=redis.call(ARGV[1], KEYS[1], start, stop)})
	`,

	"ZCARD_EX": `
	local count = redis.call('ZCARD', KEYS[1])
	if count == 0 then
		return count
	end
	redis.call('PEXPIRE', KEYS[1], ARGV[1])
	return count
	`,

	"ZCOUNT_EX": `
	local isExist = redis.call('EXISTS', KEYS[1])
	if isExist == 0 then
		return 0
	end
	redis.call('PEXPIRE', KEYS[1], ARGV[3])
	return redis.call('ZCOUNT', KEYS[1], ARGV[1], ARGV[2])
	`,

	"ZINCRBY_EX": `
	local score = redis.call('ZINCRBY', KEYS[1], ARGV[1], ARGV[2])
	redis.call('PEXPIRE', KEYS[1], ARGV[3])
	return score
	`,

	"RANGEBYSCORE_EX": `
	local isExist = redis.call('EXISTS', KEYS[1])
	if isExist == 0 then
		return nil
	end
	redis.call('PEXPIRE', KEYS[1], ARGV[#ARGV])
	return redis.call(ARGV[1], KEYS[1], unpack(ARGV, 2, #ARGV - 1))
	`,

	"RANK_EX": `
	local rank = redis.call(ARGV[1], KEYS[1], ARGV[2])
	if not rank then
		return -1
	end
	redis.call('PEXPIRE', KEYS[1], ARGV[3])
	return rank
	`,

	"ZREM_EX": `
	local count = redis.call('ZREM', KEYS[1], unpack(ARGV, 1, #ARGV - 1))
	local isExist = redis.call('EXISTS', KEYS[1])
	if isExist == 0 then
		return count
	end
	redis.call('PEXPIRE', KEYS[1], ARGV[#ARGV])
	return count
	`,

	"ZREMRANGEBYRANK_EX": `
	local count = redis.call('ZREMRANGEBYRANK', KEYS[1], ARGV[1], ARGV[2])
	local isExist = redis.call('EXISTS', KEYS[1])
	if isExist == 0 then
		return count
	end
	redis.call('PEXPIRE', KEYS[1], ARGV[3])
	return count
	`,

	"ZREMRANGEBYSCORE_EX": `
	local count = redis.call('ZREMRANGEBYSCORE', KEYS[1], ARGV[1], ARGV[2])
	local isExist = redis.call('EXISTS', KEYS[1])
	if isExist == 0 then
		return count
	end
	redis.call('PEXPIRE', KEYS[1], ARGV[3])
	return count
	`,

	"ZSCORE_EX": `
	local score = redis.call('ZSCORE', KEYS[1], ARGV[1])
	if not score then
		return -1
	end
	redis.call('PEXPIRE', KEYS[1], ARGV[2])
	return score
	`,
}

// NewRedisCache 创建 RedisCache
func NewRedisCache(ctx context.Context, cfg *gtkredis.ClientConfig) (rc *RedisCache, err error) {
	var client *gtkredis.RedisClient
	if client, err = gtkredis.NewClient(ctx, cfg); err != nil {
		return
	}
	rc = &RedisCache{
		ctx:    ctx,
		client: client,
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

// add 添加缓存
//
//	当`key`已存在且未过期时，返回现有值（不修改）
//	当`key`不存在或已过期时，设置新值并返回新值
func (rc *RedisCache) add(ctx context.Context, key string, val any, timeout ...time.Duration) (newVal any, err error) {
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		newVal, err = rc.client.EvalSha(ctx, "ADD_EX", []string{key}, val, 0)
	} else {
		newVal, err = rc.client.EvalSha(ctx, "ADD_EX", []string{key}, val, timeout[0].Milliseconds())
	}
	return
}

// Get 获取缓存
//
//	当`timeout > 0`且缓存命中时，设置/重置`key`的过期时间
func (rc *RedisCache) Get(ctx context.Context, key string, timeout ...time.Duration) (val any, err error) {
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		val, err = rc.client.Do(ctx, "GET", key)
	} else {
		val, err = rc.client.Do(ctx, "GETEX", key, "PX", timeout[0].Milliseconds())
	}
	return
}

// GetMap 批量获取缓存
//
//	当`timeout > 0`且所有缓存都命中时，设置/重置所有`key`的过期时间，所有`key`过期时间相同
//	注意：如需为每个`key`设置/重置不同的过期时间，请使用`BatchGet`
func (rc *RedisCache) GetMap(ctx context.Context, keys []string, timeout ...time.Duration) (data map[string]any, err error) {
	if len(keys) == 0 {
		return
	}

	var result any
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		args := make([]any, 0, len(keys))
		for _, v := range keys {
			args = append(args, v)
		}
		result, err = rc.client.Do(ctx, "MGET", args...)
	} else {
		result, err = rc.client.EvalSha(ctx, "MGET_EX", keys, timeout[0].Milliseconds())
	}
	if err != nil {
		return
	}

	resultList := gtkconv.ToSlice(result)
	data = make(map[string]any)
	for k, v := range keys {
		data[v] = resultList[k]
	}
	return
}

// BatchGet 创建批量获取构建器
//
//	支持为每个`key`设置/重置不同的过期时间
//	当所有`key`使用相同过期时间时，可以使用更简洁的`GetMap`方法
//	当`capacity > 0`时，预分配指定容量以优化性能
//	返回构建器实例，支持链式调用
func (rc *RedisCache) BatchGet(ctx context.Context, capacity ...int) (batchGetter IBatchGetter) {
	cap := 0
	if len(capacity) > 0 && capacity[0] > 0 {
		cap = capacity[0]
	}
	return &redisBatchGetter{
		rc:    rc,
		items: make([]batchGetItem, 0, cap),
	}
}

// GetOrSet 检索并返回`key`的值，或者当`key`不存在时，则使用`newVal`设置`key`的值
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
func (rc *RedisCache) GetOrSet(ctx context.Context, key string, newVal any, timeout ...time.Duration) (val any, err error) {
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		val, err = rc.client.EvalSha(ctx, "GETORSET", []string{key}, newVal, 0)
	} else {
		val, err = rc.client.EvalSha(ctx, "GETORSET", []string{key}, newVal, timeout[0].Milliseconds())
	}
	return
}

// GetOrSetFunc 检索并返回`key`的值，或者当`key`不存在时，则使用函数`f`的结果设置`key`的值
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
//	当`force = true`时，可防止缓存穿透（即使`f`返回`nil`也会缓存）
//	注意：使用`singleflight`机制确保相同`key`的函数`f`只执行一次，其他并发请求等待并共享第一个请求的执行结果，有效防止缓存击穿
func (rc *RedisCache) GetOrSetFunc(ctx context.Context, key string, f Func, force bool, timeout ...time.Duration) (val any, err error) {
	// 获取缓存
	if val, err = rc.Get(ctx, key, timeout...); err != nil {
		return
	}
	if val != nil {
		return
	}
	// 使用 singleflight 确保函数只执行一次
	var result any
	if result, err, _ = rc.group.Do(key, func() (v any, e error) {
		// 获取缓存（double-check）
		var cVal any
		if cVal, e = rc.Get(ctx, key, timeout...); e != nil {
			return
		}
		if cVal != nil {
			v = singleflightValue{val: cVal, fromCache: true}
			return
		}
		// 执行函数获取新值
		var fVal any
		if fVal, e = f(ctx); e != nil {
			return
		}
		v = singleflightValue{val: fVal, fromCache: false}
		return
	}); err != nil {
		return
	}
	sfVal := result.(singleflightValue)
	if sfVal.fromCache {
		val = sfVal.val
		return
	}
	if utils.IsNil(sfVal.val) && !force {
		return
	}
	return rc.add(ctx, key, sfVal.val, timeout...)
}

// CustomGetOrSetFunc 从缓存中获取指定键`keys`的值，如果缓存未命中，则使用函数`f`的结果设置`keys`的值
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
//	当`force = true`时，可防止缓存穿透（即使`f`返回`nil`也会缓存）
//	注意：使用`singleflight`机制确保相同`key`的函数`f`只执行一次，其他并发请求等待并共享第一个请求的执行结果，有效防止缓存击穿
func (rc *RedisCache) CustomGetOrSetFunc(ctx context.Context, keys []string, args []any, cc ICustomCache, f Func, force bool, timeout ...time.Duration) (val any, err error) {
	// 获取缓存
	if val, err = cc.Get(ctx, keys, args, timeout...); err != nil {
		return
	}
	if val != nil {
		return
	}
	// 生成 singleflight 的唯一 key
	var sfKey string
	if sfKey, err = generateSingleflightKey(keys, args); err != nil {
		return
	}
	// 使用 singleflight 确保函数只执行一次
	var result any
	if result, err, _ = rc.group.Do(sfKey, func() (v any, e error) {
		// 获取缓存（double-check）
		var cVal any
		if cVal, e = cc.Get(ctx, keys, args, timeout...); e != nil {
			return
		}
		if cVal != nil {
			v = singleflightValue{val: cVal, fromCache: true}
			return
		}
		// 执行函数获取新值
		var fVal any
		if fVal, e = f(ctx); e != nil {
			return
		}
		v = singleflightValue{val: fVal, fromCache: false}
		return
	}); err != nil {
		return
	}
	sfVal := result.(singleflightValue)
	if sfVal.fromCache {
		val = sfVal.val
		return
	}
	if utils.IsNil(sfVal.val) && !force {
		return
	}
	return cc.Add(ctx, keys, args, sfVal.val, timeout...)
}

// Set 设置缓存
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
func (rc *RedisCache) Set(ctx context.Context, key string, val any, timeout ...time.Duration) (err error) {
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		_, err = rc.client.Do(ctx, "SET", key, val, "KEEPTTL")
	} else {
		_, err = rc.client.Do(ctx, "PSETEX", key, timeout[0].Milliseconds(), val)
	}
	return
}

// SetMap 批量设置缓存，所有`key`的过期时间相同
//
//	当`timeout > 0`时，设置/重置所有`key`的过期时间，所有`key`过期时间相同
//	注意：如需为每个`key`设置不同的过期时间，请使用`BatchSet`
func (rc *RedisCache) SetMap(ctx context.Context, data map[string]any, timeout ...time.Duration) (err error) {
	if len(data) == 0 {
		return
	}

	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		keys := make([]string, 0, len(data))
		args := make([]any, 0, len(data))
		for k, v := range data {
			keys = append(keys, k)
			args = append(args, v)
		}
		_, err = rc.client.EvalSha(ctx, "MSET_KEEPTTL", keys, args...)
	} else {
		keys := make([]string, 0, len(data))
		args := make([]any, 0, len(data)+1)
		for k, v := range data {
			keys = append(keys, k)
			args = append(args, v)
		}
		args = append(args, timeout[0].Milliseconds())
		_, err = rc.client.EvalSha(ctx, "MSET_EX", keys, args...)
	}
	return
}

// BatchSet 创建批量设置构建器
//
//	支持为每个`key`设置不同的过期时间
//	当所有`key`使用相同过期时间时，可以使用更简洁的`SetMap`方法
//	当`capacity > 0`时，预分配指定容量以优化性能
//	返回构建器实例，支持链式调用
func (rc *RedisCache) BatchSet(ctx context.Context, capacity ...int) (batchSetter IBatchSetter) {
	cap := 0
	if len(capacity) > 0 && capacity[0] > 0 {
		cap = capacity[0]
	}
	return &redisBatchSetter{
		rc:    rc,
		items: make([]batchSetItem, 0, cap),
	}
}

// SetIfNotExist 当`key`不存在时，则使用`val`设置`key`的值，返回是否设置成功
//
//	当`timeout > 0`且`key`设置成功时，设置`key`的过期时间
func (rc *RedisCache) SetIfNotExist(ctx context.Context, key string, val any, timeout ...time.Duration) (ok bool, err error) {
	var result any
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		result, err = rc.client.Do(ctx, "SETNX", key, val)
	} else {
		result, err = rc.client.Do(ctx, "SET", key, val, "PX", timeout[0].Milliseconds(), "NX")
	}
	ok = gtkconv.ToBool(result)
	return
}

// SetIfNotExistFunc 当`key`不存在时，则使用函数`f`的结果设置`key`的值，返回是否设置成功
//
//	当`timeout > 0`且`key`设置成功时，设置`key`的过期时间
//	当`force = true`时，可防止缓存穿透（即使`f`返回`nil`也会缓存）
//	注意：使用`singleflight`机制确保相同`key`的函数`f`只执行一次，其他并发请求等待并共享第一个请求的执行结果，有效防止缓存击穿
func (rc *RedisCache) SetIfNotExistFunc(ctx context.Context, key string, f Func, force bool, timeout ...time.Duration) (ok bool, err error) {
	// 缓存是否存在
	var isExist bool
	if isExist, err = rc.IsExist(ctx, key); err != nil {
		return
	}
	if isExist {
		return
	}
	// 使用 singleflight 确保函数只执行一次
	var result any
	if result, err, _ = rc.group.Do(key, func() (v any, e error) {
		// 缓存是否存在（double-check）
		var cIsExist bool
		if cIsExist, e = rc.IsExist(ctx, key); e != nil {
			return
		}
		if cIsExist {
			v = singleflightValue{val: nil, fromCache: true}
			return
		}
		// 执行函数获取新值
		var fVal any
		if fVal, e = f(ctx); e != nil {
			return
		}
		v = singleflightValue{val: fVal, fromCache: false}
		return
	}); err != nil {
		return
	}
	sfVal := result.(singleflightValue)
	if sfVal.fromCache {
		return
	}
	if utils.IsNil(sfVal.val) && !force {
		return
	}
	return rc.SetIfNotExist(ctx, key, sfVal.val, timeout...)
}

// Update 当`key`存在时，则使用`val`更新`key`的值，返回`key`的旧值
//
//	当`timeout > 0`且`key`更新成功时，更新`key`的过期时间
func (rc *RedisCache) Update(ctx context.Context, key string, val any, timeout ...time.Duration) (oldVal any, isExist bool, err error) {
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		oldVal, err = rc.client.Do(ctx, "SET", key, val, "XX", "GET", "KEEPTTL")
	} else {
		oldVal, err = rc.client.Do(ctx, "SET", key, val, "XX", "GET", "PX", timeout[0].Milliseconds())
	}
	if err != nil {
		return
	}
	if utils.IsNil(oldVal) {
		return
	}
	isExist = true
	return
}

// UpdateExpire 当`key`存在时，则更新`key`的过期时间，返回`key`的旧的过期时间值
//
//	当`key`不存在时，则返回-1
//	当`key`存在但没有设置过期时间时，则返回0
//	当`key`存在且设置了过期时间时，则返回过期时间
//	当`timeout > 0`且`key`存在时，更新`key`的过期时间
func (rc *RedisCache) UpdateExpire(ctx context.Context, key string, timeout time.Duration) (oldTimeout time.Duration, err error) {
	if timeout.Milliseconds() <= 0 {
		err = errors.New("timeout must be greater than 0")
		return
	}

	var value any
	if value, err = rc.client.EvalSha(ctx, "UPDATE_EX", []string{key}, timeout.Milliseconds()); err != nil {
		return
	}

	pttl := gtkconv.ToInt64(value)
	if pttl <= 0 {
		return time.Duration(pttl), nil
	}
	return time.Duration(pttl) * time.Millisecond, nil
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

// Size 缓存中的key数量
func (rc *RedisCache) Size(ctx context.Context) (size int, err error) {
	var val any
	if val, err = rc.client.Do(ctx, "DBSIZE"); err != nil {
		return
	}
	size = gtkconv.ToInt(val)
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

// GetExpire 获取缓存`key`的过期时间
//
//	当`key`不存在时，则返回-1
//	当`key`存在但没有设置过期时间时，则返回0
//	当`key`存在且设置了过期时间时，则返回过期时间
func (rc *RedisCache) GetExpire(ctx context.Context, key string) (timeout time.Duration, err error) {
	var val any
	if val, err = rc.client.Do(ctx, "PTTL", key); err != nil {
		return
	}

	pttl := gtkconv.ToInt64(val)
	switch pttl {
	case -2: // key不存在
		return -1, nil
	case -1: // key存在但没有设置剩余生存时间
		return 0, nil
	default:
		return time.Duration(pttl) * time.Millisecond, nil
	}
}

// Close 关闭缓存服务
func (rc *RedisCache) Close(ctx context.Context) (err error) {
	err = rc.client.Close()
	return
}

// SAdd 向集合添加一个或多个成员
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
func (rc *RedisCache) SAdd(ctx context.Context, key string, members []any, timeout ...time.Duration) (addCount int, err error) {
	if len(members) == 0 {
		return
	}
	var (
		args   = make([]any, 0, len(members)+1)
		result any
	)
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		args = append(args, key)
		args = append(args, members...)
		result, err = rc.client.Do(ctx, "SADD", args...)
	} else {
		args = append(args, members...)
		args = append(args, timeout[0].Milliseconds())
		result, err = rc.client.EvalSha(ctx, "SADD_EX", []string{key}, args...)
	}
	if err != nil {
		return
	}
	addCount = gtkconv.ToInt(result)
	return
}

// SIsMember 判断 member 元素是否是集合 key 的成员
//
//	当`timeout > 0`且`key`存在时，设置/重置`key`的过期时间
func (rc *RedisCache) SIsMember(ctx context.Context, key string, member any, timeout ...time.Duration) (isMember bool, err error) {
	var result any
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		result, err = rc.client.Do(ctx, "SISMEMBER", key, member)
	} else {
		result, err = rc.client.EvalSha(ctx, "SISMEMBER_EX", []string{key}, member, timeout[0].Milliseconds())
	}
	if err != nil {
		return
	}
	isMember = gtkconv.ToBool(result)
	return
}

// SMembers 返回集合中的所有成员
//
//	当`timeout > 0`且`key`存在时，设置/重置`key`的过期时间
func (rc *RedisCache) SMembers(ctx context.Context, key string, timeout ...time.Duration) (members []any, err error) {
	var result any
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		result, err = rc.client.Do(ctx, "SMEMBERS", key)
	} else {
		result, err = rc.client.EvalSha(ctx, "SMEMBERS_EX", []string{key}, timeout[0].Milliseconds())
	}
	if err != nil {
		return
	}
	members = gtkconv.ToSlice(result)
	return
}

// SPop 移除并返回集合中的一个或多个随机元素
//
//	当`timeout > 0`且更新后的`key`存在时，设置/重置`key`的过期时间
func (rc *RedisCache) SPop(ctx context.Context, key string, count int, timeout ...time.Duration) (members []any, err error) {
	var result any
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		result, err = rc.client.Do(ctx, "SPOP", key, count)
	} else {
		result, err = rc.client.EvalSha(ctx, "SPOP_EX", []string{key}, count, timeout[0].Milliseconds())
	}
	if err != nil {
		return
	}
	members = gtkconv.ToSlice(result)
	return
}

// SUnion 返回所有给定集合的并集
//
//	当`timeout > 0`且所有`key`都存在时，设置/重置所有`key`的过期时间，所有`key`过期时间相同
func (rc *RedisCache) SUnion(ctx context.Context, keys []string, timeout ...time.Duration) (members []any, err error) {
	if len(keys) == 0 {
		return
	}
	var result any
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		args := make([]any, 0, len(keys))
		for _, v := range keys {
			args = append(args, v)
		}
		result, err = rc.client.Do(ctx, "SUNION", args...)
	} else {
		result, err = rc.client.EvalSha(ctx, "SUNION_EX", keys, timeout[0].Milliseconds())
	}
	if err != nil {
		return
	}
	members = gtkconv.ToSlice(result)
	return
}

// SCard 获取集合的成员数
//
//	当`timeout > 0`且`key`存在时，设置/重置`key`的过期时间
func (rc *RedisCache) SCard(ctx context.Context, key string, timeout ...time.Duration) (count int, err error) {
	var result any
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		result, err = rc.client.Do(ctx, "SCARD", key)
	} else {
		result, err = rc.client.EvalSha(ctx, "SCARD_EX", []string{key}, timeout[0].Milliseconds())
	}
	if err != nil {
		return
	}
	count = gtkconv.ToInt(result)
	return
}

// SRem 移除集合中一个或多个成员
//
//	当`timeout > 0`且更新后的`key`存在时，设置/重置`key`的过期时间
func (rc *RedisCache) SRem(ctx context.Context, key string, members []any, timeout ...time.Duration) (remCount int, err error) {
	if len(members) == 0 {
		return
	}
	var (
		args   = make([]any, 0, len(members)+1)
		result any
	)
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		args = append(args, key)
		args = append(args, members...)
		result, err = rc.client.Do(ctx, "SREM", args...)
	} else {
		args = append(args, members...)
		args = append(args, timeout[0].Milliseconds())
		result, err = rc.client.EvalSha(ctx, "SREM_EX", []string{key}, args...)
	}
	if err != nil {
		return
	}
	remCount = gtkconv.ToInt(result)
	return
}

// SSAdd 向有序集合添加一个或多个成员，或者更新已存在成员的分数
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
func (rc *RedisCache) SSAdd(ctx context.Context, key string, data map[any]float64, timeout ...time.Duration) (addCount int, err error) {
	if len(data) == 0 {
		return
	}
	var (
		args   = make([]any, 0, len(data)*2+1)
		result any
	)
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		args = append(args, key)
		for k, v := range data {
			args = append(args, v, k)
		}
		result, err = rc.client.Do(ctx, "ZADD", args...)
	} else {
		for k, v := range data {
			args = append(args, v, k)
		}
		args = append(args, timeout[0].Milliseconds())
		result, err = rc.client.EvalSha(ctx, "ZADD_EX", []string{key}, args...)
	}
	if err != nil {
		return
	}
	addCount = gtkconv.ToInt(result)
	return
}

// SSRange 返回有序集合中指定区间内的成员
//
//	当`timeout > 0`且`key`存在时，设置/重置`key`的过期时间
func (rc *RedisCache) SSRange(ctx context.Context, key string, start, stop int, isDescOrder, withScores bool, timeout ...time.Duration) (members []map[any]float64, err error) {
	var (
		args   = make([]any, 0, 5)
		result any
	)
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		args = append(args, key, start, stop)
		if withScores {
			args = append(args, "WITHSCORES")
		}
		if isDescOrder {
			// score 值递减
			result, err = rc.client.Do(ctx, "ZREVRANGE", args...)
		} else {
			// score 值递增
			result, err = rc.client.Do(ctx, "ZRANGE", args...)
		}
	} else {
		if isDescOrder {
			// score 值递减
			args = append(args, "ZREVRANGE", start, stop)
		} else {
			// score 值递增
			args = append(args, "ZRANGE", start, stop)
		}
		if withScores {
			args = append(args, "WITHSCORES")
		}
		args = append(args, timeout[0].Milliseconds())
		result, err = rc.client.EvalSha(ctx, "RANGE_EX", []string{key}, args...)
	}
	if err != nil {
		return
	}
	resultList := gtkconv.ToSlice(result)
	if withScores {
		members = make([]map[any]float64, 0, len(resultList)/2)
		for i := 0; i < len(resultList); i += 2 {
			members = append(members, map[any]float64{resultList[i]: gtkconv.ToFloat64(resultList[i+1])})
		}
	} else {
		members = make([]map[any]float64, 0, len(resultList))
		for _, v := range resultList {
			members = append(members, map[any]float64{v: 0})
		}
	}
	return
}

// SSPage 有序集合分页查询
//
//	当`timeout > 0`且`key`存在时，设置/重置`key`的过期时间
func (rc *RedisCache) SSPage(ctx context.Context, key string, page, pageSize int, isDescOrder, withScores bool, timeout ...time.Duration) (total int, members []map[any]float64, err error) {
	var (
		args   = make([]any, 0, 5)
		result any
	)
	if isDescOrder {
		// score 值递减
		args = append(args, "ZREVRANGE", page, pageSize)
	} else {
		// score 值递增
		args = append(args, "ZRANGE", page, pageSize)
	}
	if withScores {
		args = append(args, "WITHSCORES")
	} else {
		args = append(args, "")
	}
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		args = append(args, 0)
	} else {
		args = append(args, timeout[0].Milliseconds())
	}
	if result, err = rc.client.EvalSha(ctx, "PAGE_EX", []string{key}, args...); err != nil {
		return
	}
	resultMap := gtkconv.ToStringMap(result)
	total = gtkconv.ToInt(resultMap["total"])
	memberList := gtkconv.ToSlice(resultMap["members"])
	if withScores {
		members = make([]map[any]float64, 0, len(memberList)/2)
		for i := 0; i < len(memberList); i += 2 {
			members = append(members, map[any]float64{memberList[i]: gtkconv.ToFloat64(memberList[i+1])})
		}
	} else {
		members = make([]map[any]float64, 0, len(memberList))
		for _, v := range memberList {
			members = append(members, map[any]float64{v: 0})
		}
	}
	return
}

// SSCard 获取有序集合的成员数
//
//	当`timeout > 0`且`key`存在时，设置/重置`key`的过期时间
func (rc *RedisCache) SSCard(ctx context.Context, key string, timeout ...time.Duration) (count int, err error) {
	var result any
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		result, err = rc.client.Do(ctx, "ZCARD", key)
	} else {
		result, err = rc.client.EvalSha(ctx, "ZCARD_EX", []string{key}, timeout[0].Milliseconds())
	}
	if err != nil {
		return
	}
	count = gtkconv.ToInt(result)
	return
}

// SSCount 计算在有序集合中指定区间分数的成员数
//
//	关于参数`min`和`max`的详细使用方法，请参考`redis`的`ZRANGEBYSCORE`命令
//	当`timeout > 0`且`key`存在时，设置/重置`key`的过期时间
func (rc *RedisCache) SSCount(ctx context.Context, key, min, max string, timeout ...time.Duration) (count int, err error) {
	var result any
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		result, err = rc.client.Do(ctx, "ZCOUNT", key, min, max)
	} else {
		result, err = rc.client.EvalSha(ctx, "ZCOUNT_EX", []string{key}, min, max, timeout[0].Milliseconds())
	}
	if err != nil {
		return
	}
	count = gtkconv.ToInt(result)
	return
}

// SSIncrby 有序集合中对指定成员的分数加上增量 increment
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
func (rc *RedisCache) SSIncrby(ctx context.Context, key string, increment float64, member any, timeout ...time.Duration) (score float64, err error) {
	var result any
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		result, err = rc.client.Do(ctx, "ZINCRBY", key, increment, member)
	} else {
		result, err = rc.client.EvalSha(ctx, "ZINCRBY_EX", []string{key}, increment, member, timeout[0].Milliseconds())
	}
	if err != nil {
		return
	}
	score = gtkconv.ToFloat64(result)
	return
}

// SSRangeByScore 通过分数返回有序集合指定区间内的成员
//
//	关于参数`min`和`max`的详细使用方法，请参考`redis`的`ZRANGEBYSCORE`命令
//	当`timeout > 0`且`key`存在时，设置/重置`key`的过期时间
func (rc *RedisCache) SSRangeByScore(ctx context.Context, key, min, max string, isDescOrder, withScores bool, limit []int, timeout ...time.Duration) (members []map[any]float64, err error) {
	var (
		args   = make([]any, 0, 8)
		result any
	)
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		args = append(args, key, min, max)
		if withScores {
			args = append(args, "WITHSCORES")
		}
		if len(limit) >= 2 {
			args = append(args, "LIMIT", limit[0], limit[1])
		}
		if isDescOrder {
			// score 值递减
			result, err = rc.client.Do(ctx, "ZREVRANGEBYSCORE", args...)
		} else {
			// score 值递增
			result, err = rc.client.Do(ctx, "ZRANGEBYSCORE", args...)
		}
	} else {
		if isDescOrder {
			// score 值递减
			args = append(args, "ZREVRANGEBYSCORE", max, min)
		} else {
			// score 值递增
			args = append(args, "ZRANGEBYSCORE", min, max)
		}
		if withScores {
			args = append(args, "WITHSCORES")
		}
		if len(limit) >= 2 {
			args = append(args, "LIMIT", limit[0], limit[1])
		}
		args = append(args, timeout[0].Milliseconds())
		result, err = rc.client.EvalSha(ctx, "RANGEBYSCORE_EX", []string{key}, args...)
	}
	if err != nil {
		return
	}
	resultList := gtkconv.ToSlice(result)
	if withScores {
		members = make([]map[any]float64, 0, len(resultList)/2)
		for i := 0; i < len(resultList); i += 2 {
			members = append(members, map[any]float64{resultList[i]: gtkconv.ToFloat64(resultList[i+1])})
		}
	} else {
		members = make([]map[any]float64, 0, len(resultList))
		for _, v := range resultList {
			members = append(members, map[any]float64{v: 0})
		}
	}
	return
}

// SSRank 返回有序集合中指定成员的排名
//
//	当`timeout > 0`且`key`存在时，设置/重置`key`的过期时间
func (rc *RedisCache) SSRank(ctx context.Context, key string, member any, isDescOrder bool, timeout ...time.Duration) (rank int, err error) {
	var result any
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		if isDescOrder {
			// score 值递减
			result, err = rc.client.Do(ctx, "ZREVRANK", key, member)
		} else {
			// score 值递增
			result, err = rc.client.Do(ctx, "ZRANK", key, member)
		}
	} else {
		if isDescOrder {
			// score 值递减
			result, err = rc.client.EvalSha(ctx, "RANK_EX", []string{key}, "ZREVRANK", member, timeout[0].Milliseconds())
		} else {
			// score 值递增
			result, err = rc.client.EvalSha(ctx, "RANK_EX", []string{key}, "ZRANK", member, timeout[0].Milliseconds())
		}
	}
	if err != nil {
		return
	}
	rank = gtkconv.ToInt(result)
	return
}

// SSRem 移除有序集合中的一个或多个成员
//
//	当`timeout > 0`且更新后的`key`存在时，设置/重置`key`的过期时间
func (rc *RedisCache) SSRem(ctx context.Context, key string, members []any, timeout ...time.Duration) (remCount int, err error) {
	if len(members) == 0 {
		return
	}
	var (
		args   = make([]any, 0, len(members)+1)
		result any
	)
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		args = append(args, key)
		args = append(args, members...)
		result, err = rc.client.Do(ctx, "ZREM", args...)
	} else {
		args = append(args, members...)
		args = append(args, timeout[0].Milliseconds())
		result, err = rc.client.EvalSha(ctx, "ZREM_EX", []string{key}, args...)
	}
	if err != nil {
		return
	}
	remCount = gtkconv.ToInt(result)
	return
}

// SSRemRangeByRank 移除有序集合中给定的排名区间的所有成员
//
//	当`timeout > 0`且更新后的`key`存在时，设置/重置`key`的过期时间
func (rc *RedisCache) SSRemRangeByRank(ctx context.Context, key string, start, stop int, timeout ...time.Duration) (remCount int, err error) {
	var result any
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		result, err = rc.client.Do(ctx, "ZREMRANGEBYRANK", key, start, stop)
	} else {
		result, err = rc.client.EvalSha(ctx, "ZREMRANGEBYRANK_EX", []string{key}, start, stop, timeout[0].Milliseconds())
	}
	if err != nil {
		return
	}
	remCount = gtkconv.ToInt(result)
	return
}

// SSRemRangeByScore 移除有序集合中给定的分数区间的所有成员
//
//	关于参数`min`和`max`的详细使用方法，请参考`redis`的`ZRANGEBYSCORE`命令
//	当`timeout > 0`且更新后的`key`存在时，设置/重置`key`的过期时间
func (rc *RedisCache) SSRemRangeByScore(ctx context.Context, key, min, max string, timeout ...time.Duration) (remCount int, err error) {
	var result any
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		result, err = rc.client.Do(ctx, "ZREMRANGEBYSCORE", key, min, max)
	} else {
		result, err = rc.client.EvalSha(ctx, "ZREMRANGEBYSCORE_EX", []string{key}, min, max, timeout[0].Milliseconds())
	}
	if err != nil {
		return
	}
	remCount = gtkconv.ToInt(result)
	return
}

// SSScore 返回有序集中，成员的分数值
//
//	当`timeout > 0`且`key`存在时，设置/重置`key`的过期时间
func (rc *RedisCache) SSScore(ctx context.Context, key string, member any, timeout ...time.Duration) (score float64, err error) {
	var result any
	if len(timeout) == 0 || timeout[0].Milliseconds() <= 0 {
		result, err = rc.client.Do(ctx, "ZSCORE", key, member)
	} else {
		result, err = rc.client.EvalSha(ctx, "ZSCORE_EX", []string{key}, member, timeout[0].Milliseconds())
	}
	if err != nil {
		return
	}
	score = gtkconv.ToFloat64(result)
	return
}
