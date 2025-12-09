/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-27 21:30:00
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-16 15:12:45
 * @Description: 负载均衡器测试
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import (
	"sync"
	"testing"
)

// TestNewLoadBalancer tests creating load balancer
func TestNewLoadBalancer(t *testing.T) {
	t.Run("create empty load balancer", func(t *testing.T) {
		lb := NewLoadBalancer([]string{})
		if lb == nil {
			t.Fatal("load balancer should not be nil")
		}
		if len(lb.apiKeyList) != 0 {
			t.Errorf("expected API key list length to be 0, got %d", len(lb.apiKeyList))
		}
	})

	t.Run("create load balancer with API keys", func(t *testing.T) {
		keys := []string{"key1", "key2", "key3"}
		lb := NewLoadBalancer(keys)

		if lb == nil {
			t.Fatal("load balancer should not be nil")
		}
		if len(lb.apiKeyList) != 3 {
			t.Errorf("expected API key list length to be 3, got %d", len(lb.apiKeyList))
		}

		// Check initial state of each API key
		for i, apiKey := range lb.apiKeyList {
			if apiKey.Key != keys[i] {
				t.Errorf("expected key to be %s, got %s", keys[i], apiKey.Key)
			}
			if apiKey.Times != 0 {
				t.Errorf("expected usage count to be 0, got %d", apiKey.Times)
			}
			if !apiKey.Available {
				t.Error("expected API key to be available")
			}
			if apiKey.Weight != 1 {
				t.Errorf("expected weight to be 1, got %d", apiKey.Weight)
			}
		}
	})
}

// TestGetAPIKey tests getting API key
func TestGetAPIKey(t *testing.T) {
	t.Run("get API key from empty list", func(t *testing.T) {
		lb := NewLoadBalancer([]string{})
		apiKey, err := lb.GetAPIKey()

		if err != errEmptyAPIKeyList {
			t.Errorf("expected error %v, got %v", errEmptyAPIKeyList, err)
		}
		if apiKey != nil {
			t.Error("expected API key to be nil")
		}
	})

	t.Run("all API keys unavailable", func(t *testing.T) {
		lb := NewLoadBalancer([]string{"key1", "key2"})
		lb.SetAvailabilityForAll(false)

		apiKey, err := lb.GetAPIKey()
		if err != errNoAPIKeyAvailable {
			t.Errorf("expected error %v, got %v", errNoAPIKeyAvailable, err)
		}
		if apiKey != nil {
			t.Error("expected API key to be nil")
		}
	})

	t.Run("least connections algorithm test", func(t *testing.T) {
		lb := NewLoadBalancer([]string{"key1", "key2", "key3"})

		// First get should return any key (all have 0 usage)
		apiKey1, err := lb.GetAPIKey()
		if err != nil {
			t.Fatalf("failed to get API key: %v", err)
		}
		if apiKey1 == nil {
			t.Fatal("API key should not be nil")
		}

		// Second get should return a different key (first one has 1 usage)
		apiKey2, err := lb.GetAPIKey()
		if err != nil {
			t.Fatalf("failed to get API key: %v", err)
		}

		// Verify usage counts
		if apiKey1.Times != 1 {
			t.Errorf("expected first key usage count to be 1, got %d", apiKey1.Times)
		}
		if apiKey2.Times != 1 {
			t.Errorf("expected second key usage count to be 1, got %d", apiKey2.Times)
		}
	})

	t.Run("weight affects selection", func(t *testing.T) {
		lb := NewLoadBalancer([]string{"key1", "key2"})

		// Set different weights
		lb.SetWeight("key1", 1)
		lb.SetWeight("key2", 2)

		// Get multiple times to verify weight effect
		counts := make(map[string]int)
		for range 10 {
			apiKey, err := lb.GetAPIKey()
			if err != nil {
				t.Fatalf("failed to get API key: %v", err)
			}
			counts[apiKey.Key]++
		}

		// key2 has higher weight, should be selected more often
		if counts["key2"] <= counts["key1"] {
			t.Errorf("key2 with higher weight should be selected more often, key1: %d, key2: %d", counts["key1"], counts["key2"])
		}

		var wg sync.WaitGroup
		for range 10 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				lb.GetAPIKey()
			}()
		}
		wg.Wait()
		stats := lb.GetStats()
		if stats["total_requests"].(uint32) != 20 {
			t.Errorf("expected total request count to be 20, got %d", stats["total_requests"])
		}
		if stats["available_api_key"].(int) != 2 {
			t.Errorf("expected available API key count to be 2, got %d", stats["available_api_key"])
		}
		if stats["total_api_key"].(int) != 2 {
			t.Errorf("expected total API key count to be 2, got %d", stats["total_api_key"])
		}
	})
}

