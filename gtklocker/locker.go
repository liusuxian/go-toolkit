/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-04-07 19:48:18
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-04-07 22:43:39
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtklocker

import "sync"

// Locker 锁接口
type Locker interface {
	Lock()
	RLock()
	TryLock() (ok bool)
	TryRLock() (ok bool)
	Unlock()
	RUnlock()
}

// MutexLocker 互斥锁
type MutexLocker struct {
	mu sync.Mutex
}

// Lock 加锁
func (m *MutexLocker) Lock() {
	m.mu.Lock()
}

// RLock 加读锁
func (m *MutexLocker) RLock() {
	m.mu.Lock() // 对于互斥锁，读锁行为与普通锁相同
}

// TryLock 尝试加锁
func (m *MutexLocker) TryLock() (ok bool) {
	return m.mu.TryLock()
}

// TryRLock 尝试加读锁
func (m *MutexLocker) TryRLock() (ok bool) {
	return m.mu.TryLock() // 对于互斥锁，读锁行为与普通锁相同
}

// Unlock 释放锁
func (m *MutexLocker) Unlock() {
	m.mu.Unlock()
}

// RUnlock 释放读锁
func (m *MutexLocker) RUnlock() {
	m.mu.Unlock() // 对于互斥锁，读锁行为与普通锁相同
}

// RWMutexLocker 读写锁
type RWMutexLocker struct {
	mu sync.RWMutex
}

// Lock 加锁
func (m *RWMutexLocker) Lock() {
	m.mu.Lock()
}

// RLock 加读锁
func (m *RWMutexLocker) RLock() {
	m.mu.RLock()
}

// TryLock 尝试加锁
func (m *RWMutexLocker) TryLock() (ok bool) {
	return m.mu.TryLock()
}

// TryRLock 尝试加读锁
func (m *RWMutexLocker) TryRLock() (ok bool) {
	return m.mu.TryRLock()
}

// Unlock 释放锁
func (m *RWMutexLocker) Unlock() {
	m.mu.Unlock()
}

// RUnlock 释放读锁
func (m *RWMutexLocker) RUnlock() {
	m.mu.RUnlock()
}
