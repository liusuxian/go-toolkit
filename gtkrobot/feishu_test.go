/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-20 21:04:55
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-20 22:07:28
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkrobot_test

import (
	"github.com/liusuxian/go-toolkit/gtkrobot"
	"testing"
)

func TestSendTextMessage(t *testing.T) {
	fr := gtkrobot.NewFeishuRobot("https://open.feishu.cn/open-apis/bot/v2/hook/203ff12b-f912-466a-836c-fd7004c148cb")
	fr.SendTextMessage("我是测试日志上报内容，并不是真的服务器报错，请悉知")
}
