/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-19 23:42:12
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-03-21 16:13:14
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkkafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/liusuxian/go-toolkit/gtkarr"
	"github.com/liusuxian/go-toolkit/gtkconv"
	"github.com/liusuxian/go-toolkit/gtkjson"
	"github.com/liusuxian/go-toolkit/gtklog"
	"github.com/pkg/errors"
	"hash/fnv"
	"math"
	"time"
)

// ProducerConsumerConfig 生产者和消费者配置
type ProducerConsumerConfig struct {
	Topic         string `json:"topic" dc:"topic 名称"`      // topic 名称
	StartProducer bool   `json:"startProducer" dc:"启动生产者"` // 启动生产者
	StartConsumer bool   `json:"startConsumer" dc:"启动消费者"` // 启动消费者
	StartAll      bool   `json:"startAll" dc:"启动生产者和消费者"`  // 启动生产者和消费者
}

// Config kafka 客户端配置
type Config struct {
	Servers                   string                            `json:"servers" dc:"SSL接入点的IP地址以及端口"`                              // SSL接入点的IP地址以及端口
	Protocol                  string                            `json:"protocol" dc:"SASL用户认证协议"`                                  // SASL用户认证协议
	RetryDelay                time.Duration                     `json:"retryDelay" dc:"当消费失败时重试的间隔时间，默认 10s"`                      // 当消费失败时重试的间隔时间，默认 10s
	RetryMaxCount             int                               `json:"retryMaxCount" dc:"当消费失败时重试的最大次数，默认 0无限重试"`                 // 当消费失败时重试的最大次数，默认 0无限重试
	BatchSize                 int                               `json:"batchSize" dc:"批量消费的条数，默认 200"`                             // 批量消费的条数，默认 200
	BatchInterval             time.Duration                     `json:"batchInterval" dc:"批量消费的间隔时间，默认 5s"`                        // 批量消费的间隔时间，默认 5s
	IsClose                   bool                              `json:"isClose" dc:"是否不启动 Kafka 客户端（适用于本地调试有时候没有kafka环境的情况）"`      // 是否不启动 Kafka 客户端（适用于本地调试有时候没有kafka环境的情况）
	Env                       string                            `json:"env" dc:"当前服务环境，默认 local"`                                  // 当前服务环境，默认 local
	ProducerConsumerConfigMap map[string]ProducerConsumerConfig `json:"producerConsumerConfigMap" dc:"生产者和消费者配置"`                  // 生产者和消费者配置
	ExcludeEnvTopicMap        map[string][]string               `json:"excludeEnvTopicMap" dc:"指定哪些服务环境下对应的哪些 Topic 不发送 Kafka 消息"` // 指定哪些服务环境下对应的哪些 Topic 不发送 Kafka 消息
	LogConfig                 *gtklog.Config                    `json:"logConfig" dc:"日志配置"`                                       // 日志配置
}

// ConfigOption kafka 客户端配置选项
type ConfigOption func(c *Config)

// KafkaClient kafka 客户端结构
type KafkaClient struct {
	producerMap       map[string]*kafka.Producer
	consumerMap       map[string][]*kafka.Consumer
	topicPartitionMap map[string]uint32 // topic 分区数
	config            *Config           // kafka 客户端配置
	logger            *gtklog.Logger    // 日志对象
}

// NewWithOption 创建 kafka 客户端
func NewWithOption(opts ...ConfigOption) (client *KafkaClient, err error) {
	client = &KafkaClient{
		config: &Config{
			ProducerConsumerConfigMap: make(map[string]ProducerConsumerConfig),
			ExcludeEnvTopicMap:        make(map[string][]string),
			LogConfig:                 &gtklog.Config{},
		},
	}
	for _, opt := range opts {
		opt(client.config)
	}
	// SSL接入点的IP地址以及端口
	if client.config.Servers == "" {
		client.config.Servers = "host.docker.internal:9092"
	}
	// SASL用户认证协议
	if client.config.Protocol == "" {
		client.config.Protocol = "PLAINTEXT"
	}
	// 当消费失败时重试的间隔时间，默认 10s
	if client.config.RetryDelay == time.Duration(0) {
		client.config.RetryDelay = time.Second * 10
	}
	// 批量消费的条数，默认 200
	if client.config.BatchSize == 0 {
		client.config.BatchSize = 200
	}
	// 批量消费的间隔时间，默认 5s
	if client.config.BatchInterval == time.Duration(0) {
		client.config.BatchInterval = time.Second * 5
	}
	// 当前服务环境，默认 local
	if client.config.Env == "" {
		client.config.Env = "local"
	}
	if client.logger, err = gtklog.NewWithConfig(client.config.LogConfig); err != nil {
		return
	}
	client.producerMap = make(map[string]*kafka.Producer)
	client.consumerMap = make(map[string][]*kafka.Consumer)
	client.topicPartitionMap = make(map[string]uint32)
	return
}

