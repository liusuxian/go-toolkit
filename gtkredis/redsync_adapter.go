/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2026-01-26 17:44:24
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2026-01-26 18:59:46
 * @Description: https://github.com/go-redsync/redsync 适配器，让 gtkredis 可以使用 redsync 分布式锁
 *
 * Copyright (c) 2026 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkredis

import (
	"context"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis"
	"strings"
	"time"
)

// redisPool 实现 redsync 的 Pool 接口
type redisPool struct {
	rc *RedisClient
}

// redisConn 实现 redsync 的 Conn 接口
type redisConn struct {
	ctx context.Context
	rc  *RedisClient
}

// Get 实现 Pool 接口的 Get 方法
func (p *redisPool) Get(ctx context.Context) (redis.Conn, error) {
	return &redisConn{
		ctx: ctx,
		rc:  p.rc,
	}, nil
}

// Get 实现 Conn 接口的 Get 方法
func (c *redisConn) Get(name string) (string, error) {
	value, err := c.rc.client.Get(c.ctx, name).Result()
	return value, noErrNil(err)
}

// Set 实现 Conn 接口的 Set 方法
func (c *redisConn) Set(name, value string) (bool, error) {
	reply, err := c.rc.client.Set(c.ctx, name, value, 0).Result()
	return reply == "OK", err
}

// SetNX 实现 Conn 接口的 SetNX 方法
func (c *redisConn) SetNX(name, value string, expiry time.Duration) (result bool, err error) {
	return c.rc.client.SetNX(c.ctx, name, value, expiry).Result()
}

// Eval 实现 Conn 接口的 Eval 方法
func (c *redisConn) Eval(script *redis.Script, keysAndArgs ...any) (any, error) {
	var (
		keys = make([]string, script.KeyCount)
		args = keysAndArgs
	)

	if script.KeyCount > 0 {
		for i := range script.KeyCount {
			keys[i] = keysAndArgs[i].(string)
		}
		args = keysAndArgs[script.KeyCount:]
	}

	value, err := c.rc.client.EvalSha(c.ctx, script.Hash, keys, args...).Result()
	if err != nil && strings.Contains(err.Error(), "NOSCRIPT ") {
		value, err = c.rc.client.Eval(c.ctx, script.Src, keys, args...).Result()
	}
	return value, noErrNil(err)
}

// ScriptLoad 实现 Conn 接口的 ScriptLoad 方法
func (c *redisConn) ScriptLoad(script *redis.Script) error {
	value, err := c.rc.client.ScriptLoad(c.ctx, script.Src).Result()
	if err == nil {
		script.Hash = value
	}
	return err
}

// PTTL 实现 Conn 接口的 PTTL 方法
func (c *redisConn) PTTL(name string) (time.Duration, error) {
	return c.rc.client.PTTL(c.ctx, name).Result()
}

// Close 实现 Conn 接口的 Close 方法
func (c *redisConn) Close() error {
	return nil
}

// NewRedsync 创建一个 Redsync 实例
func (rc *RedisClient) NewRedsync() *redsync.Redsync {
	return redsync.New(&redisPool{rc: rc})
}

// NewMutex 创建一个 Redsync 分布式锁
func (rc *RedisClient) NewMutex(name string, options ...redsync.Option) *redsync.Mutex {
	return rc.NewRedsync().NewMutex(name, options...)
}
