/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-19 23:42:12
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2026-01-28 16:38:25
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
	"github.com/liusuxian/go-toolkit/gtkjson"
	"github.com/liusuxian/go-toolkit/gtklog"
	"github.com/liusuxian/go-toolkit/gtkretry"
	"hash/fnv"
	"math"
	"slices"
	"strings"
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
	PartitionNum uint32 `json:"partition_num"`
	// 启动模式 0:不启动生产者或消费者 1:仅启动生产者 2:仅启动消费者 3:同时启动生产者和消费者
	Mode ProducerConsumerStartMode `json:"mode"`
	// 指定消费者组名称列表。如果未指定，将使用默认格式："$consumerEnv_group_$topic"，其中`$consumerEnv_group_`是系统根据当前环境自动添加的前缀
	// 可以配置多个消费者组名称，系统会自动在每个名称前添加"$consumerEnv_group_"前缀
	Groups []string `json:"groups,omitempty"`
	// 批量消费的条数，默认 200
	BatchConsumeSize int `json:"batch_consume_size"`
	// 批量消费的间隔时间，默认 5s
	BatchConsumeInterval time.Duration `json:"batch_consume_interval"`
	// 当消费失败时的重试配置，默认不重试
	RetryConfig gtkretry.RetryConfig `json:"retry_config"`
}

// ProducerMessage 生产者消息
type ProducerMessage struct {
	Key       string `json:"key,omitempty"` // 键
	Data      any    `json:"data"`          // 数据
	dataBytes []byte // 数据字节数组
}

// Config kafka 客户端配置
type Config struct {
	BootstrapServers           string                 `json:"bootstrap_servers"`             // Kafka 服务器的地址列表，格式为 host1:port1,host2:port2
	SecurityProtocol           string                 `json:"security_protocol"`             // Kafka 通信的安全协议，如 PLAINTEXT、SSL、SASL_PLAINTEXT、SASL_SSL
	SaslMechanism              string                 `json:"sasl_mechanism"`                // SASL 认证机制，如 GSSAPI、PLAIN、OAUTHBEARER、SCRAM-SHA-256、SCRAM-SHA-512
	SaslUsername               string                 `json:"sasl_username"`                 // SASL 认证的用户名
	SaslPassword               string                 `json:"sasl_password"`                 // SASL 认证的密码
	StickyPartitioningLingerMs int                    `json:"sticky_partitioning_linger_ms"` // 黏性分区策略的延迟时间，此设置允许生产者在指定时间内将消息发送到同一个分区，以增加消息批次的大小，提高压缩效率和吞吐量。设置为 0 时，生产者不会等待，消息会立即发送。默认 100ms
	BatchSize                  int                    `json:"batch_size"`                    // 批量发送大小，默认 10485760 字节
	MessageMaxBytes            int                    `json:"message_max_bytes"`             // 最大消息大小，默认 16384 字节
	Retries                    int                    `json:"retries"`                       // 发送消息失败后允许重试的次数，默认 2147483647
	RetryBackoffMs             int                    `json:"retry_backoff_ms"`              // 发送消息失败后，下一次重试发送前的等待时间，默认 100ms
	LingerMs                   int                    `json:"linger_ms"`                     // 发送延迟时间，默认 100ms
	QueueBufferingMaxKbytes    int                    `json:"queue_buffering_max_kbytes"`    // Producer 攒批发送中，默认 1048576kb
	WaitTimeout                time.Duration          `json:"wait_timeout"`                  // 指定等待消息的最大时间，默认 -1，表示无限期等待消息，直到有消息到达
	OffsetReset                string                 `json:"offset_reset"`                  // 重置消费者偏移量的策略，可选值: earliest 最早位置，latest 最新位置，none 找不到之前的偏移量，消费者将抛出一个异常，停止工作，默认 earliest
	IsClose                    bool                   `json:"is_close"`                      // 是否不启动 Kafka 客户端（适用于本地调试有时候没有kafka环境的情况）
	Env                        string                 `json:"env"`                           // topic 服务环境，默认 local
	ConsumerEnv                string                 `json:"consumer_env"`                  // 消费者服务环境，默认和 topic 服务环境一致
	GlobalProducer             string                 `json:"global_producer"`               // 全局生产者名称，配置此项时，客户端将使用全局生产者，不再创建新的生产者，默认为空
	TopicConfig                map[string]TopicConfig `json:"topic_config"`                  // topic 配置，key 为 topic 名称
	ExcludeTopics              []string               `json:"exclude_topics"`                // 指定哪些 topic 不发送 Kafka 消息
}

