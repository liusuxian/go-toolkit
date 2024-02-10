/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-27 20:46:12
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-11 01:59:58
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkcache

import (
	"context"
	"time"
)

type Func func(ctx context.Context) (val any, err error)

// CustomCache 自定义缓存接口
type CustomCache interface {
	// Get 获取缓存
	Get(ctx context.Context, keys []string, timeout ...time.Duration) (val any, err error)
	// Set 设置缓存
	Set(ctx context.Context, keys []string, newVal any, timeout ...time.Duration) (val any, err error)
}

// ICache 缓存接口
type ICache interface {
	// Get 获取缓存
	//   当`timeout > 0`且缓存命中时，设置/重置`key`的过期时间
	Get(ctx context.Context, key string, timeout ...time.Duration) (val any, err error)
	// GetMap 批量获取缓存
	//   当`timeout > 0`且所有缓存都命中时，设置/重置所有`key`的过期时间，所有`key`过期时间相同
	GetMap(ctx context.Context, keys []string, timeout ...time.Duration) (data map[string]any, err error)
	// GetOrSet 检索并返回`key`的值，或者当`key`不存在时，则使用`newVal`设置`key`的值
	//   当`timeout > 0`时，设置/重置`key`的过期时间
	GetOrSet(ctx context.Context, key string, newVal any, timeout ...time.Duration) (val any, err error)
	// GetOrSetFunc 检索并返回`key`的值，或者当`key`不存在时，则使用函数`f`的结果设置`key`的值
	//   当`timeout > 0`时，设置/重置`key`的过期时间
	//   当`force = true`时，可防止缓存穿透
	GetOrSetFunc(ctx context.Context, key string, f Func, force bool, timeout ...time.Duration) (val any, err error)
	// GetOrSetFuncLock 检索并返回`key`的值，或者当`key`不存在时，则使用函数`f`的结果设置`key`的值，函数`f`是在读写互斥锁中执行的
	//   当`timeout > 0`时，设置/重置`key`的过期时间
	//   当`force = true`时，可防止缓存穿透
	GetOrSetFuncLock(ctx context.Context, key string, f Func, force bool, timeout ...time.Duration) (val any, err error)
	// CustomGetOrSetFunc 从缓存中获取指定键`keys`的值，如果缓存未命中，则使用函数`f`的结果设置`keys`的值
	//   当`timeout > 0`时，设置/重置`key`的过期时间
	//   当`force = true`时，可防止缓存穿透
	CustomGetOrSetFunc(ctx context.Context, keys []string, cc CustomCache, f Func, force bool, timeout ...time.Duration) (val any, err error)
	// CustomGetOrSetFuncLock 从缓存中获取指定键`keys`的值，如果缓存未命中，则使用函数`f`的结果设置`keys`的值，函数`f`是在读写互斥锁中执行的
	//   当`timeout > 0`时，设置/重置`key`的过期时间
	//   当`force = true`时，可防止缓存穿透
	CustomGetOrSetFuncLock(ctx context.Context, keys []string, cc CustomCache, f Func, force bool, timeout ...time.Duration) (val any, err error)
	// Set 设置缓存
	//   当`timeout > 0`时，设置/重置`key`的过期时间
	Set(ctx context.Context, key string, val any, timeout ...time.Duration) (err error)
	// SetMap 批量设置缓存，所有`key`的过期时间相同
	//   当`timeout > 0`时，设置/重置所有`key`的过期时间，所有`key`过期时间相同
	SetMap(ctx context.Context, data map[string]any, timeout ...time.Duration) (err error)
	// SetIfNotExist 当`key`不存在时，则使用`val`设置`key`的值，返回是否设置成功
	//   当`timeout > 0`且`key`设置成功时，设置`key`的过期时间
	SetIfNotExist(ctx context.Context, key string, val any, timeout ...time.Duration) (ok bool, err error)
	// SetIfNotExistFunc 当`key`不存在时，则使用函数`f`的结果设置`key`的值，返回是否设置成功
	//   当`timeout > 0`且`key`设置成功时，设置`key`的过期时间
	//   当`force = true`时，可防止缓存穿透
	SetIfNotExistFunc(ctx context.Context, key string, f Func, force bool, timeout ...time.Duration) (ok bool, err error)
	// SetIfNotExistFuncLock 当`key`不存在时，则使用函数`f`的结果设置`key`的值，返回是否设置成功，函数`f`是在读写互斥锁中执行的
	//   当`timeout > 0`且`key`设置成功时，设置`key`的过期时间
	//   当`force = true`时，可防止缓存穿透
	SetIfNotExistFuncLock(ctx context.Context, key string, f Func, force bool, timeout ...time.Duration) (ok bool, err error)
	// IsExist 缓存是否存在
	IsExist(ctx context.Context, key string) (isExist bool, err error)
	// Delete 删除缓存
	Delete(ctx context.Context, keys ...string) (err error)
	// GetExpire 获取缓存过期时间
	GetExpire(ctx context.Context, key string) (timeout time.Duration, err error)
	// Close 关闭缓存服务
	Close(ctx context.Context) (err error)
}

