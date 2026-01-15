/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-27 20:46:12
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2026-01-15 18:22:16
 * @Description: 缓存接口
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkcache

import (
	"context"
	"time"
)

// Func 函数类型
type Func func(ctx context.Context) (val any, err error)

// ICustomCache 自定义缓存接口
type ICustomCache interface {
	// Get 获取缓存
	//   当`timeout > 0`且缓存命中时，设置/重置`keys`的过期时间
	Get(ctx context.Context, keys []string, args []any, timeout ...time.Duration) (val any, err error)
	// Add 添加缓存
	//   当`keys`已存在且未过期时，返回现有值（不修改）
	//   当`keys`不存在或已过期时，设置新值并返回新值
	//   当`timeout > 0`时，设置/重置`keys`的过期时间
	Add(ctx context.Context, keys []string, args []any, newVal any, timeout ...time.Duration) (val any, err error)
}

// IBatchGetter 批量获取构建器接口
type IBatchGetter interface {
	// Add 添加一个 key 到批量获取队列
	//   当`timeout > 0`且该`key`缓存命中时，设置/重置该`key`的过期时间（会覆盖默认过期时间）
	//   当`timeout`未指定或`timeout <= 0`时，该`key`使用`SetDefaultTimeout`设置的默认过期时间
	//   返回构建器自身，支持链式调用
	Add(ctx context.Context, key string, timeout ...time.Duration) (batchGetter IBatchGetter)
	// SetDefaultTimeout 设置默认过期时间（对所有未单独设置过期时间的 key 生效）
	//   当`timeout > 0`且缓存命中时，所有未单独指定过期时间的`key`将使用此默认过期时间
	//   当`timeout <= 0`时，所有未单独指定过期时间的`key`将保持原有的过期时间
	//   返回构建器自身，支持链式调用
	SetDefaultTimeout(ctx context.Context, timeout time.Duration) (batchGetter IBatchGetter)
	// Execute 执行批量获取操作
	//   返回 map[key]value，不存在或已过期的`key`不会出现在结果`map`中
	//   执行成功后，自动清空构建器中的数据（不建议继续使用该构建器）
	//   执行失败时，保留构建器中的数据，可以直接再次调用本方法进行重试
	//   建议：为每次批量操作创建新的构建器实例
	Execute(ctx context.Context) (values map[string]any, err error)
}

// IBatchSetter 批量设置构建器接口
type IBatchSetter interface {
	// Add 添加一个 key-value 对到批量设置队列
	//   当`timeout > 0`时，设置该`key`的过期时间（会覆盖默认过期时间）
	//   当`timeout`未指定或`timeout <= 0`时，该`key`使用`SetDefaultTimeout`设置的默认过期时间
	//   返回构建器自身，支持链式调用
	Add(ctx context.Context, key string, val any, timeout ...time.Duration) (batchSetter IBatchSetter)
	// SetDefaultTimeout 设置默认过期时间（对所有未单独设置过期时间的 key 生效）
	//   当`timeout > 0`时，所有未单独指定过期时间的`key`将使用此默认过期时间
	//   当`timeout <= 0`时，所有未单独指定过期时间的`key`将保持原有的过期时间
	//   返回构建器自身，支持链式调用
	SetDefaultTimeout(ctx context.Context, timeout time.Duration) (batchSetter IBatchSetter)
	// Execute 执行批量设置操作
	//   执行成功后，自动清空构建器中的数据（不建议继续使用该构建器）
	//   执行失败时，保留构建器中的数据，可以直接再次调用本方法进行重试
	//   建议：为每次批量操作创建新的构建器实例
	Execute(ctx context.Context) (err error)
}

