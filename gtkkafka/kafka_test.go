/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-20 00:06:58
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-23 18:04:38
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkkafka_test

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/liusuxian/go-toolkit/gtkconf"
	"github.com/liusuxian/go-toolkit/gtkkafka"
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

func TestNewClient(t *testing.T) {
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
	kafkaClient, err = gtkkafka.NewClient(config)
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
	}, "testname1"); err != nil {
		t.Fatal("Subscribe Error: ", err)
	}
	if err = kafkaClient.Subscribe(ctx, "topic_100", func(message *kafka.Message) error {
		t.Logf("topic_100 receive message: %v", string(message.Value))
		return nil
	}, "testname2"); err != nil {
		t.Fatal("Subscribe Error: ", err)
	}
	if err = kafkaClient.BatchSubscribe(ctx, "topic_101", func(messages []*kafka.Message) error {
		return nil
	}, "testname1"); err != nil {
		t.Fatal("BatchSubscribe Error: ", err)
	}
	if err = kafkaClient.BatchSubscribe(ctx, "topic_101", func(messages []*kafka.Message) error {
		return nil
	}, "testname2"); err != nil {
		t.Fatal("BatchSubscribe Error: ", err)
	}
	if err = kafkaClient.BatchSubscribe(ctx, "topic_102", func(messages []*kafka.Message) error {
		return nil
	}); err != nil {
		t.Fatal("BatchSubscribe Error: ", err)
	}
	time.Sleep(5 * time.Second)
}
