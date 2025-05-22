/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-04-23 00:30:12
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-22 10:31:13
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkmq

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/liusuxian/go-toolkit/gtkconv"
	"github.com/liusuxian/go-toolkit/gtkhttp"
	"github.com/liusuxian/go-toolkit/gtkjson"
	"github.com/liusuxian/go-toolkit/gtkredis"
	"hash/fnv"
	"math"
	"slices"
	"strings"
	"time"
)

// RedisMQConfig Redis 消息队列配置
type RedisMQConfig struct {
	Addr                  string              `json:"addr"`                  // redis 地址
	Username              string              `json:"username"`              // redis 用户名
	Password              string              `json:"password"`              // redis 密码
	DB                    int                 `json:"db"`                    // redis 数据库
	PoolSize              int                 `json:"poolSize"`              // redis 连接池大小，默认 20
	TLSConfig             *tls.Config         `json:"tlsConfig"`             // tls 配置
	Retries               int                 `json:"retries"`               // 发送消息失败后允许重试的次数，默认 2147483647
	RetryBackoff          time.Duration       `json:"retryBackoff"`          // 发送消息失败后，下一次重试发送前的等待时间，默认 100ms
	ExpiredTime           time.Duration       `json:"expiredTime"`           // 消息过期时间，默认 90天
	DelExpiredMsgInterval time.Duration       `json:"delExpiredMsgInterval"` // 删除过期消息的时间间隔，默认 1天
	WaitTimeout           time.Duration       `json:"waitTimeout"`           // 指定等待消息的最大时间，默认最大 2500ms
	OffsetReset           string              `json:"offsetReset"`           // 重置消费者偏移量的策略，可选值: 0-0 最早位置，$ 最新位置，默认 0-0
	Env                   string              `json:"env"`                   // 消息队列服务环境，默认 local
	ConsumerEnv           string              `json:"consumerEnv"`           // 消费者服务环境，默认和消息队列服务环境一致
	GlobalProducer        string              `json:"globalProducer"`        // 全局生产者名称，配置此项时，客户端将使用全局生产者，不再创建新的生产者，默认为空
	MQConfig              map[string]MQConfig `json:"mqConfig"`              // 消息队列配置，key 为消息队列名称
	ExcludeMQList         []string            `json:"excludeMQList"`         // 指定哪些消息队列不发送消息
}

// RedisMQConfigOption Redis 消息队列配置选项
type RedisMQConfigOption func(c *RedisMQConfig)

// RedisMQClient Redis 消息队列客户端
type RedisMQClient struct {
	rc          *gtkredis.RedisClient // redis 客户端
	config      *RedisMQConfig        // Redis 消息队列配置
	producerMap map[string]bool
	consumerMap map[string]bool
	logger      ILogger   // 日志接口
	quitChan    chan bool // 退出信号
}

// 内置 lua 脚本
var internalScriptMap = map[string]string{
	"XGROUP_CREATE": `
		for i = 0, tonumber(ARGV[1], 10) - 1 do
			local partitionQueue = KEYS[1] .. "@" .. i
    	local partitionGroup = KEYS[2] .. "@" .. i
			-- 检查指定的流是否存在
			local partitionQueueExists = tonumber(redis.call('EXISTS', partitionQueue), 10)
			if partitionQueueExists == 0 then
				-- 如果流不存在，创建流和消费者组
				redis.call("XGROUP", "CREATE", partitionQueue, partitionGroup, ARGV[2], "MKSTREAM")
			else
				-- 如果流存在，检查消费者组是否存在
				local partitionGroupExists = false
				local partitionGroups = redis.call("XINFO", "GROUPS", partitionQueue)
				for j, group in ipairs(partitionGroups) do
					if group[2] == partitionGroup then
						partitionGroupExists = true
          	break
      		end
    		end
				-- 如果消费者组不存在，创建消费者组
				if not partitionGroupExists then
					redis.call("XGROUP", "CREATE", partitionQueue, partitionGroup, ARGV[2], "MKSTREAM")
				end
			end
		end
    `,

	"SEND_MESSAGE": `
		local partition = tonumber(ARGV[1], 10)
		local targetPartition = 0
		-- 寻找目标分区
		if partition < 0 then
			-- 如果 partition 为负数，则选择分区长度最短的队列
			local partitionNum = tonumber(ARGV[2], 10)
			local minPartitionQueueLen = math.huge

			for i = 0, partitionNum - 1 do
				local partitionQueue = KEYS[1] .. "@" .. i
				local rawLen = redis.call('XLEN', partitionQueue)
				local partitionQueueLen = 0
				-- 检查rawLen的类型来处理不同的返回值
				if type(rawLen) == "table" then
        	if #rawLen == 0 then
            partitionQueueLen = 0
        	else
            partitionQueueLen = tonumber(rawLen[1], 10) or 0
        	end
    		else
        	partitionQueueLen = tonumber(rawLen, 10) or 0
    		end
				-- 在第一次循环时, minPartitionQueueLen 会被设置为第一个队列的长度
				if partitionQueueLen < minPartitionQueueLen then
					minPartitionQueueLen = partitionQueueLen
					targetPartition = i
				end
			end
		else
			-- 如果 partition 为非负整数，则选择指定的分区
			targetPartition = partition
		end
		-- 发送消息到目标分区
		local targetPartitionQueue = KEYS[1] .. "@" .. targetPartition
		redis.call("XADD", targetPartitionQueue, "*", "key", ARGV[3], "value", ARGV[4], "timestamp", ARGV[5], "expire_time", ARGV[6])
		return targetPartition
		`,

	"GET_DELAY_MESSAGES": `
		local messages = redis.call('ZRANGEBYSCORE', KEYS[1], 0, tonumber(ARGV[1], 10), 'LIMIT', 0, tonumber(ARGV[2], 10))
		if #messages > 0 then
			redis.call('ZREM', KEYS[1], unpack(messages))
		end
		return messages
		`,
}

