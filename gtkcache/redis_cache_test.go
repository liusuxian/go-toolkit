/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-27 20:53:08
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-28 03:00:57
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
	r := miniredis.RunT(t)
	cache := gtkcache.NewRedisCache([]string{}, func(cc *gtkredis.ClientConfig) {
		cc.Addr = r.Addr()
		cc.DB = 1
		cc.Password = ""
	})
	var (
		ctx     = context.Background()
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
	isExist = cache.IsExist(ctx, "test_key_1")
	assert.False(isExist)
	timeout, err = cache.GetExpire(ctx, "test_key_1")
	assert.NoError(err)
	assert.Equal(float64(-2), timeout.Seconds())

	err = cache.Set(ctx, "test_key_1", 100, time.Second*10)
	assert.NoError(err)
	val, err = cache.Get(ctx, "test_key_1")
	assert.NoError(err)
	assert.Equal(100, gtkconv.ToInt(val))
	isExist = cache.IsExist(ctx, "test_key_1")
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
	isExist = cache.IsExist(ctx, "test_key_1")
	assert.False(isExist)

	err = cache.Set(ctx, "test_key_1", 100, time.Second*10)
	assert.NoError(err)
	data, err = cache.GetMap(ctx)
	assert.NoError(err)
	assert.Equal(map[string]any{}, data)
	data, err = cache.GetMap(ctx, "test_key_1", "test_key_2")
	assert.NoError(err)
	assert.Equal(map[string]any{"test_key_1": "100", "test_key_2": nil}, data)

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

	val, err = cache.CustomCache(ctx, func(ctx context.Context) (val any, err error) {
		if val, err = cache.Client().Do(ctx, "SET", "test_key_3", 300); err != nil {
			return
		}
		val, err = cache.Client().Do(ctx, "GET", "test_key_3")
		return
	})
	assert.NoError(err)
	assert.Equal(300, gtkconv.ToInt(val))
}
