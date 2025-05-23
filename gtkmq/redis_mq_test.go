/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-04-23 00:30:12
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-23 18:06:53
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkmq_test

import (
	"context"
	"errors"
	"github.com/liusuxian/go-toolkit/gtkconf"
	"github.com/liusuxian/go-toolkit/gtkjson"
	"github.com/liusuxian/go-toolkit/gtkmq"
	"github.com/liusuxian/go-toolkit/gtkredis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisMQProducer(t *testing.T) {
	var (
		ctx    = context.Background()
		assert = assert.New(t)
		config *gtkmq.RedisMQConfig
		client *gtkmq.RedisMQClient
		err    error
	)
	if err = gtkconf.StructKey("redis_mq", &config); err != nil {
		t.Fatal("Get Logger Config Error: ", err)
	}
	client, err = gtkmq.NewRedisMQClient(ctx, &gtkredis.ClientConfig{
		Addr:     "127.0.0.1:6379",
		Username: "default",
		Password: "redis!@#$%",
		DB:       0,
		PoolSize: 20,
	}, config)
	assert.NoError(err)
	defer client.Close()

	for queue := range config.MQConfig {
		// ============================= 以下为生产者 =============================
		if err := client.NewProducer(ctx, queue); err != nil {
			t.Fatalf("create %s producer error: %v", queue, err)
		}
	}
	client.PrintClientConfig(ctx)
	err = client.SendMessage(ctx, "queue_100", &gtkmq.ProducerMessage{Data: map[string]any{"a": "hello world"}})
	assert.NoError(err)
	err = client.SendMessage(ctx, "queue_200", &gtkmq.ProducerMessage{Data: map[string]any{"a": "hello world"}})
	assert.NoError(err)

	for i := 0; i < 24; i++ {
		err = client.SendMessage(ctx, "queue", &gtkmq.ProducerMessage{
			Data: map[string]any{
				"a": i,
				"b": []int{i, i + 1, i + 2},
			}})
		assert.NoError(err)
	}

	err = client.SendMessage(ctx, "queue_delay", &gtkmq.ProducerMessage{
		Data: map[string]any{
			"a": "hello world",
			"b": []int{1, 2, 3},
		},
		DelayTime: time.Now().Add(time.Second * 9),
	})
	assert.NoError(err)
	err = client.SendMessage(ctx, "queue_delay", &gtkmq.ProducerMessage{
		Data: map[string]any{
			"a": "hello world",
			"b": []int{1, 2, 3},
		},
		DelayTime: time.Now().Add(time.Second * 9),
	})
	assert.NoError(err)
	time.Sleep(time.Second * 60)
}

func TestRedisMQConsumerSubscribe(t *testing.T) {
	var (
		ctx    = context.Background()
		assert = assert.New(t)
		config *gtkmq.RedisMQConfig
		client *gtkmq.RedisMQClient
		err    error
	)
	if err = gtkconf.StructKey("redis_mq", &config); err != nil {
		t.Fatal("Get Logger Config Error: ", err)
	}
	client, err = gtkmq.NewRedisMQClient(ctx, &gtkredis.ClientConfig{
		Addr:     "127.0.0.1:6379",
		Username: "default",
		Password: "redis!@#$%",
		DB:       0,
		PoolSize: 20,
	}, config)
	assert.NoError(err)
	defer client.Close()

	for queue := range config.MQConfig {
		// ============================= 以下为消费者 =============================
		if err := client.NewConsumer(ctx, queue); err != nil {
			t.Fatalf("create %s consumer error: %v", queue, err)
		}
	}

	count1 := 0
	err = client.Subscribe(ctx, "queue", func(message *gtkmq.MQMessage) error {
		for count1 < 2 {
			count1++
			return errors.New("test error")
		}
		t.Logf("subscribe111 message: %s\n", gtkjson.MustString(message))
		return nil
	}, "testname1")
	assert.NoError(err)

	count2 := 0
	err = client.Subscribe(ctx, "queue", func(message *gtkmq.MQMessage) error {
		for count2 < 2 {
			count2++
			return errors.New("test error")
		}
		t.Logf("subscribe222 message: %s\n", gtkjson.MustString(message))
		return nil
	}, "testname2")
	assert.NoError(err)
	time.Sleep(time.Second * 5)
}