const (
	defaultPartitionNum uint32 = 12                        // 默认分区数
	partitionAny        int32  = -1                        // 任意分区
	delayQueueKey       string = "gtkmq:delay:queue:%s_%s" // 延迟队列redis key
)

// delayMessage 延迟消息
type delayMessage struct {
	UUID      string    `json:"uuid"`          // 消息唯一标识
	Queue     string    `json:"queue"`         // 队列名称
	Key       string    `json:"key,omitempty"` // 键
	Data      any       `json:"data"`          // 数据
	Timestamp time.Time `json:"timestamp"`     // 发送消息的时间戳
}

// NewRedisMQClientWithOption 创建 Redis 消息队列客户端
func NewRedisMQClientWithOption(ctx context.Context, opts ...RedisMQConfigOption) (client *RedisMQClient, err error) {
	client = &RedisMQClient{
		config: &RedisMQConfig{
			MQConfig:      make(map[string]MQConfig),
			ExcludeMQList: make([]string, 0),
		},
		producerMap: make(map[string]bool),
		consumerMap: make(map[string]bool),
		quitChan:    make(chan bool),
	}
	for _, opt := range opts {
		opt(client.config)
	}
	// redis 地址
	if client.config.Addr == "" {
		client.config.Addr = "127.0.0.1:6379"
	}
	// redis 连接池大小，默认 20
	if client.config.PoolSize <= 0 {
		client.config.PoolSize = 20
	}
	// 发送消息失败后允许重试的次数，默认 2147483647
	if client.config.Retries <= 0 {
		client.config.Retries = math.MaxInt32
	}
	// 发送消息失败后，下一次重试发送前的等待时间，默认 100ms
	if client.config.RetryBackoff <= time.Duration(0) {
		client.config.RetryBackoff = time.Millisecond * 100
	}
	// 消息过期时间，默认 90天
	if client.config.ExpiredTime <= time.Duration(0) {
		client.config.ExpiredTime = time.Hour * 24 * 90
	}
	// 删除过期消息的时间间隔，默认 1天
	if client.config.DelExpiredMsgInterval <= time.Duration(0) {
		client.config.DelExpiredMsgInterval = time.Hour * 24 * 1
	}
	// 指定等待消息的最大时间，默认最大 2500ms
	if client.config.WaitTimeout <= time.Duration(0) || client.config.WaitTimeout > time.Millisecond*2500 {
		client.config.WaitTimeout = time.Millisecond * 2500
	}
	// 重置消费者偏移量的策略，可选值: 0-0 最早位置，$ 最新位置，默认 0-0
	if client.config.OffsetReset == "" {
		client.config.OffsetReset = "0-0"
	}
	// 消息队列服务环境，默认 local
	if client.config.Env == "" {
		client.config.Env = "local"
	}
	// 消费者服务环境，默认和消息队列服务环境一致
	if client.config.ConsumerEnv == "" {
		client.config.ConsumerEnv = client.config.Env
	}
	// redis 客户端
	client.rc = gtkredis.NewClientWithOption(ctx, func(cc *gtkredis.ClientConfig) {
		cc.Addr = client.config.Addr
		cc.Password = client.config.Password
		cc.DB = client.config.DB
		cc.PoolSize = client.config.PoolSize
		cc.TLSConfig = client.config.TLSConfig
	})
	for k, v := range internalScriptMap {
		if err = client.rc.ScriptLoad(ctx, k, v); err != nil {
			client.rc.Close()
			return
		}
	}
	// 默认日志对象
	client.logger = newDefaultLogger()
	// 启动所有队列删除过期消息处理器
	client.startAllQueueDelExpiredMsgProcessor(ctx)
	// 启动所有延迟队列处理器
	client.startAllDelayQueueProcessor(ctx)
	return
}