// TestSetAvailability tests setting API key availability
func TestSetAvailability(t *testing.T) {
	t.Run("set availability for existing API key", func(t *testing.T) {
		lb := NewLoadBalancer([]string{"key1", "key2"})

		err := lb.SetAvailability("key1", false)
		if err != nil {
			t.Errorf("failed to set availability: %v", err)
		}

		// Verify the setting takes effect
		apiKeyList := lb.GetAPIKeyList()
		for _, apiKey := range apiKeyList {
			if apiKey.Key == "key1" && apiKey.Available {
				t.Error("expected key1 to be unavailable")
			}
			if apiKey.Key == "key2" && !apiKey.Available {
				t.Error("expected key2 to be available")
			}
		}
	})

	t.Run("set availability for non-existent API key", func(t *testing.T) {
		lb := NewLoadBalancer([]string{"key1"})

		err := lb.SetAvailability("nonexistent", false)
		if err != errAPIKeyNotFound {
			t.Errorf("expected error %v, got %v", errAPIKeyNotFound, err)
		}
	})
}

// TestRegisterAPIKey tests registering API key
func TestRegisterAPIKey(t *testing.T) {
	t.Run("register new API key", func(t *testing.T) {
		lb := NewLoadBalancer([]string{"key1"})

		err := lb.RegisterAPIKey("key2")
		if err != nil {
			t.Errorf("failed to register API key: %v", err)
		}

		// Verify registration success
		apiKeyList := lb.GetAPIKeyList()
		if len(apiKeyList) != 2 {
			t.Errorf("expected API key list length to be 2, got %d", len(apiKeyList))
		}

		found := false
		for _, apiKey := range apiKeyList {
			if apiKey.Key == "key2" {
				found = true
				if !apiKey.Available {
					t.Error("newly registered API key should be available")
				}
				if apiKey.Weight != 1 {
					t.Errorf("newly registered API key weight should be 1, got %d", apiKey.Weight)
				}
				break
			}
		}
		if !found {
			t.Error("newly registered API key not found")
		}
	})

	t.Run("register existing API key", func(t *testing.T) {
		lb := NewLoadBalancer([]string{"key1"})

		err := lb.RegisterAPIKey("key1")
		if err != errAPIKeyAlreadyExists {
			t.Errorf("expected error %v, got %v", errAPIKeyAlreadyExists, err)
		}
	})

	t.Run("register API key in empty list", func(t *testing.T) {
		lb := NewLoadBalancer([]string{})

		err := lb.RegisterAPIKey("key1")
		if err != nil {
			t.Errorf("failed to register API key: %v", err)
		}

		apiKeyList := lb.GetAPIKeyList()
		if len(apiKeyList) != 1 {
			t.Errorf("expected API key list length to be 1, got %d", len(apiKeyList))
		}
	})
}

