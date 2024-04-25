/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-20 00:06:58
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-04-25 23:06:33
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkkafka_test

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/liusuxian/go-toolkit/gtkconf"
	"github.com/liusuxian/go-toolkit/gtkkafka"
	"github.com/liusuxian/go-toolkit/gtklog"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type TestRecordInfo struct {
	Uid       int64  `json:"uid" dc:"uid"`   // 用户ID
	Sid       int64  `json:"sid" dc:"sid"`   // 剧id
	Time      string `json:"time" dc:"time"` // 时间
	Str       string `json:"str" dc:"str"`   // 字符串
	CreatedAt string `json:"created_at" dc:"created_at"`
}

func TestNewWithOption(t *testing.T) {
	var (
		ctx         = context.Background()
		assert      = assert.New(t)
		kafkaClient *gtkkafka.KafkaClient
		err         error
	)
	// 创建 kafka 客户端
	kafkaClient, err = gtkkafka.NewWithOption(func(cc *gtkkafka.Config) {
		cc.IsClose = true
		cc.Env = "test"
		cc.TopicConfig = map[string]gtkkafka.TopicConfig{
			"topic_100": {
				PartitionNum: 12,
				Mode:         gtkkafka.ModeBoth,
			},
		}
		cc.ExcludeEnvTopicMap = map[string][]string{
			"test": {
				"topic_100",
			},
		}
		cc.LogConfig.LogPath = "logs/kafka"
		cc.LogConfig.LogLevelFileName = map[gtklog.Level]string{
			gtklog.TraceLevel: "access.log",
			gtklog.DebugLevel: "access.log",
			gtklog.InfoLevel:  "access.log",
			gtklog.WarnLevel:  "access.log",
			gtklog.ErrorLevel: "error.log",
			gtklog.FatalLevel: "error.log",
			gtklog.PanicLevel: "error.log",
		}
		cc.LogConfig.Stdout = true
	})
	assert.NoError(err)
	// 打印消息队列客户端配置
	kafkaClient.PrintClientConfig(ctx)
	// 创建生产者
	if err = kafkaClient.NewProducer(ctx, "topic_100"); err != nil {
		t.Fatal("NewProducer Error: ", err)
	}
	if err = kafkaClient.NewConsumer(ctx, "topic_100"); err != nil {
		t.Fatal("NewConsumer Error: ", err)
	}
	if err = kafkaClient.SendMessage(ctx, "topic_100", &gtkkafka.ProducerMessage{
		Data: TestRecordInfo{
			Uid:  1111,
			Sid:  2222,
			Time: "2023-11-28 12:00:00",
			Str:  "<iphone 12.12.12, 1234...>"},
	}); err != nil {
		t.Fatal("SendMessage Error: ", err)
	}
	if err = kafkaClient.SendMessage(ctx, "topic_100", &gtkkafka.ProducerMessage{
		Key: "1",
		Data: TestRecordInfo{
			Uid:  1111,
			Sid:  2222,
			Time: "2023-11-28 12:00:00",
			Str:  "<iphone 12.12.12, 1234...>"},
	}); err != nil {
		t.Fatal("SendMessage Error: ", err)
	}
	if err = kafkaClient.Subscribe(ctx, "topic_100", func(message *kafka.Message) error {
		return nil
	}); err != nil {
		t.Fatal("Subscribe Error: ", err)
	}
	if err = kafkaClient.BatchSubscribe(ctx, "topic_100", func(messages []*kafka.Message) error {
		return nil
	}); err != nil {
		t.Fatal("BatchSubscribe Error: ", err)
	}
}

func TestNewWithConfig(t *testing.T) {
	var (
		ctx         = context.Background()
		assert      = assert.New(t)
		config      *gtkkafka.Config
		kafkaClient *gtkkafka.KafkaClient
		err         error
	)
	if err = gtkconf.StructKey("kafka", &config); err != nil {
		t.Fatal("Get Logger Config Error: ", err)
	}
	// 创建 kafka 客户端
	kafkaClient, err = gtkkafka.NewWithConfig(config)
	assert.NoError(err)
	// 打印消息队列客户端配置
	kafkaClient.PrintClientConfig(ctx)

	for topic := range config.TopicConfig {
		// ============================= 以下为生产者 =============================
		if err := kafkaClient.NewProducer(ctx, topic); err != nil {
			t.Fatalf("create %s producer error: %v", topic, err)
		}
		// ============================= 以下为消费者 =============================
		if err := kafkaClient.NewConsumer(ctx, topic); err != nil {
			t.Fatalf("create %s consumer error: %v", topic, err)
		}
	}
	if err = kafkaClient.SendMessage(ctx, "topic_100", &gtkkafka.ProducerMessage{
		Data: TestRecordInfo{
			Uid:  1111,
			Sid:  2222,
			Time: "2023-11-28 12:00:00",
			Str:  "<iphone 12.12.12, 1234...>"},
	}); err != nil {
		t.Fatal("SendMessage Error: ", err)
	}
	if err = kafkaClient.SendMessage(ctx, "topic_101", &gtkkafka.ProducerMessage{
		Key: "1",
		Data: TestRecordInfo{
			Uid:  1111,
			Sid:  2222,
			Time: "2023-11-28 12:00:00",
			Str:  "<iphone 12.12.12, 1234...>"},
	}); err != nil {
		t.Fatal("SendMessage Error: ", err)
	}
	if err = kafkaClient.SendMessage(ctx, "topic_102", &gtkkafka.ProducerMessage{
		Key: "2",
		Data: TestRecordInfo{
			Uid:  1111,
			Sid:  2222,
			Time: "2023-11-28 12:00:00",
			Str:  "<iphone 12.12.12, 1234...>"},
	}); err != nil {
		t.Fatal("SendMessage Error: ", err)
	}
	if err = kafkaClient.Subscribe(ctx, "topic_100", func(message *kafka.Message) error {
		t.Logf("topic_100 receive message: %v", string(message.Value))
		return nil
	}); err != nil {
		t.Fatal("Subscribe Error: ", err)
	}
	if err = kafkaClient.BatchSubscribe(ctx, "topic_101", func(messages []*kafka.Message) error {
		return nil
	}); err != nil {
		t.Fatal("BatchSubscribe Error: ", err)
	}
	if err = kafkaClient.BatchSubscribe(ctx, "topic_102", func(messages []*kafka.Message) error {
		return nil
	}); err != nil {
		t.Fatal("BatchSubscribe Error: ", err)
	}
	time.Sleep(5 * time.Second)
}
