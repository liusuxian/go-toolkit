/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-27 20:53:08
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-29 16:39:48
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

type AAA struct {
	A     int
	B     float64
	C     string
	cache *gtkcache.RedisCache
}

func (a *AAA) Get(ctx context.Context, keys []string, args []any, timeout ...time.Duration) (val any, err error) {
	return
}

func (a *AAA) Set(ctx context.Context, keys []string, args []any, newVal any, timeout ...time.Duration) (val any, err error) {
	script := `
		local result = redis.call('SETEX', KEYS[1], ARGV[3], ARGV[1])
		if not result['ok'] then
			return false
		end
		result = redis.call('SETEX', KEYS[2], ARGV[3], ARGV[2])
		if not result['ok'] then
			return false
		end
		return redis.call('MGET', KEYS[1], KEYS[2])
		`
	val, err = a.cache.Client().Eval(ctx, script, keys, 1000, 2000, 120)
	return
}

func TestRedisCacheString(t *testing.T) {
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
		ok      bool
		timeout time.Duration
	)
	val, err = cache.Get(ctx, "test_key_1", time.Second*10)
	assert.NoError(err)
	assert.Nil(val)
	isExist, err = cache.IsExist(ctx, "test_key_1")
	assert.NoError(err)
	assert.False(isExist)
	timeout, err = cache.GetExpire(ctx, "test_key_1")
	assert.NoError(err)
	assert.Equal(float64(-2), timeout.Seconds())

	err = cache.Set(ctx, "test_key_2", nil, time.Second)
	assert.NoError(err)
	val, err = cache.GetOrSet(ctx, "test_key_2", 200, time.Second*2)
	assert.NoError(err)
	assert.Equal("", val)
	a1 := AAA{}
	gtkconv.ToStruct(val, &a1)
	assert.Equal(AAA{A: 0, B: 0, C: ""}, a1)
	timeout, err = cache.GetExpire(ctx, "test_key_2")
	assert.NoError(err)
	assert.Equal(time.Second*2, timeout)
	val, err = cache.GetOrSet(ctx, "test_key_22", map[string]any{"a": 1}, time.Second)
	assert.NoError(err)
	assert.Equal("{\"a\":1}", val)
	val, err = cache.GetOrSetFunc(ctx, "test_key_3", func(ctx context.Context) (val any, err error) {
		return
	}, true, time.Second)
	assert.NoError(err)
	assert.Equal("", val)
	err = cache.SetMap(ctx, map[string]any{"a": 1, "b": map[string]any{"b": 100}}, time.Second)
	assert.NoError(err)
	ok, err = cache.SetIfNotExist(ctx, "test_key_3", 100, time.Second)
	assert.NoError(err)
	assert.False(ok)
	ok, err = cache.SetIfNotExist(ctx, "test_key_4", nil, time.Second)
	assert.NoError(err)
	assert.True(ok)
	ok, err = cache.SetIfNotExistFunc(ctx, "test_key_4", func(ctx context.Context) (val any, err error) {
		return
	}, true, time.Second)
	assert.NoError(err)
	assert.False(ok)
	ok, err = cache.SetIfNotExistFunc(ctx, "test_key_5", func(ctx context.Context) (val any, err error) {
		return
	}, true, time.Second)
	assert.NoError(err)
	assert.True(ok)

	val, err = cache.CustomGetOrSetFunc(ctx, []string{"test_key_10", "test_key_11"}, []any{}, &AAA{cache: cache}, func(ctx context.Context) (val any, err error) {
		return map[string]any{
			"test_key_10": 1,
			"test_key_11": 2,
		}, nil
	}, true, time.Second)
	assert.NoError(err)
	assert.Equal([]any{"1000", "2000"}, val)
}

func TestRedisCacheString2(t *testing.T) {
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
	)
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
	assert.Equal(time.Second*10, timeout)
}

func TestRedisCacheString3(t *testing.T) {
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
		err     error
		isExist bool
		timeout time.Duration
	)
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
}