func TestRedisMQConsumerBatchSubscribe(t *testing.T) {
	var (
		ctx    = context.Background()
		assert = assert.New(t)
		config *gtkmq.RedisMQConfig
		client *gtkmq.RedisMQClient
		err    error
	)
	if err = gtkconf.StructKey("redis_mq", &config); err != nil {
		t.Fatal("Get Logger Config Error: ", err)
	}
	client, err = gtkmq.NewRedisMQClient(ctx, &gtkredis.ClientConfig{
		Addr:     "127.0.0.1:6379",
		Username: "default",
		Password: "redis!@#$%",
		DB:       0,
		PoolSize: 20,
	}, config)
	assert.NoError(err)
	defer client.Close()

	for queue := range config.MQConfig {
		// ============================= 以下为消费者 =============================
		if err := client.NewConsumer(ctx, queue); err != nil {
			t.Fatalf("create %s consumer error: %v", queue, err)
		}
	}

	count1 := 0
	err = client.BatchSubscribe(ctx, "queue", func(message []*gtkmq.MQMessage) error {
		for count1 < 2 {
			count1++
			return errors.New("test error")
		}
		t.Logf("subscribe111 message: %s\n", gtkjson.MustString(message))
		return nil
	}, "testname1")
	assert.NoError(err)

	count2 := 0
	err = client.BatchSubscribe(ctx, "queue", func(message []*gtkmq.MQMessage) error {
		for count2 < 2 {
			count2++
			return errors.New("test error")
		}
		t.Logf("subscribe222 message: %s\n", gtkjson.MustString(message))
		return nil
	}, "testname2")
	assert.NoError(err)
	time.Sleep(time.Second * 5)
}

func TestRedisMQResetConsumerOffset(t *testing.T) {
	var (
		ctx    = context.Background()
		assert = assert.New(t)
		config *gtkmq.RedisMQConfig
		client *gtkmq.RedisMQClient
		err    error
	)
	if err = gtkconf.StructKey("redis_mq", &config); err != nil {
		t.Fatal("Get Logger Config Error: ", err)
	}
	client, err = gtkmq.NewRedisMQClient(ctx, &gtkredis.ClientConfig{
		Addr:     "127.0.0.1:6379",
		Username: "default",
		Password: "redis!@#$%",
		DB:       0,
		PoolSize: 20,
	}, config)
	assert.NoError(err)
	defer client.Close()

	err = client.ResetConsumerOffset(ctx, "queue", "0-0", "testname1")
	assert.NoError(err)
}

func TestRedisMQResetConsumerOffsetByPartition(t *testing.T) {
	var (
		ctx    = context.Background()
		assert = assert.New(t)
		config *gtkmq.RedisMQConfig
		client *gtkmq.RedisMQClient
		err    error
	)
	if err = gtkconf.StructKey("redis_mq", &config); err != nil {
		t.Fatal("Get Logger Config Error: ", err)
	}
	client, err = gtkmq.NewRedisMQClient(ctx, &gtkredis.ClientConfig{
		Addr:     "127.0.0.1:6379",
		Username: "default",
		Password: "redis!@#$%",
		DB:       0,
		PoolSize: 20,
	}, config)
	assert.NoError(err)
	defer client.Close()

	err = client.ResetConsumerOffsetByPartition(ctx, "queue", 0, "0-0", "testname1")
	assert.NoError(err)
}

