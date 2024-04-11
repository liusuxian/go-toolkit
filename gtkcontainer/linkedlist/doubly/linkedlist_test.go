/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-04-08 14:00:12
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-04-11 15:33:07
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package doubly_test

import (
	"github.com/liusuxian/go-toolkit/gtkcontainer/linkedlist/doubly"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestLinkedList_PushFront(t *testing.T) {
	var (
		assert = assert.New(t)
		index  = 0
		list   = doubly.NewLinkedList()
	)

	assert.Equal(0, list.Len())
	list.PushFront(doubly.Node{Uuid: "10", Value: 10}, doubly.Node{Uuid: "10", Value: 100}, doubly.Node{Uuid: "9", Value: 9}, doubly.Node{Uuid: "8", Value: 8})
	for i := 7; i >= 0; i-- {
		list.PushFront(doubly.Node{Uuid: strconv.Itoa(i), Value: i})
	}
	list.PushFront(doubly.Node{Uuid: "9", Value: 100})
	assert.Equal(11, list.Len())

	iterator := list.NewIterator()
	for node := iterator.Next(); node != nil; node = iterator.Next() {
		assert.Equal(strconv.Itoa(index), node.Uuid)
		t.Logf("node: %+v", node)
		index++
	}
}

func TestLinkedList_PushBack(t *testing.T) {
	var (
		assert = assert.New(t)
		index  = 0
		list   = doubly.NewLinkedList()
	)

	assert.Equal(0, list.Len())
	list.PushBack(doubly.Node{Uuid: "0", Value: 0}, doubly.Node{Uuid: "0", Value: 100}, doubly.Node{Uuid: "1", Value: 1}, doubly.Node{Uuid: "2", Value: 2})
	for i := 3; i < 10; i++ {
		list.PushBack(doubly.Node{Uuid: strconv.Itoa(i), Value: i})
	}
	list.PushBack(doubly.Node{Uuid: "9", Value: 100})
	assert.Equal(10, list.Len())

	iterator := list.NewIterator()
	for node := iterator.Next(); node != nil; node = iterator.Next() {
		assert.Equal(strconv.Itoa(index), node.Uuid)
		t.Logf("node: %+v", node)
		index++
	}
}

func TestLinkedList_InsertBefore(t *testing.T) {
	var (
		assert = assert.New(t)
		index  = 0
		list   = doubly.NewLinkedList()
		err    error
	)

	assert.Equal(0, list.Len())
	err = list.InsertBefore("100", doubly.Node{Uuid: "9", Value: 9}, doubly.Node{Uuid: "9", Value: 100}, doubly.Node{Uuid: "8", Value: 8}, doubly.Node{Uuid: "7", Value: 7})
	assert.NoError(err)
	assert.Equal(3, list.Len())
	err = list.InsertBefore("100", doubly.Node{Uuid: "9", Value: 200})
	assert.Error(err)

	err = list.InsertBefore("7", []doubly.Node{
		{Uuid: "6", Value: 6},
		{Uuid: "6", Value: 100},
		{Uuid: "5", Value: 5},
		{Uuid: "4", Value: 4},
		{Uuid: "3", Value: 3},
		{Uuid: "2", Value: 2},
		{Uuid: "1", Value: 1},
		{Uuid: "0", Value: 0},
	}...)
	assert.NoError(err)
	assert.Equal(10, list.Len())

	iterator := list.NewIterator()
	for node := iterator.Next(); node != nil; node = iterator.Next() {
		assert.Equal(strconv.Itoa(index), node.Uuid)
		t.Logf("node: %+v", node)
		index++
	}
}

func TestLinkedList_InsertAfter(t *testing.T) {
	var (
		assert = assert.New(t)
		index  = 0
		list   = doubly.NewLinkedList()
		err    error
	)

	assert.Equal(0, list.Len())
	err = list.InsertAfter("100", doubly.Node{Uuid: "0", Value: 0}, doubly.Node{Uuid: "0", Value: 100}, doubly.Node{Uuid: "1", Value: 1}, doubly.Node{Uuid: "2", Value: 2})
	assert.NoError(err)
	assert.Equal(3, list.Len())
	err = list.InsertAfter("100", doubly.Node{Uuid: "0", Value: 200})
	assert.Error(err)

	err = list.InsertAfter("2", []doubly.Node{
		{Uuid: "3", Value: 3},
		{Uuid: "4", Value: 4},
		{Uuid: "5", Value: 5},
		{Uuid: "6", Value: 6},
		{Uuid: "7", Value: 7},
		{Uuid: "8", Value: 8},
		{Uuid: "9", Value: 9},
		{Uuid: "9", Value: 100},
	}...)
	assert.NoError(err)
	assert.Equal(10, list.Len())

	iterator := list.NewIterator()
	for node := iterator.Next(); node != nil; node = iterator.Next() {
		assert.Equal(strconv.Itoa(index), node.Uuid)
		t.Logf("node: %+v", node)
		index++
	}
}

func TestLinkedList_Remove(t *testing.T) {
	var (
		assert = assert.New(t)
		index  = 0
		list   = doubly.NewLinkedList()
	)

	assert.Equal(0, list.Len())
	list.PushBack([]doubly.Node{
		{Uuid: "0", Value: 0},
		{Uuid: "1", Value: 1},
		{Uuid: "2", Value: 2},
		{Uuid: "3", Value: 3},
		{Uuid: "4", Value: 4},
		{Uuid: "5", Value: 5},
		{Uuid: "6", Value: 6},
		{Uuid: "7", Value: 7},
		{Uuid: "8", Value: 8},
		{Uuid: "9", Value: 9},
	}...)
	assert.Equal(10, list.Len())

	iterator := list.NewIterator()
	for node := iterator.Next(); node != nil; node = iterator.Next() {
		assert.Equal(strconv.Itoa(index), node.Uuid)
		index++
	}

	list.Remove("3", "5")
	index = 0
	iterator = list.NewIterator()
	for node := iterator.Next(); node != nil; node = iterator.Next() {
		assert.Equal(strconv.Itoa(index), node.Uuid)
		t.Logf("node: %+v", node)
		if index == 2 || index == 4 {
			index += 2
		} else {
			index++
		}
	}
}

func TestLinkedList_Poll(t *testing.T) {
	var (
		assert = assert.New(t)
		list   = doubly.NewLinkedList()
		node   *doubly.Node
		err    error
	)

	assert.Equal(0, list.Len())
	list.PushBack([]doubly.Node{
		{Uuid: "0", Value: 0},
		{Uuid: "1", Value: 1},
		{Uuid: "2", Value: 2},
		{Uuid: "3", Value: 3},
		{Uuid: "4", Value: 4},
		{Uuid: "5", Value: 5},
		{Uuid: "6", Value: 6},
		{Uuid: "7", Value: 7},
		{Uuid: "8", Value: 8},
		{Uuid: "9", Value: 9},
	}...)
	assert.Equal(10, list.Len())

	for i := 0; i < 10; i++ {
		node, err = list.Poll()
		assert.NoError(err)
		assert.Equal(strconv.Itoa(i), node.Uuid)
		t.Logf("node: %+v", node)
	}
	for i := 0; i < 10; i++ {
		node, err = list.Poll()
		assert.NoError(err)
		assert.Equal(strconv.Itoa(i), node.Uuid)
		t.Logf("node: %+v", node)
	}
}

func TestLinkedList_GetCurrentAndMoveToNext(t *testing.T) {
	var (
		assert = assert.New(t)
		list   = doubly.NewLinkedList()
		node   *doubly.Node
		err    error
		value  any
	)

	assert.Equal(0, list.Len())
	list.PushBack([]doubly.Node{
		{Uuid: "0", Value: 0},
		{Uuid: "1", Value: 1},
		{Uuid: "2", Value: 2},
		{Uuid: "3", Value: 3},
		{Uuid: "4", Value: 4},
		{Uuid: "5", Value: 5},
		{Uuid: "6", Value: 6},
		{Uuid: "7", Value: 7},
		{Uuid: "8", Value: 8},
		{Uuid: "9", Value: 9},
	}...)
	assert.Equal(10, list.Len())

	for i := 0; i < 10; i++ {
		node, err = list.GetCurrentAndMoveToNext()
		assert.NoError(err)
		assert.Equal(strconv.Itoa(i), node.Uuid)
		t.Logf("node: %+v", node)
	}
	node, err = list.GetCurrentAndMoveToNext()
	assert.Error(err)
	assert.Nil(node)

	err = list.SetCurrent("5")
	assert.NoError(err)
	value, err = list.GetNodeValue("5")
	assert.NoError(err)
	assert.Equal(5, value)
	list.UpdateNodeValue("5", 100)
	value, err = list.GetNodeValue("5")
	assert.NoError(err)
	assert.Equal(100, value)

	for i := 5; i < 10; i++ {
		node, err := list.GetCurrentAndMoveToNext()
		assert.NoError(err)
		assert.Equal(strconv.Itoa(i), node.Uuid)
		t.Logf("node: %+v", node)
	}
	node, err = list.GetCurrentAndMoveToNext()
	assert.Error(err)
	assert.Nil(node)
}
