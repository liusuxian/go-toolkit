/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-04-08 14:00:12
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-13 13:49:38
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package doubly

import (
	"fmt"
	"sync"
)

// Node 链表节点
type Node struct {
	Uuid  string // 节点唯一标识
	Value any    // 节点存储的值
	next  *Node  // 指向下一个节点
	prev  *Node  // 指向上一个节点
}

// LinkedList 双向链表
type LinkedList struct {
	head    *Node            // 链表头
	tail    *Node            // 链表尾
	lock    *sync.RWMutex    // 用于确保并发安全的锁接口
	current *Node            // 用于操作如轮询的当前节点指针
	nodes   map[string]*Node // 通过值快速查找节点的映射
	length  int              // 链表的长度
}

// Iterator 双向链表迭代器
type Iterator struct {
	start    *Node       // 迭代开始的节点
	current  *Node       // 当前迭代的节点
	finished bool        // 是否完成迭代
	lock     *sync.Mutex // 迭代器的锁
}

// NewLinkedList 新建双向链表
func NewLinkedList() (list *LinkedList) {
	return &LinkedList{
		lock:  &sync.RWMutex{},
		nodes: make(map[string]*Node),
	}
}

// PushFront 将一个或多个节点插入链表头部，如果 uuid 已存在则更新节点值
func (list *LinkedList) PushFront(items ...Node) {
	if len(items) == 0 {
		return
	}

	list.lock.Lock()
	defer list.lock.Unlock()

	list.doPushFront(items...)
}

// PushBack 将一个或多个节点插入链表尾部，如果 uuid 已存在则更新节点值
func (list *LinkedList) PushBack(items ...Node) {
	if len(items) == 0 {
		return
	}

	list.lock.Lock()
	defer list.lock.Unlock()

	list.doPushBack(items...)
}

// InsertBefore 在指定节点前插入一个或多个节点，如果指定节点不存在则返回错误，如果 uuid 已存在则更新节点值
func (list *LinkedList) InsertBefore(targetUuid string, items ...Node) (err error) {
	if len(items) == 0 {
		return
	}

	list.lock.Lock()
	defer list.lock.Unlock()

	// 处理链表为空的情况
	if list.length == 0 {
		list.doPushFront(items...)
		return
	}

	// 处理链表非空的情况
	return list.doInsertBefore(targetUuid, items...)
}

// InsertAfter 在指定节点后插入一个或多个节点，如果指定节点不存在则返回错误，如果 uuid 已存在则更新节点值
func (list *LinkedList) InsertAfter(targetUuid string, items ...Node) (err error) {
	if len(items) == 0 {
		return
	}

	list.lock.Lock()
	defer list.lock.Unlock()

	// 处理链表为空的情况
	if list.length == 0 {
		list.doPushBack(items...)
		return
	}

	// 处理链表非空的情况
	return list.doInsertAfter(targetUuid, items...)
}

// Remove 移除一个或多个节点
func (list *LinkedList) Remove(uuids ...string) {
	if len(uuids) == 0 {
		return
	}

	list.lock.Lock()
	defer list.lock.Unlock()

	list.doRemove(uuids...)
}

// Len 链表的当前长度
func (list *LinkedList) Len() (length int) {
	list.lock.RLock()
	defer list.lock.RUnlock()

	return list.length
}

// SetCurrent 设置当前节点
func (list *LinkedList) SetCurrent(uuid string) (err error) {
	list.lock.Lock()
	defer list.lock.Unlock()

	if node, exists := list.nodes[uuid]; exists {
		list.current = node
		return nil
	}
	return fmt.Errorf("node [%s] not found", uuid)
}

// GetNodeValue 获取节点的值
func (list *LinkedList) GetNodeValue(uuid string) (value any, err error) {
	list.lock.Lock()
	defer list.lock.Unlock()

	if node, exists := list.nodes[uuid]; exists {
		return node.Value, nil
	}
	return nil, fmt.Errorf("node [%s] not found", uuid)
}

// UpdateNodeValue 更新节点的值
func (list *LinkedList) UpdateNodeValue(uuid string, value any) (err error) {
	list.lock.Lock()
	defer list.lock.Unlock()

	if node, exists := list.nodes[uuid]; exists {
		node.Value = value
		return nil
	}
	return fmt.Errorf("node [%s] not found", uuid)
}

// Poll 获取当前节点并移动到下一个节点，当到达链表末尾时，移动到链表头部
func (list *LinkedList) Poll() (node *Node, err error) {
	list.lock.Lock()
	defer list.lock.Unlock()

	// 空链表
	if list.length == 0 {
		return nil, fmt.Errorf("empty linked list")
	}

	if list.current == nil {
		return nil, fmt.Errorf("current node is nil")
	}

	node = &Node{}
	*node = *list.current

	if list.current.next == nil {
		list.current = list.head
	} else {
		list.current = list.current.next
	}

	return node, nil
}

// GetCurrentAndMoveToNext 获取当前节点并移动到下一个节点，当到达链表末尾时，停止移动
func (list *LinkedList) GetCurrentAndMoveToNext() (node *Node, err error) {
	list.lock.Lock()
	defer list.lock.Unlock()

	// 空链表
	if list.length == 0 {
		return nil, fmt.Errorf("empty linked list")
	}

	if list.current == nil {
		return nil, fmt.Errorf("current node is nil")
	}

	node = &Node{}
	*node = *list.current

	list.current = list.current.next
	return node, nil
}

// NewIterator 新建环形链表迭代器
func (list *LinkedList) NewIterator() (it *Iterator) {
	list.lock.RLock()
	defer list.lock.RUnlock()

	return &Iterator{
		start:    list.head,
		current:  nil,
		finished: false,
		lock:     &sync.Mutex{},
	}
}

// Next 返回迭代器的下一个元素
func (it *Iterator) Next() (node *Node) {
	it.lock.Lock()
	defer it.lock.Unlock()

	// 如果迭代已经完成或链表为空，则直接返回
	if it.finished || it.start == nil {
		return nil
	}

	if it.current == nil {
		it.current = it.start
	} else {
		it.current = it.current.next
	}

	// 检查是否到达链表末尾或已经遍历完成
	if it.current == nil {
		it.finished = true
		return nil
	}
	return it.current
}