func TestRedisCacheString4(t *testing.T) {
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
		timeout time.Duration
		data    map[string]any
	)
	err = cache.Set(ctx, "test_key_1", 100, time.Second*10)
	assert.NoError(err)
	data, err = cache.GetMap(ctx, []string{})
	assert.NoError(err)
	assert.Nil(data)
	data, err = cache.GetMap(ctx, []string{"test_key_1", "test_key_2"}, 0)
	assert.NoError(err)
	assert.Equal(map[string]any{"test_key_1": "100", "test_key_2": nil}, data)
	data, err = cache.GetMap(ctx, []string{"test_key_1", "test_key_2"})
	assert.NoError(err)
	assert.Equal(map[string]any{"test_key_1": "100", "test_key_2": nil}, data)
	data, err = cache.GetMap(ctx, []string{"test_key_1", "test_key_2"}, time.Second*20)
	assert.NoError(err)
	assert.Equal(map[string]any{"test_key_1": "100", "test_key_2": nil}, data)
	timeout, err = cache.GetExpire(ctx, "test_key_1")
	assert.NoError(err)
	assert.Equal(time.Second*10, timeout)
	err = cache.SetMap(ctx, map[string]any{"test_key_1": 100, "test_key_2": 200}, time.Second*5)
	assert.NoError(err)
	data, err = cache.GetMap(ctx, []string{"test_key_1", "test_key_2"})
	assert.NoError(err)
	assert.Equal(map[string]any{"test_key_1": "100", "test_key_2": "200"}, data)
	timeout, err = cache.GetExpire(ctx, "test_key_1")
	assert.NoError(err)
	assert.Equal(time.Second*5, timeout)
	timeout, err = cache.GetExpire(ctx, "test_key_2")
	assert.NoError(err)
	assert.Equal(time.Second*5, timeout)
	val, err = cache.Get(ctx, "test_key_1", time.Second*10)
	assert.NoError(err)
	assert.Equal(100, gtkconv.ToInt(val))
	timeout, err = cache.GetExpire(ctx, "test_key_1")
	assert.NoError(err)
	assert.Equal(time.Second*10, timeout)
	data, err = cache.GetMap(ctx, []string{"test_key_1", "test_key_2"}, time.Second*60)
	assert.NoError(err)
	assert.Equal(map[string]any{"test_key_1": "100", "test_key_2": "200"}, data)
	timeout, err = cache.GetExpire(ctx, "test_key_1")
	assert.NoError(err)
	assert.Equal(time.Second*60, timeout)
	timeout, err = cache.GetExpire(ctx, "test_key_2")
	assert.NoError(err)
	assert.Equal(time.Second*60, timeout)
}

func TestRedisCacheSet(t *testing.T) {
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
		assert   = assert.New(t)
		val      any
		err      error
		timeout  time.Duration
		isMember bool
		members  []any
	)
	members, err = cache.SMembers(ctx, "test_key_0", time.Second*10)
	assert.NoError(err)
	assert.Equal([]any{}, members)
	err = cache.Set(ctx, "test_key_1", 100, time.Second*10)
	assert.NoError(err)
	val, err = cache.SAdd(ctx, "test_key_1", []any{100, "hello", 1.11}, time.Second*5)
	assert.Error(err)
	assert.Equal(0, val)
	val, err = cache.SAdd(ctx, "test_key_2", []any{100, "hello", 1.11})
	assert.NoError(err)
	assert.Equal(3, val)
	timeout, err = cache.GetExpire(ctx, "test_key_2")
	assert.NoError(err)
	assert.Equal(time.Second*-1, timeout)
	val, err = cache.SAdd(ctx, "test_key_2", []any{100, "hello", 1.11}, time.Second*5)
	assert.NoError(err)
	assert.Equal(0, val)
	timeout, err = cache.GetExpire(ctx, "test_key_2")
	assert.NoError(err)
	assert.Equal(time.Second*5, timeout)
	val, err = cache.SAdd(ctx, "test_key_2", []any{200}, time.Second*10)
	assert.NoError(err)
	assert.Equal(1, val)
	timeout, err = cache.GetExpire(ctx, "test_key_2")
	assert.NoError(err)
	assert.Equal(time.Second*10, timeout)
	isMember, err = cache.SIsMember(ctx, "test_key_2", 100, time.Second*20)
	assert.NoError(err)
	assert.True(isMember)
	timeout, err = cache.GetExpire(ctx, "test_key_2")
	assert.NoError(err)
	assert.Equal(time.Second*20, timeout)
	members, err = cache.SMembers(ctx, "test_key_3", time.Second*20)
	assert.NoError(err)
	assert.Equal([]any{}, members)
	timeout, err = cache.GetExpire(ctx, "test_key_3")
	assert.NoError(err)
	assert.Equal(time.Second*-2, timeout)
	members, err = cache.SMembers(ctx, "test_key_2", time.Second*30)
	assert.NoError(err)
	assert.Equal([]any{"1.11", "100", "200", "hello"}, members)
	timeout, err = cache.GetExpire(ctx, "test_key_2")
	assert.NoError(err)
	assert.Equal(time.Second*30, timeout)
	isMember, err = cache.SIsMember(ctx, "test_key_2", 1000, time.Second*5)
	assert.NoError(err)
	assert.False(isMember)
	timeout, err = cache.GetExpire(ctx, "test_key_2")
	assert.NoError(err)
	assert.Equal(time.Second*5, timeout)
	members, err = cache.SPop(ctx, "test_key_2", 10, time.Second*10)
	assert.NoError(err)
	assert.Equal(4, len(members))
	timeout, err = cache.GetExpire(ctx, "test_key_2")
	assert.NoError(err)
	assert.Equal(time.Second*-2, timeout)
	val, err = cache.SAdd(ctx, "test_key_3", []any{100, "world", 1.11}, time.Second*5)
	assert.NoError(err)
	assert.Equal(3, val)
	timeout, err = cache.GetExpire(ctx, "test_key_3")
	assert.NoError(err)
	assert.Equal(time.Second*5, timeout)
	members, err = cache.SUnion(ctx, []string{"test_key_2", "test_key_3"}, time.Second*20)
	assert.NoError(err)
	assert.Equal(3, len(members))
	timeout, err = cache.GetExpire(ctx, "test_key_3")
	assert.NoError(err)
	assert.Equal(time.Second*5, timeout)
	val, err = cache.SRem(ctx, "test_key_2", []any{100})
	assert.NoError(err)
	assert.Equal(0, val)
	val, err = cache.SCard(ctx, "test_key_2")
	assert.NoError(err)
	assert.Equal(0, val)
}

