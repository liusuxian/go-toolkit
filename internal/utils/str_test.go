/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-12-24 15:03:05
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-24 15:05:13
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package utils_test

import (
	"github.com/liusuxian/go-toolkit/internal/utils"
	"testing"
)

func TestGenerateRandomString(t *testing.T) {
	for range 10 {
		t.Log(utils.GenerateRandomString(12))
	}
	for range 10 {
		t.Log(utils.GenerateRandomString(16))
	}
}
