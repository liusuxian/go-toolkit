/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-04-01 13:15:12
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-13 14:32:50
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtktask

import (
	"fmt"
	"github.com/liusuxian/go-toolkit/gtkcontainer/linkedlist/doubly"
	"github.com/orcaman/concurrent-map/v2"
	"strconv"
	"time"
)

// PollInfo 用于存储轮询信息的结构体
type PollInfo struct {
	availableList   *doubly.LinkedList                    // 可用的轮询对象标识列表
	unavailableMap  cmap.ConcurrentMap[string, time.Time] // 不可用的轮询对象标识映射
	unavailableTime time.Duration                         // 不可用的轮询对象的冷却时间
	interval        time.Duration                         // 可用性检测时间间隔
	quitChan        chan bool                             // 退出信号
}

// NewPoll 新建轮询对象
func NewPoll(unavailableTime, interval time.Duration) (p *PollInfo) {
	if unavailableTime <= 0 {
		unavailableTime = time.Minute * 10
	}
	if interval <= 0 {
		interval = time.Minute * 10
	}
	p = &PollInfo{
		availableList:   doubly.NewLinkedList(),
		unavailableMap:  cmap.New[time.Time](),
		unavailableTime: unavailableTime,
		interval:        interval,
		quitChan:        make(chan bool),
	}
	// 启动可用性检测
	go p.start()
	return
}

// start 启动可用性检测
func (p *PollInfo) start() {
	ticker := time.NewTicker(p.interval)
	for {
		select {
		case <-ticker.C:
			// 遍历不可用的轮询对象
			now := time.Now()
			for uuid, unavailableTime := range p.unavailableMap.Items() {
				if now.After(unavailableTime) || now.Equal(unavailableTime) {
					// 超过冷却时间，将不可用对象从不可用列表中移除
					p.unavailableMap.Remove(uuid)
					p.availableList.PushBack(doubly.Node{Uuid: uuid, Value: nil})
				}
			}
		case <-p.quitChan:
			ticker.Stop()
			return
		}
	}
}

// Stop 停止可用性检测
func (p *PollInfo) Stop() {
	p.quitChan <- true
}

// Init 初始化轮询对象
func (p *PollInfo) Init(ids ...int) {
	items := make([]doubly.Node, 0, len(ids))
	for _, id := range ids {
		items = append(items, doubly.Node{Uuid: strconv.Itoa(id), Value: nil})
	}
	p.availableList.PushBack(items...)
}

// IsAvailable 是否有可用的轮询对象标识
func (p *PollInfo) IsAvailable() (isAvailable bool) {
	return p.availableList.Len() > 0
}

// SetUnavailable 设置不可用的轮询对象标识
func (p *PollInfo) SetUnavailable(ids ...int) {
	uuids := make([]string, 0, len(ids))
	for _, id := range ids {
		uuids = append(uuids, strconv.Itoa(id))
	}
	p.availableList.Remove(uuids...)

	unavailableIds := make(map[string]time.Time, len(uuids))
	for _, uuid := range uuids {
		unavailableIds[uuid] = time.Now().Add(p.unavailableTime)
	}
	p.unavailableMap.MSet(unavailableIds)
}

// Poll 轮询
func (p *PollInfo) Poll() (id int, err error) {
	var node *doubly.Node
	if node, err = p.availableList.Poll(); err != nil {
		err = fmt.Errorf("no available id found")
		return
	}
	id, err = strconv.Atoi(node.Uuid)
	return
}

// GetGoroutinesAndTasks 根据任务总数计算所需协程数和每个协程处理的任务数量
//
//	total: 任务总数
//	expected: 每个协程期望处理的任务数量
func GetGoroutinesAndTasks(total, expected uint) (goroutineNum uint, tasks []uint) {
	// 如果总任务数小于或等于每个协程预期的任务数，则只需要一个协程来处理所有任务
	if total <= expected {
		return 1, []uint{total}
	}
	// 根据总任务数除以每个协程预期的任务数来计算需要的协程数
	// 这里添加`expected - 1`是为了在除法中实现向上取整，确保即使有余数也能分配足够的协程
	goroutineNum = (total + expected - 1) / expected
	// 初始化一个切片来保存每个协程将要处理的任务数
	tasks = make([]uint, goroutineNum)
	// 尽可能均匀地将任务分配给每个协程
	for k := range tasks {
		tasks[k] = total / goroutineNum
	}
	// 如果有余数，则将剩余的任务逐一分配给部分协程，以保证所有任务都能被处理
	remainder := total % goroutineNum
	for i := uint(0); i < remainder; i++ {
		tasks[i]++
	}
	return
}
