/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-20 21:04:55
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-28 01:14:11
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkrobot_test

import (
	"context"
	"github.com/liusuxian/go-toolkit/gtkconf"
	"github.com/liusuxian/go-toolkit/gtkrobot"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Config struct {
	WebHookURL string `json:"webHookURL" dc:"webHookURL"`
}

func TestSendTextMessage(t *testing.T) {
	var (
		assert   = assert.New(t)
		config   = Config{}
		localCfg *gtkconf.Config
		err      error
	)
	localCfg, err = gtkconf.NewConfig("../test_config/feishu_robot.json")
	assert.NoError(err)
	err = localCfg.StructKey("test", &config)
	assert.NoError(err)
	fr := gtkrobot.NewFeishuRobot(config.WebHookURL)
	fr.SendTextMessage(context.Background(), "我是测试日志上报内容，并不是真的服务器报错，请悉知")
}
