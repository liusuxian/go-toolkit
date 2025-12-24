/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-12-24 15:01:30
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-24 15:27:00
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package utils

import "math/rand/v2"

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int, charset ...string) (str string) {
	chars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if len(charset) > 0 {
		chars = charset[0] // 使用自定义字符集
	}

	result := make([]byte, length)
	for i := range length {
		result[i] = chars[rand.IntN(len(chars))]
	}
	return string(result)
}
