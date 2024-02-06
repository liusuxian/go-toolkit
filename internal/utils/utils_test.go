/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-07 02:45:07
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-07 02:45:12
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package utils_test

import (
	"github.com/liusuxian/go-toolkit/internal/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDoRedisCmdArgs(t *testing.T) {
	var (
		assert = assert.New(t)
		args   = []any{1, 1.2, "hello", []byte{1}, map[string]any{"a": 1}}
		err    error
	)
	err = utils.DoRedisArgs(0, args...)
	assert.NoError(err)
	assert.Equal([]any{1, 1.2, "hello", []uint8{0x1}, []uint8{0x7b, 0x22, 0x61, 0x22, 0x3a, 0x31, 0x7d}}, args)
}
