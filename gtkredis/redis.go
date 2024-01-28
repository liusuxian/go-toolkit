/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-04-15 02:58:43
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-28 17:05:53
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkredis

import (
	"context"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/liusuxian/go-toolkit/gtkconv"
	"github.com/liusuxian/go-toolkit/gtkfile"
	"github.com/liusuxian/go-toolkit/gtkreflection"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"reflect"
	"strings"
	"time"
)

// ClientConfig redis 客户端配置
type ClientConfig = redis.Options

// ClientConfigOption redis 客户端配置选项
type ClientConfigOption func(cc *ClientConfig)

// RedisClient redis 客户端结构
type RedisClient struct {
	client        *redis.Client // redis 客户端
	luaEvalShaMap map[string]string
}

// PipelineResult 管道返回值
type PipelineResult struct {
	Val any
	Err error
}

// RedisLock redis 分布式锁
type RedisLock struct {
	client     *RedisClient
	key        string
	uuid       string
	cancelFunc context.CancelFunc
}

const (
	defaultExpiration = 10 // 单位，秒
	sleepDur          = 10 * time.Millisecond
)

// 内置 lua 脚本
var internalScriptMap = map[string]string{
	"compareAndDelete": `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
		`,
}

// NewClient 创建 redis 客户端
func NewClient(ctx context.Context, opts ...ClientConfigOption) (client *RedisClient) {
	ro := &redis.Options{}
	for _, opt := range opts {
		opt(ro)
	}
	client = &RedisClient{
		client:        redis.NewClient(ro),
		luaEvalShaMap: make(map[string]string),
	}
	for k, v := range internalScriptMap {
		if err := client.ScriptLoad(ctx, k, v); err != nil {
			panic(err)
		}
	}
	return
}