// KafkaClient kafka 客户端
type KafkaClient struct {
	producerMap map[string]*kafka.Producer
	consumerMap map[string][]*kafka.Consumer
	config      *Config        // kafka 客户端配置
	logger      gtklog.ILogger // 日志接口
}

const (
	defaultPartitionNum uint32 = 12 // 默认分区数
)

// NewClient 创建 kafka 客户端
func NewClient(cfg *Config) (client *KafkaClient, err error) {
	if cfg == nil {
		err = fmt.Errorf("kafka client config is nil")
		return
	}
	client = &KafkaClient{
		producerMap: make(map[string]*kafka.Producer),
		consumerMap: make(map[string][]*kafka.Consumer),
		config:      cfg,
		logger:      gtklog.NewDefaultLogger(gtklog.TraceLevel),
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
	// 重置消费者偏移量的策略，可选值: earliest 最早位置，latest 最新位置，none 找不到之前的偏移量，消费者将抛出一个异常，停止工作，默认 earliest
	if client.config.OffsetReset == "" {
		client.config.OffsetReset = "earliest"
	}
	// topic 服务环境，默认 local
	if client.config.Env == "" {
		client.config.Env = "local"
	}
	// 消费者服务环境，默认和 topic 服务环境一致
	if client.config.ConsumerEnv == "" {
		client.config.ConsumerEnv = client.config.Env
	}
	// 处理每个 topic 的默认配置
	for topic, topicCfg := range client.config.TopicConfig {
		// topic 分区数量，默认 12 个分区
		if topicCfg.PartitionNum == 0 {
			topicCfg.PartitionNum = defaultPartitionNum
		}
		// 批量消费的条数，默认 200
		if topicCfg.BatchConsumeSize <= 0 {
			topicCfg.BatchConsumeSize = 200
		}
		// 批量消费的间隔时间，默认 5s
		if topicCfg.BatchConsumeInterval <= time.Duration(0) {
			topicCfg.BatchConsumeInterval = time.Second * 5
		}
		// 填充重试配置的默认值
		topicCfg.RetryConfig = gtkretry.WithDefaults(topicCfg.RetryConfig)
		// 更新回配置（因为 map 中存的是值类型，需要重新赋值）
		client.config.TopicConfig[topic] = topicCfg
	}
	return
}

// SetLogger 设置日志对象
func (kc *KafkaClient) SetLogger(logger gtklog.ILogger) {
	kc.logger = logger
}

// PrintClientConfig 打印消息队列客户端配置
func (kc *KafkaClient) PrintClientConfig(ctx context.Context) {
	kc.logger.Debugf(ctx, "client config: %s\n", gtkjson.MustString(kc.config))
}

// NewProducer 创建生产者
func (kc *KafkaClient) NewProducer(ctx context.Context, topic string) (err error) {
	// 获取生产者配置
	var (
		isStart     bool
		topicConfig *TopicConfig
	)
	if isStart, topicConfig, err = kc.getProducerConfig(topic); err != nil {
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
	globalProducerName := strings.Trim(kc.config.GlobalProducer, " ")
	if globalProducerName != "" {
		producerName = kc.getGlobalProducerName(globalProducerName)
	}
	if kc.config.IsClose {
		kc.logger.Infof(ctx, "new producer: %s, topic: %s, partitionNum: %d success (isClosed)", producerName, fullTopicName, topicConfig.PartitionNum)
		return
	}
	// 判断全局生产者是否已存在
	if globalProducerName != "" {
		if _, ok := kc.producerMap[producerName]; ok {
			kc.logger.Infof(ctx, "new producer: %s, topic: %s, partitionNum: %d success", producerName, fullTopicName, topicConfig.PartitionNum)
			return
		}
	} else {
		if _, ok := kc.producerMap[producerName]; ok {
			return fmt.Errorf("new producer: %s, topic: %s, partitionNum: %d already exists", producerName, fullTopicName, topicConfig.PartitionNum)
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
		"bootstrap.servers": kc.config.BootstrapServers,
		"security.protocol": kc.config.SecurityProtocol,
		"sasl.mechanism":    kc.config.SaslMechanism,
	}); err != nil {
		return
	}

	var producer *kafka.Producer
	if producer, err = kafka.NewProducer(kafkaCnf); err != nil {
		return
	}
	if producer == nil {
		err = fmt.Errorf("new producer %s failed", producerName)
		return
	}
	if kc.config.SaslUsername != "" || kc.config.SaslPassword != "" {
		if err = producer.SetSaslCredentials(kc.config.SaslUsername, kc.config.SaslPassword); err != nil {
			return
		}
	}

	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					kc.logger.Errorf(ctx, "producer: %s, failed to write access log entry: %v, topic: %v, partition: %v, offset: %v, key: %s, content: %s, timestamp: %v",
						producerName, ev.TopicPartition.Error, *ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset, ev.Key, ev.Value, ev.Timestamp)
				} else {
					kc.logger.Debugf(ctx, "producer: %s, send ok topic: %v, partition: %v, offset: %v, key: %s, content: %s, timestamp: %v",
						producerName, *ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset, ev.Key, ev.Value, ev.Timestamp)
				}
			}
		}
	}()

	kc.producerMap[producerName] = producer
	kc.logger.Infof(ctx, "new producer: %s, topic: %s, partitionNum: %d success", producerName, fullTopicName, topicConfig.PartitionNum)
	return
}

