/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-12-16 20:38:43
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-17 19:22:52
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkcache

import (
	"context"
	"crypto/rand"
	"math"
	"math/big"
	insecurerand "math/rand/v2"
	"runtime"
	"time"
)

const (
	defaultShardCount = 32 // 默认分片数量
)

// MemoryCache 内存缓存
type MemoryCache struct {
	*memoryCache
}

// memoryCache 内存缓存
type memoryCache struct {
	shards     []*cacheShard // 分片
	shardCount uint32        // 分片数量
	seed       uint32        // 随机种子
	janitor    *janitor      // 清理器
}

// NewMemoryCache 创建内存缓存
//
//	cleanupInterval: 清理间隔，0 表示懒惰删除
func NewMemoryCache(cleanupInterval ...time.Duration) *MemoryCache {
	return NewMemoryCacheWithShards(defaultShardCount, cleanupInterval...)
}

// NewMemoryCacheWithShards 创建指定分片数的内存缓存
//
//	shardCount: 分片数量，建议设置为 2 的幂次（如 8, 16, 32, 64）
//	cleanupInterval: 清理间隔，0 表示懒惰删除
func NewMemoryCacheWithShards(shardCount int, cleanupInterval ...time.Duration) *MemoryCache {
	var (
		mc = newMemoryCacheWithShards(shardCount)
		MC = &MemoryCache{mc}
	)
	// 启动清理器
	if len(cleanupInterval) > 0 && cleanupInterval[0] > 0 {
		runJanitor(mc, cleanupInterval[0])
		runtime.SetFinalizer(MC, stopJanitor)
	}
	return MC
}

// DeleteExpired 删除过期缓存项
func (mc *memoryCache) DeleteExpired() {
	for _, v := range mc.shards {
		v.DeleteExpired()
	}
}

// Get 获取缓存
//
//	当`timeout > 0`且缓存命中时，设置/重置`key`的过期时间
func (mc *memoryCache) Get(ctx context.Context, key string, timeout ...time.Duration) (val any, err error) {
	shard := mc.getShard(key)
	// 获取缓存项并重置过期时间
	if expiration := mc.getExpiration(timeout...); expiration > 0 {
		return shard.getItemWithExpiration(key, expiration), nil
	}
	// 获取缓存项
	return shard.getItem(key), nil
}

// GetMap 批量获取缓存
//
//	当`timeout > 0`且所有缓存都命中时，设置/重置所有`key`的过期时间，所有`key`过期时间相同
func (mc *memoryCache) GetMap(ctx context.Context, keys []string, timeout ...time.Duration) (data map[string]any, err error) {
	if len(keys) == 0 {
		return
	}
	// 批量获取缓存项并重置过期时间
	if expiration := mc.getExpiration(timeout...); expiration > 0 {
		return mc.getItemMapWithExpiration(keys, expiration), nil
	}
	// 批量获取缓存项
	return mc.getItemMap(keys), nil
}

// GetOrSet 检索并返回`key`的值，或者当`key`不存在时，则使用`newVal`设置`key`的值
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
func (mc *memoryCache) GetOrSet(ctx context.Context, key string, newVal any, timeout ...time.Duration) (val any, err error) {
	shard := mc.getShard(key)
	// 获取缓存项并重置过期时间或设置新值并设置过期时间
	return shard.getOrSetItemWithExpiration(key, newVal, mc.getExpiration(timeout...)), nil
}

// Set 设置缓存
//
//	当`timeout > 0`时，设置/重置`key`的过期时间
func (mc *memoryCache) Set(ctx context.Context, key string, val any, timeout ...time.Duration) (err error) {
	shard := mc.getShard(key)
	// 设置缓存项并设置过期时间
	shard.setItemWithExpiration(key, val, mc.getExpiration(timeout...))
	return
}

// SetMap 批量设置缓存，所有`key`的过期时间相同
//
//	当`timeout > 0`时，设置/重置所有`key`的过期时间，所有`key`过期时间相同
func (mc *memoryCache) SetMap(ctx context.Context, data map[string]any, timeout ...time.Duration) (err error) {
	if len(data) == 0 {
		return
	}
	// 批量设置缓存项并设置过期时间
	mc.setItemMapWithExpiration(data, mc.getExpiration(timeout...))
	return
}

// getItemMap 批量获取缓存项
func (mc *memoryCache) getItemMap(keys []string) (data map[string]any) {
	data = make(map[string]any)
	// 按分片分组 keys
	shardMap := make(map[*cacheShard][]string)
	for _, k := range keys {
		shard := mc.getShard(k)
		shardMap[shard] = append(shardMap[shard], k)
	}
	for shard, shardKeyList := range shardMap {
		shard.mu.RLock()
		for _, k := range shardKeyList {
			item, found := shard.items[k]
			if !found || item.isExpired() {
				data[k] = nil
				continue
			}
			data[k] = item.Object
		}
		shard.mu.RUnlock()
	}
	return
}

