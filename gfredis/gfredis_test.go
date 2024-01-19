/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-19 21:51:32
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-19 21:53:51
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gfredis_test

import (
	"context"
	"encoding/json"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"go-toolkit/gfredis"
	"testing"
)

type A struct {
	A int
	B float64
	C string
	D []any
}

func TestRedis(t *testing.T) {
	r := miniredis.RunT(t)
	client := gfredis.NewClient(func(cc *gfredis.ClientConfig) {
		cc.Addr = r.Addr()
		cc.Password = ""
		cc.DB = 1
	})
	defer client.Close()

	ctx := context.Background()
	assert := assert.New(t)
	actualObj, err := client.Do(ctx, "FLUSHDB")
	if assert.NoError(err) {
		assert.Equal("OK", actualObj.String())
	}
	actualObj, err = client.Do(ctx, "SET", "aaa", 1)
	if assert.NoError(err) {
		assert.Equal("OK", actualObj.String())
	}
	actualObj, err = client.Do(ctx, "GET", "aaa")
	if assert.NoError(err) {
		assert.Equal(1, actualObj.Int())
	}
	var actualPipelineObj []*gfredis.PipelineResult
	actualPipelineObj, err = client.Pipeline(ctx, []any{"SET", "bbb", 2}, []any{"SADD", "ccc", 3})
	if assert.NoError(err) {
		for k, v := range actualPipelineObj {
			assert.IsType(&gfredis.PipelineResult{}, v)
			assert.NoError(v.Err)
			if k == 0 {
				assert.Equal("OK", v.Val.String())
			} else {
				assert.Equal(1, v.Val.Int())
			}
		}
	}
	actualPipelineObj, err = client.Pipeline(ctx, []any{"GET", "bbb"}, []any{"SMEMBERS", "ccc"})
	if assert.NoError(err) {
		for k, v := range actualPipelineObj {
			assert.IsType(&gfredis.PipelineResult{}, v)
			assert.NoError(v.Err)
			if k == 0 {
				assert.Equal(2, v.Val.Int())
			} else {
				assert.Equal([]any{"3"}, v.Val.Slice())
			}
		}
	}
	actualObj, err = client.Do(ctx, "SET", "ddd", &A{A: 1, B: 1.2, C: "hello", D: []any{1, 1.2, "hello", true}})
	if assert.NoError(err) {
		assert.Equal("OK", actualObj.String())
	}
	actualObj, err = client.Do(ctx, "GET", "ddd")
	if assert.NoError(err) {
		val := &A{}
		err = actualObj.Struct(&val)
		if assert.NoError(err) {
			assert.IsType(&A{}, val)
			assert.Equal(&A{A: 1, B: 1.2, C: "hello", D: []any{json.Number("1"), json.Number("1.2"), "hello", true}}, val)
			assert.Equal(map[string]any{"A": json.Number("1"), "B": json.Number("1.2"), "C": "hello", "D": []any{json.Number("1"), json.Number("1.2"), "hello", true}}, actualObj.MapStrAny())
		}
	}
	actualObj, err = client.Do(ctx, "SET", "eee", []any{1, 1.2, "hello", true})
	if assert.NoError(err) {
		assert.Equal("OK", actualObj.String())
	}
	actualObj, err = client.Do(ctx, "GET", "eee")
	if assert.NoError(err) {
		assert.ElementsMatch([]any{json.Number("1"), json.Number("1.2"), "hello", true}, actualObj.Slice())
	}
	var rl *gfredis.RedisLock
	rl, err = client.NewRedisLock("test_redis_lock")
	if assert.NoError(err) {
		t.Log("new redis lock succ")
	}
	var ok bool
	ok, err = rl.TryLock(ctx)
	if assert.NoError(err) {
		t.Log("try lock1: ", ok)
	}
	if ok {
		rl.Unlock(ctx)
	}
	ok, err = rl.TryLock(ctx)
	if assert.NoError(err) {
		t.Log("try lock2: ", ok)
	}
	if ok {
		rl.Unlock(ctx)
	}
	ok, err = rl.SpinLock(ctx, 10)
	if assert.NoError(err) {
		t.Log("try lock3: ", ok)
	}
	if ok {
		rl.Unlock(ctx)
	}
	err = client.ScriptLoad(ctx, "lua_script/test1.lua")
	if assert.Error(err) {
		t.Log("lua script load failed: ", err)
	}
	err = client.ScriptLoad(ctx, "lua_script/test.lua")
	if assert.NoError(err) {
		t.Log("lua script load succ")
	}
	actualObj, err = client.EvalSha(ctx, "test", []string{"lua_key"}, 1)
	if assert.NoError(err) {
		assert.Equal(1, actualObj.Int())
	}
}