func TestRedisMQExpiredMessages(t *testing.T) {
	var (
		ctx      = context.Background()
		assert   = assert.New(t)
		config   *gtkmq.RedisMQConfig
		client   *gtkmq.RedisMQClient
		messages map[int32][]*gtkmq.MQMessage
		err      error
	)
	if err = gtkconf.StructKey("redis_mq", &config); err != nil {
		t.Fatal("Get Logger Config Error: ", err)
	}
	client, err = gtkmq.NewRedisMQClient(ctx, &gtkredis.ClientConfig{
		Addr:     "127.0.0.1:6379",
		Username: "default",
		Password: "redis!@#$%",
		DB:       0,
		PoolSize: 20,
	}, config)
	assert.NoError(err)
	defer client.Close()

	messages, err = client.GetExpiredMessages(ctx, "queue", true)
	assert.NoError(err)
	for _, v := range messages {
		for _, m := range v {
			t.Logf("expired message: %s %s\n", m.MQPartition.PartitionName, m.MQPartition.Offset)
		}
	}
}

func TestRedisMQDelGroup(t *testing.T) {
	var (
		ctx    = context.Background()
		assert = assert.New(t)
		config *gtkmq.RedisMQConfig
		client *gtkmq.RedisMQClient
		err    error
	)
	if err = gtkconf.StructKey("redis_mq", &config); err != nil {
		t.Fatal("Get Logger Config Error: ", err)
	}
	client, err = gtkmq.NewRedisMQClient(ctx, &gtkredis.ClientConfig{
		Addr:     "127.0.0.1:6379",
		Username: "default",
		Password: "redis!@#$%",
		DB:       0,
		PoolSize: 20,
	}, config)
	assert.NoError(err)
	defer client.Close()

	err = client.DelGroup(ctx, "queue_100")
	assert.NoError(err)
	err = client.DelGroup(ctx, "queue_200")
	assert.NoError(err)
	err = client.DelGroup(ctx, "queue", "testname2")
	assert.NoError(err)
}

func TestRedisMQDelQueue(t *testing.T) {
	var (
		ctx    = context.Background()
		assert = assert.New(t)
		config *gtkmq.RedisMQConfig
		client *gtkmq.RedisMQClient
		err    error
	)
	if err = gtkconf.StructKey("redis_mq", &config); err != nil {
		t.Fatal("Get Logger Config Error: ", err)
	}
	client, err = gtkmq.NewRedisMQClient(ctx, &gtkredis.ClientConfig{
		Addr:     "127.0.0.1:6379",
		Username: "default",
		Password: "redis!@#$%",
		DB:       0,
		PoolSize: 20,
	}, config)
	assert.NoError(err)
	defer client.Close()

	err = client.DelQueue(ctx, "queue_100")
	assert.NoError(err)
	err = client.DelQueue(ctx, "queue_200")
	assert.NoError(err)
	err = client.DelQueue(ctx, "queue")
	assert.NoError(err)
}

func BenchmarkSendMessage(b *testing.B) {
	var (
		ctx    = context.Background()
		assert = assert.New(b)
		config *gtkmq.RedisMQConfig
		client *gtkmq.RedisMQClient
		err    error
	)
	if err = gtkconf.StructKey("redis_mq", &config); err != nil {
		b.Fatal("Get Logger Config Error: ", err)
	}
	client, err = gtkmq.NewRedisMQClient(ctx, &gtkredis.ClientConfig{
		Addr:     "127.0.0.1:6379",
		Username: "default",
		Password: "redis!@#$%",
		DB:       0,
		PoolSize: 20,
	}, config)
	assert.NoError(err)
	defer client.Close()

	err = client.NewProducer(ctx, "queue")
	assert.NoError(err)

	message := &gtkmq.ProducerMessage{
		Data: map[string]any{"data": "test data"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err = client.SendMessage(ctx, "queue", message); err != nil {
			b.Error("发送消息失败:", err)
			return
		}
	}
	client.DelGroup(ctx, "queue")
	client.DelQueue(ctx, "queue")
}