// NewWithConfig 创建 kafka 客户端
func NewWithConfig(cfg *Config) (client *KafkaClient, err error) {
	if cfg == nil {
		cfg = &Config{
			ProducerConsumerConfigMap: make(map[string]ProducerConsumerConfig),
			ExcludeEnvTopicMap:        make(map[string][]string),
			LogConfig:                 &gtklog.Config{},
		}
	}
	client = &KafkaClient{
		config: cfg,
	}
	// SSL接入点的IP地址以及端口
	if client.config.Servers == "" {
		client.config.Servers = "host.docker.internal:9092"
	}
	// SASL用户认证协议
	if client.config.Protocol == "" {
		client.config.Protocol = "PLAINTEXT"
	}
	// 当消费失败时重试的间隔时间，默认 10s
	if client.config.RetryDelay == time.Duration(0) {
		client.config.RetryDelay = time.Second * 10
	}
	// 批量消费的条数，默认 200
	if client.config.BatchSize == 0 {
		client.config.BatchSize = 200
	}
	// 批量消费的间隔时间，默认 5s
	if client.config.BatchInterval == time.Duration(0) {
		client.config.BatchInterval = time.Second * 5
	}
	// 当前服务环境，默认 local
	if client.config.Env == "" {
		client.config.Env = "local"
	}
	if client.logger, err = gtklog.NewWithConfig(client.config.LogConfig); err != nil {
		return
	}
	client.producerMap = make(map[string]*kafka.Producer)
	client.consumerMap = make(map[string][]*kafka.Consumer)
	client.topicPartitionMap = make(map[string]uint32)
	return
}

// GetConfig 获取 kafka 客户端配置
func (kc *KafkaClient) GetConfig() (config Config) {
	return *(kc.config)
}

// GetTopicPartitionNum 获取 topic 分区数
func (kc *KafkaClient) GetTopicPartitionNum(producerConsumerName string) (partitionNum uint32, err error) {
	if config, ok := kc.config.ProducerConsumerConfigMap[producerConsumerName]; ok {
		topic := fmt.Sprintf("%s_%s", kc.config.Env, config.Topic)
		partitionNum = kc.getTopicPartitionNum(topic)
		return
	}
	err = errors.Errorf("ProducerConsumerName `%s` Not Found", producerConsumerName)
	return
}

// SetTopicPartitionNum 设置 topic 分区数
func (kc *KafkaClient) SetTopicPartitionNum(producerConsumerName string, partitionNum uint32) (err error) {
	if config, ok := kc.config.ProducerConsumerConfigMap[producerConsumerName]; ok {
		topic := fmt.Sprintf("%s_%s", kc.config.Env, config.Topic)
		kc.setTopicPartitionNum(topic, partitionNum)
		return
	}
	err = errors.Errorf("ProducerConsumerName `%s` Not Found", producerConsumerName)
	return
}