// TestUnregisterAPIKey tests unregistering API key
func TestUnregisterAPIKey(t *testing.T) {
	t.Run("unregister existing API key", func(t *testing.T) {
		lb := NewLoadBalancer([]string{"key1", "key2", "key3"})

		err := lb.UnregisterAPIKey("key2")
		if err != nil {
			t.Errorf("failed to unregister API key: %v", err)
		}

		// Verify unregistration success
		apiKeyList := lb.GetAPIKeyList()
		if len(apiKeyList) != 2 {
			t.Errorf("expected API key list length to be 2, got %d", len(apiKeyList))
		}

		for _, apiKey := range apiKeyList {
			if apiKey.Key == "key2" {
				t.Error("unregistered API key should not exist")
			}
		}
	})

	t.Run("unregister non-existent API key", func(t *testing.T) {
		lb := NewLoadBalancer([]string{"key1"})

		err := lb.UnregisterAPIKey("nonexistent")
		if err != errAPIKeyNotFound {
			t.Errorf("expected error %v, got %v", errAPIKeyNotFound, err)
		}
	})
}

// TestSetWeight tests setting weight
func TestSetWeight(t *testing.T) {
	t.Run("set valid weight", func(t *testing.T) {
		lb := NewLoadBalancer([]string{"key1"})

		err := lb.SetWeight("key1", 5)
		if err != nil {
			t.Errorf("failed to set weight: %v", err)
		}

		// Verify weight setting success
		apiKeyList := lb.GetAPIKeyList()
		if apiKeyList[0].Weight != 5 {
			t.Errorf("expected weight to be 5, got %d", apiKeyList[0].Weight)
		}
	})

	t.Run("set weight to zero", func(t *testing.T) {
		lb := NewLoadBalancer([]string{"key1"})

		err := lb.SetWeight("key1", 0)
		if err != errWeightMustBeGreaterThan0 {
			t.Errorf("expected error %v, got %v", errWeightMustBeGreaterThan0, err)
		}
	})

	t.Run("set weight for non-existent API key", func(t *testing.T) {
		lb := NewLoadBalancer([]string{"key1"})

		err := lb.SetWeight("nonexistent", 5)
		if err != errAPIKeyNotFound {
			t.Errorf("expected error %v, got %v", errAPIKeyNotFound, err)
		}
	})
}

// TestSetAvailabilityForAll tests setting availability for all API keys
func TestSetAvailabilityForAll(t *testing.T) {
	t.Run("set all API keys unavailable", func(t *testing.T) {
		lb := NewLoadBalancer([]string{"key1", "key2", "key3"})

		lb.SetAvailabilityForAll(false)

		// Verify all API keys are unavailable
		apiKeyList := lb.GetAPIKeyList()
		for _, apiKey := range apiKeyList {
			if apiKey.Available {
				t.Errorf("API key %s should be unavailable", apiKey.Key)
			}
		}
	})

	t.Run("set all API keys available", func(t *testing.T) {
		lb := NewLoadBalancer([]string{"key1", "key2"})
		lb.SetAvailabilityForAll(false) // First set to unavailable
		lb.SetAvailabilityForAll(true)  // Then set to available

		// Verify all API keys are available
		apiKeyList := lb.GetAPIKeyList()
		for _, apiKey := range apiKeyList {
			if !apiKey.Available {
				t.Errorf("API key %s should be available", apiKey.Key)
			}
		}
	})
}

// TestGetAPIKeyList tests getting API key list
func TestGetAPIKeyList(t *testing.T) {
	t.Run("get API key list copy", func(t *testing.T) {
		lb := NewLoadBalancer([]string{"key1", "key2"})

		// Modify some state
		lb.GetAPIKey() // Increase usage count
		lb.SetWeight("key1", 3)
		lb.SetAvailability("key2", false)

		apiKeyList := lb.GetAPIKeyList()

		// Verify returned is a copy
		if len(apiKeyList) != 2 {
			t.Errorf("expected API key list length to be 2, got %d", len(apiKeyList))
		}

		// Modify returned list, should not affect original data
		apiKeyList[0].Times = 999
		apiKeyList[0].Available = false

		// Get again, verify original data is not modified
		newList := lb.GetAPIKeyList()
		if newList[0].Times == 999 {
			t.Error("returned should be a copy, modification should not affect original data")
		}
	})
}