// NewRedisMQClientWithConfig 创建 Redis 消息队列客户端
func NewRedisMQClientWithConfig(ctx context.Context, cfg *RedisMQConfig) (client *RedisMQClient, err error) {
	if cfg == nil {
		cfg = &RedisMQConfig{
			MQConfig:      make(map[string]MQConfig),
			ExcludeMQList: make([]string, 0),
		}
	}
	client = &RedisMQClient{
		config:      cfg,
		producerMap: make(map[string]bool),
		consumerMap: make(map[string]bool),
		quitChan:    make(chan bool),
	}
	// redis 地址
	if client.config.Addr == "" {
		client.config.Addr = "127.0.0.1:6379"
	}
	// redis 连接池大小，默认 20
	if client.config.PoolSize <= 0 {
		client.config.PoolSize = 20
	}
	// 发送消息失败后允许重试的次数，默认 2147483647
	if client.config.Retries <= 0 {
		client.config.Retries = math.MaxInt32
	}
	// 发送消息失败后，下一次重试发送前的等待时间，默认 100ms
	if client.config.RetryBackoff <= time.Duration(0) {
		client.config.RetryBackoff = time.Millisecond * 100
	}
	// 消息过期时间，默认 90天
	if client.config.ExpiredTime <= time.Duration(0) {
		client.config.ExpiredTime = time.Hour * 24 * 90
	}
	// 删除过期消息的时间间隔，默认 1天
	if client.config.DelExpiredMsgInterval <= time.Duration(0) {
		client.config.DelExpiredMsgInterval = time.Hour * 24 * 1
	}
	// 指定等待消息的最大时间，默认最大 2500ms
	if client.config.WaitTimeout <= time.Duration(0) || client.config.WaitTimeout > time.Millisecond*2500 {
		client.config.WaitTimeout = time.Millisecond * 2500
	}
	// 重置消费者偏移量的策略，可选值: 0-0 最早位置，$ 最新位置，默认 0-0
	if client.config.OffsetReset == "" {
		client.config.OffsetReset = "0-0"
	}
	// 消息队列服务环境，默认 local
	if client.config.Env == "" {
		client.config.Env = "local"
	}
	// 消费者服务环境，默认和消息队列服务环境一致
	if client.config.ConsumerEnv == "" {
		client.config.ConsumerEnv = client.config.Env
	}
	// redis 客户端
	client.rc = gtkredis.NewClientWithOption(ctx, func(cc *gtkredis.ClientConfig) {
		cc.Addr = client.config.Addr
		cc.Password = client.config.Password
		cc.DB = client.config.DB
		cc.PoolSize = client.config.PoolSize
		cc.TLSConfig = client.config.TLSConfig
	})
	for k, v := range internalScriptMap {
		if err = client.rc.ScriptLoad(ctx, k, v); err != nil {
			client.rc.Close()
			return
		}
	}
	// 默认日志对象
	client.logger = newDefaultLogger()
	// 启动所有队列删除过期消息处理器
	client.startAllQueueDelExpiredMsgProcessor(ctx)
	// 启动所有延迟队列处理器
	client.startAllDelayQueueProcessor(ctx)
	return
}

// SetLogger 设置日志对象
func (mq *RedisMQClient) SetLogger(logger ILogger) {
	mq.logger = logger
}

// PrintClientConfig 打印消息队列客户端配置
func (mq *RedisMQClient) PrintClientConfig(ctx context.Context) {
	mq.logger.Debugf(ctx, "client config: %s\n", gtkjson.MustString(mq.config))
}

// NewProducer 创建生产者
func (mq *RedisMQClient) NewProducer(ctx context.Context, queue string) (err error) {
	// 获取生产者配置
	var (
		isStart  bool
		mqConfig *MQConfig
	)
	if isStart, mqConfig, err = mq.getProducerConfig(queue); err != nil {
		return
	}
	if !isStart {
		return
	}
	// 创建生产者
	var (
		producerName  = mq.getProducerName(queue)
		fullQueueName = mq.getFullQueueName(queue)
	)
	// 判断是否配置了全局生产者名称
	globalProducerName := strings.Trim(mq.config.GlobalProducer, " ")
	if globalProducerName != "" {
		producerName = mq.getGlobalProducerName(globalProducerName)
		if _, ok := mq.producerMap[producerName]; ok {
			mq.logger.Infof(ctx, "new producer: %s, queue: %s, partitionNum: %d success", producerName, fullQueueName, mqConfig.PartitionNum)
			return
		}
	} else {
		if _, ok := mq.producerMap[producerName]; ok {
			return fmt.Errorf("new producer: %s, queue: %s, partitionNum: %d already exists", producerName, fullQueueName, mqConfig.PartitionNum)
		}
	}
	mq.producerMap[producerName] = true
	mq.logger.Infof(ctx, "new producer: %s, queue: %s, partitionNum: %d success", producerName, fullQueueName, mqConfig.PartitionNum)
	return
}

// NewConsumer 创建消费者
func (mq *RedisMQClient) NewConsumer(ctx context.Context, queue string) (err error) {
	// 获取消费者配置
	var (
		isStart  bool
		mqConfig *MQConfig
	)
	if isStart, mqConfig, err = mq.getConsumerConfig(queue); err != nil {
		return
	}
	if !isStart {
		return
	}
	// 创建消费者
	var (
		consumerNameList = []string{mq.getConsumerName(queue)}
		groupList        = []string{mq.getConsumerGroupName(queue)}
		fullQueueName    = mq.getFullQueueName(queue)
	)
	if len(mqConfig.Groups) > 0 {
		consumerNameList = make([]string, 0, len(mqConfig.Groups))
		groupList = make([]string, 0, len(mqConfig.Groups))
		for _, g := range mqConfig.Groups {
			consumerNameList = append(consumerNameList, mq.getConsumerName(g))
			groupList = append(groupList, mq.getConsumerGroupName(g))
		}
	}

	for i := 0; i < len(consumerNameList); i++ {
		var (
			consumerName = consumerNameList[i]
			group        = groupList[i]
		)
		if _, ok := mq.consumerMap[consumerName]; ok {
			return fmt.Errorf("new consumer: %s, queue: %s, group: %s, partitionNum: %d already exists", consumerName, fullQueueName, group, mqConfig.PartitionNum)
		}
		if _, err = mq.rc.EvalSha(ctx, "XGROUP_CREATE", []string{fullQueueName, group}, mqConfig.PartitionNum, mq.config.OffsetReset); err != nil {
			return
		}
		mq.consumerMap[consumerName] = true
		mq.logger.Infof(ctx, "new consumer: %s, queue: %s, group: %s, partitionNum: %d success", consumerName, fullQueueName, group, mqConfig.PartitionNum)
	}
	return
}