// NewProducer 创建生产者
func (kc *KafkaClient) NewProducer(ctx context.Context, producerConsumerName string) (err error) {
	// 获取生产者配置
	var (
		isStart      bool
		producerName string
		topic        string
		originTopic  string
	)
	if isStart, producerName, topic, originTopic, err = kc.getProducerConfig(producerConsumerName); err != nil {
		return
	}
	if !isStart {
		return
	}
	if kc.config.IsClose {
		kc.logger.Infof(ctx, "init kafka producer: %s, topic: %s, originTopic: %s success (isClosed)", producerName, topic, originTopic)
		return
	}
	var kafkaCnf = &kafka.ConfigMap{
		"api.version.request":                   "true",
		"message.max.bytes":                     1048576,       // 最大消息大小
		"compression.type":                      "gzip",        // 消息压缩方式
		"max.in.flight.requests.per.connection": 1,             // 生产者在收到服务器响应之前可以发送多少个消息，设置为1可以保证消息是按照发送的顺序写入服务器，即使发生了重试。
		"sticky.partitioning.linger.ms":         1000,          // 黏性分区策略
		"batch.size":                            16384,         // 批量发送大小
		"linger.ms":                             100,           // 发送延迟时间
		"retries":                               math.MaxInt32, // 发送消息失败后允许重试的次数
		"retry.backoff.ms":                      1000,          // 重试间隔次数
		"acks":                                  "1",           // 回复
	}
	if err = kafkaSetKey(kafkaCnf, map[string]string{
		"bootstrap.servers": kc.config.Servers,
		"security.protocol": kc.config.Protocol,
	}); err != nil {
		return
	}

	var producer *kafka.Producer
	if producer, err = kafka.NewProducer(kafkaCnf); err != nil {
		return
	}
	if producer == nil {
		err = errors.Errorf("NewProducer %s Failed", producerName)
		return
	}

	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					kc.logger.Errorf(ctx, "%s producer failed to write access log entry: %v, topic: %v, partition: %v, offset: %v, key: %s, content: %s",
						producerName, ev.TopicPartition.Error, *ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset, ev.Key, ev.Value)
				} else {
					kc.logger.Debugf(ctx, "%s producer send ok topic: %v, partition: %v, offset: %v, key: %s, content: %s",
						producerName, *ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset, ev.Key, ev.Value)
				}
			}
		}
	}()

	kc.producerMap[producerName] = producer
	kc.logger.Infof(ctx, "init kafka producer: %s, topic: %s, originTopic: %s success", producerName, topic, originTopic)
	return
}

// NewConsumer 创建消费者
func (kc *KafkaClient) NewConsumer(ctx context.Context, producerConsumerName string) (err error) {
	// 获取消费者配置
	var (
		isStart      bool
		consumerName string
		topic        string
		group        string
		partitionNum uint32
	)
	if isStart, consumerName, topic, group, partitionNum, err = kc.getConsumerConfig(producerConsumerName); err != nil {
		return
	}
	if !isStart {
		return
	}
	if kc.config.IsClose {
		kc.logger.Infof(ctx, "init kafka consumer: %s, topic: %s, group: %s, partitionNum: %v success (isClosed)", consumerName, topic, group, partitionNum)
		return
	}
	var kafkaCnf = &kafka.ConfigMap{
		"api.version.request":       "true",
		"auto.offset.reset":         "earliest", // 消费偏移量
		"heartbeat.interval.ms":     3000,       // 心跳间隔时间
		"session.timeout.ms":        30000,      // 会话超时时间
		"max.poll.interval.ms":      60000,      // 最大拉取间隔时间
		"fetch.max.bytes":           52428800,   // 一次fetch请求，从一个broker中取得的records最大大小
		"max.partition.fetch.bytes": 104857600,  // 服务器从每个分区里返回给消费者的最大字节数
		"enable.auto.commit":        false,      // 关闭自动提交
	}
	if err = kafkaSetKey(kafkaCnf, map[string]string{
		"bootstrap.servers": kc.config.Servers,
		"security.protocol": kc.config.Protocol,
		"group.id":          group,
	}); err != nil {
		return
	}

	consumerList := make([]*kafka.Consumer, 0, int(partitionNum))
	for i := 0; i < int(partitionNum); i++ {
		consumer, cErr := kafka.NewConsumer(kafkaCnf)
		if cErr != nil {
			return cErr
		}
		if consumer == nil {
			return errors.Errorf("NewConsumer %s Failed", consumerName)
		}
		consumerList = append(consumerList, consumer)
	}

	kc.consumerMap[consumerName] = consumerList
	kc.logger.Infof(ctx, "init kafka consumer: %s, topic: %s, group: %s, partitionNum: %v success", consumerName, topic, group, partitionNum)
	return
}

// SendJsonData 发送 Json 格式的数据
func (kc *KafkaClient) SendJsonData(ctx context.Context, producerConsumerName string, data any, key ...string) (err error) {
	var dataMap map[string]any
	if dataMap, err = gtkconv.ToStringMapE(data); err != nil {
		return
	}
	delete(dataMap, "created_at")
	delete(dataMap, "updated_at")
	delete(dataMap, "deleted_at")
	// 处理数据
	var buf []byte
	if buf, err = json.Marshal(dataMap); err != nil {
		return
	}
	return kc.sendData(ctx, producerConsumerName, buf, key...)
}