// getItemMapWithExpiration 批量获取缓存项并重置过期时间
func (mc *memoryCache) getItemMapWithExpiration(keys []string, expiration int64) (data map[string]any) {
	data = make(map[string]any)
	// 按分片分组 keys
	shardMap := make(map[*cacheShard][]string)
	for _, k := range keys {
		shard := mc.getShard(k)
		shardMap[shard] = append(shardMap[shard], k)
	}
	allHit := true
	// 读取数据
	for shard, shardKeyList := range shardMap {
		shard.mu.RLock()
		for _, k := range shardKeyList {
			item, found := shard.items[k]
			if !found || item.isExpired() {
				data[k] = nil
				allHit = false
				continue
			}
			data[k] = item.Object
		}
		shard.mu.RUnlock()
	}
	// 只有所有 key 都命中时才重置过期时间
	if allHit {
		for shard, shardKeyList := range shardMap {
			shard.mu.Lock()
			for _, k := range shardKeyList {
				item, found := shard.items[k]
				if !found {
					continue
				}
				if item.isExpired() {
					delete(shard.items, k)
					continue
				}
				item.Expiration = expiration
			}
			shard.mu.Unlock()
		}
	}
	return
}

// setItemMapWithExpiration 批量设置缓存项并设置过期时间
func (mc *memoryCache) setItemMapWithExpiration(data map[string]any, expiration int64) {
	// 按分片分组
	shardMap := make(map[*cacheShard]map[string]any)
	for k, v := range data {
		shard := mc.getShard(k)
		if shardMap[shard] == nil {
			shardMap[shard] = make(map[string]any)
		}
		shardMap[shard][k] = v
	}
	// 批量设置每个分片
	for shard, shardData := range shardMap {
		shard.mu.Lock()
		for k, v := range shardData {
			shard.items[k] = &cacheItem{
				Object:     v,
				Expiration: expiration,
			}
		}
		shard.mu.Unlock()
	}
}

// getShard 根据 key 获取对应的分片
func (mc *memoryCache) getShard(k string) *cacheShard {
	return mc.shards[djb33(mc.seed, k)%mc.shardCount]
}

// getExpiration 获取过期时间戳（Unix纳秒时间戳），0 表示永不过期
func (mc *memoryCache) getExpiration(timeout ...time.Duration) int64 {
	if len(timeout) > 0 && timeout[0] > 0 {
		return time.Now().Add(timeout[0]).UnixNano()
	}
	return 0
}

// newMemoryCacheWithShards 创建指定分片数的内存缓存
func newMemoryCacheWithShards(shardCount int) (mc *memoryCache) {
	// 分片数量
	if shardCount <= 0 {
		shardCount = defaultShardCount
	}
	// 生成随机种子
	var (
		seed uint32
		max  = big.NewInt(0).SetUint64(uint64(math.MaxUint32))
	)
	rnd, err := rand.Int(rand.Reader, max)
	if err != nil {
		seed = insecurerand.Uint32()
	} else {
		seed = uint32(rnd.Uint64())
	}
	// 创建内存缓存对象
	mc = &memoryCache{
		shards:     make([]*cacheShard, shardCount),
		shardCount: uint32(shardCount),
		seed:       seed,
	}
	// 创建分片
	for i := 0; i < shardCount; i++ {
		mc.shards[i] = &cacheShard{
			items: make(map[string]*cacheItem),
		}
	}
	return
}

// djb33 哈希算法
func djb33(seed uint32, k string) uint32 {
	var (
		l = uint32(len(k))
		d = 5381 + seed + l
		i = uint32(0)
	)
	if l >= 4 {
		for i < l-4 {
			d = (d * 33) ^ uint32(k[i])
			d = (d * 33) ^ uint32(k[i+1])
			d = (d * 33) ^ uint32(k[i+2])
			d = (d * 33) ^ uint32(k[i+3])
			i += 4
		}
	}
	switch l - i {
	case 1:
	case 2:
		d = (d * 33) ^ uint32(k[i])
	case 3:
		d = (d * 33) ^ uint32(k[i])
		d = (d * 33) ^ uint32(k[i+1])
	case 4:
		d = (d * 33) ^ uint32(k[i])
		d = (d * 33) ^ uint32(k[i+1])
		d = (d * 33) ^ uint32(k[i+2])
	}
	return d ^ (d >> 16)
}
