/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-27 15:03:40
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-09 13:11:37
 * @Description: 负载均衡器
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"sync"
	"time"
)

var (
	errEmptyAPIKeyList          = errors.New("api key list is empty")         // API密钥列表为空
	errNoAPIKeyAvailable        = errors.New("no api key available")          // 没有可用的API密钥
	errAPIKeyNotFound           = errors.New("api key not found")             // API密钥不存在
	errAPIKeyAlreadyExists      = errors.New("api key already exists")        // API密钥已存在
	errWeightMustBeGreaterThan0 = errors.New("weight must be greater than 0") // 权重必须大于0
)

// APIKey API密钥
type APIKey struct {
	Key       string // 密钥
	Times     uint32 // 请求次数
	Available bool   // 是否可用
	Weight    uint32 // 权重
}

// LoadBalancer 负载均衡器
type LoadBalancer struct {
	apiKeyList []*APIKey    // API密钥列表
	rng        *rand.Rand   // 随机数生成器
	mu         sync.RWMutex // 读写锁
}

// NewLoadBalancer 创建负载均衡器
func NewLoadBalancer(keyList []string) (lb *LoadBalancer) {
	now := time.Now().UnixNano()
	lb = &LoadBalancer{
		rng: rand.New(rand.NewPCG(uint64(now), uint64(now>>32))),
	}
	// 初始化API密钥列表
	for _, key := range keyList {
		lb.apiKeyList = append(lb.apiKeyList, &APIKey{
			Key:       key,
			Available: true,
			Weight:    1, // 默认权重为1
		})
	}
	return
}

// GetAPIKey 获取一个APIKey，使用最少连接算法
func (lb *LoadBalancer) GetAPIKey() (apiKey *APIKey, err error) {
	if len(lb.apiKeyList) == 0 {
		return nil, errEmptyAPIKeyList
	}
	// 选择使用次数最少的APIKey
	lb.mu.RLock()
	var (
		selectedAPIKey *APIKey
		minScore       = math.MaxFloat64
	)
	for _, v := range lb.apiKeyList {
		if v.Available {
			score := float64(v.Times) / float64(v.Weight)
			if score < minScore {
				selectedAPIKey = v
				minScore = score
			}
		}
	}
	lb.mu.RUnlock()
	// 如果未找到可用的APIKey，则返回错误
	if selectedAPIKey == nil {
		return nil, errNoAPIKeyAvailable
	}
	// 增加使用次数（需要写锁）
	lb.mu.Lock()
	selectedAPIKey.Times++
	lb.mu.Unlock()
	return selectedAPIKey, nil
}

// SetAvailability 设置指定APIKey的可用性
func (lb *LoadBalancer) SetAvailability(key string, available bool) (err error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	// 获取APIKey的索引
	index := slices.IndexFunc(lb.apiKeyList, func(apiKey *APIKey) bool {
		return apiKey.Key == key
	})
	// 如果APIKey不存在，则返回错误
	if index == -1 {
		return errAPIKeyNotFound
	}
	// 设置APIKey的可用性
	lb.apiKeyList[index].Available = available
	return
}

// RegisterAPIKey 注册新的APIKey
func (lb *LoadBalancer) RegisterAPIKey(key string) (err error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	// 检查是否已存在
	if isExists := slices.ContainsFunc(lb.apiKeyList, func(apiKey *APIKey) bool {
		return apiKey.Key == key
	}); isExists {
		return errAPIKeyAlreadyExists
	}

	if lb.apiKeyList == nil {
		lb.apiKeyList = make([]*APIKey, 0)
	}

	lb.apiKeyList = append(lb.apiKeyList, &APIKey{
		Key:       key,
		Available: true,
		Weight:    1, // 默认权重为1
	})
	return
}

// UnregisterAPIKey 注销APIKey
func (lb *LoadBalancer) UnregisterAPIKey(key string) (err error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	originalLen := len(lb.apiKeyList)
	// 删除该APIKey
	lb.apiKeyList = slices.DeleteFunc(lb.apiKeyList, func(apiKey *APIKey) bool {
		return apiKey.Key == key
	})

	// 如果长度没有变化，说明APIKey不存在
	if len(lb.apiKeyList) == originalLen {
		return errAPIKeyNotFound
	}
	return
}

// SetWeight 设置API权重，权重必须大于0
func (lb *LoadBalancer) SetWeight(key string, weight uint32) (err error) {
	if weight == 0 {
		return errWeightMustBeGreaterThan0
	}

	lb.mu.Lock()
	defer lb.mu.Unlock()

	// 获取APIKey的索引
	index := slices.IndexFunc(lb.apiKeyList, func(apiKey *APIKey) bool {
		return apiKey.Key == key
	})
	// 如果APIKey不存在，则返回错误
	if index == -1 {
		return errAPIKeyNotFound
	}
	// 设置APIKey的权重
	lb.apiKeyList[index].Weight = weight
	return
}

// SetAvailabilityForAll 设置所有APIKey的可用性
func (lb *LoadBalancer) SetAvailabilityForAll(available bool) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for _, apiKey := range lb.apiKeyList {
		apiKey.Available = available
	}
}

// GetAPIKeyList 获取所有APIKey的副本
func (lb *LoadBalancer) GetAPIKeyList() (apiKeyList []*APIKey) {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	apiKeyList = make([]*APIKey, len(lb.apiKeyList))
	for i, apiKey := range lb.apiKeyList {
		// 深拷贝：创建新的APIKey对象
		apiKeyList[i] = &APIKey{
			Key:       apiKey.Key,
			Times:     apiKey.Times,
			Available: apiKey.Available,
			Weight:    apiKey.Weight,
		}
	}
	return
}

// GetStats 获取负载均衡器统计信息
func (lb *LoadBalancer) GetStats() (stats map[string]any) {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	stats = make(map[string]any)
	var (
		totalAPIKey     = len(lb.apiKeyList)
		availableAPIKey = 0
		totalRequests   = uint32(0)
	)

	for _, apiKey := range lb.apiKeyList {
		if apiKey.Available {
			availableAPIKey++
		}
		totalRequests += apiKey.Times
	}

	stats["total_api_key"] = totalAPIKey
	stats["available_api_key"] = availableAPIKey
	stats["total_requests"] = totalRequests
	return stats
}