// NewConsumer 创建消费者
func (kc *KafkaClient) NewConsumer(ctx context.Context, topic string) (err error) {
	// 获取消费者配置
	var (
		isStart     bool
		topicConfig *TopicConfig
	)
	if isStart, topicConfig, err = kc.getConsumerConfig(topic); err != nil {
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
	if len(topicConfig.Groups) > 0 {
		consumerNameList = make([]string, 0, len(topicConfig.Groups))
		groupList = make([]string, 0, len(topicConfig.Groups))
		for _, g := range topicConfig.Groups {
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
			return fmt.Errorf("new consumer: %s, topic: %s, group: %s, partitionNum: %d already exists", consumerName, fullTopicName, group, topicConfig.PartitionNum)
		}
		if kc.config.IsClose {
			kc.logger.Infof(ctx, "new consumer: %s, topic: %s, group: %s, partitionNum: %d success (isClosed)", consumerName, fullTopicName, group, topicConfig.PartitionNum)
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
			"bootstrap.servers": kc.config.BootstrapServers,
			"security.protocol": kc.config.SecurityProtocol,
			"sasl.mechanism":    kc.config.SaslMechanism,
			"group.id":          group,
		}); err != nil {
			return
		}

		consumerList := make([]*kafka.Consumer, 0, topicConfig.PartitionNum)
		for i := 0; i < int(topicConfig.PartitionNum); i++ {
			consumer, cErr := kafka.NewConsumer(kafkaCnf)
			if cErr != nil {
				return cErr
			}
			if consumer == nil {
				return fmt.Errorf("new consumer %s failed", consumerName)
			}
			if kc.config.SaslUsername != "" || kc.config.SaslPassword != "" {
				if err = consumer.SetSaslCredentials(kc.config.SaslUsername, kc.config.SaslPassword); err != nil {
					return
				}
			}
			consumerList = append(consumerList, consumer)
		}

		kc.consumerMap[consumerName] = consumerList
		kc.logger.Infof(ctx, "new consumer: %s, topic: %s, group: %s, partitionNum: %d success", consumerName, fullTopicName, group, topicConfig.PartitionNum)
	}
	return
}

