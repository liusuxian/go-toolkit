/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-19 23:42:12
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-01-16 20:02:02
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkkafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/liusuxian/go-toolkit/gtkarr"
	"github.com/liusuxian/go-toolkit/gtkconv"
	"github.com/liusuxian/go-toolkit/gtkjson"
	"github.com/liusuxian/go-toolkit/gtklog"
	"github.com/liusuxian/go-toolkit/gtkstr"
	"github.com/pkg/errors"
	"hash/fnv"
	"math"
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

// TopicConfig topic 配置
type TopicConfig struct {
	// topic 分区数量，默认 12 个分区
	PartitionNum uint32 `json:"partitionNum"`
	// 启动模式 0:不启动生产者或消费者 1:仅启动生产者 2:仅启动消费者 3:同时启动生产者和消费者
	Mode ProducerConsumerStartMode `json:"mode"`
	// 指定消费者组名称列表。如果未指定，将使用默认格式："$consumerEnv_group_$topic"，其中`$consumerEnv_group_`是系统根据当前环境自动添加的前缀
	// 可以配置多个消费者组名称，系统会自动在每个名称前添加"$consumerEnv_group_"前缀
	Groups []string `json:"groups"`
}

// ProducerMessage 生产者消息
type ProducerMessage struct {
	Key       string    `json:"key" dc:"键"`              // 键
	Data      any       `json:"data" dc:"数据"`            // 数据
	Timestamp time.Time `json:"timestamp" dc:"发送消息的时间戳"` // 发送消息的时间戳
	dataBytes []byte    // 数据字节数组
}

// Config kafka 客户端配置
type Config struct {
	BootstrapServers           string                 `json:"bootstrapServers"`           // Kafka 服务器的地址列表，格式为 host1:port1,host2:port2
	SecurityProtocol           string                 `json:"securityProtocol"`           // Kafka 通信的安全协议，如 PLAINTEXT、SSL、SASL_PLAINTEXT、SASL_SSL
	SaslMechanism              string                 `json:"saslMechanism"`              // SASL 认证机制，如 GSSAPI、PLAIN、OAUTHBEARER、SCRAM-SHA-256、SCRAM-SHA-512
	SaslUsername               string                 `json:"saslUsername"`               // SASL 认证的用户名
	SaslPassword               string                 `json:"saslPassword"`               // SASL 认证的密码
	StickyPartitioningLingerMs int                    `json:"stickyPartitioningLingerMs"` // 黏性分区策略的延迟时间，此设置允许生产者在指定时间内将消息发送到同一个分区，以增加消息批次的大小，提高压缩效率和吞吐量。设置为 0 时，生产者不会等待，消息会立即发送。默认 100ms
	BatchSize                  int                    `json:"batchSize"`                  // 批量发送大小，默认 10485760 字节
	MessageMaxBytes            int                    `json:"messageMaxBytes"`            // 最大消息大小，默认 16384 字节
	Retries                    int                    `json:"retries"`                    // 发送消息失败后允许重试的次数，默认 2147483647
	RetryBackoffMs             int                    `json:"retryBackoffMs"`             // 发送消息失败后，下一次重试发送前的等待时间，默认 100ms
	LingerMs                   int                    `json:"lingerMs"`                   // 发送延迟时间，默认 100ms
	QueueBufferingMaxKbytes    int                    `json:"queueBufferingMaxKbytes"`    // Producer 攒批发送中，默认 1048576kb
	WaitTimeout                time.Duration          `json:"waitTimeout"`                // 指定等待消息的最大时间，默认 -1，表示无限期等待消息，直到有消息到达
	RetryDelay                 time.Duration          `json:"retryDelay"`                 // 当消费失败时重试的间隔时间，默认 10s
	RetryMaxCount              int                    `json:"retryMaxCount"`              // 当消费失败时重试的最大次数，默认 0，无限重试
	OffsetReset                string                 `json:"offsetReset"`                // 重置消费者偏移量的策略，可选值: earliest 最早位置，latest 最新位置，none 找不到之前的偏移量，消费者将抛出一个异常，停止工作，默认 earliest
	BatchConsumeSize           int                    `json:"batchConsumeSize"`           // 批量消费的条数，默认 200
	BatchConsumeInterval       time.Duration          `json:"batchConsumeInterval"`       // 批量消费的间隔时间，默认 5s
	IsClose                    bool                   `json:"isClose"`                    // 是否不启动 Kafka 客户端（适用于本地调试有时候没有kafka环境的情况）
	Env                        string                 `json:"env"`                        // topic 服务环境，默认 local
	ConsumerEnv                string                 `json:"consumerEnv"`                // 消费者服务环境，默认和 topic 服务环境一致
	GlobalProducer             string                 `json:"globalProducer"`             // 全局生产者名称，配置此项时，客户端将使用全局生产者，不再创建新的生产者，默认为空
	TopicConfig                map[string]TopicConfig `json:"topicConfig"`                // topic 配置，key 为 topic 名称
	ExcludeTopics              []string               `json:"excludeTopics"`              // 指定哪些 topic 不发送 Kafka 消息
	LogConfig                  *gtklog.Config         `json:"logConfig"`                  // 日志配置
}

