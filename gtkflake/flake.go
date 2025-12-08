/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-05 15:48:50
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-08 23:42:52
 * @Description: 分布式唯一ID生成器
 * 默认情况下，ID由以下部分组成：
 * 39位时间（以10毫秒为单位）
 * 8位序列号
 * 16位机器ID
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkflake

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

var (
	errInvalidBitsTime      = errors.New("bit length for time must be 32 or more") // 时间位长度必须大于32
	errInvalidBitsSequence  = errors.New("invalid bit length for sequence number") // 序列号位长度无效
	errInvalidBitsMachineID = errors.New("invalid bit length for machine id")      // 机器ID位长度无效
	errInvalidTimeUnit      = errors.New("invalid time unit")                      // 时间单位无效
	errInvalidSequence      = errors.New("invalid sequence number")                // 序列号无效
	errInvalidMachineID     = errors.New("invalid machine id")                     // 机器ID无效
	errStartTimeAhead       = errors.New("start time is ahead")                    // 起始时间早于当前时间
	errOverTimeLimit        = errors.New("over the time limit")                    // 超过时间限制
	errNoPrivateAddress     = errors.New("no private ip address")                  // 没有私有IP地址
)

// InterfaceAddrs 用于获取网络地址的接口
type InterfaceAddrs func() ([]net.Addr, error)

// Settings 配置
type Settings struct {
	BitsSequence   int                       // 序列号的位长度
	BitsMachineID  int                       // 机器ID的位长度
	TimeUnit       time.Duration             // 时间单位
	StartTime      time.Time                 // 起始时间
	MachineID      func() (n int, err error) // 机器ID
	CheckMachineID func(int) (ok bool)       // 验证机器ID的唯一性
}

// Flake 分布式唯一ID生成器
type Flake struct {
	bitsTime     int
	bitsSequence int
	bitsMachine  int
	timeUnit     int64
	startTime    int64
	elapsedTime  int64
	sequence     int
	machine      int
	mutex        *sync.Mutex
}

const (
	defaultTimeUnit = 1e7 // nsec, i.e. 10 msec
	// defaultBitsTime     = 39
	defaultBitsSequence = 8
	defaultBitsMachine  = 16
)

var defaultInterfaceAddrs = net.InterfaceAddrs

// New 创建一个分布式唯一ID生成器
func New(st Settings) (flake *Flake, err error) {
	if st.BitsSequence < 0 || st.BitsSequence > 30 {
		return nil, errInvalidBitsSequence
	}
	if st.BitsMachineID < 0 || st.BitsMachineID > 30 {
		return nil, errInvalidBitsMachineID
	}
	if st.TimeUnit < 0 || (st.TimeUnit > 0 && st.TimeUnit < time.Millisecond) {
		return nil, errInvalidTimeUnit
	}
	if st.StartTime.After(time.Now()) {
		return nil, errStartTimeAhead
	}

	f := new(Flake)
	f.mutex = new(sync.Mutex)

	if st.BitsSequence == 0 {
		f.bitsSequence = defaultBitsSequence
	} else {
		f.bitsSequence = st.BitsSequence
	}

	if st.BitsMachineID == 0 {
		f.bitsMachine = defaultBitsMachine
	} else {
		f.bitsMachine = st.BitsMachineID
	}

	f.bitsTime = 63 - f.bitsSequence - f.bitsMachine
	if f.bitsTime < 32 {
		return nil, errInvalidBitsTime
	}

	if st.TimeUnit == 0 {
		f.timeUnit = defaultTimeUnit
	} else {
		f.timeUnit = int64(st.TimeUnit)
	}

	if st.StartTime.IsZero() {
		f.startTime = f.toInternalTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
	} else {
		f.startTime = f.toInternalTime(st.StartTime)
	}

	f.sequence = 1<<f.bitsSequence - 1

	var e error
	if st.MachineID == nil {
		f.machine, e = lower16BitPrivateIP(defaultInterfaceAddrs)
	} else {
		f.machine, e = st.MachineID()
	}
	if e != nil {
		return nil, e
	}

	if f.machine < 0 || f.machine >= 1<<f.bitsMachine {
		return nil, errInvalidMachineID
	}

	if st.CheckMachineID != nil && !st.CheckMachineID(f.machine) {
		return nil, errInvalidMachineID
	}

	return f, nil
}

// ID 生成唯一ID
func (f *Flake) ID() (id int64, err error) {
	maskSequence := 1<<f.bitsSequence - 1

	f.mutex.Lock()
	defer f.mutex.Unlock()

	current := f.currentElapsedTime()
	if f.elapsedTime < current {
		f.elapsedTime = current
		f.sequence = 0
	} else {
		f.sequence = (f.sequence + 1) & maskSequence
		if f.sequence == 0 {
			f.elapsedTime++
			overtime := f.elapsedTime - current
			f.sleep(overtime)
		}
	}

	return f.toID()
}

