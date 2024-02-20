/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-20 21:04:55
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-20 21:59:20
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkrobot

import (
	"bytes"
	"github.com/liusuxian/go-toolkit/gtkjson"
	"github.com/liusuxian/go-toolkit/gtkstr"
	"github.com/pkg/errors"
	"net/http"
)

// FeishuRobot
type FeishuRobot struct {
	WebHookURL string
	data       FeiShuMessage
}

// FeiShuMessage 飞书消息
type FeiShuMessage struct {
	MsgType string        `json:"msg_type" dc:"消息类型"` // 消息类型
	Content FeishuContent `json:"content" dc:"消息内容"`  // 消息内容
}

// FeishuContent 飞书消息内容
type FeishuContent struct {
	Text string `json:"text" dc:"文本内容"` // 文本内容
}

// NewFeishuRobot 新建飞书机器人
func NewFeishuRobot(webHookURL string) (fr *FeishuRobot) {
	return &FeishuRobot{
		WebHookURL: webHookURL,
	}
}

// SendTextMessage 发送文本消息
func (fr *FeishuRobot) SendTextMessage(content string) (err error) {
	if gtkstr.TrimAll(fr.WebHookURL) == "" {
		return
	}
	fr.data = FeiShuMessage{
		MsgType: "text",
		Content: FeishuContent{
			Text: content,
		},
	}
	return fr.send()
}

// send 发送
func (fr *FeishuRobot) send() (err error) {
	var message []byte
	if message, err = gtkjson.Bytes(fr.data); err != nil {
		return
	}
	var buffer bytes.Buffer
	if _, err = buffer.Write(message); err != nil {
		return
	}
	// 创建`HTTP`请求
	var req *http.Request
	if req, err = http.NewRequest("POST", fr.WebHookURL, &buffer); err != nil {
		return
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	// 发送请求
	var (
		client = &http.Client{}
		resp   *http.Response
	)
	if resp, err = client.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.Errorf("Request Failed With Status Code: %d", resp.StatusCode)
		return
	}
	return
}
