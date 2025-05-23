/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-29 16:15:07
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-23 17:41:17
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkcache_test

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/liusuxian/go-toolkit/gtkcache"
	"github.com/liusuxian/go-toolkit/gtkredis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewWechatCache(t *testing.T) {
	var (
		ctx    = context.Background()
		r      = miniredis.RunT(t)
		assert = assert.New(t)
		cache  *gtkcache.WechatCache
		err    error
	)
	cache, err = gtkcache.NewWechatCache(ctx, &gtkredis.ClientConfig{
		Addr:     r.Addr(),
		DB:       1,
		Password: "",
	})
	assert.NoError(err)
	assert.NotNil(cache)

	var (
		val     any
		isExist bool
	)
	val = cache.Get("test_key_1")
	assert.Nil(val)
	err = cache.Set("test_key_1", 100, time.Second*10)
	assert.NoError(err)
	isExist = cache.IsExist("test_key_1")
	assert.True(isExist)
	err = cache.Delete("test_key_1")
	assert.NoError(err)
	isExist = cache.IsExist("test_key_1")
	assert.False(isExist)
}