// SubscribeTopics 订阅数据
func (kc *KafkaClient) SubscribeTopics(ctx context.Context, producerConsumerName string, fn func(message *kafka.Message) error) (err error) {
	// 获取消费者配置
	var (
		isStart      bool
		consumerName string
		topic        string
	)
	if isStart, consumerName, topic, _, _, err = kc.getConsumerConfig(producerConsumerName); err != nil {
		return
	}
	if !isStart {
		return
	}
	topics := []string{topic}
	if kc.config.IsClose {
		kc.logger.Infof(ctx, "subscribeTopics consumer: %s, topics: %v (isClosed)", consumerName, topics)
		return
	}
	// 订阅topics
	var consumerList []*kafka.Consumer
	var ok bool
	if consumerList, ok = kc.consumerMap[consumerName]; !ok {
		err = errors.Errorf("consumer: %s, topics: %v not found", consumerName, topics)
		return
	}
	for _, v := range consumerList {
		if err = v.SubscribeTopics(topics, nil); err != nil {
			return
		}

		// 接收消息
		go func(consumer *kafka.Consumer) {
			// 添加对 panic 的处理
			defer func() {
				if r := recover(); r != nil {
					kc.logger.Errorf(ctx, "%s %v consumer panic: %+v", consumerName, consumer, r)
				}
			}()

			for {
				select {
				case <-ctx.Done():
					return
				default:
					msg, pErr := consumer.ReadMessage(100 * time.Millisecond)
					if pErr == nil {
						// 处理数据
						kc.handelData(ctx, consumerName, consumer, msg, fn)
					} else {
						if kafkaErr, ok := pErr.(kafka.Error); ok {
							if kafkaErr.Code() != kafka.ErrTimedOut {
								kc.logger.Errorf(ctx, "%s %v consumer error: %+v, msg: %+v", consumerName, consumer, pErr, msg)
							}
						} else {
							kc.logger.Errorf(ctx, "%s %v consumer error: %+v, msg: %+v", consumerName, consumer, pErr, msg)
						}
					}
				}
			}
		}(v)
	}

	return
}

// BatchSubscribeTopics 批量订阅数据
func (kc *KafkaClient) BatchSubscribeTopics(ctx context.Context, producerConsumerName string, fn func(messages []*kafka.Message) error) (err error) {
	// 获取消费者配置
	var (
		isStart      bool
		consumerName string
		topic        string
	)
	if isStart, consumerName, topic, _, _, err = kc.getConsumerConfig(producerConsumerName); err != nil {
		return
	}
	if !isStart {
		return
	}
	topics := []string{topic}
	if kc.config.IsClose {
		kc.logger.Infof(ctx, "batchSubscribeTopics consumer: %s, topics: %v (isClosed)", consumerName, topics)
		return
	}
	// 订阅topics
	var consumerList []*kafka.Consumer
	var ok bool
	if consumerList, ok = kc.consumerMap[consumerName]; !ok {
		err = errors.Errorf("consumer: %s, topics: %v not found", consumerName, topics)
		return
	}
	for _, v := range consumerList {
		if err = v.SubscribeTopics(topics, nil); err != nil {
			return
		}

		// 接收消息
		go func(consumer *kafka.Consumer) {
			// 添加对 panic 的处理
			defer func() {
				if r := recover(); r != nil {
					kc.logger.Errorf(ctx, "%s %v consumer panic: %+v", consumerName, consumer, r)
				}
			}()

			// 批量数据
			msgList := make([]*kafka.Message, 0, kc.config.BatchSize)
			// 定时器
			ticker := time.NewTicker(kc.config.BatchInterval)
			defer ticker.Stop()

			// 读取消息
			readMsg := make(chan *kafka.Message)
			defer close(readMsg)
			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					default:
						msg, pErr := consumer.ReadMessage(100 * time.Millisecond)
						if pErr == nil {
							readMsg <- msg
						} else {
							if kafkaErr, ok := pErr.(kafka.Error); ok {
								if kafkaErr.Code() != kafka.ErrTimedOut {
									kc.logger.Errorf(ctx, "%s %v consumer error: %+v, msg: %+v", consumerName, consumer, pErr, msg)
								}
							} else {
								kc.logger.Errorf(ctx, "%s %v consumer error: %+v, msg: %+v", consumerName, consumer, pErr, msg)
							}
						}
					}
				}
			}()

			for {
				select {
				case <-ticker.C:
					if len(msgList) > 0 {
						// 批量处理数据
						kc.handelBatchData(ctx, consumerName, consumer, msgList, fn)
						// 重新创建一个新的 msgList
						msgList = make([]*kafka.Message, 0, kc.config.BatchSize)
					}
				case msg, ok := <-readMsg:
					if ok {
						msgList = append(msgList, msg)
						if len(msgList) == kc.config.BatchSize {
							// 批量处理数据
							kc.handelBatchData(ctx, consumerName, consumer, msgList, fn)
							// 重新创建一个新的 msgList
							msgList = make([]*kafka.Message, 0, kc.config.BatchSize)
						}
					}
				case <-ctx.Done():
					if len(msgList) > 0 {
						// 批量处理数据
						kc.handelBatchData(ctx, consumerName, consumer, msgList, fn)
					}
					return
				}
			}
		}(v)
	}

	return
}

