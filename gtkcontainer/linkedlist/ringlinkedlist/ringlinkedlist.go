/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-04-07 19:22:14
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-04-07 23:14:40
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package ringlinkedlist

import "github.com/liusuxian/go-toolkit/gtklocker"

// Node 链表节点
type Node struct {
	value int64
	next  *Node // 指向下一个节点
	prev  *Node // 指向上一个节点
}

// LinkedList 环形链表
type LinkedList struct {
	head    *Node            // 链表头
	tail    *Node            // 链表尾
	lock    gtklocker.Locker // 用于确保并发安全的锁接口
	current *Node            // 用于操作如轮询的当前节点指针
	nodes   map[int64]*Node  // 通过值快速查找节点的映射
}

// Iterator 环形链表迭代器
type Iterator struct {
	start    *Node // 迭代开始的节点
	current  *Node // 当前迭代的节点
	finished bool  // 是否完成迭代
}

// NewLinkedList 新建环形链表
func NewLinkedList() (list *LinkedList) {
	return &LinkedList{
		lock:  &gtklocker.MutexLocker{},
		nodes: make(map[int64]*Node),
	}
}

// Append 追加
func (list *LinkedList) Append(value int64) {
	list.lock.Lock()
	defer list.lock.Unlock()

	// 检查待添加的值是否已存在于链表中，存在则直接返回，不添加重复的值
	if _, exists := list.nodes[value]; exists {
		return
	}
	// 创建一个新的节点
	newNode := &Node{value: value}
	if list.head == nil {
		// 列表中的第一个节点
		list.head = newNode
		list.tail = newNode
		newNode.next = newNode // 环形链表自指向，形成一个环
		newNode.prev = newNode // 自指向，形成一个完整的环
	} else {
		// 将新节点追加到列表中
		newNode.prev = list.tail // 设置新节点的`prev`为当前的尾节点
		newNode.next = list.head // 设置新节点的`next`为头节点，因为是环形链表
		list.tail.next = newNode // 将当前尾节点的`next`指向新节点
		list.head.prev = newNode // 更新头节点的`prev`为新节点
		list.tail = newNode      // 更新链表的尾节点为新节点
	}
	// 将新节点添加到`nodes`映射中，以便于快速查找
	list.nodes[value] = newNode
	if list.current == nil {
		list.current = list.head
	}
}

// Insert 在链表中找到`value`值的节点，并在其前面或后面插入新值`newValue`
func (list *LinkedList) Insert(value, newValue int64, before bool) (result bool) {
	list.lock.Lock()
	defer list.lock.Unlock()

	// 检查`newValue`是否已存在
	if _, exists := list.nodes[newValue]; exists {
		return false
	}
	// 查找`value`值的节点
	node, exists := list.nodes[value]
	if !exists {
		return false
	}

	// 创建新节点
	newNode := &Node{value: newValue}
	if before {
		// 插入到找到的节点前面
		newNode.prev = node.prev
		newNode.next = node
		node.prev.next = newNode
		node.prev = newNode
	} else {
		// 插入到找到的节点后面
		newNode.next = node.next
		newNode.prev = node
		node.next.prev = newNode
		node.next = newNode
	}
	// 更新头尾节点
	if node == list.head && before {
		list.head = newNode
	}
	if node == list.tail && !before {
		list.tail = newNode
	}
	// 将新节点加入到`nodes`映射中
	list.nodes[newValue] = newNode
	if list.current == nil {
		list.current = list.head
	}
	return true
}

// Remove 移除
func (list *LinkedList) Remove(value int64) {
	list.lock.Lock()
	defer list.lock.Unlock()

	// 尝试从`nodes`映射中找到要删除的节点
	node, exists := list.nodes[value]
	if !exists {
		return
	}

	if node == node.next {
		// 处理只有一个节点的情况
		list.head = nil
		list.tail = nil
		list.current = nil
	} else {
		// 处理链表中有多个节点的情况，从链表中移除目标节点
		node.prev.next = node.next
		node.next.prev = node.prev
		// 更新头节点
		if node == list.head {
			list.head = node.next
		}
		// 更新尾节点
		if node == list.tail {
			list.tail = node.prev
		}
		// 如果当前节点是要删除的节点，则将`current`更新为下一个节点
		if node == list.current {
			list.current = node.next
		}
	}
	// 从`nodes`映射中删除节点，完成整个删除操作
	delete(list.nodes, value)
}

// Poll 轮询，返回当前节点并将当前指针移到下一个节点
func (list *LinkedList) Poll() (value int64, ok bool) {
	list.lock.Lock()
	defer list.lock.Unlock()

	if list.current == nil && list.head == nil {
		return
	}
	value = list.current.value
	ok = true
	list.current = list.current.next
	return
}

// NewIterator 新建环形链表迭代器（非线程安全）
func (list *LinkedList) NewIterator() (it *Iterator) {
	return &Iterator{
		start:    list.head,
		current:  nil,
		finished: false,
	}
}

// Next 返回迭代器的下一个元素（非线程安全）
func (it *Iterator) Next() (value int64, ok bool) {
	if it.finished || (it.current != nil && it.current.next == it.start) {
		// 如果已完成迭代，或者回到了开始的节点，则结束迭代
		it.finished = true
		return
	}

	if it.current == nil {
		it.current = it.start
	} else {
		it.current = it.current.next
	}

	if it.current != nil {
		return it.current.value, true
	}
	return
}
