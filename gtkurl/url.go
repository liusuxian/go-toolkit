/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-06 18:30:38
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-22 18:31:49
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkurl

import (
	"net/url"
	"strings"
)

// IsUrlEncoded 检查字符串是否已经被`URL`编码
func IsUrlEncoded(str string) (ok bool, err error) {
	if len(str) == 0 {
		return
	}
	var decodeStr string
	if decodeStr, err = url.QueryUnescape(str); err != nil {
		return
	}
	encodeStr := url.QueryEscape(decodeStr)
	ok = encodeStr == str
	return
}

// QueryDecode 字符串`URL`解码
func QueryDecode(str string) (decodeStr string, err error) {
	if len(str) == 0 {
		return
	}
	var tmpDecodeStr string
	if tmpDecodeStr, err = url.QueryUnescape(str); err != nil {
		return
	}
	encodeStr := url.QueryEscape(tmpDecodeStr)
	if encodeStr == str {
		decodeStr = tmpDecodeStr
		return
	}
	decodeStr = str
	return
}

// IsUrl 检查字符串是否是`URL`链接
func IsUrl(str string) (ok bool) {
	if len(str) == 0 {
		return
	}
	ok = strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://")
	return
}
