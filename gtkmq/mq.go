/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-04-23 00:35:41
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-22 10:28:39
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkmq

import (
	"context"
	"time"
)

// ProducerConsumerStartMode 定义生产者和消费者的启动模式
type ProducerConsumerStartMode int

const (
	ModeNone     ProducerConsumerStartMode = iota // 不启动生产者或消费者
	ModeProducer                                  // 仅启动生产者
	ModeConsumer                                  // 仅启动消费者
	ModeBoth                                      // 同时启动生产者和消费者
)

// MQConfig 消息队列配置
type MQConfig struct {
	// 消息队列分区数量，默认 12 个分区
	PartitionNum uint32 `json:"partitionNum,omitempty"`
	// 启动模式 0:不启动生产者或消费者 1:仅启动生产者 2:仅启动消费者 3:同时启动生产者和消费者
	Mode ProducerConsumerStartMode `json:"mode"`
	// 指定消费者组名称列表。如果未指定，将使用默认格式："$consumerEnv_group_$topic"，其中`$consumerEnv_group_`是系统根据当前环境自动添加的前缀
	// 可以配置多个消费者组名称，系统会自动在每个名称前添加"$consumerEnv_group_"前缀
	Groups []string `json:"groups,omitempty"`
	// 当消费失败时重试的间隔时间，默认 10s
	RetryDelay time.Duration `json:"retryDelay,omitempty"`
	// 当消费失败时重试的最大次数，默认 0，无限重试
	RetryMaxCount int `json:"retryMaxCount,omitempty"`
	// 批量消费的条数，默认 200
	BatchSize int `json:"batchSize,omitempty"`
	// 批量消费的间隔时间，默认 5s
	BatchInterval time.Duration `json:"batchInterval,omitempty"`
	// 是否开启延迟队列
	IsDelayQueue bool `json:"isDelayQueue,omitempty"`
	// 延迟队列检查间隔，默认 10s
	DelayQueueCheckInterval time.Duration `json:"delayQueueCheckInterval,omitempty"`
	// 延迟队列批处理大小，默认 100
	DelayQueueBatchSize int `json:"delayQueueBatchSize,omitempty"`
}

// ProducerMessage 生产者消息
type ProducerMessage struct {
	Key       string    `json:"key,omitempty"` // 键
	Data      any       `json:"data"`          // 数据
	DelayTime time.Time `json:"delay_time"`    // 延迟投递的截止时间点
	dataBytes []byte    // 数据字节数组
}

// MQPartition 消息队列分区
type MQPartition struct {
	Queue         string `json:"queue"`          // 队列名称
	PartitionName string `json:"partition_name"` // 分区名称
	Partition     int32  `json:"partition"`      // 分区号
	Offset        string `json:"offset"`         // 偏移量
}

// MQMessage 消息队列消息
type MQMessage struct {
	MQPartition MQPartition `json:"mq_partition"` // 消息队列分区
	Key         []byte      `json:"key"`          // 键
	Value       []byte      `json:"value"`        // 值
	Timestamp   time.Time   `json:"timestamp"`    // 发送消息的时间戳
	ExpireTime  time.Time   `json:"expire_time"`  // 消息过期时间
}

// MQClient 消息队列客户端接口
type MQClient interface {
	// NewProducer 创建生产者
	NewProducer(ctx context.Context, queue string) (err error)
	// NewConsumer 创建消费者
	NewConsumer(ctx context.Context, queue string) (err error)
	// SendMessage 发送消息
	SendMessage(ctx context.Context, queue string, producerMessage *ProducerMessage) (err error)
	// Subscribe 订阅数据
	Subscribe(ctx context.Context, queue string, fn func(message *MQMessage) error, group ...string) (err error)
	// BatchSubscribe 批量订阅数据
	BatchSubscribe(ctx context.Context, queue string, fn func(messages []*MQMessage) error, group ...string) (err error)
	// GetExpiredMessages 获取过期消息，每个分区每次最多返回 100 条
	//
	//	isDelete: 是否删除过期消息
	GetExpiredMessages(ctx context.Context, queue string, isDelete bool) (messages map[int32][]*MQMessage, err error)
	// ResetConsumerOffset 重置消费起点，所有分区（请谨慎使用）
	//
	//	offset: 0-0 重置为最早位置
	//	offset: $ 重置为最新位置
	ResetConsumerOffset(ctx context.Context, queue string, offset string, group ...string) (err error)
	// ResetConsumerOffsetByPartition 重置消费起点，指定分区（请谨慎使用）
	//
	//	offset: 0-0 重置为最早位置
	//	offset: $ 重置为最新位置
	//	offset: <ID> 重置为指定位置
	ResetConsumerOffsetByPartition(ctx context.Context, queue string, partition int32, offset string, group ...string) (err error)
	// DelGroup 删除消费者组（请谨慎使用）
	DelGroup(ctx context.Context, queue string, group ...string) (err error)
	// DelQueue 删除队列（请谨慎使用）
	DelQueue(ctx context.Context, queue string) (err error)
	// 关闭客户端
	Close() (err error)
}
