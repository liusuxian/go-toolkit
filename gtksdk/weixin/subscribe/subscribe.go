/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-13 18:51:12
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-13 18:52:20
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package subscribe

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strings"
)

// SubscribeSign 验证消息订阅的签名
func SubscribeSign(token, signature string, params map[string]string) (err error) {
	paramsNum := len(params)
	if paramsNum == 0 {
		err = fmt.Errorf("params is empty")
		return
	}
	// 将参数进行字典序排序
	keys := make([]string, 0, paramsNum+1)
	keys = append(keys, token)
	for _, v := range params {
		keys = append(keys, v)
	}
	sort.Strings(keys)
	// 将参数拼接成一个字符串进行SHA1加密
	hash := sha1.New()
	hash.Write([]byte(strings.Join(keys, "")))
	sign := fmt.Sprintf("%x", hash.Sum(nil))
	// 将计算出的签名与URL中的签名参数进行比较，如果相等，则验证通过
	if sign != signature {
		err = fmt.Errorf("signature verification failed, sign: %s, signature: %s", sign, signature)
		return
	}
	return
}
