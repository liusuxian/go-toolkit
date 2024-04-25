/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-04-23 00:35:41
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-04-25 17:27:44
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
	PartitionNum uint32                    `json:"partitionNum" dc:"消息队列分区数量，默认 12 个分区"`                        // 消息队列分区数量，默认 12 个分区
	Mode         ProducerConsumerStartMode `json:"mode" dc:"启动模式 0:不启动生产者或消费者 1:仅启动生产者 2:仅启动消费者 3:同时启动生产者和消费者"` // 启动模式 0:不启动生产者或消费者 1:仅启动生产者 2:仅启动消费者 3:同时启动生产者和消费者
}

// ProducerMessage 生产者消息
type ProducerMessage struct {
	Key       string    `json:"key" dc:"键"`              // 键
	Data      any       `json:"data" dc:"数据"`            // 数据
	Timestamp time.Time `json:"timestamp" dc:"发送消息的时间戳"` // 发送消息的时间戳
	dataBytes []byte    // 数据字节数组
}

// MQPartition 消息队列分区
type MQPartition struct {
	Queue         string `json:"queue" dc:"队列名称"`         // 队列名称
	PartitionName string `json:"partitionName" dc:"分区名称"` // 分区名称
	Partition     int32  `json:"partition" dc:"分区号"`      // 分区号
	Offset        string `json:"offset" dc:"偏移量"`         // 偏移量
}

// MQMessage 消息队列消息
type MQMessage struct {
	MQPartition MQPartition `json:"mqPartition" dc:"消息队列分区"` // 消息队列分区
	Key         []byte      `json:"key" dc:"键"`              // 键
	Value       []byte      `json:"value" dc:"值"`            // 值
	Timestamp   time.Time   `json:"timestamp" dc:"发送消息的时间戳"` // 发送消息的时间戳
	ExpireTime  time.Time   `json:"expireTime" dc:"消息过期时间"`  // 消息过期时间
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
	Subscribe(ctx context.Context, queue string, fn func(message *MQMessage) error) (err error)
	// BatchSubscribe 批量订阅数据
	BatchSubscribe(ctx context.Context, queue string, fn func(messages []*MQMessage) error) (err error)
	// GetExpiredMessages 获取过期消息，每个分区每次最多返回 100 条
	//
	//	isDelete: 是否删除过期消息
	GetExpiredMessages(ctx context.Context, queue string, isDelete bool) (messages map[int32][]*MQMessage, err error)
	// ResetConsumerOffset 重置消费起点，所有分区（请谨慎使用）
	//
	//	offset: 0-0 重置为最早位置
	//	offset: $ 重置为最新位置
	ResetConsumerOffset(ctx context.Context, queue string, offset string) (err error)
	// ResetConsumerOffsetByPartition 重置消费起点，指定分区（请谨慎使用）
	//
	//	offset: 0-0 重置为最早位置
	//	offset: $ 重置为最新位置
	//	offset: <ID> 重置为指定位置
	ResetConsumerOffsetByPartition(ctx context.Context, queue string, partition int32, offset string) (err error)
	// DelGroup 删除消费者组（请谨慎使用）
	DelGroup(ctx context.Context, queue string) (err error)
	// DelQueue 删除队列（请谨慎使用）
	DelQueue(ctx context.Context, queue string) (err error)
	// 关闭客户端
	Close() (err error)
}