// SendMessage 发送消息
func (mq *RedisMQClient) SendMessage(ctx context.Context, queue string, producerMessage *ProducerMessage) (err error) {
	// 获取生产者配置
	var (
		isStart  bool
		mqConfig *MQConfig
	)
	if isStart, mqConfig, err = mq.getProducerConfig(queue); err != nil {
		return
	}
	if !isStart {
		return
	}
	if mqConfig.IsDelayQueue && !producerMessage.DelayTime.IsZero() {
		// 发送延迟消息
		return mq.sendDelayMessage(ctx, queue, producerMessage)
	}

	// 序列化消息
	var dataMap map[string]any
	if dataMap, err = gtkconv.ToStringMapE(producerMessage.Data); err != nil {
		return
	}
	// 处理数据
	var dataBytes []byte
	if dataBytes, err = json.Marshal(dataMap); err != nil {
		return
	}
	producerMessage.dataBytes = dataBytes
	// 发送消息
	return mq.sendMessage(ctx, queue, mqConfig, producerMessage)
}

// Subscribe 订阅数据
func (mq *RedisMQClient) Subscribe(ctx context.Context, queue string, fn func(message *MQMessage) error, group ...string) (err error) {
	return mq.handelSubscribe(ctx, queue, false, func(messages []*MQMessage) error {
		return fn(messages[0])
	}, group...)
}

// BatchSubscribe 批量订阅数据
func (mq *RedisMQClient) BatchSubscribe(ctx context.Context, queue string, fn func(messages []*MQMessage) error, group ...string) (err error) {
	return mq.handelSubscribe(ctx, queue, true, fn, group...)
}

// GetExpiredMessages 获取过期消息，每个分区每次最多返回 100 条
//
//	isDelete: 是否删除过期消息
func (mq *RedisMQClient) GetExpiredMessages(ctx context.Context, queue string, isDelete bool) (messages map[int32][]*MQMessage, err error) {
	// 获取消息队列的分区数量
	var partitionNum uint32
	if partitionNum, err = mq.getPartitionNum(queue); err != nil {
		return
	}
	// 组装命令参数
	cmdArgsList := make([][]any, 0, partitionNum)
	for i := uint32(0); i < partitionNum; i++ {
		partition := int32(i)
		cmdArgsList = append(cmdArgsList, []any{"XRANGE", mq.getPartitionQueueName(queue, partition), "-", "+", "COUNT", 100})
	}
	// 执行 redis 管道命令
	var results []*gtkredis.PipelineResult
	if results, err = mq.rc.Pipeline(ctx, cmdArgsList...); err != nil {
		return
	}
	// 组装结果数据
	messages = make(map[int32][]*MQMessage)
	for i, result := range results {
		if result.Err != nil {
			err = result.Err
			return
		}

		partition := int32(i)
		resultSliceSlice := gtkconv.ToSlice(result.Val)
		mqMessageList := make([]*MQMessage, 0, len(resultSliceSlice))
		// 遍历结果数据
		for _, resultSliceAny := range resultSliceSlice {
			resultSlice := gtkconv.ToSlice(resultSliceAny)
			dataSlice := gtkconv.ToSlice(resultSlice[1])
			expireTime := gtkconv.ToInt64(dataSlice[7])
			// 判断是否过期
			if time.Now().Unix() >= expireTime {
				offset := gtkconv.ToString(resultSlice[0])
				mqMessage := &MQMessage{
					MQPartition: MQPartition{
						Queue:         queue,
						PartitionName: mq.getPartitionQueueName(queue, partition),
						Partition:     partition,
						Offset:        offset,
					},
					Key:        gtkconv.ToBytes(dataSlice[1]),
					Value:      gtkconv.ToBytes(dataSlice[3]),
					Timestamp:  time.UnixMilli(gtkconv.ToInt64(dataSlice[5])),
					ExpireTime: time.Unix(gtkconv.ToInt64(dataSlice[7]), 0),
				}
				mqMessageList = append(mqMessageList, mqMessage)
			}
		}
		if len(mqMessageList) > 0 {
			messages[partition] = mqMessageList
		}
	}
	// 删除过期消息
	if isDelete {
		err = mq.delExpiredMessages(ctx, messages)
	}
	return
}

// ResetConsumerOffset 重置消费起点，所有分区（请谨慎使用）
//
//	offset: 0-0 重置为最早位置
//	offset: $ 重置为最新位置
func (mq *RedisMQClient) ResetConsumerOffset(ctx context.Context, queue string, offset string, group ...string) (err error) {
	// 检查 offset 参数
	if offset != "0-0" && offset != "$" {
		err = fmt.Errorf("offset must be 0-0 or $")
		return
	}
	// 获取消费者配置
	var mqConfig *MQConfig
	if _, mqConfig, err = mq.getConsumerConfig(queue); err != nil {
		return
	}
	if len(mqConfig.Groups) > 0 && len(group) > 0 {
		if !slices.Contains(mqConfig.Groups, group[0]) {
			return fmt.Errorf("group: %s not found in groups: %s", group[0], mqConfig.Groups)
		}
	}
	// 组装命令参数
	cmdArgsList := make([][]any, 0, mqConfig.PartitionNum)
	for i := uint32(0); i < mqConfig.PartitionNum; i++ {
		var (
			partition          = int32(i)
			partitionGroupName = mq.getPartitionGroupName(queue, partition)
		)
		if len(group) > 0 {
			partitionGroupName = mq.getPartitionGroupName(group[0], partition)
		}
		cmdArgsList = append(cmdArgsList, []any{"XGROUP", "SETID", mq.getPartitionQueueName(queue, partition), partitionGroupName, offset})
	}
	// 执行 redis 管道命令
	var results []*gtkredis.PipelineResult
	if results, err = mq.rc.Pipeline(ctx, cmdArgsList...); err != nil {
		return
	}
	for _, result := range results {
		if result.Err != nil {
			err = result.Err
			return
		}
	}
	return
}

