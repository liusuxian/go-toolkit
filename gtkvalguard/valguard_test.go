/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-05-09 21:15:02
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-05-09 21:25:37
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkvalguard_test

import (
	"github.com/liusuxian/go-toolkit/gtkvalguard"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetValue(t *testing.T) {
	assert := assert.New(t)
	var intPtr *int
	assert.Equal(0, gtkvalguard.GetValue(intPtr))
	assert.Equal(1, gtkvalguard.GetValue(intPtr, 1))
	var floatPtr *float64
	assert.Equal(float64(0), gtkvalguard.GetValue(floatPtr))
	var strPtr *string
	assert.Equal("", gtkvalguard.GetValue(strPtr))
	var boolPtr *bool
	assert.Equal(false, gtkvalguard.GetValue(boolPtr))
}