// IRedisCache Redis 缓存接口
type IRedisCache interface {
	ICache
	/* Set（集合）*/
	// SAdd 向集合添加一个或多个成员
	//   当`timeout > 0`时，设置/重置`key`的过期时间
	SAdd(ctx context.Context, key string, members []any, timeout ...time.Duration) (addCount int, err error)
	// SIsMember 判断 member 元素是否是集合 key 的成员
	//   当`timeout > 0`且`key`存在时，设置/重置`key`的过期时间
	SIsMember(ctx context.Context, key string, member any, timeout ...time.Duration) (isMember bool, err error)
	// SMembers 返回集合中的所有成员
	//   当`timeout > 0`且`key`存在时，设置/重置`key`的过期时间
	SMembers(ctx context.Context, key string, timeout ...time.Duration) (members []any, err error)
	// SPop 移除并返回集合中的一个或多个随机元素
	//   当`timeout > 0`且更新后的`key`存在时，设置/重置`key`的过期时间
	SPop(ctx context.Context, key string, count int, timeout ...time.Duration) (members []any, err error)
	// SUnion 返回所有给定集合的并集
	//   当`timeout > 0`且所有`key`都存在时，设置/重置所有`key`的过期时间，所有`key`过期时间相同
	SUnion(ctx context.Context, keys []string, timeout ...time.Duration) (members []any, err error)
	// SCard 获取集合的成员数
	//   当`timeout > 0`且`key`存在时，设置/重置`key`的过期时间
	SCard(ctx context.Context, key string, timeout ...time.Duration) (count int, err error)
	// SRem 移除集合中一个或多个成员
	//   当`timeout > 0`且更新后的`key`存在时，设置/重置`key`的过期时间
	SRem(ctx context.Context, key string, members []any, timeout ...time.Duration) (remCount int, err error)

	/* SortedSet（有序集合）*/
	// SSAdd 向有序集合添加一个或多个成员，或者更新已存在成员的分数
	//   当`timeout > 0`时，设置/重置`key`的过期时间
	SSAdd(ctx context.Context, key string, data map[any]float64, timeout ...time.Duration) (addCount int, err error)
	// SSRange 返回有序集合中指定区间内的成员
	//   当`timeout > 0`且`key`存在时，设置/重置`key`的过期时间
	SSRange(ctx context.Context, key string, start, stop int, isDescOrder, withScores bool, timeout ...time.Duration) (members []map[any]float64, err error)
	// SSPage 有序集合分页查询
	//   当`timeout > 0`且`key`存在时，设置/重置`key`的过期时间
	SSPage(ctx context.Context, key string, page, pageSize int, isDescOrder, withScores bool, timeout ...time.Duration) (total int, members []map[any]float64, err error)
	// SSCard 获取有序集合的成员数
	//   当`timeout > 0`且`key`存在时，设置/重置`key`的过期时间
	SSCard(ctx context.Context, key string, timeout ...time.Duration) (count int, err error)
	// SSCount 计算在有序集合中指定区间分数的成员数
	//   关于参数`min`和`max`的详细使用方法，请参考`redis`的`ZRANGEBYSCORE`命令
	//   当`timeout > 0`且`key`存在时，设置/重置`key`的过期时间
	SSCount(ctx context.Context, key, min, max string, timeout ...time.Duration) (count int, err error)
	// SSIncrby 有序集合中对指定成员的分数加上增量 increment
	//   当`timeout > 0`时，设置/重置`key`的过期时间
	SSIncrby(ctx context.Context, key string, increment float64, member any, timeout ...time.Duration) (score float64, err error)
	// SSRangeByScore 通过分数返回有序集合指定区间内的成员
	//   关于参数`min`和`max`的详细使用方法，请参考`redis`的`ZRANGEBYSCORE`命令
	//   当`timeout > 0`且`key`存在时，设置/重置`key`的过期时间
	SSRangeByScore(ctx context.Context, key, min, max string, isDescOrder, withScores bool, limit []int, timeout ...time.Duration) (members []map[any]float64, err error)
	// SSRank 返回有序集合中指定成员的排名
	//   当`timeout > 0`且`key`存在时，设置/重置`key`的过期时间
	SSRank(ctx context.Context, key string, member any, isDescOrder bool, timeout ...time.Duration) (rank int, err error)
	// SSRem 移除有序集合中的一个或多个成员
	//   当`timeout > 0`且更新后的`key`存在时，设置/重置`key`的过期时间
	SSRem(ctx context.Context, key string, members []any, timeout ...time.Duration) (remCount int, err error)
	// SSRemRangeByRank 移除有序集合中给定的排名区间的所有成员
	//   当`timeout > 0`且更新后的`key`存在时，设置/重置`key`的过期时间
	SSRemRangeByRank(ctx context.Context, key string, start, stop int, timeout ...time.Duration) (remCount int, err error)
	// SSRemRangeByScore 移除有序集合中给定的分数区间的所有成员
	//   关于参数`min`和`max`的详细使用方法，请参考`redis`的`ZRANGEBYSCORE`命令
	//   当`timeout > 0`且更新后的`key`存在时，设置/重置`key`的过期时间
	SSRemRangeByScore(ctx context.Context, key, min, max string, timeout ...time.Duration) (remCount int, err error)
	// SSScore 返回有序集中，成员的分数值
	//   当`timeout > 0`且`key`存在时，设置/重置`key`的过期时间
	SSScore(ctx context.Context, key string, member any, timeout ...time.Duration) (score float64, err error)

	/* TODO Hash（哈希表）*/

	/* TODO List（列表）*/
}
