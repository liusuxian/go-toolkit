/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-12-16 23:11:19
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-18 01:33:32
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkcache

import (
	"sync"
	"time"
)

// cacheItem 缓存项
type cacheItem struct {
	Object     any   // 缓存的值
	Expiration int64 // 过期时间（Unix纳秒时间戳），0 表示永不过期
}

// isExpired 检查是否过期
func (item *cacheItem) isExpired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

// keyAndValue 键值对
type keyAndValue struct {
	key   string
	value any
}

// cacheShard 缓存分片
type cacheShard struct {
	items     map[string]*cacheItem
	mu        sync.RWMutex
	onEvicted func(string, any)
}

// OnEvicted 设置删除回调函数
func (cs *cacheShard) OnEvicted(f func(string, any)) {
	cs.mu.Lock()
	cs.onEvicted = f
	cs.mu.Unlock()
}

// DeleteExpired 删除过期缓存
func (cs *cacheShard) DeleteExpired() {
	var evictedItems []keyAndValue
	now := time.Now().UnixNano()
	cs.mu.Lock()
	for k, v := range cs.items {
		if v.Expiration > 0 && now > v.Expiration {
			ov, evicted := cs.delete(k)
			if evicted {
				evictedItems = append(evictedItems, keyAndValue{k, ov})
			}
		}
	}
	cs.mu.Unlock()
	for _, v := range evictedItems {
		cs.onEvicted(v.key, v.value)
	}
}

// delete 删除缓存
func (cs *cacheShard) delete(k string) (any, bool) {
	if cs.onEvicted != nil {
		if v, found := cs.items[k]; found {
			delete(cs.items, k)
			return v.Object, true
		}
	}
	delete(cs.items, k)
	return nil, false
}

// get 获取缓存
func (cs *cacheShard) get(k string) any {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	item, found := cs.items[k]
	if !found || item.isExpired() {
		return nil
	}
	return item.Object
}

// getWithExpiration 获取缓存并重置过期时间
func (cs *cacheShard) getWithExpiration(k string, expiration int64) any {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	item, found := cs.items[k]
	if !found || item.isExpired() {
		return nil
	}

	item.Expiration = expiration
	return item.Object
}

// getOrSetWithExpiration 获取缓存并重置过期时间或设置新值并设置过期时间
func (cs *cacheShard) getOrSetWithExpiration(k string, v any, expiration int64) any {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	item, found := cs.items[k]
	if found && !item.isExpired() {
		if expiration > 0 {
			item.Expiration = expiration
		}
		return item.Object
	}

	cs.items[k] = &cacheItem{
		Object:     v,
		Expiration: expiration,
	}
	return v
}

// getOrSetWithValue 获取缓存或设置指定值
func (cs *cacheShard) getOrSetWithValue(k string, v any, expiration int64) any {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	item, found := cs.items[k]
	if found && !item.isExpired() {
		return item.Object
	}

	cs.items[k] = &cacheItem{
		Object:     v,
		Expiration: expiration,
	}
	return v
}

// setWithExpiration 设置缓存并设置过期时间
func (cs *cacheShard) setWithExpiration(k string, v any, expiration int64) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.items[k] = &cacheItem{
		Object:     v,
		Expiration: expiration,
	}
}

// janitor 清理器
type janitor struct {
	interval time.Duration
	stop     chan bool
}

// Run 启动清理任务
func (j *janitor) Run(mc *memoryCache) {
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mc.DeleteExpired()
		case <-j.stop:
			return
		}
	}
}

// stopJanitor 停止清理任务
func stopJanitor(mc *MemoryCache) {
	mc.janitor.stop <- true
}

// runJanitor 启动清理器
func runJanitor(mc *memoryCache, interval time.Duration) {
	j := &janitor{
		interval: interval,
		stop:     make(chan bool, 1),
	}
	mc.janitor = j
	go j.Run(mc)
}