// SendMessage 发送消息
func (kc *KafkaClient) SendMessage(ctx context.Context, topic string, producerMessage *ProducerMessage) (err error) {
	// 处理数据
	var dataBytes []byte
	if dataBytes, err = json.Marshal(producerMessage.Data); err != nil {
		return
	}
	producerMessage.dataBytes = dataBytes
	return kc.sendMessage(ctx, topic, producerMessage)
}

// Subscribe 订阅数据
func (kc *KafkaClient) Subscribe(ctx context.Context, topic string, fn func(message *kafka.Message) error, group ...string) (err error) {
	return kc.handelSubscribe(ctx, topic, false, func(messages []*kafka.Message) error {
		return fn(messages[0])
	}, group...)
}

// BatchSubscribe 批量订阅数据
func (kc *KafkaClient) BatchSubscribe(ctx context.Context, topic string, fn func(messages []*kafka.Message) error, group ...string) (err error) {
	return kc.handelSubscribe(ctx, topic, true, fn, group...)
}

// handelSubscribe 处理订阅数据
func (kc *KafkaClient) handelSubscribe(ctx context.Context, topic string, isBatch bool, fn func(messages []*kafka.Message) error, group ...string) (err error) {
	// 获取消费者配置
	var (
		isStart     bool
		topicConfig *TopicConfig
	)
	if isStart, topicConfig, err = kc.getConsumerConfig(topic); err != nil {
		return
	}
	if !isStart {
		return
	}
	if len(topicConfig.Groups) > 0 && len(group) > 0 {
		if !slices.Contains(topicConfig.Groups, group[0]) {
			return fmt.Errorf("group: %s not found in groups: %s", group[0], topicConfig.Groups)
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
		kc.logger.Infof(ctx, "handelSubscribe consumer: %s, topics: %v (is closed)", consumerName, topics)
		return
	}
	// 订阅topics
	var consumerList []*kafka.Consumer
	var ok bool
	if consumerList, ok = kc.consumerMap[consumerName]; !ok {
		err = fmt.Errorf("consumer: %s, topics: %v not found", consumerName, topics)
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
					kc.logger.Errorf(ctx, "handelSubscribe consumer: %s %v, panic: %+v", consumerName, consumer, r)
				}
			}()
			// 统一处理批量和单条数据
			var (
				batchSize = 1 // 单条模式默认为1
				ticker    *time.Ticker
				tickerC   <-chan time.Time // nil channel 在 select 中会永久阻塞，实现单条模式下不触发定时器
			)
			if isBatch {
				batchSize = topicConfig.BatchConsumeSize
				ticker = time.NewTicker(topicConfig.BatchConsumeInterval)
				tickerC = ticker.C
				defer ticker.Stop()
			}
			msgList := make([]*kafka.Message, 0, batchSize)
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
									kc.logger.Errorf(ctx, "handelSubscribe consumer: %s %v, message: %+v, error: %+v", consumerName, consumer, msg, pErr)
								}
							} else {
								kc.logger.Errorf(ctx, "handelSubscribe consumer: %s %v, message: %+v, error: %+v", consumerName, consumer, msg, pErr)
							}
						}
					}
				}
			}()
			// 处理消息
			for {
				select {
				case <-tickerC: // tickerC 为 nil 时此分支永远不会被选中（单条模式）
					if len(msgList) > 0 {
						// 处理数据
						kc.handelData(ctx, topicConfig, consumerName, consumer, msgList, fn)
						// 重新创建一个新的 msgList
						msgList = make([]*kafka.Message, 0, batchSize)
					}
				case msg, ok := <-readMsg:
					if ok {
						msgList = append(msgList, msg)
						if len(msgList) == batchSize {
							// 处理数据
							kc.handelData(ctx, topicConfig, consumerName, consumer, msgList, fn)
							// 重新创建一个新的 msgList
							msgList = make([]*kafka.Message, 0, batchSize)
						}
					}
				case <-ctx.Done():
					if len(msgList) > 0 {
						// 处理数据
						kc.handelData(ctx, topicConfig, consumerName, consumer, msgList, fn)
					}
					return
				}
			}
		}(v)
	}
	return
}

