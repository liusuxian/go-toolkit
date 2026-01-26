/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-04-15 02:58:43
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2026-01-26 18:55:23
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkredis

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/liusuxian/go-toolkit/gtkconv"
	"github.com/liusuxian/go-toolkit/internal/utils"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

// ClientConfig redis 客户端配置
type ClientConfig struct {
	Addr            string        `json:"addr"`               // 地址:端口
	ClientName      string        `json:"client_name"`        // 执行 CLIENT SETNAME 命令所用的客户端名称
	Protocol        int           `json:"protocol"`           // 设置与 Redis Server 通信的 RESP 协议版本，默认 3，可选 2 或 3
	Username        string        `json:"username"`           // 访问授权用户
	Password        string        `json:"password"`           // 访问授权密码
	DB              int           `json:"db"`                 // 数据库索引，默认 0
	MaxRetries      int           `json:"max_retries"`        // 最大重试次数，默认 3，-1 表示禁用重试
	MinRetryBackoff time.Duration `json:"min_retry_backoff"`  // 每次重试之间的最小退避时间，默认 8ms，-1 表示禁用退避
	MaxRetryBackoff time.Duration `json:"max_retry_backoff"`  // 每次重试之间的最大退避时间，默认 512ms，-1 表示禁用退避
	DialTimeout     time.Duration `json:"dial_timeout"`       // 连接的超时时间，默认 5s
	ReadTimeout     time.Duration `json:"read_timeout"`       // Read 操作超时时间，默认 3s，-1 表示无超时，-2 表示完全禁用 SetReadDeadline 调用
	WriteTimeout    time.Duration `json:"write_timeout"`      // Write 操作超时时间，默认 3s，-1 表示无超时，-2 表示完全禁用 SetWriteDeadline 调用
	PoolFIFO        bool          `json:"pool_fifo"`          // 连接池类型，true 表示 FIFO（先进先出），false 表示 LIFO（后进先出），默认 false，FIFO 相比 LIFO 有略高的开销，但它有助于更快地关闭空闲连接，减少池大小
	PoolSize        int           `json:"pool_size"`          // 连接池大小，默认每个可用 CPU 10 个连接，如果池中没有足够的连接，将分配超出 PoolSize 的新连接，您可以通过 MaxActiveConns 进行限制
	PoolTimeout     time.Duration `json:"pool_timeout"`       // 如果所有连接都忙，客户端在返回错误前等待连接的时间，默认为 ReadTimeout + 1s
	MinIdleConns    int           `json:"min_idle_conns"`     // 允许闲置的最小连接数，默认 0
	MaxIdleConns    int           `json:"max_idle_conns"`     // 允许闲置的最大连接数，默认 0，0 表示不限制
	MaxActiveConns  int           `json:"max_active_conns"`   // 最大连接数量限制，默认 0，0 表示不限制
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time"` // 连接最大空闲时间，默认 30m，-1 表示禁用空闲超时检查
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`  // 连接最长存活时间，默认 0 表示不关闭空闲连接
	TLSConfig       *tls.Config   `json:"-"`                  // tls 配置
	DisableIdentity bool          `json:"disable_identity"`   // 用于在连接时禁用 CLIENT SETINFO 命令，默认 false
	IdentitySuffix  string        `json:"identity_suffix"`    // 添加客户端名称后缀
	UnstableResp3   bool          `json:"unstable_resp_3"`    // 为 Redis Search 模块启用 RESP3 的不稳定模式，默认 false
}

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

// 内置 lua 脚本
var internalScriptMap = map[string]string{
	"COMPARE_AND_DELETE": `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
	`,

	"POLLING": `
	local next = redis.call("INCRBY", KEYS[1], 1)
	if next > tonumber(ARGV[1], 10) then
		redis.call("SET", KEYS[1], 1)
		return 0
	end
	return next-1
	`,
}

// NewClient 创建 redis 客户端
func NewClient(ctx context.Context, cfg *ClientConfig) (client *RedisClient, err error) {
	if cfg == nil {
		err = fmt.Errorf("redis client config is nil")
		return
	}
	client = &RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr:            cfg.Addr,
			ClientName:      cfg.ClientName,
			Protocol:        cfg.Protocol,
			Username:        cfg.Username,
			Password:        cfg.Password,
			DB:              cfg.DB,
			MaxRetries:      cfg.MaxRetries,
			MinRetryBackoff: cfg.MinRetryBackoff,
			MaxRetryBackoff: cfg.MaxRetryBackoff,
			DialTimeout:     cfg.DialTimeout,
			ReadTimeout:     cfg.ReadTimeout,
			WriteTimeout:    cfg.WriteTimeout,
			PoolFIFO:        cfg.PoolFIFO,
			PoolSize:        cfg.PoolSize,
			PoolTimeout:     cfg.PoolTimeout,
			MinIdleConns:    cfg.MinIdleConns,
			MaxIdleConns:    cfg.MaxIdleConns,
			MaxActiveConns:  cfg.MaxActiveConns,
			ConnMaxIdleTime: cfg.ConnMaxIdleTime,
			ConnMaxLifetime: cfg.ConnMaxLifetime,
			TLSConfig:       cfg.TLSConfig,
			DisableIdentity: cfg.DisableIdentity,
			IdentitySuffix:  cfg.IdentitySuffix,
			UnstableResp3:   cfg.UnstableResp3,
		}),
		luaEvalShaMap: make(map[string]string),
	}
	for k, v := range internalScriptMap {
		if err = client.ScriptLoad(ctx, k, v); err != nil {
			return
		}
	}
	return
}

// Do 执行 redis 命令
func (rc *RedisClient) Do(ctx context.Context, cmd string, args ...any) (value any, err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	// 处理`redis`命令参数
	if err = utils.DoRedisArgs(0, args...); err != nil {
		return
	}
	// 执行`redis`命令
	cmdArgs := make([]any, 0, len(args)+1)
	cmdArgs = append(cmdArgs, cmd)
	cmdArgs = append(cmdArgs, args...)
	value, err = rc.client.Do(ctx, cmdArgs...).Result()
	err = noErrNil(err)
	return
}

// Pipeline 执行 redis 管道命令
func (rc *RedisClient) Pipeline(ctx context.Context, cmdArgsList ...[]any) (results []*PipelineResult, err error) {
	if len(cmdArgsList) == 0 {
		err = fmt.Errorf("pipeline cmd args list is empty")
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}
	// 执行`redis`管道命令
	p := rc.client.Pipeline()
	// 处理`redis`命令参数
	for _, cmdArgs := range cmdArgsList {
		if err = utils.DoRedisArgs(1, cmdArgs...); err != nil {
			return
		}
		// 执行`redis`命令
		p.Do(ctx, cmdArgs...)
	}
	var resList []redis.Cmder
	resList, err = p.Exec(ctx)
	if err = noErrNil(err); err != nil {
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
	script := utils.GetContents(scriptPath)
	if strings.EqualFold("", script) {
		err = fmt.Errorf("[%s] script not found", scriptPath)
		return
	}
	var evalsha string
	if evalsha, err = rc.client.ScriptLoad(ctx, script).Result(); err != nil {
		return
	}
	name := utils.Name(scriptPath)
	rc.luaEvalShaMap[name] = evalsha
	return
}

// Eval 执行 lua 脚本
func (rc *RedisClient) Eval(ctx context.Context, script string, keys []string, args ...any) (value any, err error) {
	// 处理`redis`命令参数
	if err = utils.DoRedisArgs(0, args...); err != nil {
		return
	}
	value, err = rc.client.Eval(ctx, script, keys, args...).Result()
	err = noErrNil(err)
	return
}

// EvalSha 执行 lua 脚本
func (rc *RedisClient) EvalSha(ctx context.Context, name string, keys []string, args ...any) (value any, err error) {
	evalsha, ok := rc.luaEvalShaMap[name]
	if !ok {
		err = fmt.Errorf("[%s] Script Not Found", name)
		return
	}
	// 处理`redis`命令参数
	if err = utils.DoRedisArgs(0, args...); err != nil {
		return
	}
	value, err = rc.client.EvalSha(ctx, evalsha, keys, args...).Result()
	err = noErrNil(err)
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
	if value, err = rc.Do(ctx, "SET", key, 1, "PX", cd.Milliseconds(), "NX"); err != nil {
		return
	}
	ok = gtkconv.ToBool(value)
	return
}

// CompareAndDelete 比较并删除
func (rc *RedisClient) CompareAndDelete(ctx context.Context, key string, value any) (ok bool, err error) {
	var result any
	if result, err = rc.EvalSha(ctx, "COMPARE_AND_DELETE", []string{key}, value); err != nil {
		return
	}
	ok = gtkconv.ToBool(result)
	return
}

// Polling 轮询
func (rc *RedisClient) Polling(ctx context.Context, key string, max int) (index int, err error) {
	var result any
	if result, err = rc.EvalSha(ctx, "POLLING", []string{key}, max); err != nil {
		return
	}
	index = gtkconv.ToInt(result)
	return
}

// Close 关闭 redis
func (rc *RedisClient) Close() (err error) {
	return rc.client.Close()
}

// noErrNil 处理 redis.Nil 错误
func noErrNil(err error) error {
	if err == redis.Nil {
		return nil
	}
	return err
}