// ResetConsumerOffsetByPartition 重置消费起点，指定分区（请谨慎使用）
//
//	offset: 0-0 重置为最早位置
//	offset: $ 重置为最新位置
//	offset: <ID> 重置为指定位置
func (mq *RedisMQClient) ResetConsumerOffsetByPartition(ctx context.Context, queue string, partition int32, offset string, group ...string) (err error) {
	// 获取消费者配置
	var mqConfig *MQConfig
	if _, mqConfig, err = mq.getConsumerConfig(queue); err != nil {
		return
	}
	if len(mqConfig.Groups) > 0 && len(group) > 0 {
		if !slices.Contains(mqConfig.Groups, group[0]) {
			return fmt.Errorf("group: %s not found in groups: %s", group[0], mqConfig.Groups)
		}
	}
	partitionGroupName := mq.getPartitionGroupName(queue, partition)
	if len(group) > 0 {
		partitionGroupName = mq.getPartitionGroupName(group[0], partition)
	}
	_, err = mq.rc.Do(ctx, "XGROUP", "SETID", mq.getPartitionQueueName(queue, partition), partitionGroupName, offset)
	return
}

// DelGroup 删除消费者组（请谨慎使用）
func (mq *RedisMQClient) DelGroup(ctx context.Context, queue string, group ...string) (err error) {
	// 获取消费者配置
	var mqConfig *MQConfig
	if _, mqConfig, err = mq.getConsumerConfig(queue); err != nil {
		return
	}
	if len(mqConfig.Groups) > 0 && len(group) > 0 {
		if !slices.Contains(mqConfig.Groups, group[0]) {
			return fmt.Errorf("group: %s not found in groups: %s", group[0], mqConfig.Groups)
		}
	}
	// 组装命令参数
	cmdArgsList := make([][]any, 0, mqConfig.PartitionNum)
	for i := uint32(0); i < mqConfig.PartitionNum; i++ {
		var (
			partition          = int32(i)
			partitionGroupName = mq.getPartitionGroupName(queue, partition)
		)
		if len(group) > 0 {
			partitionGroupName = mq.getPartitionGroupName(group[0], partition)
		}
		cmdArgsList = append(cmdArgsList, []any{"XGROUP", "DESTROY", mq.getPartitionQueueName(queue, partition), partitionGroupName})
	}
	// 执行 redis 管道命令
	var results []*gtkredis.PipelineResult
	if results, err = mq.rc.Pipeline(ctx, cmdArgsList...); err != nil {
		return
	}
	for _, result := range results {
		if result.Err != nil {
			err = result.Err
			return
		}
	}
	return
}

// DelQueue 删除队列（请谨慎使用）
func (mq *RedisMQClient) DelQueue(ctx context.Context, queue string) (err error) {
	// 获取消息队列的分区数量
	var partitionNum uint32
	if partitionNum, err = mq.getPartitionNum(queue); err != nil {
		return
	}
	// 组装命令参数
	cmdArgs := make([]any, 0, partitionNum)
	for i := uint32(0); i < partitionNum; i++ {
		partition := int32(i)
		cmdArgs = append(cmdArgs, mq.getPartitionQueueName(queue, partition))
	}
	// 执行 redis 命令
	_, err = mq.rc.Do(ctx, "DEL", cmdArgs...)
	return
}

// Close 关闭客户端
func (mq *RedisMQClient) Close() (err error) {
	mq.quitChan <- true
	return mq.rc.Close()
}

// startAllQueueDelExpiredMsgProcessor 启动所有队列删除过期消息处理器
func (mq *RedisMQClient) startAllQueueDelExpiredMsgProcessor(ctx context.Context) {
	for queue, mqConfig := range mq.config.MQConfig {
		if mqConfig.Mode == ModeBoth || mqConfig.Mode == ModeProducer {
			go mq.startQueueDelExpiredMsgProcessor(ctx, queue)
		}
	}
}

// startQueueDelExpiredMsgProcessor 启动队列删除过期消息处理器
func (mq *RedisMQClient) startQueueDelExpiredMsgProcessor(ctx context.Context, queue string) {
	ticker := time.NewTicker(mq.config.DelExpiredMsgInterval)
	for {
		select {
		case <-ticker.C:
			if _, err := mq.GetExpiredMessages(ctx, queue, true); err != nil {
				mq.logger.Errorf(ctx, "delete expired messages, queue: %s error: %+v", queue, err)
			}
		case <-mq.quitChan:
			ticker.Stop()
			return
		}
	}
}

// startAllDelayQueueProcessor 启动所有延迟队列处理器
func (mq *RedisMQClient) startAllDelayQueueProcessor(ctx context.Context) {
	for queue, mqConfig := range mq.config.MQConfig {
		if mqConfig.Mode == ModeBoth || mqConfig.Mode == ModeProducer {
			if mqConfig.IsDelayQueue {
				var (
					interval  = mqConfig.DelayQueueCheckInterval
					batchSize = mqConfig.DelayQueueBatchSize
				)
				// 延迟队列检查间隔，默认 10s
				if interval <= time.Duration(0) {
					interval = time.Second * 10
				}
				// 延迟队列批处理大小，默认 100
				if batchSize <= 0 {
					batchSize = 100
				}
				go mq.startDelayQueueProcessor(ctx, queue, interval, batchSize)
			}
		}
	}
}

