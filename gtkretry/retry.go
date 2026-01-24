/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-12-09 17:23:44
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2026-01-24 15:02:47
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkretry

import (
	"context"
	"math"
	"math/rand/v2"
	"time"
)

// RetryFunc 重试函数的类型
type RetryFunc func(ctx context.Context) (err error)

// RetryCondition 重试条件函数
type RetryCondition func(attempt int, err error) (ok bool)

// RetryStrategy 重试策略
type RetryStrategy string

const (
	RetryStrategyFixed       RetryStrategy = "fixed"       // 固定间隔
	RetryStrategyLinear      RetryStrategy = "linear"      // 线性递增
	RetryStrategyExponential RetryStrategy = "exponential" // 指数退避
	RetryStrategyJitter      RetryStrategy = "jitter"      // 带抖动的指数退避
)

// RetryConfig 重试配置
type RetryConfig struct {
	MaxAttempts   int            `json:"max_attempts"`   // 最大重试次数（-1表示无限重试，0表示不重试只执行一次，>0表示重试指定次数），默认不重试
	Strategy      RetryStrategy  `json:"strategy"`       // 重试策略，可选值: fixed 固定间隔，linear 线性递增，exponential 指数退避，jitter 带抖动的指数退避，默认 exponential
	BaseDelay     time.Duration  `json:"base_delay"`     // 基础延迟时间，默认 1s
	MaxDelay      time.Duration  `json:"max_delay"`      // 最大延迟时间，默认 10s
	Multiplier    float64        `json:"multiplier"`     // 重试间隔倍数（用于指数退避），默认 2.0
	JitterPercent float64        `json:"jitter_percent"` // 抖动百分比（用于抖动策略，范围0-1，如0.1表示±10%），默认 0.1
	Condition     RetryCondition `json:"-"`              // 重试条件，默认总是重试
}

// Retry 重试
type Retry struct {
	config RetryConfig
}

// WithDefaults 返回填充了默认值的重试配置
func WithDefaults(config RetryConfig) RetryConfig {
	// 设置最大重试次数（-1表示无限重试，0表示不重试只执行一次，>0表示重试指定次数），默认不重试
	// 只有小于-1的值才会被设置为3
	if config.MaxAttempts < -1 {
		config.MaxAttempts = 3
	}
	// 设置重试策略，可选值: fixed 固定间隔，linear 线性递增，exponential 指数退避，jitter 带抖动的指数退避，默认 exponential
	if config.Strategy == "" {
		config.Strategy = RetryStrategyExponential
	}
	// 设置基础延迟时间，默认 1s
	if config.BaseDelay <= 0 {
		config.BaseDelay = 1 * time.Second
	}
	// 设置最大延迟时间，默认 10s
	if config.MaxDelay <= 0 {
		config.MaxDelay = 10 * time.Second
	}
	// 确保最大延迟不小于基础延迟
	if config.MaxDelay < config.BaseDelay {
		config.MaxDelay = config.BaseDelay * 10 // 设置为基础延迟的10倍
	}
	// 设置重试间隔倍数（用于指数退避），默认 2.0
	if config.Multiplier <= 0 {
		config.Multiplier = 2.0
	}
	// 设置抖动百分比（用于抖动策略，范围0-1，如0.1表示±10%），默认 0.1
	if config.JitterPercent <= 0 || config.JitterPercent > 1 {
		config.JitterPercent = 0.1 // 默认±10%抖动
	}
	// 设置重试条件，默认总是重试
	if config.Condition == nil {
		config.Condition = func(attempt int, err error) (ok bool) {
			return true
		}
	}
	return config
}

// NewRetry 创建重试实例
func NewRetry(config RetryConfig) (r *Retry) {
	return &Retry{
		config: WithDefaults(config),
	}
}

// Do 执行重试
func (r *Retry) Do(ctx context.Context, f RetryFunc) (err error) {
	attempt := 0
	for {
		// 执行函数
		err = f(ctx)
		// 如果成功或者不需要重试，直接返回
		if err == nil || !r.config.Condition(attempt, err) {
			return
		}
		// MaxAttempts = 0: 只执行一次，不重试
		if r.config.MaxAttempts == 0 {
			break
		}
		// MaxAttempts > 0: 已达到最大重试次数
		if r.config.MaxAttempts > 0 && attempt >= r.config.MaxAttempts {
			break
		}
		// MaxAttempts = -1: 无限重试，继续执行
		// 计算延迟时间
		delay := r.calculateDelay(attempt + 1)
		// 等待延迟时间
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			// 继续下次重试
		}
		attempt++
	}
	return
}

// calculateDelay 计算延迟时间
func (r *Retry) calculateDelay(attempt int) (delay time.Duration) {
	switch r.config.Strategy {
	case RetryStrategyFixed:
		delay = r.config.BaseDelay
	case RetryStrategyLinear:
		if attempt > 0 && r.config.BaseDelay > 0 {
			// 防止整数溢出，先检查是否会超过最大延迟
			maxAttempts := int64(r.config.MaxDelay / r.config.BaseDelay)
			if int64(attempt) > maxAttempts {
				delay = r.config.MaxDelay
			} else {
				delay = min(time.Duration(attempt)*r.config.BaseDelay, r.config.MaxDelay)
			}
		} else {
			delay = r.config.BaseDelay
		}
	case RetryStrategyExponential:
		delay = r.calculateExponentialDelay(attempt)
	case RetryStrategyJitter:
		exponentialDelay := r.calculateExponentialDelay(attempt)
		// 添加双向随机抖动（基于配置的百分比）
		jitterRange := float64(exponentialDelay) * r.config.JitterPercent
		jitter := time.Duration((rand.Float64() - 0.5) * 2 * jitterRange)
		// 确保抖动后的延迟在合理范围内
		delay = max(exponentialDelay+jitter, exponentialDelay/2)
	default:
		delay = r.config.BaseDelay
	}
	// 确保不超过最大延迟
	delay = min(delay, r.config.MaxDelay)
	// 确保延迟不为负数或过小
	delay = max(delay, r.config.BaseDelay/2)
	return
}

// calculateExponentialDelay 计算指数退避延迟
func (r *Retry) calculateExponentialDelay(attempt int) (delay time.Duration) {
	// 如果基础延迟、最大延迟小于等于0或重试间隔倍数小于等于1，返回基础延迟
	if r.config.BaseDelay <= 0 || r.config.MaxDelay <= 0 || r.config.Multiplier <= 1 {
		return r.config.BaseDelay
	}
	// 计算最大尝试次数
	maxExponent := math.Log(float64(r.config.MaxDelay/r.config.BaseDelay)) / math.Log(r.config.Multiplier)
	// 如果最大尝试次数为无穷大或NaN，或者尝试次数大于最大尝试次数，返回最大延迟
	if math.IsInf(maxExponent, 0) || math.IsNaN(maxExponent) || float64(attempt-1) > maxExponent {
		return r.config.MaxDelay
	}
	// 计算延迟时间
	delayFloat := float64(r.config.BaseDelay) * math.Pow(r.config.Multiplier, float64(attempt-1))
	// 如果延迟时间为无穷大或NaN，返回最大延迟
	if math.IsInf(delayFloat, 0) || math.IsNaN(delayFloat) {
		return r.config.MaxDelay
	}
	// 返回延迟时间
	return time.Duration(delayFloat)
}
