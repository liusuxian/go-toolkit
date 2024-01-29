/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-27 20:53:08
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-29 16:32:47
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkcache_test

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/liusuxian/go-toolkit/gtkcache"
	"github.com/liusuxian/go-toolkit/gtkconv"
	"github.com/liusuxian/go-toolkit/gtkredis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisCache(t *testing.T) {
	var (
		ctx   = context.Background()
		r     = miniredis.RunT(t)
		cache *gtkcache.RedisCache
	)
	cache = gtkcache.NewRedisCache(ctx, func(cc *gtkredis.ClientConfig) {
		cc.Addr = r.Addr()
		cc.DB = 1
		cc.Password = ""
	})
	var (
		assert  = assert.New(t)
		val     any
		err     error
		isExist bool
		timeout time.Duration
		data    map[string]any
	)
	val, err = cache.Get(ctx, "test_key_1")
	assert.NoError(err)
	assert.Nil(val)
	isExist, err = cache.IsExist(ctx, "test_key_1")
	assert.NoError(err)
	assert.False(isExist)
	timeout, err = cache.GetExpire(ctx, "test_key_1")
	assert.NoError(err)
	assert.Equal(float64(-2), timeout.Seconds())

	err = cache.Set(ctx, "test_key_1", 100, time.Second*10)
	assert.NoError(err)
	val, err = cache.Get(ctx, "test_key_1")
	assert.NoError(err)
	assert.Equal(100, gtkconv.ToInt(val))
	isExist, err = cache.IsExist(ctx, "test_key_1")
	assert.NoError(err)
	assert.True(isExist)
	timeout, err = cache.GetExpire(ctx, "test_key_1")
	assert.NoError(err)
	assert.Equal(time.Second*time.Duration(10), timeout)

	err = cache.Set(ctx, "test_key_1", 200)
	assert.NoError(err)
	timeout, err = cache.GetExpire(ctx, "test_key_1")
	assert.NoError(err)
	assert.Equal(float64(-1), timeout.Seconds())
	err = cache.Delete(ctx, "test_key_1", "test_key_2")
	assert.NoError(err)
	isExist, err = cache.IsExist(ctx, "test_key_1")
	assert.NoError(err)
	assert.False(isExist)

	err = cache.Set(ctx, "test_key_1", 100, time.Second*10)
	assert.NoError(err)
	data, err = cache.GetMap(ctx)
	assert.NoError(err)
	assert.Equal(map[string]any{}, data)
	data, err = cache.GetMapReset(ctx, 0, "test_key_1", "test_key_2")
	assert.Error(err)
	assert.Equal(map[string]any{}, data)
	data, err = cache.GetMap(ctx, "test_key_1", "test_key_2")
	assert.NoError(err)
	assert.Equal(map[string]any{"test_key_1": "100", "test_key_2": nil}, data)
	data, err = cache.GetMapReset(ctx, time.Second*20, "test_key_1", "test_key_2")
	assert.NoError(err)
	assert.Equal(map[string]any{"test_key_1": "100", "test_key_2": nil}, data)
	timeout, err = cache.GetExpire(ctx, "test_key_1")
	assert.NoError(err)
	assert.Equal(time.Second*time.Duration(10), timeout)
	err = cache.SetMap(ctx, map[string]any{"test_key_1": 100, "test_key_2": 200}, time.Second*5)
	assert.NoError(err)
	data, err = cache.GetMap(ctx, "test_key_1", "test_key_2")
	assert.NoError(err)
	assert.Equal(map[string]any{"test_key_1": "100", "test_key_2": "200"}, data)
	timeout, err = cache.GetExpire(ctx, "test_key_1")
	assert.NoError(err)
	assert.Equal(time.Second*time.Duration(5), timeout)
	timeout, err = cache.GetExpire(ctx, "test_key_2")
	assert.NoError(err)
	assert.Equal(time.Second*time.Duration(5), timeout)
	val, err = cache.GetReset(ctx, "test_key_1", time.Second*10)
	assert.NoError(err)
	assert.Equal(100, gtkconv.ToInt(val))
	timeout, err = cache.GetExpire(ctx, "test_key_1")
	assert.NoError(err)
	assert.Equal(time.Second*time.Duration(10), timeout)
	data, err = cache.GetMapReset(ctx, time.Second*time.Duration(60), "test_key_1", "test_key_2")
	assert.NoError(err)
	assert.Equal(map[string]any{"test_key_1": "100", "test_key_2": "200"}, data)
	timeout, err = cache.GetExpire(ctx, "test_key_1")
	assert.NoError(err)
	assert.Equal(time.Second*time.Duration(60), timeout)
	timeout, err = cache.GetExpire(ctx, "test_key_2")
	assert.NoError(err)
	assert.Equal(time.Second*time.Duration(60), timeout)

	val, err = cache.CustomCache(ctx, func(ctx context.Context) (val any, err error) {
		script := `
		local result = redis.call('SETEX', KEYS[1], ARGV[3], ARGV[1])
		if not result['ok'] then
			return false
		end
		result = redis.call('SETEX', KEYS[2], ARGV[3], ARGV[2])
		if not result['ok'] then
			return false
		end
		return true
		`
		val, err = cache.Client().Eval(ctx, script, []string{"test_key_1", "test_key_2"}, 1000, 2000, 120)
		return
	})
	assert.NoError(err)
	assert.True(gtkconv.ToBool(val))
}