func TestRedisCacheSortedSet(t *testing.T) {
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
		assert   = assert.New(t)
		val      any
		err      error
		timeout  time.Duration
		total    int
		remCount int
		members  []map[any]float64
		score    float64
	)
	members, err = cache.SSRange(ctx, "test_key_0", 0, -1, false, false, time.Second*10)
	assert.NoError(err)
	assert.Equal([]map[any]float64{}, members)
	err = cache.Set(ctx, "test_key_1", 100, time.Second*10)
	assert.NoError(err)
	val, err = cache.SSAdd(ctx, "test_key_1", map[any]float64{"1": 100, "2": 200}, time.Second*10)
	assert.Error(err)
	assert.Equal(0, val)
	timeout, err = cache.GetExpire(ctx, "test_key_1")
	assert.NoError(err)
	assert.Equal(time.Second*10, timeout)
	val, err = cache.SSAdd(ctx, "test_key_2", map[any]float64{"1": 100, "2": 200}, time.Second*20)
	assert.NoError(err)
	assert.Equal(2, val)
	timeout, err = cache.GetExpire(ctx, "test_key_2")
	assert.NoError(err)
	assert.Equal(time.Second*20, timeout)
	members, err = cache.SSRange(ctx, "test_key_2", 0, -1, false, false, time.Second*30)
	assert.NoError(err)
	assert.Equal([]map[any]float64{{"1": 0}, {"2": 0}}, members)
	members, err = cache.SSRange(ctx, "test_key_2", 0, -1, true, false, time.Second*40)
	assert.NoError(err)
	assert.Equal([]map[any]float64{{"2": 0}, {"1": 0}}, members)
	members, err = cache.SSRange(ctx, "test_key_2", 0, -1, false, true, time.Second*50)
	assert.NoError(err)
	assert.Equal([]map[any]float64{{"1": 100}, {"2": 200}}, members)
	members, err = cache.SSRange(ctx, "test_key_2", 0, -1, true, true, time.Second*60)
	assert.NoError(err)
	assert.Equal([]map[any]float64{{"2": 200}, {"1": 100}}, members)
	timeout, err = cache.GetExpire(ctx, "test_key_2")
	assert.NoError(err)
	assert.Equal(time.Second*60, timeout)
	val, err = cache.SSAdd(ctx, "test_key_3", map[any]float64{"1": 100, "2": 200, "3": 300, "4": 400, "5": 500, "6": 600}, time.Second*5)
	assert.NoError(err)
	assert.Equal(6, val)
	total, members, err = cache.SSPage(ctx, "test_key_3", 1, 10, false, false, time.Second*10)
	assert.NoError(err)
	assert.Equal(6, total)
	assert.Equal([]map[any]float64{{"1": 0}, {"2": 0}, {"3": 0}, {"4": 0}, {"5": 0}, {"6": 0}}, members)
	total, members, err = cache.SSPage(ctx, "test_key_3", 2, 10, false, false, time.Second*10)
	assert.NoError(err)
	assert.Equal(6, total)
	assert.Equal([]map[any]float64{}, members)
	total, members, err = cache.SSPage(ctx, "test_key_3", 1, 10, false, true, time.Second*20)
	assert.NoError(err)
	assert.Equal(6, total)
	assert.Equal([]map[any]float64{{"1": 100}, {"2": 200}, {"3": 300}, {"4": 400}, {"5": 500}, {"6": 600}}, members)
	timeout, err = cache.GetExpire(ctx, "test_key_3")
	assert.NoError(err)
	assert.Equal(time.Second*20, timeout)
	total, err = cache.SSCard(ctx, "test_key_3", time.Second*30)
	assert.NoError(err)
	assert.Equal(6, total)
	timeout, err = cache.GetExpire(ctx, "test_key_3")
	assert.NoError(err)
	assert.Equal(time.Second*30, timeout)
	total, err = cache.SSCount(ctx, "test_key_3", "(300", "(500", time.Second*40)
	assert.NoError(err)
	assert.Equal(1, total)
	timeout, err = cache.GetExpire(ctx, "test_key_3")
	assert.NoError(err)
	assert.Equal(time.Second*40, timeout)
	total, err = cache.SSCount(ctx, "test_key_3", "-inf", "+inf", time.Second*40)
	assert.NoError(err)
	assert.Equal(6, total)
	score, err = cache.SSIncrby(ctx, "test_key_4", 100, 1, time.Second)
	assert.NoError(err)
	assert.Equal(float64(100), score)
	timeout, err = cache.GetExpire(ctx, "test_key_4")
	assert.NoError(err)
	assert.Equal(time.Second, timeout)
	score, err = cache.SSIncrby(ctx, "test_key_3", 200, 1, time.Second)
	assert.NoError(err)
	assert.Equal(float64(300), score)
	timeout, err = cache.GetExpire(ctx, "test_key_3")
	assert.NoError(err)
	assert.Equal(time.Second, timeout)
	members, err = cache.SSRangeByScore(ctx, "test_key_3", "(100", "(500", false, true, []int{}, time.Second)
	assert.NoError(err)
	assert.Equal([]map[any]float64{{"2": 200}, {"1": 300}, {"3": 300}, {"4": 400}}, members)
	members, err = cache.SSRangeByScore(ctx, "test_key_3", "(100", "(500", true, true, []int{}, time.Second)
	assert.NoError(err)
	assert.Equal([]map[any]float64{{"4": 400}, {"3": 300}, {"1": 300}, {"2": 200}}, members)
	members, err = cache.SSRangeByScore(ctx, "test_key_3", "(100", "(500", true, true, []int{0, 2}, time.Second)
	assert.NoError(err)
	assert.Equal([]map[any]float64{{"4": 400}, {"3": 300}}, members)
	total, err = cache.SSRank(ctx, "test_key_3", 6, false, time.Second)
	assert.NoError(err)
	assert.Equal(5, total)
	total, err = cache.SSRank(ctx, "test_key_3", 6, true, time.Second)
	assert.NoError(err)
	assert.Equal(0, total)
	total, err = cache.SSRank(ctx, "test_key_3", 7, true, time.Second)
	assert.NoError(err)
	assert.Equal(-1, total)
	remCount, err = cache.SSRem(ctx, "test_key_3", []any{"1"}, time.Second)
	assert.NoError(err)
	assert.Equal(1, remCount)
	remCount, err = cache.SSRem(ctx, "test_key_3", []any{"1"}, time.Second)
	assert.NoError(err)
	assert.Equal(0, remCount)
	remCount, err = cache.SSRemRangeByRank(ctx, "test_key_3", 0, 0, time.Second)
	assert.NoError(err)
	assert.Equal(1, remCount)
	remCount, err = cache.SSRemRangeByScore(ctx, "test_key_3", "300", "400", time.Second)
	assert.NoError(err)
	assert.Equal(2, remCount)
	score, err = cache.SSScore(ctx, "test_key_3", 1, time.Second)
	assert.NoError(err)
	assert.Equal(float64(-1), score)
}
