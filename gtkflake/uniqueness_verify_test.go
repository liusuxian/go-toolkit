/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-12-24
 * @Description: 唯一性验证测试
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkflake

import (
	"sync"
	"testing"
)

// TestRequestIDUniqueness 测试 RequestID 的唯一性（并发场景）
func TestRequestIDUniqueness(t *testing.T) {
	flake, err := New(Settings{})
	if err != nil {
		t.Fatalf("创建 Flake 失败: %v", err)
	}

	const goroutines = 10        // 并发协程数
	const idsPerGoroutine = 1000 // 每个协程生成的ID数

	idChan := make(chan string, goroutines*idsPerGoroutine)
	var wg sync.WaitGroup

	// 启动多个协程并发生成 ID
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < idsPerGoroutine; j++ {
				id, err := flake.RequestID()
				if err != nil {
					t.Errorf("生成 RequestID 失败: %v", err)
					return
				}
				idChan <- id
			}
		}()
	}

	// 等待所有协程完成
	wg.Wait()
	close(idChan)

	// 收集所有 ID 并检查重复
	idMap := make(map[string]bool)
	totalCount := 0
	for id := range idChan {
		totalCount++
		if idMap[id] {
			t.Errorf("发现重复的 RequestID: %s", id)
		}
		idMap[id] = true
	}

	expectedCount := goroutines * idsPerGoroutine
	if totalCount != expectedCount {
		t.Errorf("ID 数量不符: 期望 %d, 实际 %d", expectedCount, totalCount)
	}

	t.Logf("✅ 并发生成 %d 个 RequestID，无重复", totalCount)
}

// TestShortIDUniqueness 测试 ShortID 的唯一性（并发场景）
func TestShortIDUniqueness(t *testing.T) {
	flake, err := New(Settings{})
	if err != nil {
		t.Fatalf("创建 Flake 失败: %v", err)
	}

	const goroutines = 10
	const idsPerGoroutine = 1000

	idChan := make(chan string, goroutines*idsPerGoroutine)
	var wg sync.WaitGroup

	// 启动多个协程并发生成 ID
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < idsPerGoroutine; j++ {
				id, err := flake.ShortID()
				if err != nil {
					t.Errorf("生成 ShortID 失败: %v", err)
					return
				}
				idChan <- id
			}
		}()
	}

	// 等待所有协程完成
	wg.Wait()
	close(idChan)

	// 收集所有 ID 并检查重复
	idMap := make(map[string]bool)
	totalCount := 0
	for id := range idChan {
		totalCount++
		if idMap[id] {
			t.Errorf("发现重复的 ShortID: %s", id)
		}
		idMap[id] = true
	}

	expectedCount := goroutines * idsPerGoroutine
	if totalCount != expectedCount {
		t.Errorf("ID 数量不符: 期望 %d, 实际 %d", expectedCount, totalCount)
	}

	t.Logf("✅ 并发生成 %d 个 ShortID，无重复", totalCount)
}

// TestEncodingConsistency 测试编码的一致性（相同输入→相同输出）
func TestEncodingConsistency(t *testing.T) {
	flake, err := New(Settings{})
	if err != nil {
		t.Fatalf("创建 Flake 失败: %v", err)
	}

	// 生成一个 Snowflake ID
	id1, err := flake.ID()
	if err != nil {
		t.Fatalf("生成 ID 失败: %v", err)
	}

	// 测试 formatHexUpper 的一致性
	hex1 := formatHexUpper(id1)
	hex2 := formatHexUpper(id1)
	if hex1 != hex2 {
		t.Errorf("formatHexUpper 不一致: %s != %s", hex1, hex2)
	}

	// 测试 encodeBase62 的一致性
	b62_1 := encodeBase62(id1)
	b62_2 := encodeBase62(id1)
	if b62_1 != b62_2 {
		t.Errorf("encodeBase62 不一致: %s != %s", b62_1, b62_2)
	}

	t.Logf("✅ 编码函数一致性验证通过")
	t.Logf("   Snowflake ID: %d", id1)
	t.Logf("   Hex: %s", hex1)
	t.Logf("   Base62: %s", b62_1)
}

// TestDifferentIDsDifferentEncodings 测试不同ID生成不同编码
func TestDifferentIDsDifferentEncodings(t *testing.T) {
	flake, err := New(Settings{})
	if err != nil {
		t.Fatalf("创建 Flake 失败: %v", err)
	}

	// 生成多个不同的 ID
	const count = 100
	ids := make([]int64, count)
	for i := 0; i < count; i++ {
		id, err := flake.ID()
		if err != nil {
			t.Fatalf("生成 ID 失败: %v", err)
		}
		ids[i] = id
	}

	// 检查 Hex 编码的唯一性
	hexSet := make(map[string]bool)
	for _, id := range ids {
		hex := formatHexUpper(id)
		if hexSet[hex] {
			t.Errorf("发现重复的 Hex 编码: %s (ID: %d)", hex, id)
		}
		hexSet[hex] = true
	}

	// 检查 Base62 编码的唯一性
	b62Set := make(map[string]bool)
	for _, id := range ids {
		b62 := encodeBase62(id)
		if b62Set[b62] {
			t.Errorf("发现重复的 Base62 编码: %s (ID: %d)", b62, id)
		}
		b62Set[b62] = true
	}

	t.Logf("✅ %d 个不同的 ID 生成了 %d 个唯一的 Hex 编码", count, len(hexSet))
	t.Logf("✅ %d 个不同的 ID 生成了 %d 个唯一的 Base62 编码", count, len(b62Set))
}
