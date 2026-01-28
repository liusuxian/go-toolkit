/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2026-01-24 20:17:14
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2026-01-28 16:46:12
 * @Description:
 *
 * Copyright (c) 2026 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkmq_test

import (
	"context"
	"fmt"
	"github.com/liusuxian/go-toolkit/gtkmq"
	"github.com/liusuxian/go-toolkit/gtkredis"
	"github.com/liusuxian/go-toolkit/gtkretry"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

// getPodTestConfig 获取 Pod 测试配置
func getPodTestConfig() *gtkmq.RedisMQConfig {
	return &gtkmq.RedisMQConfig{
		ExpiredTime:           1 * time.Minute,
		DelExpiredMsgInterval: 2 * time.Second,
		Env:                   "test",
		MQConfig: map[string]gtkmq.MQConfig{
			"pod_queue": {
				PartitionNum: 1,
				Mode:         gtkmq.ModeBoth,
				Groups:       []string{"testname1", "testname2"},
				RetryConfig: gtkretry.RetryConfig{
					MaxAttempts: -1,
				},
			},
		},
	}
}

func TestRedisMQPodProducer(t *testing.T) {
	var (
		ctx    = context.Background()
		assert = assert.New(t)
		config = getPodTestConfig()
		client *gtkmq.RedisMQClient
		err    error
	)
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
	err = client.SendMessage(ctx, "pod_queue", &gtkmq.ProducerMessage{
		Data: map[string]any{
			"a": "hello world",
			"b": []int{1, 2, 3},
		}})
	assert.NoError(err)
}

func TestRedisMQPodConsumerSubscribe(t *testing.T) {
	var (
		ctx    = context.Background()
		assert = assert.New(t)
		config = getPodTestConfig()
		client *gtkmq.RedisMQClient
		err    error
	)
	client, err = gtkmq.NewRedisMQClient(ctx, &gtkredis.ClientConfig{
		Addr:     "127.0.0.1:6379",
		Username: "default",
		Password: "redis!@#$%",
		DB:       0,
		PoolSize: 20,
	}, config)
	assert.NoError(err)
	defer func() {
		client.Close()
	}()

	for queue := range config.MQConfig {
		// ============================= 以下为消费者 =============================
		if err := client.NewConsumer(ctx, queue); err != nil {
			t.Fatalf("create %s consumer error: %v", queue, err)
		}
	}

	ctxOut, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(10)
	for i := range 10 {
		go func(i int) {
			defer wg.Done()

			err = client.Subscribe(ctx, "pod_queue", func(message *gtkmq.MQMessage) error {
				return fmt.Errorf("test error: %d", i)
			}, "testname1")
			assert.NoError(err)
			<-ctxOut.Done()
		}(i)
	}
	wg.Wait()
}