// ConfigOption kafka 客户端配置选项
type ConfigOption func(c *Config)

// KafkaClient kafka 客户端
type KafkaClient struct {
	producerMap map[string]*kafka.Producer
	consumerMap map[string][]*kafka.Consumer
	config      *Config        // kafka 客户端配置
	logger      *gtklog.Logger // 日志对象
}

const (
	defaultPartitionNum uint32 = 12 // 默认分区数
)

// NewWithOption 创建 kafka 客户端
func NewWithOption(opts ...ConfigOption) (client *KafkaClient, err error) {
	client = &KafkaClient{
		config: &Config{
			TopicConfig:   make(map[string]TopicConfig),
			ExcludeTopics: make([]string, 0),
			LogConfig:     &gtklog.Config{},
		},
	}
	for _, opt := range opts {
		opt(client.config)
	}
	// SSL接入点的IP地址以及端口
	if client.config.BootstrapServers == "" {
		client.config.BootstrapServers = "127.0.0.1:9092"
	}
	// SASL用户认证协议
	if client.config.SecurityProtocol == "" {
		client.config.SecurityProtocol = "PLAINTEXT"
	}
	// 黏性分区策略的延迟时间，此设置允许生产者在指定时间内将消息发送到同一个分区，以增加消息批次的大小，提高压缩效率和吞吐量。设置为 0 时，生产者不会等待，消息会立即发送。默认 100ms
	if client.config.StickyPartitioningLingerMs <= 0 {
		client.config.StickyPartitioningLingerMs = 100
	}
	// 批量发送大小，默认 10485760 字节
	if client.config.BatchSize <= 0 {
		client.config.BatchSize = 10485760
	}
	// 最大消息大小，默认 16384 字节
	if client.config.MessageMaxBytes <= 0 {
		client.config.MessageMaxBytes = 16384
	}
	// 发送消息失败后允许重试的次数，默认 2147483647
	if client.config.Retries <= 0 {
		client.config.Retries = math.MaxInt32
	}
	// 发送消息失败后，下一次重试发送前的等待时间，默认 100ms
	if client.config.RetryBackoffMs <= 0 {
		client.config.RetryBackoffMs = 100
	}
	// 发送延迟时间，默认 100ms
	if client.config.LingerMs <= 0 {
		client.config.LingerMs = 100
	}
	// Producer 攒批发送中，默认 1048576kb
	if client.config.QueueBufferingMaxKbytes <= 0 {
		client.config.QueueBufferingMaxKbytes = 1048576
	}
	// 指定等待消息的最大时间，默认 -1，表示无限期等待消息，直到有消息到达
	if client.config.WaitTimeout <= time.Duration(0) {
		client.config.WaitTimeout = time.Duration(-1)
	}
	// 当消费失败时重试的间隔时间，默认 10s
	if client.config.RetryDelay == time.Duration(0) {
		client.config.RetryDelay = time.Second * 10
	}
	// 重置消费者偏移量的策略，可选值: earliest 最早位置，latest 最新位置，none 找不到之前的偏移量，消费者将抛出一个异常，停止工作，默认 earliest
	if client.config.OffsetReset == "" {
		client.config.OffsetReset = "earliest"
	}
	// 批量消费的条数，默认 200
	if client.config.BatchConsumeSize == 0 {
		client.config.BatchConsumeSize = 200
	}
	// 批量消费的间隔时间，默认 5s
	if client.config.BatchConsumeInterval == time.Duration(0) {
		client.config.BatchConsumeInterval = time.Second * 5
	}
	// topic 服务环境，默认 local
	if client.config.Env == "" {
		client.config.Env = "local"
	}
	// 消费者服务环境，默认和 topic 服务环境一致
	if client.config.ConsumerEnv == "" {
		client.config.ConsumerEnv = client.config.Env
	}
	if client.logger, err = gtklog.NewWithConfig(client.config.LogConfig); err != nil {
		return
	}
	client.producerMap = make(map[string]*kafka.Producer)
	client.consumerMap = make(map[string][]*kafka.Consumer)
	return
}

