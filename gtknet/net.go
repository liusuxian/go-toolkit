/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-05-26 15:33:37
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-25 10:29:24
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtknet

import (
	"errors"
	"net"
	"net/http"
	"strings"
)

var (
	ErrNoPrivateAddress = errors.New("no private ip address") // 没有私有IP地址
)

// IsPrivateIPv4 判断是否是私有 IPv4 地址
func IsPrivateIPv4(ip net.IP) (ok bool) {
	return ip != nil && (ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168 || ip[0] == 169 && ip[1] == 254)
}

// PrivateIPv4 获取私有 IPv4 地址
func PrivateIPv4() (ip net.IP, err error) {
	var as []net.Addr
	if as, err = net.InterfaceAddrs(); err != nil {
		return
	}

	for _, a := range as {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}

		ip = ipnet.IP.To4()
		if IsPrivateIPv4(ip) {
			return
		}
	}

	err = ErrNoPrivateAddress
	return
}

// GetClientIp 获取客户端的`IP`地址
func GetClientIp(r *http.Request) (clientIp string) {
	realIps := r.Header.Get("X-Forwarded-For")
	if realIps != "" && !strings.EqualFold("unknown", realIps) {
		ipArray := strings.Split(realIps, ",")
		for i := range ipArray {
			ipArray[i] = strings.TrimSpace(ipArray[i])
		}
		clientIp = ipArray[0]
	}
	if clientIp == "" {
		clientIp = r.Header.Get("Proxy-Client-IP")
	}
	if clientIp == "" {
		clientIp = r.Header.Get("WL-Proxy-Client-IP")
	}
	if clientIp == "" {
		clientIp = r.Header.Get("HTTP_CLIENT_IP")
	}
	if clientIp == "" {
		clientIp = r.Header.Get("HTTP_X_FORWARDED_FOR")
	}
	if clientIp == "" {
		clientIp = r.Header.Get("X-Real-IP")
	}
	if clientIp == "" {
		clientIp, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return
}