// handelData 处理数据
func (kc *KafkaClient) handelData(ctx context.Context, topicConfig *TopicConfig, consumerName string, consumer *kafka.Consumer, msgList []*kafka.Message, fn func(messages []*kafka.Message) error) {
	// 判断 msgList 是否为空
	msgSize := len(msgList)
	if msgSize == 0 {
		return
	}
	// 获取最后一条消息
	lastMessage := msgList[msgSize-1]
	// 重试条件
	retryConfig := topicConfig.RetryConfig
	retryConfig.Condition = func(attempt int, err error) bool {
		// 打印错误日志
		kc.logger.Errorf(ctx, "consumer: %s %v, error: %+v, attempt: %d, topic: %v, partition: %v, offset: %v, key: %s, content: %s, timestamp: %v", consumerName, consumer,
			err, attempt, *lastMessage.TopicPartition.Topic, lastMessage.TopicPartition.Partition, lastMessage.TopicPartition.Offset, string(lastMessage.Key), string(lastMessage.Value), lastMessage.Timestamp)
		// 如果是首次失败（attempt == 0），暂停该分区的消费
		if attempt == 0 {
			if err := consumer.Pause([]kafka.TopicPartition{
				{
					Topic:     lastMessage.TopicPartition.Topic,
					Partition: lastMessage.TopicPartition.Partition,
					Offset:    lastMessage.TopicPartition.Offset,
				},
			}); err != nil {
				kc.logger.Errorf(ctx, "handelData pause consumer: %s %v, error: %+v, attempt: %d, topic: %v, partition: %v, offset: %v, key: %s, content: %s, timestamp: %v", consumerName, consumer,
					err, attempt, *lastMessage.TopicPartition.Topic, lastMessage.TopicPartition.Partition, lastMessage.TopicPartition.Offset, string(lastMessage.Key), string(lastMessage.Value), lastMessage.Timestamp)
			}
		}
		return true
	}
	// 创建重试实例，并且立即执行重试
	if err := gtkretry.NewRetry(retryConfig).Do(ctx, func(ctx context.Context) error {
		// 发送心跳
		consumer.Poll(0)
		// 执行业务函数
		return fn(msgList)
	}); err != nil {
		kc.logger.Errorf(ctx, "handelData finished, consumer: %s %v, error: %+v, topic: %v, partition: %v, offset: %v, key: %s, content: %s, timestamp: %v", consumerName, consumer,
			err, *lastMessage.TopicPartition.Topic, lastMessage.TopicPartition.Partition, lastMessage.TopicPartition.Offset, string(lastMessage.Key), string(lastMessage.Value), lastMessage.Timestamp)
		// 检查是否是因为 context 被取消（退出信号）
		if ctx.Err() != nil {
			return
		}
	}
	// 无论成功还是失败，重试结束后都应该恢复该分区的消费
	if err := consumer.Resume([]kafka.TopicPartition{
		{
			Topic:     lastMessage.TopicPartition.Topic,
			Partition: lastMessage.TopicPartition.Partition,
			Offset:    lastMessage.TopicPartition.Offset,
		},
	}); err != nil {
		kc.logger.Errorf(ctx, "handelData resume consumer: %s %v, error: %+v, topic: %v, partition: %v, offset: %v, key: %s, content: %s, timestamp: %v", consumerName, consumer,
			err, *lastMessage.TopicPartition.Topic, lastMessage.TopicPartition.Partition, lastMessage.TopicPartition.Offset, string(lastMessage.Key), string(lastMessage.Value), lastMessage.Timestamp)
	}
	// 提交
	if _, err := consumer.CommitMessage(lastMessage); err != nil {
		kc.logger.Errorf(ctx, "handelData submit consumer: %s %v, error: %+v, topic: %v, partition: %v, offset: %v, key: %s, content: %s, timestamp: %v", consumerName, consumer,
			err, *lastMessage.TopicPartition.Topic, lastMessage.TopicPartition.Partition, lastMessage.TopicPartition.Offset, string(lastMessage.Key), string(lastMessage.Value), lastMessage.Timestamp)
	}
}

