/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-05-26 15:33:37
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2023-05-26 16:19:41
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtknet_test

import (
	"github.com/liusuxian/go-toolkit/gtknet"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestIsPrivateIPv4(t *testing.T) {
	assert := assert.New(t)
	assert.False(gtknet.IsPrivateIPv4(net.ParseIP("121, 199, 16, 7")))
	assert.False(gtknet.IsPrivateIPv4(net.ParseIP("127, 0, 0, 1")))
}

func TestPrivateIPv4(t *testing.T) {
	assert := assert.New(t)
	actualObj, err := gtknet.PrivateIPv4()
	if assert.NoError(err) {
		assert.True(gtknet.IsPrivateIPv4(actualObj))
	}
}