// handelData 处理数据
func (kc *KafkaClient) handelData(ctx context.Context, consumerName string, consumer *kafka.Consumer, msg *kafka.Message, fn func(msg *kafka.Message) error) {
	// 执行处理函数
	if fErr := fn(msg); fErr != nil {
		kc.logger.Errorf(ctx, "%s %v consumer error: %+v, topic: %v, partition: %v, offset: %v, content: %s", consumerName, consumer,
			fErr, *msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset, string(msg.Value))
		// 暂停该分区的消费
		consumer.Pause([]kafka.TopicPartition{
			{
				Topic:     msg.TopicPartition.Topic,
				Partition: msg.TopicPartition.Partition,
				Offset:    msg.TopicPartition.Offset,
			},
		})
		// 重试处理函数
		var (
			retryMaxCount = kc.config.RetryMaxCount
			count         = 0
		)
		for retryMaxCount == 0 || (retryMaxCount > 0 && count < retryMaxCount) {
			count++
			consumer.Poll(0) // 用以给kafka发送心跳
			time.Sleep(kc.config.RetryDelay)
			if sfErr := fn(msg); sfErr != nil {
				kc.logger.Errorf(ctx, "%s %v consumer error: %+v, count: %d, topic: %v, partition: %v, offset: %v, content: %s", consumerName, consumer,
					sfErr, count, *msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset, string(msg.Value))
				continue
			}
			// 恢复该分区的消费
			consumer.Resume([]kafka.TopicPartition{
				{
					Topic:     msg.TopicPartition.Topic,
					Partition: msg.TopicPartition.Partition,
					Offset:    msg.TopicPartition.Offset,
				},
			})
			break
		}
	}
	// 提交
	consumer.CommitMessage(msg)
}

// handelBatchData 批量处理数据
func (kc *KafkaClient) handelBatchData(ctx context.Context, consumerName string, consumer *kafka.Consumer, msgList []*kafka.Message, fn func(messages []*kafka.Message) error) {
	// 判断 msgList 是否为空
	msgSize := len(msgList)
	if msgSize == 0 {
		return
	}
	// 获取最后一条消息
	endMsg := msgList[msgSize-1]
	// 执行处理函数
	if fErr := fn(msgList); fErr != nil {
		kc.logger.Errorf(ctx, "%s %v consumer error: %+v, topic: %v, partition: %v, offset: %v", consumerName, consumer,
			fErr, *endMsg.TopicPartition.Topic, endMsg.TopicPartition.Partition, endMsg.TopicPartition.Offset)
		// 暂停该分区的消费
		consumer.Pause([]kafka.TopicPartition{
			{
				Topic:     endMsg.TopicPartition.Topic,
				Partition: endMsg.TopicPartition.Partition,
				Offset:    endMsg.TopicPartition.Offset,
			},
		})
		// 重试处理函数
		var (
			retryMaxCount = kc.config.RetryMaxCount
			count         = 0
		)
		for retryMaxCount == 0 || (retryMaxCount > 0 && count > retryMaxCount) {
			count++
			consumer.Poll(0) // 用以给kafka发送心跳
			time.Sleep(kc.config.RetryDelay)
			if sfErr := fn(msgList); sfErr != nil {
				kc.logger.Errorf(ctx, "%s %v consumer error: %+v, count: %d, topic: %v, partition: %v, offset: %v", consumerName, consumer,
					sfErr, count, *endMsg.TopicPartition.Topic, endMsg.TopicPartition.Partition, endMsg.TopicPartition.Offset)
				continue
			}
			// 恢复该分区的消费
			consumer.Resume([]kafka.TopicPartition{
				{
					Topic:     endMsg.TopicPartition.Topic,
					Partition: endMsg.TopicPartition.Partition,
					Offset:    endMsg.TopicPartition.Offset,
				},
			})
			break
		}
	}
	// 提交
	consumer.CommitMessage(endMsg)
}

