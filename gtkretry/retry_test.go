/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-12-09 17:23:44
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-11 12:39:15
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkretry_test

import (
	"context"
	"errors"
	"github.com/liusuxian/go-toolkit/gtkretry"
	"testing"
	"time"
)

func TestRetryDo(t *testing.T) {
	ctx := context.Background()

	gtkretry.NewRetry(gtkretry.RetryConfig{
		MaxAttempts: 0,
		Strategy:    gtkretry.RetryStrategyFixed,
		BaseDelay:   1 * time.Second,
		Condition: func(attempt int, err error) (ok bool) {
			t.Logf("attempt: %d, time: %v", attempt, time.Now().Unix())
			return true
		},
	}).Do(ctx, func(ctx context.Context) (err error) {
		return errors.New("test error")
	})

	gtkretry.NewRetry(gtkretry.RetryConfig{
		MaxAttempts:   3,
		Strategy:      gtkretry.RetryStrategyExponential,
		BaseDelay:     1 * time.Second,
		MaxDelay:      10 * time.Second,
		Multiplier:    2.0,
		JitterPercent: 0.1,
		Condition: func(attempt int, err error) (ok bool) {
			t.Logf("attempt: %d, time: %v", attempt, time.Now().Unix())
			return true
		},
	}).Do(ctx, func(ctx context.Context) (err error) {
		return errors.New("test error")
	})

	gtkretry.NewRetry(gtkretry.RetryConfig{
		MaxAttempts: 3,
		Strategy:    gtkretry.RetryStrategyFixed,
		BaseDelay:   1 * time.Second,
		Condition: func(attempt int, err error) (ok bool) {
			t.Logf("attempt: %d, time: %v", attempt, time.Now().Unix())
			return true
		},
	}).Do(ctx, func(ctx context.Context) (err error) {
		return errors.New("test error")
	})

	gtkretry.NewRetry(gtkretry.RetryConfig{
		MaxAttempts:   -1,
		Strategy:      gtkretry.RetryStrategyExponential,
		BaseDelay:     1 * time.Second,
		MaxDelay:      10 * time.Second,
		Multiplier:    2.0,
		JitterPercent: 0.1,
		Condition: func(attempt int, err error) (ok bool) {
			t.Logf("attempt: %d, time: %v", attempt, time.Now().Unix())
			return attempt != 5
		},
	}).Do(ctx, func(ctx context.Context) (err error) {
		return errors.New("test error")
	})
}