// sendMessage 发送消息
func (kc *KafkaClient) sendMessage(ctx context.Context, topic string, producerMessage *ProducerMessage) (err error) {
	// 获取生产者配置
	var (
		isStart     bool
		topicConfig *TopicConfig
	)
	if isStart, topicConfig, err = kc.getProducerConfig(topic); err != nil {
		return
	}
	if !isStart {
		return
	}
	// 组装消息
	var (
		fullTopicName = kc.getFullTopicName(topic)
		msg           = &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &fullTopicName, Partition: kc.calcPartition(producerMessage.Key, topicConfig.PartitionNum)},
			Value:          producerMessage.dataBytes,
			Key:            []byte(producerMessage.Key),
		}
		producerName = kc.getProducerName(topic)
	)
	// 判断是否配置了全局生产者名称
	globalProducerName := strings.Trim(kc.config.GlobalProducer, " ")
	if globalProducerName != "" {
		producerName = kc.getGlobalProducerName(globalProducerName)
	}
	if kc.config.IsClose {
		kc.logger.Debugf(ctx, "producer: %s, send message(is closed): %s, data: %s", producerName, gtkjson.MustString(msg), string(producerMessage.dataBytes))
		return
	}
	// 检测哪些 topic 不发送 Kafka 消息
	if slices.Contains(kc.config.ExcludeTopics, topic) {
		kc.logger.Debugf(ctx, "producer: %s, send message(exclude topic): %s, data: %s", producerName, gtkjson.MustString(msg), string(producerMessage.dataBytes))
		return
	}
	// 发送消息
	var producer *kafka.Producer
	var ok bool
	if producer, ok = kc.producerMap[producerName]; !ok {
		err = fmt.Errorf("producer: %s, send message: %s, data: %s, producer not found", producerName, gtkjson.MustString(msg), string(producerMessage.dataBytes))
		return
	}
	if producer == nil {
		err = fmt.Errorf("producer: %s, send message: %s, data: %s, producer nil", producerName, gtkjson.MustString(msg), string(producerMessage.dataBytes))
		return
	}
	if err = producer.Produce(msg, nil); err != nil {
		return
	}
	producer.Flush(10 * 1000)
	return
}

// getProducerConfig 获取生产者配置
func (kc *KafkaClient) getProducerConfig(topic string) (isStart bool, topicConfig *TopicConfig, err error) {
	if config, ok := kc.config.TopicConfig[topic]; ok {
		isStart = (config.Mode == ModeBoth || config.Mode == ModeProducer)
		topicConfig = &config
		return
	}
	err = fmt.Errorf("topic `%s` Not Found", topic)
	return
}

// getConsumerConfig 获取消费者配置
func (kc *KafkaClient) getConsumerConfig(topic string) (isStart bool, topicConfig *TopicConfig, err error) {
	if config, ok := kc.config.TopicConfig[topic]; ok {
		isStart = (config.Mode == ModeBoth || config.Mode == ModeConsumer)
		topicConfig = &config
		return
	}
	err = fmt.Errorf("topic `%s` Not Found", topic)
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
	return "producer_" + globalProducer
}

// getProducerName 获取生产者名称
func (kc *KafkaClient) getProducerName(topic string) (producerName string) {
	return "producer_" + topic
}

// getConsumerName 获取消费者名称
func (kc *KafkaClient) getConsumerName(topic string) (consumerName string) {
	return "consumer_" + topic
}

// getFullTopicName 获取完整的 topic 名称
func (kc *KafkaClient) getFullTopicName(topic string) (fullTopicName string) {
	return kc.config.Env + "_" + topic
}

// getConsumerGroupName 获取消费者组名称
func (kc *KafkaClient) getConsumerGroupName(topic string) (group string) {
	return kc.config.ConsumerEnv + "_group_" + topic
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