// ICache 缓存接口
type ICache interface {
	// Get 获取缓存
	//   当`timeout > 0`且缓存命中时，设置/重置`key`的过期时间
	Get(ctx context.Context, key string, timeout ...time.Duration) (val any, err error)
	// GetMap 批量获取缓存
	//   当`timeout > 0`且所有缓存都命中时，设置/重置所有`key`的过期时间，所有`key`过期时间相同
	//   注意：如需为每个`key`设置/重置不同的过期时间，请使用`BatchGet`
	GetMap(ctx context.Context, keys []string, timeout ...time.Duration) (data map[string]any, err error)
	// BatchGet 创建批量获取构建器
	//   支持为每个`key`设置/重置不同的过期时间
	//   当所有`key`使用相同过期时间时，可以使用更简洁的`GetMap`方法
	//   当`capacity > 0`时，预分配指定容量以优化性能
	//   返回构建器实例，支持链式调用
	BatchGet(ctx context.Context, capacity ...int) (batchGetter IBatchGetter)
	// GetOrSet 检索并返回`key`的值，或者当`key`不存在时，则使用`newVal`设置`key`的值
	//   当`timeout > 0`时，设置/重置`key`的过期时间
	GetOrSet(ctx context.Context, key string, newVal any, timeout ...time.Duration) (val any, err error)
	// GetOrSetFunc 检索并返回`key`的值，或者当`key`不存在时，则使用函数`f`的结果设置`key`的值
	//	 当`timeout > 0`时，设置/重置`key`的过期时间
	//	 当`force = true`时，可防止缓存穿透（即使`f`返回`nil`也会缓存）
	//	 注意：使用`singleflight`机制确保相同`key`的函数`f`只执行一次，其他并发请求等待并共享第一个请求的执行结果，有效防止缓存击穿
	GetOrSetFunc(ctx context.Context, key string, f Func, force bool, timeout ...time.Duration) (val any, err error)
	// CustomGetOrSetFunc 从缓存中获取指定键`keys`的值，如果缓存未命中，则使用函数`f`的结果设置`keys`的值
	//	 当`timeout > 0`时，设置/重置`key`的过期时间
	//	 当`force = true`时，可防止缓存穿透（即使`f`返回`nil`也会缓存）
	//	 注意：使用`singleflight`机制确保相同`key`的函数`f`只执行一次，其他并发请求等待并共享第一个请求的执行结果，有效防止缓存击穿
	CustomGetOrSetFunc(ctx context.Context, keys []string, args []any, cc ICustomCache, f Func, force bool, timeout ...time.Duration) (val any, err error)
	// Set 设置缓存
	//   当`timeout > 0`时，设置/重置`key`的过期时间
	Set(ctx context.Context, key string, val any, timeout ...time.Duration) (err error)
	// SetMap 批量设置缓存，所有`key`的过期时间相同
	//   当`timeout > 0`时，设置/重置所有`key`的过期时间，所有`key`过期时间相同
	//   注意：如需为每个`key`设置不同的过期时间，请使用`BatchSet`
	SetMap(ctx context.Context, data map[string]any, timeout ...time.Duration) (err error)
	// BatchSet 创建批量设置构建器
	//   支持为每个`key`设置不同的过期时间
	//   当所有`key`使用相同过期时间时，可以使用更简洁的`SetMap`方法
	//   当`capacity > 0`时，预分配指定容量以优化性能
	//   返回构建器实例，支持链式调用
	BatchSet(ctx context.Context, capacity ...int) (batchSetter IBatchSetter)
	// SetIfNotExist 当`key`不存在时，则使用`val`设置`key`的值，返回是否设置成功
	//   当`timeout > 0`且`key`设置成功时，设置`key`的过期时间
	SetIfNotExist(ctx context.Context, key string, val any, timeout ...time.Duration) (ok bool, err error)
	// SetIfNotExistFunc 当`key`不存在时，则使用函数`f`的结果设置`key`的值，返回是否设置成功
	//	 当`timeout > 0`且`key`设置成功时，设置`key`的过期时间
	//	 当`force = true`时，可防止缓存穿透（即使`f`返回`nil`也会缓存）
	//	 注意：使用`singleflight`机制确保相同`key`的函数`f`只执行一次，其他并发请求等待并共享第一个请求的执行结果，有效防止缓存击穿
	SetIfNotExistFunc(ctx context.Context, key string, f Func, force bool, timeout ...time.Duration) (ok bool, err error)
	// Update 当`key`存在时，则使用`val`更新`key`的值，返回`key`的旧值
	//   当`timeout > 0`且`key`更新成功时，更新`key`的过期时间
	Update(ctx context.Context, key string, val any, timeout ...time.Duration) (oldVal any, isExist bool, err error)
	// UpdateExpire 当`key`存在时，则更新`key`的过期时间，返回`key`的旧的过期时间值
	//   当`key`不存在时，则返回-1
	//   当`key`存在但没有设置过期时间时，则返回0
	//   当`key`存在且设置了过期时间时，则返回过期时间
	//   当`timeout > 0`且`key`存在时，更新`key`的过期时间
	UpdateExpire(ctx context.Context, key string, timeout time.Duration) (oldTimeout time.Duration, err error)
	// IsExist 缓存是否存在
	IsExist(ctx context.Context, key string) (isExist bool, err error)
	// Size 缓存中的key数量
	Size(ctx context.Context) (size int, err error)
	// Delete 删除缓存
	Delete(ctx context.Context, keys ...string) (err error)
	// GetExpire 获取缓存`key`的过期时间
	//   当`key`不存在时，则返回-1
	//   当`key`存在但没有设置过期时间时，则返回0
	//   当`key`存在且设置了过期时间时，则返回过期时间
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

// IWechatCache 微信缓存接口（适配 github.com/silenceper/wechat/v2 库的缓存）
type IWechatCache interface {
	Get(key string) (val any)                                   // 获取缓存
	Set(key string, val any, timeout time.Duration) (err error) // 设置缓存
	IsExist(key string) (isExist bool)                          // 缓存是否存在
	Delete(key string) (err error)                              // 删除缓存
}

// IContextWechatCache 上下文微信缓存接口（适配 github.com/silenceper/wechat/v2 库的缓存）
type IContextWechatCache interface {
	IWechatCache
	GetContext(ctx context.Context, key string) (val any)                                   // 获取缓存
	SetContext(ctx context.Context, key string, val any, timeout time.Duration) (err error) // 设置缓存
	IsExistContext(ctx context.Context, key string) (isExist bool)                          // 缓存是否存在
	DeleteContext(ctx context.Context, key string) (err error)                              // 删除缓存
}