// TestGetStats tests getting statistics
func TestGetStats(t *testing.T) {
	t.Run("get statistics", func(t *testing.T) {
		lb := NewLoadBalancer([]string{"key1", "key2", "key3"})

		// Simulate some usage
		lb.GetAPIKey()
		lb.GetAPIKey()
		lb.SetAvailability("key3", false)

		stats := lb.GetStats()

		if stats["total_api_key"] != 3 {
			t.Errorf("expected total API key count to be 3, got %v", stats["total_api_key"])
		}
		if stats["available_api_key"] != 2 {
			t.Errorf("expected available API key count to be 2, got %v", stats["available_api_key"])
		}
		if stats["total_requests"] != uint32(2) {
			t.Errorf("expected total request count to be 2, got %v", stats["total_requests"])
		}
	})

	t.Run("empty list statistics", func(t *testing.T) {
		lb := NewLoadBalancer([]string{})

		stats := lb.GetStats()

		if stats["total_api_key"] != 0 {
			t.Errorf("expected total API key count to be 0, got %v", stats["total_api_key"])
		}
		if stats["available_api_key"] != 0 {
			t.Errorf("expected available API key count to be 0, got %v", stats["available_api_key"])
		}
		if stats["total_requests"] != uint32(0) {
			t.Errorf("expected total request count to be 0, got %v", stats["total_requests"])
		}
	})
}

// TestConcurrency tests concurrency safety
func TestConcurrency(t *testing.T) {
	t.Run("concurrent get API key", func(t *testing.T) {
		lb := NewLoadBalancer([]string{"key1", "key2", "key3"})

		var wg sync.WaitGroup
		numGoroutines := 100
		numRequests := 10

		for range numGoroutines {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for range numRequests {
					_, err := lb.GetAPIKey()
					if err != nil {
						t.Errorf("concurrent get API key failed: %v", err)
					}
				}
			}()
		}

		wg.Wait()

		// Verify total request count
		stats := lb.GetStats()
		expectedRequests := uint32(numGoroutines * numRequests)
		if stats["total_requests"] != expectedRequests {
			t.Errorf("expected total request count to be %d, got %v", expectedRequests, stats["total_requests"])
		}
	})

	t.Run("concurrent modify API key state", func(t *testing.T) {
		lb := NewLoadBalancer([]string{"key1", "key2", "key3"})

		var wg sync.WaitGroup
		numGoroutines := 50

		// Concurrent set availability
		for i := range numGoroutines {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				available := index%2 == 0
				lb.SetAvailability("key1", available)
			}(i)
		}

		// Concurrent register and unregister
		for i := range numGoroutines {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				keyName := "temp_key_" + string(rune(index))
				lb.RegisterAPIKey(keyName)
				lb.UnregisterAPIKey(keyName)
			}(i)
		}

		wg.Wait()

		// Verify original keys still exist
		apiKeyList := lb.GetAPIKeyList()
		originalKeys := []string{"key1", "key2", "key3"}
		for _, originalKey := range originalKeys {
			found := false
			for _, apiKey := range apiKeyList {
				if apiKey.Key == originalKey {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("original key %s should still exist", originalKey)
			}
		}
	})
}

// BenchmarkGetAPIKey benchmark test for getting API key
func BenchmarkGetAPIKey(b *testing.B) {
	lb := NewLoadBalancer([]string{"key1", "key2", "key3", "key4", "key5"})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := lb.GetAPIKey()
			if err != nil {
				b.Errorf("failed to get API key: %v", err)
			}
		}
	})
}

// BenchmarkSetAvailability benchmark test for setting availability
func BenchmarkSetAvailability(b *testing.B) {
	lb := NewLoadBalancer([]string{"key1", "key2", "key3", "key4", "key5"})

	for i := 0; b.Loop(); i++ {
		key := "key" + string(rune(i%5+1))
		available := i%2 == 0
		lb.SetAvailability(key, available)
	}
}
