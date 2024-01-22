/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-20 00:06:58
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-22 22:11:56
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gfkafka_test

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/liusuxian/go-toolkit/gf/gfkafka"
	"testing"
)

type TestRecordInfo struct {
	Uid       int64  `json:"uid" dc:"uid"`   // 用户ID
	Sid       int64  `json:"sid" dc:"sid"`   // 剧id
	Time      string `json:"time" dc:"time"` // 时间
	Str       string `json:"str" dc:"str"`   // 字符串
	CreatedAt string `json:"created_at" dc:"created_at"`
}

func TestKafka1(t *testing.T) {
	var err error
	ctx := context.Background()
	// 创建 kafka 客户端
	kafkaClient := gfkafka.NewClient(func(cc *gfkafka.ClientConfig) {
		cc.IsClose = true
		cc.Env = "test"
		cc.ProducerConsumerConfigMap = map[string]gfkafka.ProducerConsumerConfig{
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
	})
	t.Logf("NewClient kafkaClient: %+v, ClientConfig: %+v\n", kafkaClient, kafkaClient.GetClientConfig())
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
