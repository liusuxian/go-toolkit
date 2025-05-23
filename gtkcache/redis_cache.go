/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-27 20:53:08
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-23 17:33:51
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkcache

import (
	"context"
	"github.com/liusuxian/go-toolkit/gtkconv"
	"github.com/liusuxian/go-toolkit/gtkredis"
	"github.com/liusuxian/go-toolkit/internal/utils"
	"time"
)

// RedisCache Redis 缓存
type RedisCache struct {
	ctx    context.Context
	client *gtkredis.RedisClient // redis 客户端
}

// 内置 lua 脚本
var internalScriptMap = map[string]string{
	"SETGET": `
	if tonumber(ARGV[2], 10) > 0 then
		redis.call('SETEX', KEYS[1], ARGV[2], ARGV[1])
	else
		redis.call('SET', KEYS[1], ARGV[1])
	end
	return redis.call('GET', KEYS[1])
	`,

	"GET_EX": `
	local val = redis.call('GET', KEYS[1])
	if not val then
		return nil
	end
	redis.call('EXPIRE', KEYS[1], ARGV[1])
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
			redis.call('EXPIRE', KEYS[i], ARGV[1])
    end
	end
	return vals
	`,

	"GETORSET": `
	local val = redis.call('GET', KEYS[1])
	if not val then
		if tonumber(ARGV[2], 10) > 0 then
			redis.call('SETEX', KEYS[1], ARGV[2], ARGV[1])
		else
			redis.call('SET', KEYS[1], ARGV[1])
		end
		return redis.call('GET', KEYS[1])
	end
	if tonumber(ARGV[2], 10) > 0 then
		redis.call('EXPIRE', KEYS[1], ARGV[2])
	end
	return val
	`,

	"MSET_EX": `
	for i = 1, #KEYS do
		redis.call('SETEX', KEYS[i], ARGV[#ARGV], ARGV[i])
	end
	return 'OK'
	`,

	"SADD_EX": `
	local count = redis.call('SADD', KEYS[1], unpack(ARGV, 1, #ARGV - 1))
	redis.call('EXPIRE', KEYS[1], ARGV[#ARGV])
	return count
	`,

	"SISMEMBER_EX": `
	local isExist = redis.call('EXISTS', KEYS[1])
	if isExist == 0 then
    return 0
  end
	redis.call('EXPIRE', KEYS[1], ARGV[2])
	return redis.call('SISMEMBER', KEYS[1], ARGV[1])
	`,

	"SMEMBERS_EX": `
	local isExist = redis.call('EXISTS', KEYS[1])
	if isExist == 0 then
    return nil
  end
	redis.call('EXPIRE', KEYS[1], ARGV[1])
	return redis.call('SMEMBERS', KEYS[1])
	`,

	"SPOP_EX": `
	local members = redis.call('SPOP', KEYS[1], ARGV[1])
	local isExist = redis.call('EXISTS', KEYS[1])
	if isExist == 0 then
    return members
  end
	redis.call('EXPIRE', KEYS[1], ARGV[2])
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
			redis.call('EXPIRE', KEYS[i], ARGV[1])
    end
	end
	return redis.call('SUNION', unpack(KEYS))
	`,

	"SCARD_EX": `
	local count = redis.call('SCARD', KEYS[1])
	if count == 0 then
    return count
  end
	redis.call('EXPIRE', KEYS[1], ARGV[1])
	return count
	`,

	"SREM_EX": `
	local count = redis.call('SREM', KEYS[1], unpack(ARGV, 1, #ARGV - 1))
	local isExist = redis.call('EXISTS', KEYS[1])
	if isExist == 0 then
    return count
  end
	redis.call('EXPIRE', KEYS[1], ARGV[#ARGV])
	return count
	`,

	"ZADD_EX": `
	local count = redis.call('ZADD', KEYS[1], unpack(ARGV, 1, #ARGV - 1))
	redis.call('EXPIRE', KEYS[1], ARGV[#ARGV])
	return count
	`,

	"RANGE_EX": `
	local isExist = redis.call('EXISTS', KEYS[1])
	if isExist == 0 then
    return nil
  end
	redis.call('EXPIRE', KEYS[1], ARGV[#ARGV])
	return redis.call(ARGV[1], KEYS[1], unpack(ARGV, 2, #ARGV - 1))
	`,

	"PAGE_EX": `
	local total = redis.call('ZCARD', KEYS[1])
	if total == 0 then
		return cjson.encode({total=total})
	end
	if tonumber(ARGV[5], 10) > 0 then
		redis.call('EXPIRE', KEYS[1], ARGV[5])
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
	redis.call('EXPIRE', KEYS[1], ARGV[1])
	return count
	`,

	"ZCOUNT_EX": `
	local isExist = redis.call('EXISTS', KEYS[1])
	if isExist == 0 then
    return 0
  end
	redis.call('EXPIRE', KEYS[1], ARGV[3])
	return redis.call('ZCOUNT', KEYS[1], ARGV[1], ARGV[2])
	`,

	"ZINCRBY_EX": `
	local score = redis.call('ZINCRBY', KEYS[1], ARGV[1], ARGV[2])
	redis.call('EXPIRE', KEYS[1], ARGV[3])
	return score
	`,

	"RANGEBYSCORE_EX": `
	local isExist = redis.call('EXISTS', KEYS[1])
	if isExist == 0 then
    return nil
  end
	redis.call('EXPIRE', KEYS[1], ARGV[#ARGV])
	return redis.call(ARGV[1], KEYS[1], unpack(ARGV, 2, #ARGV - 1))
	`,

	"RANK_EX": `
	local isExist = redis.call('EXISTS', KEYS[1])
	if isExist == 0 then
    return -1
  end
	redis.call('EXPIRE', KEYS[1], ARGV[3])
	local rank = redis.call(ARGV[1], KEYS[1], ARGV[2])
	if not rank then
    return -1
  end
	return rank
	`,

	"ZREM_EX": `
	local count = redis.call('ZREM', KEYS[1], unpack(ARGV, 1, #ARGV - 1))
	local isExist = redis.call('EXISTS', KEYS[1])
	if isExist == 0 then
    return count
  end
	redis.call('EXPIRE', KEYS[1], ARGV[#ARGV])
	return count
	`,

	"ZREMRANGEBYRANK_EX": `
	local count = redis.call('ZREMRANGEBYRANK', KEYS[1], ARGV[1], ARGV[2])
	local isExist = redis.call('EXISTS', KEYS[1])
	if isExist == 0 then
    return count
  end
	redis.call('EXPIRE', KEYS[1], ARGV[3])
	return count
	`,

	"ZREMRANGEBYSCORE_EX": `
	local count = redis.call('ZREMRANGEBYSCORE', KEYS[1], ARGV[1], ARGV[2])
	local isExist = redis.call('EXISTS', KEYS[1])
	if isExist == 0 then
    return count
  end
	redis.call('EXPIRE', KEYS[1], ARGV[3])
	return count
	`,

	"ZSCORE_EX": `
	local isExist = redis.call('EXISTS', KEYS[1])
	if isExist == 0 then
    return -1
  end
	redis.call('EXPIRE', KEYS[1], ARGV[2])
	local score = redis.call('ZSCORE', KEYS[1], ARGV[1])
	if not score then
    return -1
  end
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

// setget 设置`key`的值，并读取`key`的新值
//
//	当`timeout > 0`时，设置`key`的过期时间
func (rc *RedisCache) setget(ctx context.Context, key string, val any, timeout ...time.Duration) (newVal any, err error) {
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
		newVal, err = rc.client.EvalSha(ctx, "SETGET", []string{key}, val, 0)
	} else {
		newVal, err = rc.client.EvalSha(ctx, "SETGET", []string{key}, val, int64(timeout[0].Seconds()))
	}
	return
}

// Get 获取缓存
//
//	当`timeout > 0`且缓存命中时，设置/重置`key`的过期时间
func (rc *RedisCache) Get(ctx context.Context, key string, timeout ...time.Duration) (val any, err error) {
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
		val, err = rc.client.Do(ctx, "GET", key)
	} else {
		val, err = rc.client.EvalSha(ctx, "GET_EX", []string{key}, int64(timeout[0].Seconds()))
	}
	return
}

// GetMap 批量获取缓存
//
//	当`timeout > 0`且所有缓存都命中时，设置/重置所有`key`的过期时间，所有`key`过期时间相同
func (rc *RedisCache) GetMap(ctx context.Context, keys []string, timeout ...time.Duration) (data map[string]any, err error) {
	if len(keys) == 0 {
		return
	}
	var result any
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
		args := make([]any, 0, len(keys))
		for _, v := range keys {
			args = append(args, v)
		}
		result, err = rc.client.Do(ctx, "MGET", args...)
	} else {
		result, err = rc.client.EvalSha(ctx, "MGET_EX", keys, int64(timeout[0].Seconds()))
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

// GetOrSet 检索并返回`key`的值，或者当`key`不存在时，则使用`newVal`设置`key`的值
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
func (rc *RedisCache) GetOrSet(ctx context.Context, key string, newVal any, timeout ...time.Duration) (val any, err error) {
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
		val, err = rc.client.EvalSha(ctx, "GETORSET", []string{key}, newVal, 0)
	} else {
		val, err = rc.client.EvalSha(ctx, "GETORSET", []string{key}, newVal, int64(timeout[0].Seconds()))
	}
	return
}

// GetOrSetFunc 检索并返回`key`的值，或者当`key`不存在时，则使用函数`f`的结果设置`key`的值
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
//	当`force = true`时，可防止缓存穿透
func (rc *RedisCache) GetOrSetFunc(ctx context.Context, key string, f Func, force bool, timeout ...time.Duration) (val any, err error) {
	if val, err = rc.Get(ctx, key, timeout...); err != nil {
		return
	}
	if val == nil {
		var newVal any
		if newVal, err = f(ctx); err != nil {
			return
		}
		if utils.IsNil(newVal) && !force {
			return
		}
		// 此处不判断`newVal == nil`是因为防止缓存穿透
		val, err = rc.setget(ctx, key, newVal, timeout...)
		return
	}
	return
}

// GetOrSetFuncLock 检索并返回`key`的值，或者当`key`不存在时，则使用函数`f`的结果设置`key`的值，函数`f`是在读写互斥锁中执行的
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
//	当`force = true`时，可防止缓存穿透
func (rc *RedisCache) GetOrSetFuncLock(ctx context.Context, key string, f Func, force bool, timeout ...time.Duration) (val any, err error) {
	return rc.GetOrSetFunc(ctx, key, f, force, timeout...)
}

// CustomGetOrSetFunc 从缓存中获取指定键`keys`的值，如果缓存未命中，则使用函数`f`的结果设置`keys`的值
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
//	当`force = true`时，可防止缓存穿透
func (rc *RedisCache) CustomGetOrSetFunc(ctx context.Context, keys []string, args []any, cc CustomCache, f Func, force bool, timeout ...time.Duration) (val any, err error) {
	// 获取缓存
	if val, err = cc.Get(ctx, keys, args, timeout...); err != nil {
		return
	}
	if val == nil {
		var newVal any
		if newVal, err = f(ctx); err != nil {
			return
		}
		if utils.IsNil(newVal) && !force {
			return
		}
		// 此处不判断`newVal == nil`是因为防止缓存穿透
		val, err = cc.Set(ctx, keys, args, newVal, timeout...)
		return
	}
	return
}

// CustomGetOrSetFuncLock 从缓存中获取指定键`keys`的值，如果缓存未命中，则使用函数`f`的结果设置`keys`的值，函数`f`是在读写互斥锁中执行的
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
//	当`force = true`时，可防止缓存穿透
func (rc *RedisCache) CustomGetOrSetFuncLock(ctx context.Context, keys []string, args []any, cc CustomCache, f Func, force bool, timeout ...time.Duration) (val any, err error) {
	return rc.CustomGetOrSetFunc(ctx, keys, args, cc, f, force, timeout...)
}

// Set 设置缓存
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
func (rc *RedisCache) Set(ctx context.Context, key string, val any, timeout ...time.Duration) (err error) {
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
		_, err = rc.client.Do(ctx, "SET", key, val)
	} else {
		_, err = rc.client.Do(ctx, "SETEX", key, int64(timeout[0].Seconds()), val)
	}
	return
}

// SetMap 批量设置缓存，所有`key`的过期时间相同
//
//	当`timeout > 0`时，设置/重置所有`key`的过期时间，所有`key`过期时间相同
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
		_, err = rc.client.EvalSha(ctx, "MSET_EX", keys, args...)
	}
	return
}

// SetIfNotExist 当`key`不存在时，则使用`val`设置`key`的值，返回是否设置成功
//
//	当`timeout > 0`且`key`设置成功时，设置`key`的过期时间
func (rc *RedisCache) SetIfNotExist(ctx context.Context, key string, val any, timeout ...time.Duration) (ok bool, err error) {
	var result any
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
		result, err = rc.client.Do(ctx, "SETNX", key, val)
	} else {
		result, err = rc.client.Do(ctx, "SET", key, val, "EX", int64(timeout[0].Seconds()), "NX")
	}
	ok = gtkconv.ToBool(result)
	return
}

// SetIfNotExistFunc 当`key`不存在时，则使用函数`f`的结果设置`key`的值，返回是否设置成功
//
//	当`timeout > 0`且`key`设置成功时，设置`key`的过期时间
//	当`force = true`时，可防止缓存穿透
func (rc *RedisCache) SetIfNotExistFunc(ctx context.Context, key string, f Func, force bool, timeout ...time.Duration) (ok bool, err error) {
	var val any
	if val, err = f(ctx); err != nil {
		return
	}
	if utils.IsNil(val) && !force {
		return
	}
	// 此处不判断`val == nil`是因为防止缓存穿透
	return rc.SetIfNotExist(ctx, key, val, timeout...)
}

// SetIfNotExistFuncLock 当`key`不存在时，则使用函数`f`的结果设置`key`的值，返回是否设置成功，函数`f`是在读写互斥锁中执行的
//
//	当`timeout > 0`且`key`设置成功时，设置`key`的过期时间
//	当`force = true`时，可防止缓存穿透
func (rc *RedisCache) SetIfNotExistFuncLock(ctx context.Context, key string, f Func, force bool, timeout ...time.Duration) (ok bool, err error) {
	return rc.SetIfNotExistFunc(ctx, key, f, force, timeout...)
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
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
		args = append(args, key)
		args = append(args, members...)
		result, err = rc.client.Do(ctx, "SADD", args...)
	} else {
		args = append(args, members...)
		args = append(args, int64(timeout[0].Seconds()))
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
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
		result, err = rc.client.Do(ctx, "SISMEMBER", key, member)
	} else {
		result, err = rc.client.EvalSha(ctx, "SISMEMBER_EX", []string{key}, member, int64(timeout[0].Seconds()))
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
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
		result, err = rc.client.Do(ctx, "SMEMBERS", key)
	} else {
		result, err = rc.client.EvalSha(ctx, "SMEMBERS_EX", []string{key}, int64(timeout[0].Seconds()))
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
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
		result, err = rc.client.Do(ctx, "SPOP", key, count)
	} else {
		result, err = rc.client.EvalSha(ctx, "SPOP_EX", []string{key}, count, int64(timeout[0].Seconds()))
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
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
		args := make([]any, 0, len(keys))
		for _, v := range keys {
			args = append(args, v)
		}
		result, err = rc.client.Do(ctx, "SUNION", args...)
	} else {
		result, err = rc.client.EvalSha(ctx, "SUNION_EX", keys, int64(timeout[0].Seconds()))
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
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
		result, err = rc.client.Do(ctx, "SCARD", key)
	} else {
		result, err = rc.client.EvalSha(ctx, "SCARD_EX", []string{key}, int64(timeout[0].Seconds()))
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
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
		args = append(args, key)
		args = append(args, members...)
		result, err = rc.client.Do(ctx, "SREM", args...)
	} else {
		args = append(args, members...)
		args = append(args, int64(timeout[0].Seconds()))
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
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
		args = append(args, key)
		for k, v := range data {
			args = append(args, v, k)
		}
		result, err = rc.client.Do(ctx, "ZADD", args...)
	} else {
		for k, v := range data {
			args = append(args, v, k)
		}
		args = append(args, int64(timeout[0].Seconds()))
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
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
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
		args = append(args, int64(timeout[0].Seconds()))
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
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
		args = append(args, 0)
	} else {
		args = append(args, int64(timeout[0].Seconds()))
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
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
		result, err = rc.client.Do(ctx, "ZCARD", key)
	} else {
		result, err = rc.client.EvalSha(ctx, "ZCARD_EX", []string{key}, int64(timeout[0].Seconds()))
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
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
		result, err = rc.client.Do(ctx, "ZCOUNT", key, min, max)
	} else {
		result, err = rc.client.EvalSha(ctx, "ZCOUNT_EX", []string{key}, min, max, int64(timeout[0].Seconds()))
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
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
		result, err = rc.client.Do(ctx, "ZINCRBY", key, increment, member)
	} else {
		result, err = rc.client.EvalSha(ctx, "ZINCRBY_EX", []string{key}, increment, member, int64(timeout[0].Seconds()))
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
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
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
		args = append(args, int64(timeout[0].Seconds()))
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
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
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
			result, err = rc.client.EvalSha(ctx, "RANK_EX", []string{key}, "ZREVRANK", member, int64(timeout[0].Seconds()))
		} else {
			// score 值递增
			result, err = rc.client.EvalSha(ctx, "RANK_EX", []string{key}, "ZRANK", member, int64(timeout[0].Seconds()))
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
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
		args = append(args, key)
		args = append(args, members...)
		result, err = rc.client.Do(ctx, "ZREM", args...)
	} else {
		args = append(args, members...)
		args = append(args, int64(timeout[0].Seconds()))
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
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
		result, err = rc.client.Do(ctx, "ZREMRANGEBYRANK", key, start, stop)
	} else {
		result, err = rc.client.EvalSha(ctx, "ZREMRANGEBYRANK_EX", []string{key}, start, stop, int64(timeout[0].Seconds()))
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
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
		result, err = rc.client.Do(ctx, "ZREMRANGEBYSCORE", key, min, max)
	} else {
		result, err = rc.client.EvalSha(ctx, "ZREMRANGEBYSCORE_EX", []string{key}, min, max, int64(timeout[0].Seconds()))
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
	if len(timeout) == 0 || int64(timeout[0].Seconds()) <= 0 {
		result, err = rc.client.Do(ctx, "ZSCORE", key, member)
	} else {
		result, err = rc.client.EvalSha(ctx, "ZSCORE_EX", []string{key}, member, int64(timeout[0].Seconds()))
	}
	if err != nil {
		return
	}
	score = gtkconv.ToFloat64(result)
	return
}