// Do 执行 redis 命令
func (rc *RedisClient) Do(ctx context.Context, cmd string, args ...any) (value any, err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	// 处理 redis 命令参数
	for k, v := range args {
		reflectInfo := gtkreflection.OriginTypeAndKind(v)
		switch reflectInfo.OriginKind {
		case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
			// 忽略切片类型为 []byte 的情况
			if _, ok := v.([]byte); !ok {
				if args[k], err = json.Marshal(v); err != nil {
					return
				}
			}
		}
	}
	// 执行 redis 命令
	cmdArgs := make([]any, 0, len(args)+1)
	cmdArgs = append(cmdArgs, cmd)
	cmdArgs = append(cmdArgs, args...)
	value, err = rc.client.Do(ctx, cmdArgs...).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

// Pipeline 执行 redis 管道命令
func (rc *RedisClient) Pipeline(ctx context.Context, cmdArgsList ...[]any) (results []*PipelineResult, err error) {
	if len(cmdArgsList) == 0 {
		err = errors.New("pipeline cmd args list is empty")
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}
	// 执行 redis 管道命令
	p := rc.client.Pipeline()
	// 处理redis命令参数
	for _, cmdArgs := range cmdArgsList {
		for k, v := range cmdArgs {
			if k > 0 {
				reflectInfo := gtkreflection.OriginTypeAndKind(v)
				switch reflectInfo.OriginKind {
				case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
					// 忽略切片类型为 []byte 的情况
					if _, ok := v.([]byte); !ok {
						if cmdArgs[k], err = json.Marshal(v); err != nil {
							return
						}
					}
				}
			}
		}
		// 执行 redis 命令
		p.Do(ctx, cmdArgs...)
	}
	var resList []redis.Cmder
	resList, err = p.Exec(ctx)
	if err == redis.Nil {
		err = nil
	}
	if err != nil {
		return
	}
	// 处理返回结果
	results = make([]*PipelineResult, 0, len(resList))
	for _, v := range resList {
		results = append(results, &PipelineResult{
			Val: v.(*redis.Cmd).Val(),
			Err: v.Err(),
		})
	}
	return
}

// ScriptLoad 加载 lua 脚本
func (rc *RedisClient) ScriptLoad(ctx context.Context, name, script string) (err error) {
	var evalsha string
	if evalsha, err = rc.client.ScriptLoad(ctx, script).Result(); err != nil {
		return
	}
	rc.luaEvalShaMap[name] = evalsha
	return
}

// ScriptLoadByPath 通过 lua 脚本文件的路径加载 lua 脚本
func (rc *RedisClient) ScriptLoadByPath(ctx context.Context, scriptPath string) (err error) {
	script := gtkfile.GetContents(scriptPath)
	if strings.EqualFold("", script) {
		err = errors.Errorf("[%s] script not found", scriptPath)
		return
	}
	var evalsha string
	if evalsha, err = rc.client.ScriptLoad(ctx, script).Result(); err != nil {
		return
	}
	name := gtkfile.Name(scriptPath)
	rc.luaEvalShaMap[name] = evalsha
	return
}

// Eval 执行 lua 脚本
func (rc *RedisClient) Eval(ctx context.Context, script string, keys []string, args ...any) (value any, err error) {
	value, err = rc.client.Eval(ctx, script, keys, args...).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

// EvalSha 执行 lua 脚本
func (rc *RedisClient) EvalSha(ctx context.Context, name string, keys []string, args ...any) (value any, err error) {
	evalsha, ok := rc.luaEvalShaMap[name]
	if !ok {
		err = errors.Errorf("[%s] Script Not Found", name)
		return
	}
	value, err = rc.client.EvalSha(ctx, evalsha, keys, args...).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

// CD 冷却时间检测
func (rc *RedisClient) CD(ctx context.Context, key string) (ok bool, err error) {
	var value any
	if value, err = rc.Do(ctx, "TTL", key); err != nil {
		return
	}
	switch gtkconv.ToInt(value) {
	case -2:
		// 不存在
		ok = true
	case -1:
		// 永久key 异常清理
		if _, err = rc.Do(ctx, "DEL", key); err != nil {
			return
		}
	}
	return
}

// SetCD 设置冷却时间
func (rc *RedisClient) SetCD(ctx context.Context, key string, cd time.Duration) (ok bool, err error) {
	var value any
	if value, err = rc.Do(ctx, "SET", key, 1, "EX", cd.Seconds(), "NX"); err != nil {
		return
	}
	ok = gtkconv.ToBool(value)
	return
}

// Cad compare and delete
func (rc *RedisClient) Cad(ctx context.Context, key string, value any) (ok bool, err error) {
	var result any
	if result, err = rc.EvalSha(ctx, "compareAndDelete", []string{key}, value); err != nil {
		return
	}
	ok = gtkconv.ToBool(result)
	return
}

// Close 关闭 redis
func (rc *RedisClient) Close() (err error) {
	return rc.client.Close()
}

// NewRedisLock 创建 redis 分布式锁
func (rc *RedisClient) NewRedisLock(key string) (rl *RedisLock, err error) {
	var id uuid.UUID
	if id, err = uuid.NewV4(); err != nil {
		return
	}
	rl = &RedisLock{
		client: rc,
		key:    key,
		uuid:   id.String(),
	}
	return
}

// TryLock 尝试加锁
func (rl *RedisLock) TryLock(ctx context.Context) (ok bool, err error) {
	var value any
	if value, err = rl.client.Do(ctx, "SET", rl.key, rl.uuid, "EX", defaultExpiration, "NX"); err != nil {
		return
	}
	ok = gtkconv.ToBool(value)
	if !ok {
		return
	}
	c, cancel := context.WithCancel(ctx)
	rl.cancelFunc = cancel
	rl.refresh(c)
	return
}

// SpinLock 自旋加锁
func (rl *RedisLock) SpinLock(ctx context.Context, retryTimes int) (ok bool, err error) {
	for i := 0; i < retryTimes; i++ {
		if ok, err = rl.TryLock(ctx); err != nil {
			return
		}
		if ok {
			return
		}
		time.Sleep(sleepDur)
	}
	return
}

// Unlock
func (rl *RedisLock) Unlock(ctx context.Context) (ok bool, err error) {
	if ok, err = rl.client.Cad(ctx, rl.key, rl.uuid); err != nil {
		return
	}
	if ok {
		rl.cancelFunc()
	}
	return
}

// refresh 刷新
func (rl *RedisLock) refresh(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(defaultExpiration / 4)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				rl.client.Do(ctx, "EXPIRE", rl.key, defaultExpiration)
			}
		}
	}()
}
