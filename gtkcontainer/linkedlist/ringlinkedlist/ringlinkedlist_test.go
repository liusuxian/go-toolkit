/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-04-07 19:22:14
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-04-07 23:13:50
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package ringlinkedlist_test

import (
	"github.com/liusuxian/go-toolkit/gtkcontainer/linkedlist/ringlinkedlist"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLinkedList(t *testing.T) {
	var (
		assert       = assert.New(t)
		index  int64 = 0
		list         = ringlinkedlist.NewLinkedList()
	)

	for i := int64(0); i < 10; i++ {
		list.Append(i)
	}
	list.Append(9)
	iterator := list.NewIterator()
	for value, ok := iterator.Next(); ok; value, ok = iterator.Next() {
		assert.Equal(index, value)
		index++
	}

	index = -2
	result := list.Insert(0, -1, true)
	assert.True(result)
	result = list.Insert(-1, -2, true)
	assert.True(result)
	result = list.Insert(9, 10, false)
	assert.True(result)
	result = list.Insert(9, 10, false)
	assert.False(result)
	iterator = list.NewIterator()
	for value, ok := iterator.Next(); ok; value, ok = iterator.Next() {
		assert.Equal(index, value)
		index++
	}

	index = 0
	list.Remove(-2)
	list.Remove(-1)
	list.Remove(10)
	iterator = list.NewIterator()
	for value, ok := iterator.Next(); ok; value, ok = iterator.Next() {
		assert.Equal(index, value)
		index++
	}

	index = 0
	for i := 0; i < 10; i++ {
		value, ok := list.Poll()
		assert.True(ok)
		assert.Equal(index, value)
		index++
	}
	index = 0
	for i := 0; i < 10; i++ {
		value, ok := list.Poll()
		assert.True(ok)
		assert.Equal(index, value)
		index++
	}

	for i := int64(0); i < 9; i++ {
		list.Remove(i)
	}
	for i := 0; i < 10; i++ {
		value, ok := list.Poll()
		assert.True(ok)
		assert.Equal(int64(9), value)
	}

	for i := int64(0); i < 10; i++ {
		list.Remove(i)
	}
	value, ok := list.Poll()
	assert.False(ok)
	assert.Equal(int64(0), value)
}