// startDelayQueueProcessor 启动延迟队列处理器
func (mq *RedisMQClient) startDelayQueueProcessor(ctx context.Context, queue string, interval time.Duration, batchSize int) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			// 执行 lua 脚本
			var (
				keys = []string{fmt.Sprintf(delayQueueKey, mq.config.Env, queue)}
				args = []any{
					time.Now().UnixMilli(),
					batchSize,
				}
				value any
				err   error
			)
			if value, err = mq.rc.EvalSha(ctx, "GET_DELAY_MESSAGES", keys, args...); err != nil {
				mq.logger.Errorf(ctx, "get delay messages, queue: %s error: %+v", queue, err)
				return
			}
			// 解析结果数据
			var (
				anyList     = gtkconv.ToSlice(value)
				messageList = make([]*delayMessage, 0, len(anyList))
			)
			for _, v := range anyList {
				message := &delayMessage{}
				if err = gtkconv.ToStructE(v, &message); err != nil {
					mq.logger.Errorf(ctx, "parse delay message, queue: %s error: %+v", queue, err)
					continue
				}
				messageList = append(messageList, message)
			}
			// 发送到实际队列
			for _, message := range messageList {
				mq.SendMessage(ctx, message.Queue, &ProducerMessage{
					Key:  message.Key,
					Data: message.Data,
				})
			}
		case <-mq.quitChan:
			ticker.Stop()
			return
		}
	}
}

// sendMessage 发送消息
func (mq *RedisMQClient) sendMessage(ctx context.Context, queue string, mqConfig *MQConfig, producerMessage *ProducerMessage) (err error) {
	// 检测哪些消息队列不发送消息
	if slices.Contains(mq.config.ExcludeMQList, queue) {
		return
	}
	// 计算分区号
	var partition int32
	if producerMessage.Key == "" {
		partition = partitionAny
	} else {
		hash := fnv.New32a()
		hash.Write([]byte(producerMessage.Key))
		partition = int32(hash.Sum32() % mqConfig.PartitionNum)
	}
	// 判断是否配置了全局生产者名称
	var (
		producerName       = mq.getProducerName(queue)
		globalProducerName = strings.Trim(mq.config.GlobalProducer, " ")
	)
	if globalProducerName != "" {
		producerName = mq.getGlobalProducerName(globalProducerName)
	}
	// 发送消息
	var now time.Time
	if err = gtkhttp.Retry(ctx, func(ctx context.Context) (e error) {
		keys := []string{mq.getFullQueueName(queue)}
		now = time.Now()
		args := []any{
			partition,
			mqConfig.PartitionNum,
			producerMessage.Key,
			producerMessage.dataBytes,
			now.UnixMilli(),
			now.Add(mq.config.ExpiredTime).Unix(),
		}
		// 执行 lua 脚本
		var value any
		if value, e = mq.rc.EvalSha(ctx, "SEND_MESSAGE", keys, args...); e != nil {
			return
		}
		partition = gtkconv.ToInt32(value)
		return
	}, mq.config.Retries, mq.config.RetryBackoff, false); err != nil {
		mq.logger.Errorf(ctx, "producer: %s, send message, queue: %s, partition: %d, data: %s, timestamp: %v, error: %+v", producerName, queue, partition, gtkjson.MustString(producerMessage), now, err)
		return
	}
	mq.logger.Debugf(ctx, "producer: %s, send message, queue: %s, partition: %d, data: %s, timestamp: %v, success", producerName, queue, partition, gtkjson.MustString(producerMessage), now)
	return
}

// sendDelayMessage 发送延迟消息
func (mq *RedisMQClient) sendDelayMessage(ctx context.Context, queue string, producerMessage *ProducerMessage) (err error) {
	// 检测哪些消息队列不发送消息
	if slices.Contains(mq.config.ExcludeMQList, queue) {
		return
	}
	// 判断是否配置了全局生产者名称
	var (
		producerName       = mq.getProducerName(queue)
		globalProducerName = strings.Trim(mq.config.GlobalProducer, " ")
	)
	if globalProducerName != "" {
		producerName = mq.getGlobalProducerName(globalProducerName)
	}
	// 构造延迟消息
	delayMsg := &delayMessage{
		UUID:  uuid.New().String(),
		Queue: queue,
		Key:   producerMessage.Key,
		Data:  producerMessage.Data,
	}
	// 将消息添加到延迟队列
	if err = gtkhttp.Retry(ctx, func(ctx context.Context) (e error) {
		delayMsg.Timestamp = time.Now()
		_, e = mq.rc.Do(ctx, "ZADD", fmt.Sprintf(delayQueueKey, mq.config.Env, queue), producerMessage.DelayTime.UnixMilli(), delayMsg)
		return
	}, mq.config.Retries, mq.config.RetryBackoff, false); err != nil {
		mq.logger.Errorf(ctx, "producer: %s send delay message, data: %s error: %+v", producerName, gtkjson.MustString(delayMsg), err)
		return
	}
	mq.logger.Debugf(ctx, "producer: %s send delay message, data: %s success", producerName, gtkjson.MustString(delayMsg))
	return
}

