/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-04-10 13:10:33
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-13 13:48:49
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package doubly

import "fmt"

// doPushFront 将一个或多个节点插入链表头部，如果 uuid 已存在则更新节点值
func (list *LinkedList) doPushFront(items ...Node) {
	for _, item := range items {
		if node, exists := list.nodes[item.Uuid]; exists {
			// uuid 已存在，更新节点值
			node.Value = item.Value
		} else {
			newNode := &Node{Uuid: item.Uuid, Value: item.Value}
			if list.length == 0 {
				list.head = newNode
				list.tail = newNode
			} else {
				newNode.next = list.head // 新节点的 next 指向当前头节点
				list.head.prev = newNode // 当前头节点的 prev 指向新节点
				list.head = newNode      // 更新链表的头节点为新节点
			}

			list.nodes[item.Uuid] = newNode
			list.length++
		}
	}

	if list.current == nil {
		list.current = list.head
	}
}

// doPushBack 将一个或多个节点插入链表尾部，如果 uuid 已存在则更新节点值
func (list *LinkedList) doPushBack(items ...Node) {
	for _, item := range items {
		if node, exists := list.nodes[item.Uuid]; exists {
			// uuid 已存在，更新节点值
			node.Value = item.Value
		} else {
			newNode := &Node{Uuid: item.Uuid, Value: item.Value}
			if list.length == 0 {
				list.head = newNode
				list.tail = newNode
			} else {
				list.tail.next = newNode // 当前尾节点的 next 指向新节点
				newNode.prev = list.tail // 新节点的 prev 指向当前尾节点
				list.tail = newNode      // 更新链表的尾节点为新节点
			}

			list.nodes[item.Uuid] = newNode
			list.length++
		}
	}

	if list.current == nil {
		list.current = list.head
	}
}

// doInsertBefore 在指定节点前插入一个或多个节点，如果指定节点不存在则返回错误，如果 uuid 已存在则更新节点值
func (list *LinkedList) doInsertBefore(targetUuid string, items ...Node) (err error) {
	var (
		targetNode   *Node
		targetExists bool
	)
	if targetNode, targetExists = list.nodes[targetUuid]; !targetExists {
		err = fmt.Errorf("target node [%s] not found", targetUuid)
		return
	}

	lastNode := targetNode
	for _, item := range items {
		if node, exists := list.nodes[item.Uuid]; exists {
			// uuid 已存在，更新节点值
			node.Value = item.Value
			lastNode = node
		} else {
			newNode := &Node{Uuid: item.Uuid, Value: item.Value}
			newNode.prev = lastNode.prev
			newNode.next = lastNode
			if lastNode.prev == nil {
				list.head = newNode
			} else {
				lastNode.prev.next = newNode
			}
			lastNode.prev = newNode
			lastNode = newNode

			list.nodes[item.Uuid] = newNode
			list.length++
		}
	}

	if list.current == nil {
		list.current = list.head
	}
	return
}

// doInsertAfter 在指定节点后插入一个或多个节点，如果指定节点不存在则返回错误，如果 uuid 已存在则更新节点值
func (list *LinkedList) doInsertAfter(targetUuid string, items ...Node) (err error) {
	var (
		targetNode   *Node
		targetExists bool
	)
	if targetNode, targetExists = list.nodes[targetUuid]; !targetExists {
		err = fmt.Errorf("target node [%s] not found", targetUuid)
		return
	}

	lastNode := targetNode
	for _, item := range items {
		if node, exists := list.nodes[item.Uuid]; exists {
			// uuid 已存在，更新节点值
			node.Value = item.Value
			lastNode = node
		} else {
			newNode := &Node{Uuid: item.Uuid, Value: item.Value}
			newNode.next = lastNode.next
			newNode.prev = lastNode
			if lastNode.next == nil {
				list.tail = newNode
			} else {
				lastNode.next.prev = newNode
			}
			lastNode.next = newNode
			lastNode = newNode

			list.nodes[item.Uuid] = newNode
			list.length++
		}
	}

	if list.current == nil {
		list.current = list.head
	}
	return
}

// doRemove 移除一个或多个节点
func (list *LinkedList) doRemove(uuids ...string) {
	for _, uuid := range uuids {
		node, exists := list.nodes[uuid]
		if !exists {
			continue
		}

		if list.length == 1 {
			list.head = nil
			list.tail = nil
			list.current = nil
		} else {
			if node == list.head {
				list.head = node.next
				list.head.prev = nil
			} else if node == list.tail {
				list.tail = node.prev
				list.tail.next = nil
			} else {
				node.prev.next = node.next
				node.next.prev = node.prev
			}

			// 如果当前节点是要删除的节点，则将 current 更新为下一个节点
			if node == list.current {
				list.current = node.next
			}
		}

		list.length--
		delete(list.nodes, uuid)
	}
}
