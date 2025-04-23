/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-20 21:04:55
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-23 19:21:11
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkrobot

import (
	"context"
	"github.com/liusuxian/go-toolkit/gtkhttp"
	"github.com/liusuxian/go-toolkit/gtkstr"
	"net/http"
)

// FeishuRobot
type FeishuRobot struct {
	webHookURL string
	httpClient *gtkhttp.HTTPClient // http 客户端
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
		webHookURL: webHookURL,
		httpClient: gtkhttp.NewHTTPClient(""),
	}
}

// SendTextMessage 发送文本消息
func (fr *FeishuRobot) SendTextMessage(ctx context.Context, content string) (err error) {
	if gtkstr.TrimAll(fr.webHookURL) == "" {
		return
	}
	return fr.send(ctx, FeiShuMessage{
		MsgType: "text",
		Content: FeishuContent{
			Text: content,
		},
	})
}

// send 发送
func (fr *FeishuRobot) send(ctx context.Context, data FeiShuMessage) (err error) {
	var (
		setters = []gtkhttp.RequestOption{
			gtkhttp.WithBody(data),
			gtkhttp.WithContentType("application/json; charset=utf-8"),
		}
		req *http.Request
	)
	if req, err = fr.httpClient.NewRequest(ctx, http.MethodPost, fr.webHookURL, setters...); err != nil {
		return
	}
	return fr.httpClient.SendRequest(req, nil)
}