// RequestID 生成唯一请求ID
func (f *Flake) RequestID() (requestId string, err error) {
	var (
		now = time.Now()
		id  int64
	)
	if id, err = f.ID(); err != nil {
		return
	}
	timeStr := now.Format("20060102150405")
	// 生成2字节随机数
	randomBytes := make([]byte, 2)
	if _, err = rand.Reader.Read(randomBytes); err != nil {
		return
	}
	requestId = fmt.Sprintf("%s%016X%04X", timeStr, id, randomBytes)
	return
}

// ToTime 返回生成给定ID时的时间
func (f *Flake) ToTime(id int64) (t time.Time) {
	return time.Unix(0, (f.startTime+f.timePart(id))*f.timeUnit)
}

// Compose 创建ID
func (f *Flake) Compose(t time.Time, sequence, machineID int) (id int64, err error) {
	elapsedTime := f.toInternalTime(t.UTC()) - f.startTime
	if elapsedTime < 0 {
		return 0, errStartTimeAhead
	}
	if elapsedTime >= 1<<f.bitsTime {
		return 0, errOverTimeLimit
	}

	if sequence < 0 || sequence >= 1<<f.bitsSequence {
		return 0, errInvalidSequence
	}

	if machineID < 0 || machineID >= 1<<f.bitsMachine {
		return 0, errInvalidMachineID
	}

	return elapsedTime<<(f.bitsSequence+f.bitsMachine) |
		int64(sequence)<<f.bitsMachine |
		int64(machineID), nil
}

// Decompose 分解ID
func (f *Flake) Decompose(id int64) (m map[string]int64) {
	time := f.timePart(id)
	sequence := f.sequencePart(id)
	machine := f.machinePart(id)
	return map[string]int64{
		"id":       id,
		"time":     time,
		"sequence": sequence,
		"machine":  machine,
	}
}

// toInternalTime 将时间转换为内部时间
func (f *Flake) toInternalTime(t time.Time) (internalTime int64) {
	return t.UTC().UnixNano() / f.timeUnit
}

// currentElapsedTime 当前经过的时间
func (f *Flake) currentElapsedTime() (elapsedTime int64) {
	return f.toInternalTime(time.Now()) - f.startTime
}

// sleep 睡眠
func (f *Flake) sleep(overtime int64) {
	sleepTime := time.Duration(overtime*f.timeUnit) -
		time.Duration(time.Now().UTC().UnixNano()%f.timeUnit)
	time.Sleep(sleepTime)
}

// toID 生成ID
func (f *Flake) toID() (id int64, err error) {
	if f.elapsedTime >= 1<<f.bitsTime {
		return 0, errOverTimeLimit
	}

	return f.elapsedTime<<(f.bitsSequence+f.bitsMachine) |
		int64(f.sequence)<<f.bitsMachine |
		int64(f.machine), nil
}

// timePart 时间部分
func (f *Flake) timePart(id int64) (time int64) {
	return id >> (f.bitsSequence + f.bitsMachine)
}

// sequencePart 序列号部分
func (f *Flake) sequencePart(id int64) (sequence int64) {
	maskSequence := int64((1<<f.bitsSequence - 1) << f.bitsMachine)
	return (id & maskSequence) >> f.bitsMachine
}

// machinePart 机器ID部分
func (f *Flake) machinePart(id int64) (machine int64) {
	maskMachine := int64(1<<f.bitsMachine - 1)
	return id & maskMachine
}

// privateIPv4 获取私有IPv4地址
func privateIPv4(interfaceAddrs InterfaceAddrs) (ip net.IP, err error) {
	var as []net.Addr
	if as, err = interfaceAddrs(); err != nil {
		return
	}

	for _, a := range as {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}

		ip = ipnet.IP.To4()
		if isPrivateIPv4(ip) {
			return
		}
	}

	err = errNoPrivateAddress
	return
}

// isPrivateIPv4 判断是否为私有IPv4地址
func isPrivateIPv4(ip net.IP) (ok bool) {
	return ip != nil && (ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168 || ip[0] == 169 && ip[1] == 254)
}

// lower16BitPrivateIP 获取私有IPv4地址的低16位
func lower16BitPrivateIP(interfaceAddrs InterfaceAddrs) (machine int, err error) {
	var ip net.IP
	if ip, err = privateIPv4(interfaceAddrs); err != nil {
		return
	}

	machine = (int(ip[2]) << 8) + int(ip[3])
	return
}