// handelSubscribe 处理订阅数据
func (mq *RedisMQClient) handelSubscribe(ctx context.Context, queue string, isBatch bool, fn func(messages []*MQMessage) error, group ...string) (err error) {
	// 获取消费者配置
	var (
		isStart  bool
		mqConfig *MQConfig
	)
	if isStart, mqConfig, err = mq.getConsumerConfig(queue); err != nil {
		return
	}
	if !isStart {
		return
	}
	if len(mqConfig.Groups) > 0 && len(group) > 0 {
		if !slices.Contains(mqConfig.Groups, group[0]) {
			return fmt.Errorf("group: %s not found in groups: %s", group[0], mqConfig.Groups)
		}
	}
	// 订阅数据
	var (
		block           = mq.config.WaitTimeout.Milliseconds()
		readWaitTimeout = time.Millisecond * 100
		count           = 1
	)
	if isBatch {
		readWaitTimeout = mqConfig.BatchInterval
		count = mqConfig.BatchSize
	}
	for i := int32(0); i < int32(mqConfig.PartitionNum); i++ {
		go func(partition int32) {
			var (
				partitionGroupName    = mq.getPartitionGroupName(queue, partition)
				partitionConsumerName = mq.getPartitionConsumerName(queue, partition)
				partitionQueueName    = mq.getPartitionQueueName(queue, partition)
			)
			if len(group) > 0 {
				partitionGroupName = mq.getPartitionGroupName(group[0], partition)
				partitionConsumerName = mq.getPartitionConsumerName(group[0], partition)
			}
			// 添加对 panic 的处理
			defer func() {
				if r := recover(); r != nil {
					mq.logger.Errorf(ctx, "partition-consumer: %s, partition-queue: %s, panic: %+v", partitionConsumerName, partitionQueueName, r)
				}
			}()

			for {
				select {
				case <-ctx.Done():
					return
				default:
					// 读取数据
					value, e := mq.rc.Do(ctx, "XREADGROUP", "GROUP", partitionGroupName, partitionConsumerName, "COUNT", count, "BLOCK", block, "STREAMS", partitionQueueName, ">")
					if e != nil {
						mq.logger.Errorf(ctx, "partition-consumer: %s, partition-queue: %s, error: %+v", partitionConsumerName, partitionQueueName, e)
						time.Sleep(readWaitTimeout)
					} else if value != nil {
						// 处理结果数据
						valueMap := gtkconv.ToStringMap(value)
						result := valueMap[partitionQueueName]
						resultSliceSlice := gtkconv.ToSlice(result)
						mqMessageList := make([]*MQMessage, 0, len(resultSliceSlice))
						// 遍历结果数据
						for _, resultSliceAny := range resultSliceSlice {
							resultSlice := gtkconv.ToSlice(resultSliceAny)
							dataSlice := gtkconv.ToSlice(resultSlice[1])
							expireTime := gtkconv.ToInt64(dataSlice[7])
							// 判断是否过期
							if time.Now().Unix() < expireTime {
								offset := gtkconv.ToString(resultSlice[0])
								mqMessage := &MQMessage{
									MQPartition: MQPartition{
										Queue:         queue,
										PartitionName: mq.getPartitionQueueName(queue, partition),
										Partition:     partition,
										Offset:        offset,
									},
									Key:        gtkconv.ToBytes(dataSlice[1]),
									Value:      gtkconv.ToBytes(dataSlice[3]),
									Timestamp:  time.UnixMilli(gtkconv.ToInt64(dataSlice[5])),
									ExpireTime: time.Unix(gtkconv.ToInt64(dataSlice[7]), 0),
								}
								mqMessageList = append(mqMessageList, mqMessage)
							}
						}
						if len(mqMessageList) > 0 {
							mq.handelData(ctx, mqConfig, partitionConsumerName, partitionGroupName, mqMessageList, fn)
						}
					}
				}
			}
		}(i)
	}
	return
}

// handelData 处理数据
func (mq *RedisMQClient) handelData(ctx context.Context, mqConfig *MQConfig, partitionConsumerName, partitionGroupName string, messages []*MQMessage, fn func(messages []*MQMessage) error) {
	// 判断是否有数据
	length := len(messages)
	if length == 0 {
		return
	}

	var (
		lastMessage        = messages[length-1]
		partitionQueueName = lastMessage.MQPartition.PartitionName
		partition          = lastMessage.MQPartition.Partition
		offset             = lastMessage.MQPartition.Offset
		key                = string(lastMessage.Key)
		content            = string(lastMessage.Value)
		timestamp          = lastMessage.Timestamp
	)
	// 执行处理函数
	if err := fn(messages); err != nil {
		mq.logger.Errorf(ctx, "partition-consumer: %s, partition-queue: %s, partition: %d, offset: %v, key: %s, content: %s, timestamp: %v, error: %+v",
			partitionConsumerName, partitionQueueName, partition, offset, key, content, timestamp, err)
		// 重试处理函数
		var (
			retryMaxCount = mqConfig.RetryMaxCount
			count         = 0
		)
		for retryMaxCount == 0 || (retryMaxCount > 0 && count < retryMaxCount) {
			count++
			time.Sleep(mqConfig.RetryDelay)
			if err := fn(messages); err != nil {
				mq.logger.Errorf(ctx, "partition-consumer: %s, partition-queue: %s, partition: %d, offset: %v, key: %s, content: %s, timestamp: %v, error: %+v",
					partitionConsumerName, partitionQueueName, partition, offset, key, content, timestamp, err)
				continue
			}
			break
		}
	}
	// 提交
	var cmdArgs = make([]any, 0, length+2)
	cmdArgs = append(cmdArgs, lastMessage.MQPartition.PartitionName, partitionGroupName)
	for _, message := range messages {
		cmdArgs = append(cmdArgs, message.MQPartition.Offset)
	}
	gtkhttp.Retry(ctx, func(ctx context.Context) (err error) {
		_, err = mq.rc.Do(ctx, "XACK", cmdArgs...)
		return
	}, mq.config.Retries, mq.config.RetryBackoff, false)
}

