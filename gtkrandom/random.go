/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-04-01 23:27:01
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-26 21:09:41
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkrandom

import (
	"crypto/rand"
	"math/big"
	"sort"
)

// RandomWeight 随机权重
func RandomWeight(weights []int) (index int) {
	// 计算前缀和
	length := len(weights)
	if length == 0 {
		return
	}
	prefixSum := make([]int, length)
	prefixSum[0] = weights[0]
	for i := 1; i < length; i++ {
		prefixSum[i] = prefixSum[i-1] + weights[i]
	}
	// 生成一个随机权重值
	randomWeight, _ := rand.Int(rand.Reader, big.NewInt(int64(prefixSum[length-1])))
	// 使用二分查找算法找到随机权重值对应的下标
	return sort.SearchInts(prefixSum, int(randomWeight.Int64()))
}