// NewWithConfig 创建 kafka 客户端
func NewWithConfig(cfg *Config) (client *KafkaClient, err error) {
	if cfg == nil {
		cfg = &Config{
			TopicConfig:   make(map[string]TopicConfig),
			ExcludeTopics: make([]string, 0),
			LogConfig:     &gtklog.Config{},
		}
	}
	client = &KafkaClient{
		config: cfg,
	}
	// SSL接入点的IP地址以及端口
	if client.config.BootstrapServers == "" {
		client.config.BootstrapServers = "127.0.0.1:9092"
	}
	// SASL用户认证协议
	if client.config.SecurityProtocol == "" {
		client.config.SecurityProtocol = "PLAINTEXT"
	}
	// 黏性分区策略的延迟时间，此设置允许生产者在指定时间内将消息发送到同一个分区，以增加消息批次的大小，提高压缩效率和吞吐量。设置为 0 时，生产者不会等待，消息会立即发送。默认 100ms
	if client.config.StickyPartitioningLingerMs <= 0 {
		client.config.StickyPartitioningLingerMs = 100
	}
	// 批量发送大小，默认 10485760 字节
	if client.config.BatchSize <= 0 {
		client.config.BatchSize = 10485760
	}
	// 最大消息大小，默认 16384 字节
	if client.config.MessageMaxBytes <= 0 {
		client.config.MessageMaxBytes = 16384
	}
	// 发送消息失败后允许重试的次数，默认 2147483647
	if client.config.Retries <= 0 {
		client.config.Retries = math.MaxInt32
	}
	// 发送消息失败后，下一次重试发送前的等待时间，默认 100ms
	if client.config.RetryBackoffMs <= 0 {
		client.config.RetryBackoffMs = 100
	}
	// 发送延迟时间，默认 100ms
	if client.config.LingerMs <= 0 {
		client.config.LingerMs = 100
	}
	// Producer 攒批发送中，默认 1048576kb
	if client.config.QueueBufferingMaxKbytes <= 0 {
		client.config.QueueBufferingMaxKbytes = 1048576
	}
	// 指定等待消息的最大时间，默认 -1，表示无限期等待消息，直到有消息到达
	if client.config.WaitTimeout <= time.Duration(0) {
		client.config.WaitTimeout = time.Duration(-1)
	}
	// 当消费失败时重试的间隔时间，默认 10s
	if client.config.RetryDelay == time.Duration(0) {
		client.config.RetryDelay = time.Second * 10
	}
	// 重置消费者偏移量的策略，可选值: earliest 最早位置，latest 最新位置，none 找不到之前的偏移量，消费者将抛出一个异常，停止工作，默认 earliest
	if client.config.OffsetReset == "" {
		client.config.OffsetReset = "earliest"
	}
	// 批量消费的条数，默认 200
	if client.config.BatchConsumeSize == 0 {
		client.config.BatchConsumeSize = 200
	}
	// 批量消费的间隔时间，默认 5s
	if client.config.BatchConsumeInterval == time.Duration(0) {
		client.config.BatchConsumeInterval = time.Second * 5
	}
	// topic 服务环境，默认 local
	if client.config.Env == "" {
		client.config.Env = "local"
	}
	// 消费者服务环境，默认和 topic 服务环境一致
	if client.config.ConsumerEnv == "" {
		client.config.ConsumerEnv = client.config.Env
	}
	if client.logger, err = gtklog.NewWithConfig(client.config.LogConfig); err != nil {
		return
	}
	client.producerMap = make(map[string]*kafka.Producer)
	client.consumerMap = make(map[string][]*kafka.Consumer)
	return
}