// delExpiredMessages 删除过期消息
func (mq *RedisMQClient) delExpiredMessages(ctx context.Context, messages map[int32][]*MQMessage) (err error) {
	if len(messages) == 0 {
		return
	}
	// 组装命令参数
	cmdArgsList := make([][]any, 0, len(messages))
	for _, mqMessageList := range messages {
		cmdArgs := make([]any, 0, len(mqMessageList)+2)
		cmdArgs = append(cmdArgs, "XDEL", mqMessageList[0].MQPartition.PartitionName)
		for _, mqMessage := range mqMessageList {
			cmdArgs = append(cmdArgs, mqMessage.MQPartition.Offset)
		}
		cmdArgsList = append(cmdArgsList, cmdArgs)
	}
	// 执行 redis 管道命令
	var results []*gtkredis.PipelineResult
	if results, err = mq.rc.Pipeline(ctx, cmdArgsList...); err != nil {
		return
	}
	for _, result := range results {
		if result.Err != nil {
			err = result.Err
			return
		}
	}
	return
}

// getProducerConfig 获取生产者配置
func (mq *RedisMQClient) getProducerConfig(queue string) (isStart bool, mqConfig *MQConfig, err error) {
	if config, ok := mq.config.MQConfig[queue]; ok {
		isStart = (config.Mode == ModeBoth || config.Mode == ModeProducer)
		mqConfig = &MQConfig{}
		if config.PartitionNum > 0 {
			mqConfig.PartitionNum = config.PartitionNum
		} else {
			// 默认分区数
			mqConfig.PartitionNum = defaultPartitionNum
		}
		// 是否开启延迟队列
		mqConfig.IsDelayQueue = config.IsDelayQueue
		return
	}
	err = fmt.Errorf("queue `%s` Not Found", queue)
	return
}

// getConsumerConfig 获取消费者配置
func (mq *RedisMQClient) getConsumerConfig(queue string) (isStart bool, mqConfig *MQConfig, err error) {
	if config, ok := mq.config.MQConfig[queue]; ok {
		isStart = (config.Mode == ModeBoth || config.Mode == ModeConsumer)
		mqConfig = &MQConfig{}
		// 消息队列分区数量，默认 12 个分区
		if config.PartitionNum > 0 {
			mqConfig.PartitionNum = config.PartitionNum
		} else {
			// 默认分区数
			mqConfig.PartitionNum = defaultPartitionNum
		}
		// 指定消费者组名称列表
		mqConfig.Groups = config.Groups
		// 当消费失败时重试的间隔时间，默认 10s
		if config.RetryDelay <= time.Duration(0) {
			mqConfig.RetryDelay = time.Second * 10
		} else {
			mqConfig.RetryDelay = config.RetryDelay
		}
		// 当消费失败时重试的最大次数，默认 0，无限重试
		if config.RetryMaxCount <= 0 {
			mqConfig.RetryMaxCount = 0
		} else {
			mqConfig.RetryMaxCount = config.RetryMaxCount
		}
		// 批量消费的条数，默认 200
		if config.BatchSize <= 0 {
			mqConfig.BatchSize = 200
		} else {
			mqConfig.BatchSize = config.BatchSize
		}
		// 批量消费的间隔时间，默认 5s
		if config.BatchInterval <= time.Duration(0) {
			mqConfig.BatchInterval = time.Second * 5
		} else {
			mqConfig.BatchInterval = config.BatchInterval
		}
		return
	}
	err = fmt.Errorf("queue `%s` Not Found", queue)
	return
}

// getPartitionNum 获取消息队列的分区数量
func (mq *RedisMQClient) getPartitionNum(queue string) (partitionNum uint32, err error) {
	if config, ok := mq.config.MQConfig[queue]; ok {
		if config.PartitionNum > 0 {
			partitionNum = config.PartitionNum
			return
		}
		// 默认分区数
		partitionNum = defaultPartitionNum
		return
	}
	err = fmt.Errorf("queue `%s` Not Found", queue)
	return
}

// getGlobalProducerName 获取全局生产者名称
func (mq *RedisMQClient) getGlobalProducerName(globalProducer string) (producerName string) {
	return fmt.Sprintf("producer_%s", globalProducer)
}

// getProducerName 获取生产者名称
func (mq *RedisMQClient) getProducerName(queue string) (producerName string) {
	return fmt.Sprintf("producer_%s", queue)
}

// getConsumerName 获取消费者名称
func (mq *RedisMQClient) getConsumerName(queue string) (consumerName string) {
	return fmt.Sprintf("consumer_%s", queue)
}

// getFullQueueName 获取完整的队列名称
func (mq *RedisMQClient) getFullQueueName(queue string) (fullQueueName string) {
	return fmt.Sprintf("%s_%s", mq.config.Env, queue)
}

// getPartitionQueueName 获取分区队列名称
func (mq *RedisMQClient) getPartitionQueueName(queue string, partition int32) (partitionQueueName string) {
	return fmt.Sprintf("%s@%d", mq.getFullQueueName(queue), partition)
}

// getConsumerGroupName 获取消费者组名称
func (mq *RedisMQClient) getConsumerGroupName(queue string) (group string) {
	return fmt.Sprintf("%s_group_%s", mq.config.ConsumerEnv, queue)
}

// getPartitionGroupName 获取分区消费者组名称
func (mq *RedisMQClient) getPartitionGroupName(queue string, partition int32) (partitionGroupName string) {
	return fmt.Sprintf("%s@%d", mq.getConsumerGroupName(queue), partition)
}

// getPartitionConsumerName 获取分区消费者名称
func (mq *RedisMQClient) getPartitionConsumerName(queue string, partition int32) (partitionConsumerName string) {
	return fmt.Sprintf("%s@%d", mq.getConsumerName(queue), partition)
}
