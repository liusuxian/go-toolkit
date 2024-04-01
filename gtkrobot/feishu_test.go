/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-20 21:04:55
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-03-31 20:33:00
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkrobot_test

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/liusuxian/go-toolkit/gtkenv"
	"github.com/liusuxian/go-toolkit/gtkrobot"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSendTextMessage(t *testing.T) {
	var (
		assert = assert.New(t)
		ctx    = context.Background()
		err    error
	)
	err = godotenv.Load(".env")
	assert.NoError(err)
	webHookUrl := gtkenv.Get("testFeishuRobotWebHookUrl")
	t.Logf("webHookUrl: %s\n", webHookUrl)
	fr := gtkrobot.NewFeishuRobot(webHookUrl)
	err = fr.SendTextMessage(ctx, "我是测试日志上报内容，并不是真的服务器报错，请悉知")
	assert.NoError(err)
}
