/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-04-04 12:14:28
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-28 00:13:46
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkredis_test

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/liusuxian/go-toolkit/gtkconv"
	"github.com/liusuxian/go-toolkit/gtkredis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type A struct {
	A int
	B float64
	C string
	D []any
}

func TestRedis(t *testing.T) {
	r := miniredis.RunT(t)
	client := gtkredis.NewClient(func(cc *gtkredis.ClientConfig) {
		cc.Addr = r.Addr()
		cc.Password = ""
		cc.DB = 1
	})
	defer client.Close()

	ctx := context.Background()
	assert := assert.New(t)
	actualObj, err := client.Do(ctx, "FLUSHDB")
	if assert.NoError(err) {
		assert.Equal("OK", actualObj)
	}
	actualObj, err = client.Do(ctx, "SET", "aaa", 1)
	if assert.NoError(err) {
		assert.Equal("OK", actualObj)
	}
	actualObj, err = client.Do(ctx, "GET", "aaa")
	if assert.NoError(err) {
		assert.Equal(1, gtkconv.ToInt(actualObj))
	}
	var actualPipelineObj []*gtkredis.PipelineResult
	actualPipelineObj, err = client.Pipeline(ctx, []any{"SET", "bbb", 2}, []any{"SADD", "ccc", 3})
	if assert.NoError(err) {
		for k, v := range actualPipelineObj {
			assert.IsType(&gtkredis.PipelineResult{}, v)
			assert.NoError(v.Err)
			if k == 0 {
				assert.Equal("OK", v.Val)
			} else {
				assert.Equal(int64(1), v.Val)
			}
		}
	}
	actualPipelineObj, err = client.Pipeline(ctx, []any{"GET", "bbb"}, []any{"SMEMBERS", "ccc"})
	if assert.NoError(err) {
		for k, v := range actualPipelineObj {
			assert.IsType(&gtkredis.PipelineResult{}, v)
			assert.NoError(v.Err)
			if k == 0 {
				assert.Equal("2", v.Val)
			} else {
				assert.Equal([]any{"3"}, v.Val)
			}
		}
	}
	actualObj, err = client.Do(ctx, "SET", "ddd", &A{A: 1, B: 1.2, C: "hello", D: []any{1, 1.2, "hello", true}})
	if assert.NoError(err) {
		assert.Equal("OK", actualObj)
	}
	actualObj, err = client.Do(ctx, "GET", "ddd")
	if assert.NoError(err) {
		val := &A{}
		err = gtkconv.ToStructE(actualObj, &val)
		if assert.NoError(err) {
			assert.IsType(&A{}, val)
			assert.Equal(&A{A: 1, B: 1.2, C: "hello", D: []any{float64(1), 1.2, "hello", true}}, val)
			assert.Equal(map[string]any{"A": float64(1), "B": 1.2, "C": "hello", "D": []any{float64(1), 1.2, "hello", true}}, gtkconv.ToStringMap(actualObj))
		}
	}
	actualObj, err = client.Do(ctx, "SET", "eee", []any{1, 1.2, "hello", true})
	if assert.NoError(err) {
		assert.Equal("OK", actualObj)
	}
	actualObj, err = client.Do(ctx, "GET", "eee")
	if assert.NoError(err) {
		assert.ElementsMatch([]any{float64(1), 1.2, "hello", true}, gtkconv.ToSlice(actualObj))
	}

	actualObj, err = client.SetCD(ctx, "test_a", time.Second*2)
	if assert.NoError(err) {
		assert.Equal(true, gtkconv.ToBool(actualObj))
	}
	actualObj, err = client.CD(ctx, "test_a")
	if assert.NoError(err) {
		assert.Equal(false, gtkconv.ToBool(actualObj))
	}
	actualObj, err = client.CD(ctx, "test_b")
	if assert.NoError(err) {
		assert.Equal(true, gtkconv.ToBool(actualObj))
	}

	var rl *gtkredis.RedisLock
	rl, err = client.NewRedisLock("test_redis_lock")
	if assert.NoError(err) {
		assert.NotNil(rl)
	}
	var ok bool
	ok, err = rl.TryLock(ctx)
	if assert.NoError(err) {
		assert.True(ok)
	}
	if ok {
		rl.Unlock(ctx)
	}
	ok, err = rl.TryLock(ctx)
	if assert.NoError(err) {
		assert.True(ok)
	}
	if ok {
		rl.Unlock(ctx)
	}
	ok, err = rl.SpinLock(ctx, 10)
	if assert.NoError(err) {
		assert.True(ok)
	}
	if ok {
		rl.Unlock(ctx)
	}

	err = client.ScriptLoad(ctx, "lua_script/test1.lua")
	assert.Error(err)
	err = client.ScriptLoad(ctx, "lua_script/test.lua")
	assert.NoError(err)
	actualObj, err = client.EvalSha(ctx, "test", []string{"lua_key"}, 1)
	if assert.NoError(err) {
		assert.Equal(1, gtkconv.ToInt(actualObj))
	}
}

func TestRedisLuaScript(t *testing.T) {
	r := miniredis.RunT(t)
	client := gtkredis.NewClient(func(cc *gtkredis.ClientConfig) {
		cc.Addr = r.Addr()
		cc.Password = ""
		cc.DB = 1
	})
	defer client.Close()

	var (
		err    error
		val    any
		ctx    = context.Background()
		assert = assert.New(t)
	)
	err = client.ScriptLoad(ctx, "lua_script/test_get.lua")
	assert.NoError(err)
	val, err = client.Do(ctx, "SADD", "test_get1", 100, 200, 300)
	assert.NoError(err)
	assert.Equal(3, gtkconv.ToInt(val))
	val, err = client.EvalSha(ctx, "test_get", []string{"test_get1"})
	assert.Error(err)
	assert.Nil(val)
	val, err = client.Do(ctx, "SET", "test_get2", 100)
	assert.NoError(err)
	assert.Equal("OK", gtkconv.ToString(val))
	val, err = client.EvalSha(ctx, "test_get", []string{"test_get2", "test_get2"})
	assert.NoError(err)
	assert.Equal(3, gtkconv.ToInt(val))
	val, err = client.EvalSha(ctx, "test_get", []string{"test_get2", "test_get3"})
	assert.NoError(err)
	assert.Equal(2, gtkconv.ToInt(val))

	err = client.ScriptLoad(ctx, "lua_script/test_set.lua")
	assert.NoError(err)
	val, err = client.EvalSha(ctx, "test_set", []string{"test_set1"}, 100)
	assert.NoError(err)
	assert.Equal(1, gtkconv.ToInt(val))
}
