/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-23 15:47:18
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-23 17:25:04
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkredis

import (
	"context"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"time"
)

// RedisLockOption redis 分布式锁选项
type RedisLockOption func(*redisLockConfig)

// redisLockConfig redis 分布式锁配置
type redisLockConfig struct {
	expiry        *time.Duration                    // 互斥锁的过期时间，默认为8秒
	tries         *int                              // 获取锁的重试次数，默认为32次
	delayFunc     func(tries int) (d time.Duration) // 重试间隔时间，默认为50ms到250ms之间的随机值
	driftFactor   *float64                          // 时钟漂移因子，默认值为0.01
	timeoutFactor *float64                          // 超时因子，默认值为0.05
	genValueFunc  func() (str string, err error)    // 自定义值生成器
	value         *string                           // 预设锁的值，支持锁的所有权转移
	shuffle       *bool                             // 是否随机化 Redis 连接池访问顺序，避免热点访问
	failFast      *bool                             // 是否快速失败模式
	setNXOnExtend *bool                             // 是否使用SETNX扩展锁
}

// WithExpiry 设置互斥锁的过期时间
func WithExpiry(expiry time.Duration) (opt RedisLockOption) {
	return func(cfg *redisLockConfig) {
		cfg.expiry = &expiry
	}
}

// WithTries 设置获取锁的重试次数
func WithTries(tries int) (opt RedisLockOption) {
	return func(cfg *redisLockConfig) {
		cfg.tries = &tries
	}
}

// WithRetryDelayFunc 设置覆盖默认的延迟行为
func WithRetryDelayFunc(delayFunc func(tries int) (d time.Duration)) (opt RedisLockOption) {
	return func(cfg *redisLockConfig) {
		cfg.delayFunc = delayFunc
	}
}

// WithSetNXOnExtend 优化了锁续期逻辑：若键存在则续期，若不存在则尝试在 Redis 中设置新键
func WithSetNXOnExtend(setNX bool) (opt RedisLockOption) {
	return func(cfg *redisLockConfig) {
		cfg.setNXOnExtend = &setNX
	}
}

// WithDriftFactor 设置时钟漂移因子
func WithDriftFactor(factor float64) (opt RedisLockOption) {
	return func(cfg *redisLockConfig) {
		cfg.driftFactor = &factor
	}
}

// WithTimeoutFactor 设置超时因子
func WithTimeoutFactor(factor float64) (opt RedisLockOption) {
	return func(cfg *redisLockConfig) {
		cfg.timeoutFactor = &factor
	}
}

// WithGenValueFunc 设置自定义值生成器
func WithGenValueFunc(genValueFunc func() (str string, err error)) (opt RedisLockOption) {
	return func(cfg *redisLockConfig) {
		cfg.genValueFunc = genValueFunc
	}
}

// WithValue 预设锁的值，支持锁的所有权转移
func WithValue(v string) (opt RedisLockOption) {
	return func(cfg *redisLockConfig) {
		cfg.value = &v
	}
}

// WithFailFast 设置是否快速失败模式
func WithFailFast(b bool) (opt RedisLockOption) {
	return func(cfg *redisLockConfig) {
		cfg.failFast = &b
	}
}

// WithShufflePools 设置是否随机化 Redis 连接池访问顺序，避免热点访问
func WithShufflePools(b bool) (opt RedisLockOption) {
	return func(cfg *redisLockConfig) {
		cfg.shuffle = &b
	}
}

// convertToRedsyncOptions 将自定义选项转换为redsync选项
func (cfg *redisLockConfig) convertToRedsyncOptions() []redsync.Option {
	var options []redsync.Option

	if cfg.expiry != nil {
		options = append(options, redsync.WithExpiry(*cfg.expiry))
	}
	if cfg.tries != nil {
		options = append(options, redsync.WithTries(*cfg.tries))
	}
	if cfg.delayFunc != nil {
		options = append(options, redsync.WithRetryDelayFunc(cfg.delayFunc))
	}
	if cfg.driftFactor != nil {
		options = append(options, redsync.WithDriftFactor(*cfg.driftFactor))
	}
	if cfg.timeoutFactor != nil {
		options = append(options, redsync.WithTimeoutFactor(*cfg.timeoutFactor))
	}
	if cfg.genValueFunc != nil {
		options = append(options, redsync.WithGenValueFunc(cfg.genValueFunc))
	}
	if cfg.value != nil {
		options = append(options, redsync.WithValue(*cfg.value))
	}
	if cfg.shuffle != nil {
		options = append(options, redsync.WithShufflePools(*cfg.shuffle))
	}
	if cfg.failFast != nil {
		options = append(options, redsync.WithFailFast(*cfg.failFast))
	}
	if cfg.setNXOnExtend != nil && *cfg.setNXOnExtend {
		options = append(options, redsync.WithSetNXOnExtend())
	}
	return options
}

// RedisLock redis 分布式锁
type RedisLock struct {
	mutex *redsync.Mutex
}

// NewRedisLock 创建 redis 分布式锁
func (rc *RedisClient) NewRedisLock(key string, options ...RedisLockOption) (rl *RedisLock) {
	cfg := &redisLockConfig{}
	for _, opt := range options {
		opt(cfg)
	}
	rs := redsync.New(goredis.NewPool(rc.client))
	redsyncOptions := cfg.convertToRedsyncOptions()
	return &RedisLock{
		mutex: rs.NewMutex(key, redsyncOptions...),
	}
}

// Extend 重置互斥锁的过期时间
func (rl *RedisLock) Extend(ctx context.Context) (ok bool, err error) {
	return rl.mutex.ExtendContext(ctx)
}

// Lock 获取锁
func (rl *RedisLock) Lock(ctx context.Context) (err error) {
	return rl.mutex.LockContext(ctx)
}

// Name 获取锁的名称
func (rl *RedisLock) Name() (name string) {
	return rl.mutex.Name()
}

// TryLock 尝试获取锁
func (rl *RedisLock) TryLock(ctx context.Context) (err error) {
	return rl.mutex.TryLockContext(ctx)
}

// Unlock 释放锁
func (rl *RedisLock) Unlock(ctx context.Context) (ok bool, err error) {
	return rl.mutex.UnlockContext(ctx)
}

// Until 获取锁的过期时间
func (rl *RedisLock) Until() (until time.Time) {
	return rl.mutex.Until()
}

// Value 获取锁的值
func (rl *RedisLock) Value() (value string) {
	return rl.mutex.Value()
}