// PrintClientConfig 打印消息队列客户端配置
func (kc *KafkaClient) PrintClientConfig(ctx context.Context) {
	kc.logger.Debugf(ctx, "client config: %s\n", gtkjson.MustString(kc.config))
}

// NewProducer 创建生产者
func (kc *KafkaClient) NewProducer(ctx context.Context, topic string) (err error) {
	// 获取生产者配置
	var (
		isStart      bool
		partitionNum uint32
	)
	if isStart, partitionNum, err = kc.getProducerConfig(topic); err != nil {
		return
	}
	if !isStart {
		return
	}
	// 获取生产者名称和完整的 topic 名称
	var (
		producerName  = kc.getProducerName(topic)
		fullTopicName = kc.getFullTopicName(topic)
	)
	// 判断是否配置了全局生产者名称
	globalProducerName := gtkstr.TrimAll(kc.config.GlobalProducer)
	if globalProducerName != "" {
		producerName = kc.getGlobalProducerName(globalProducerName)
	}
	if kc.config.IsClose {
		kc.logger.Infof(ctx, "new producer: %s, topic: %s, partitionNum: %d success (isClosed)", producerName, fullTopicName, partitionNum)
		return
	}
	// 判断全局生产者是否已存在
	if globalProducerName != "" {
		if _, ok := kc.producerMap[producerName]; ok {
			kc.logger.Infof(ctx, "new producer: %s, topic: %s, partitionNum: %d success", producerName, fullTopicName, partitionNum)
			return
		}
	} else {
		if _, ok := kc.producerMap[producerName]; ok {
			return errors.Errorf("new producer: %s, topic: %s, partitionNum: %d already exists", producerName, fullTopicName, partitionNum)
		}
	}

	var kafkaCnf = &kafka.ConfigMap{
		"compression.type":                      "none",                               // 消息压缩方式，如 none、gzip、snappy、lz4、zstd。默认 none
		"sticky.partitioning.linger.ms":         kc.config.StickyPartitioningLingerMs, // 黏性分区策略的延迟时间，此设置允许生产者在指定时间内将消息发送到同一个分区，以增加消息批次的大小，提高压缩效率和吞吐量。设置为 0 时，生产者不会等待，消息会立即发送。默认 100ms
		"batch.size":                            kc.config.BatchSize,                  // 批量发送大小，默认 10485760 字节
		"message.max.bytes":                     kc.config.MessageMaxBytes,            // 最大消息大小，默认 16384 字节
		"retries":                               kc.config.Retries,                    // 发送消息失败后允许重试的次数，默认 2147483647
		"retry.backoff.ms":                      kc.config.RetryBackoffMs,             // 发送消息失败后，下一次重试发送前的等待时间，默认 100ms
		"linger.ms":                             kc.config.LingerMs,                   // 发送延迟时间，默认 100ms
		"queue.buffering.max.kbytes":            kc.config.QueueBufferingMaxKbytes,    // Producer 攒批发送中，默认 1048576kb
		"max.in.flight.requests.per.connection": 1,                                    // 生产者在收到服务器响应之前可以发送多少个消息，设置为 1 可以保证消息是按照发送的顺序写入服务器，即使发生了重试
		"acks":                                  "1",                                  // 回复
	}
	if err = kafkaSetKey(kafkaCnf, map[string]string{
		"api.version.request": "true",
		"bootstrap.servers":   kc.config.BootstrapServers,
		"security.protocol":   kc.config.SecurityProtocol,
		"sasl.mechanism":      kc.config.SaslMechanism,
		"sasl.username":       kc.config.SaslUsername,
		"sasl.password":       kc.config.SaslPassword,
	}); err != nil {
		return
	}

	var producer *kafka.Producer
	if producer, err = kafka.NewProducer(kafkaCnf); err != nil {
		return
	}
	if producer == nil {
		err = errors.Errorf("new producer %s failed", producerName)
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
	kc.logger.Infof(ctx, "new producer: %s, topic: %s, partitionNum: %d success", producerName, fullTopicName, partitionNum)
	return
}

// NewConsumer 创建消费者
func (kc *KafkaClient) NewConsumer(ctx context.Context, topic string) (err error) {
	// 获取消费者配置
	var (
		isStart      bool
		partitionNum uint32
		groups       []string
	)
	if isStart, partitionNum, groups, err = kc.getConsumerConfig(topic); err != nil {
		return
	}
	if !isStart {
		return
	}
	// 创建消费者
	var (
		consumerNameList = []string{kc.getConsumerName(topic)}
		groupList        = []string{kc.getConsumerGroupName(topic)}
		fullTopicName    = kc.getFullTopicName(topic)
	)
	if len(groups) > 0 {
		consumerNameList = make([]string, 0, len(groups))
		groupList = make([]string, 0, len(groups))
		for _, g := range groups {
			consumerNameList = append(consumerNameList, kc.getConsumerName(g))
			groupList = append(groupList, kc.getConsumerGroupName(g))
		}
	}

	for i := 0; i < len(consumerNameList); i++ {
		var (
			consumerName = consumerNameList[i]
			group        = groupList[i]
		)
		if _, ok := kc.consumerMap[consumerName]; ok {
			return errors.Errorf("new consumer: %s, topic: %s, group: %s, partitionNum: %d already exists", consumerName, fullTopicName, group, partitionNum)
		}
		if kc.config.IsClose {
			kc.logger.Infof(ctx, "new consumer: %s, topic: %s, group: %s, partitionNum: %d success (isClosed)", consumerName, fullTopicName, group, partitionNum)
			kc.consumerMap[consumerName] = nil
			continue
		}
		var kafkaCnf = &kafka.ConfigMap{
			"auto.offset.reset":         kc.config.OffsetReset, // 重置消费者偏移量的策略，可选值: earliest 最早位置，latest 最新位置，none 找不到之前的偏移量，消费者将抛出一个异常，停止工作，默认 earliest
			"heartbeat.interval.ms":     3000,                  // 心跳间隔时间，默认3s
			"session.timeout.ms":        45000,                 // 会话超时时间，默认45s
			"max.poll.interval.ms":      300000,                // 最大拉取间隔时间，默认300s
			"fetch.max.bytes":           52428800,              // 一次 fetch 请求，从一个 broker 中取得的 records 最大大小
			"max.partition.fetch.bytes": 104857600,             // 服务器从每个分区里返回给消费者的最大字节数
			"enable.auto.commit":        false,                 // 关闭自动提交
		}
		if err = kafkaSetKey(kafkaCnf, map[string]string{
			"api.version.request": "true",
			"bootstrap.servers":   kc.config.BootstrapServers,
			"security.protocol":   kc.config.SecurityProtocol,
			"sasl.mechanism":      kc.config.SaslMechanism,
			"sasl.username":       kc.config.SaslUsername,
			"sasl.password":       kc.config.SaslPassword,
			"group.id":            group,
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
				return errors.Errorf("new consumer %s failed", consumerName)
			}
			consumerList = append(consumerList, consumer)
		}

		kc.consumerMap[consumerName] = consumerList
		kc.logger.Infof(ctx, "new consumer: %s, topic: %s, group: %s, partitionNum: %d success", consumerName, fullTopicName, group, partitionNum)
	}
	return
}

// SendMessage 发送消息
func (kc *KafkaClient) SendMessage(ctx context.Context, topic string, producerMessage *ProducerMessage) (err error) {
	var dataMap map[string]any
	if dataMap, err = gtkconv.ToStringMapE(producerMessage.Data); err != nil {
		return
	}
	delete(dataMap, "created_at")
	delete(dataMap, "updated_at")
	delete(dataMap, "deleted_at")
	// 处理数据
	var dataBytes []byte
	if dataBytes, err = json.Marshal(dataMap); err != nil {
		return
	}
	producerMessage.dataBytes = dataBytes
	return kc.sendMessage(ctx, topic, producerMessage)
}

// Subscribe 订阅数据
func (kc *KafkaClient) Subscribe(ctx context.Context, topic string, fn func(message *kafka.Message) error, group ...string) (err error) {
	// 获取消费者配置
	var (
		isStart bool
		groups  []string
	)
	if isStart, _, groups, err = kc.getConsumerConfig(topic); err != nil {
		return
	}
	if !isStart {
		return
	}
	if len(groups) > 0 && len(group) > 0 {
		if !gtkarr.ContainsStr(groups, group[0]) {
			return errors.Errorf("group: %s not found in groups: %s", group[0], groups)
		}
	}

	var (
		consumerName = kc.getConsumerName(topic)
		topics       = []string{kc.getFullTopicName(topic)}
	)
	if len(group) > 0 {
		consumerName = kc.getConsumerName(group[0])
	}
	if kc.config.IsClose {
		kc.logger.Infof(ctx, "subscribe consumer: %s, topics: %v (isClosed)", consumerName, topics)
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
					msg, pErr := consumer.ReadMessage(kc.config.WaitTimeout)
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

// BatchSubscribe 批量订阅数据
func (kc *KafkaClient) BatchSubscribe(ctx context.Context, topic string, fn func(messages []*kafka.Message) error, group ...string) (err error) {
	// 获取消费者配置
	var (
		isStart bool
		groups  []string
	)
	if isStart, _, groups, err = kc.getConsumerConfig(topic); err != nil {
		return
	}
	if !isStart {
		return
	}
	if len(groups) > 0 && len(group) > 0 {
		if !gtkarr.ContainsStr(groups, group[0]) {
			return errors.Errorf("group: %s not found in groups: %s", group[0], groups)
		}
	}

	var (
		consumerName = kc.getConsumerName(topic)
		topics       = []string{kc.getFullTopicName(topic)}
	)
	if len(group) > 0 {
		consumerName = kc.getConsumerName(group[0])
	}
	if kc.config.IsClose {
		kc.logger.Infof(ctx, "batch subscribe consumer: %s, topics: %v (isClosed)", consumerName, topics)
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
			msgList := make([]*kafka.Message, 0, kc.config.BatchConsumeSize)
			// 定时器
			ticker := time.NewTicker(kc.config.BatchConsumeInterval)
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
						msg, pErr := consumer.ReadMessage(kc.config.WaitTimeout)
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
						msgList = make([]*kafka.Message, 0, kc.config.BatchConsumeSize)
					}
				case msg, ok := <-readMsg:
					if ok {
						msgList = append(msgList, msg)
						if len(msgList) == kc.config.BatchConsumeSize {
							// 批量处理数据
							kc.handelBatchData(ctx, consumerName, consumer, msgList, fn)
							// 重新创建一个新的 msgList
							msgList = make([]*kafka.Message, 0, kc.config.BatchConsumeSize)
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
		for retryMaxCount == 0 || (retryMaxCount > 0 && count < retryMaxCount) {
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

// sendMessage 发送消息
func (kc *KafkaClient) sendMessage(ctx context.Context, topic string, producerMessage *ProducerMessage) (err error) {
	// 获取生产者配置
	var (
		isStart      bool
		partitionNum uint32
	)
	if isStart, partitionNum, err = kc.getProducerConfig(topic); err != nil {
		return
	}
	if !isStart {
		return
	}
	// 组装消息
	var (
		fullTopicName = kc.getFullTopicName(topic)
		msg           = &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &fullTopicName, Partition: kc.calcPartition(producerMessage.Key, partitionNum)},
			Value:          producerMessage.dataBytes,
			Key:            []byte(producerMessage.Key),
		}
		producerName = kc.getProducerName(topic)
	)
	// 判断是否配置了全局生产者名称
	globalProducerName := gtkstr.TrimAll(kc.config.GlobalProducer)
	if globalProducerName != "" {
		producerName = kc.getGlobalProducerName(globalProducerName)
	}
	if kc.config.IsClose {
		kc.logger.Debugf(ctx, "%s producer sendMessage(isClosed): %s, data: %s", producerName, gtkjson.MustString(msg), string(producerMessage.dataBytes))
		return
	}
	// 检测哪些 topic 不发送 Kafka 消息
	if gtkarr.ContainsStr(kc.config.ExcludeTopics, topic) {
		kc.logger.Debugf(ctx, "%s producer sendMessage: %s, data: %s", producerName, gtkjson.MustString(msg), string(producerMessage.dataBytes))
		return
	}
	// 发送消息
	var producer *kafka.Producer
	var ok bool
	if producer, ok = kc.producerMap[producerName]; !ok {
		err = errors.Errorf("sendMessage producer: %s, msg: %s, data: %s not found", producerName, gtkjson.MustString(msg), string(producerMessage.dataBytes))
		return
	}
	if producer == nil {
		err = errors.Errorf("sendMessage producer: %s, msg: %s, data: %s producer nil", producerName, gtkjson.MustString(msg), string(producerMessage.dataBytes))
		return
	}
	if err = producer.Produce(msg, nil); err != nil {
		return
	}
	producer.Flush(10 * 1000)
	return
}

// getProducerConfig 获取生产者配置
func (kc *KafkaClient) getProducerConfig(topic string) (isStart bool, partitionNum uint32, err error) {
	if config, ok := kc.config.TopicConfig[topic]; ok {
		isStart = (config.Mode == ModeBoth || config.Mode == ModeProducer)
		if config.PartitionNum > 0 {
			partitionNum = config.PartitionNum
		} else {
			// 默认分区数
			partitionNum = defaultPartitionNum
		}
		return
	}
	err = errors.Errorf("topic `%s` Not Found", topic)
	return
}

// getConsumerConfig 获取消费者配置
func (kc *KafkaClient) getConsumerConfig(topic string) (isStart bool, partitionNum uint32, groups []string, err error) {
	if config, ok := kc.config.TopicConfig[topic]; ok {
		isStart = (config.Mode == ModeBoth || config.Mode == ModeConsumer)
		if config.PartitionNum > 0 {
			partitionNum = config.PartitionNum
		} else {
			// 默认分区数
			partitionNum = defaultPartitionNum
		}
		groups = config.Groups
		return
	}
	err = errors.Errorf("topic `%s` Not Found", topic)
	return
}

// calcPartition 计算分区号
func (kc *KafkaClient) calcPartition(key string, partitionNum uint32) (partition int32) {
	if key == "" {
		partition = kafka.PartitionAny
		return
	}
	hash := fnv.New32a()
	hash.Write([]byte(key))
	partition = int32(hash.Sum32() % partitionNum)
	return
}

// getGlobalProducerName 获取全局生产者名称
func (kc *KafkaClient) getGlobalProducerName(globalProducer string) (producerName string) {
	return fmt.Sprintf("producer_%s", globalProducer)
}

// getProducerName 获取生产者名称
func (kc *KafkaClient) getProducerName(topic string) (producerName string) {
	return fmt.Sprintf("producer_%s", topic)
}

// getConsumerName 获取消费者名称
func (kc *KafkaClient) getConsumerName(topic string) (consumerName string) {
	return fmt.Sprintf("consumer_%s", topic)
}

// getFullTopicName 获取完整的 topic 名称
func (kc *KafkaClient) getFullTopicName(topic string) (fullTopicName string) {
	return fmt.Sprintf("%s_%s", kc.config.Env, topic)
}

// getConsumerGroupName 获取消费者组名称
func (kc *KafkaClient) getConsumerGroupName(topic string) (group string) {
	return fmt.Sprintf("%s_group_%s", kc.config.ConsumerEnv, topic)
}

// kafkaSetKey kafka 设置连接配置
func kafkaSetKey(kafkaCnf *kafka.ConfigMap, kafkaCnfMap map[string]string) (err error) {
	for k, v := range kafkaCnfMap {
		if v == "" {
			continue
		}
		if err = kafkaCnf.SetKey(k, v); err != nil {
			return
		}
	}

	return
}
