/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-20 00:06:58
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-03-21 16:05:23
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkkafka_test

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/liusuxian/go-toolkit/gtkconf"
	"github.com/liusuxian/go-toolkit/gtkjson"
	"github.com/liusuxian/go-toolkit/gtkkafka"
	"github.com/liusuxian/go-toolkit/gtklog"
	"github.com/stretchr/testify/assert"
	"testing"
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
		cc.ProducerConsumerConfigMap = map[string]gtkkafka.ProducerConsumerConfig{
			"test": {
				Topic:    "topic100",
				StartAll: true,
			},
		}
		cc.ExcludeEnvTopicMap = map[string][]string{
			"test": {
				"topic100",
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
	})
	assert.NoError(err)
	t.Logf("NewClient kafkaClient: %+v, Config: %s\n", kafkaClient, gtkjson.MustString(kafkaClient.GetConfig()))
	if err = kafkaClient.NewProducer(ctx, "test"); err != nil {
		t.Fatal("NewProducer Error: ", err)
	}
	if err = kafkaClient.NewConsumer(ctx, "test"); err != nil {
		t.Fatal("NewConsumer Error: ", err)
	}
	if err = kafkaClient.SendJsonData(ctx, "test", TestRecordInfo{
		Uid:  1111,
		Sid:  2222,
		Time: "2023-11-28 12:00:00",
		Str:  "<iphone 12.12.12, 1234...>"}); err != nil {
		t.Fatal("SendJsonData Error: ", err)
	}
	if err = kafkaClient.SendJsonData(ctx, "test", TestRecordInfo{
		Uid:  1111,
		Sid:  2222,
		Time: "2023-11-28 12:00:00",
		Str:  "<iphone 12.12.12, 1234...>"}, "1"); err != nil {
		t.Fatal("SendJsonData Error: ", err)
	}
	if err = kafkaClient.SubscribeTopics(ctx, "test", func(message *kafka.Message) error {
		return nil
	}); err != nil {
		t.Fatal("SubscribeTopics Error: ", err)
	}
	if err = kafkaClient.BatchSubscribeTopics(ctx, "test", func(messages []*kafka.Message) error {
		return nil
	}); err != nil {
		t.Fatal("BatchSubscribeTopics Error: ", err)
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
	t.Logf("NewClient kafkaClient: %+v, Config: %s\n", kafkaClient, gtkjson.MustString(kafkaClient.GetConfig()))
	for producerConsumerName := range config.ProducerConsumerConfigMap {
		// ============================= 以下为生产者 =============================
		if err := kafkaClient.NewProducer(ctx, producerConsumerName); err != nil {
			t.Fatalf("create %s producer error: %v", producerConsumerName, err)
		}
		// ============================= 以下为消费者 =============================
		if err := kafkaClient.NewConsumer(ctx, producerConsumerName); err != nil {
			t.Fatalf("create %s consumer error: %v", producerConsumerName, err)
		}
	}
	if err = kafkaClient.SendJsonData(ctx, "test_100", TestRecordInfo{
		Uid:  1111,
		Sid:  2222,
		Time: "2023-11-28 12:00:00",
		Str:  "<iphone 12.12.12, 1234...>"}); err != nil {
		t.Fatal("SendJsonData Error: ", err)
	}
	if err = kafkaClient.SendJsonData(ctx, "test_101", TestRecordInfo{
		Uid:  1111,
		Sid:  2222,
		Time: "2023-11-28 12:00:00",
		Str:  "<iphone 12.12.12, 1234...>"}, "1"); err != nil {
		t.Fatal("SendJsonData Error: ", err)
	}
	if err = kafkaClient.SendJsonData(ctx, "test_102", TestRecordInfo{
		Uid:  1111,
		Sid:  2222,
		Time: "2023-11-28 12:00:00",
		Str:  "<iphone 12.12.12, 1234...>"}, "2"); err != nil {
		t.Fatal("SendJsonData Error: ", err)
	}
	if err = kafkaClient.SubscribeTopics(ctx, "test_100", func(message *kafka.Message) error {
		return nil
	}); err != nil {
		t.Fatal("SubscribeTopics Error: ", err)
	}
	if err = kafkaClient.BatchSubscribeTopics(ctx, "test_101", func(messages []*kafka.Message) error {
		return nil
	}); err != nil {
		t.Fatal("BatchSubscribeTopics Error: ", err)
	}
	if err = kafkaClient.BatchSubscribeTopics(ctx, "test_102", func(messages []*kafka.Message) error {
		return nil
	}); err != nil {
		t.Fatal("BatchSubscribeTopics Error: ", err)
	}
}