// sendData 发送数据
func (kc *KafkaClient) sendData(ctx context.Context, producerConsumerName string, data []byte, key ...string) (err error) {
	// 获取生产者配置
	var (
		isStart      bool
		producerName string
		topic        string
		originTopic  string
	)
	if isStart, producerName, topic, originTopic, err = kc.getProducerConfig(producerConsumerName); err != nil {
		return
	}
	if !isStart {
		return
	}
	var dataKey string
	if len(key) > 0 {
		dataKey = key[0]
	}
	// 组装消息
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kc.calcPartition(topic, dataKey)},
		Value:          data,
	}
	if dataKey != "" {
		msg.Key = []byte(dataKey)
	}
	if kc.config.IsClose {
		kc.logger.Debugf(ctx, "%s producer sendData(isClosed) msg: %s, data: %s", producerName, gtkjson.MustString(msg), string(data))
		return
	}
	// 环境检测
	if list, ok := kc.config.ExcludeEnvTopicMap[kc.config.Env]; ok {
		if gtkarr.ContainsStr(list, originTopic) {
			kc.logger.Debugf(ctx, "%s producer sendData msg: %s, data: %s", producerName, gtkjson.MustString(msg), string(data))
			return
		}
	}
	// 发送消息
	var producer *kafka.Producer
	var ok bool
	if producer, ok = kc.producerMap[producerName]; !ok {
		err = errors.Errorf("sendData producer: %s, msg: %s, data: %s not found", producerName, gtkjson.MustString(msg), string(data))
		return
	}
	if producer == nil {
		err = errors.Errorf("sendData producer: %s, msg: %s, data: %s producer nil", producerName, gtkjson.MustString(msg), string(data))
		return
	}
	if err = producer.Produce(msg, nil); err != nil {
		return
	}
	producer.Flush(10 * 1000)
	return
}

// getProducerConfig 获取生产者配置
func (kc *KafkaClient) getProducerConfig(producerConsumerName string) (isStart bool, producerName, topic, originTopic string, err error) {
	if config, ok := kc.config.ProducerConsumerConfigMap[producerConsumerName]; ok {
		isStart = (config.StartAll || config.StartProducer)
		producerName = fmt.Sprintf("producer_%s", producerConsumerName)
		topic = fmt.Sprintf("%s_%s", kc.config.Env, config.Topic)
		originTopic = config.Topic
		return
	}
	err = errors.Errorf("ProducerConsumerName `%s` Not Found", producerConsumerName)
	return
}

// getConsumerConfig 获取消费者配置
func (kc *KafkaClient) getConsumerConfig(producerConsumerName string) (isStart bool, consumerName, topic, group string, partitionNum uint32, err error) {
	if config, ok := kc.config.ProducerConsumerConfigMap[producerConsumerName]; ok {
		isStart = (config.StartAll || config.StartConsumer)
		consumerName = fmt.Sprintf("consumer_%s", producerConsumerName)
		topic = fmt.Sprintf("%s_%s", kc.config.Env, config.Topic)
		group = fmt.Sprintf("%s_group_%s", kc.config.Env, config.Topic)
		partitionNum = kc.getTopicPartitionNum(topic)
		return
	}
	err = errors.Errorf("ProducerConsumerName `%s` Not Found", producerConsumerName)
	return
}

// getTopicPartitionNum 获取`topic`分区数
func (kc *KafkaClient) getTopicPartitionNum(topic string) (partitionNum uint32) {
	if value, ok := kc.topicPartitionMap[topic]; ok {
		return value
	}
	// 默认都是 12 个分区
	return 12
}

// setTopicPartitionNum 设置`topic`分区数
func (kc *KafkaClient) setTopicPartitionNum(topic string, partitionNum uint32) {
	kc.topicPartitionMap[topic] = partitionNum
}

// calcPartition 计算分区号
func (kc *KafkaClient) calcPartition(topic, key string) (partition int32) {
	if key == "" {
		partition = kafka.PartitionAny
		return
	}
	hash := fnv.New32a()
	hash.Write([]byte(key))
	partition = int32(hash.Sum32() % kc.getTopicPartitionNum(topic))
	return
}

// kafkaSetKey kafka 设置连接配置
func kafkaSetKey(kafkaCnf *kafka.ConfigMap, kafkaCnfMap map[string]string) (err error) {
	for k, v := range kafkaCnfMap {
		if err = kafkaCnf.SetKey(k, v); err != nil {
			return
		}
	}

	return
}
