/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-20 21:04:55
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-28 14:25:55
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkrobot

import (
	"context"
	"github.com/liusuxian/go-toolkit/gtkhttp"
	"github.com/liusuxian/go-toolkit/gtkstr"
	"github.com/pkg/errors"
	"net/http"
)

// FeishuRobot
type FeishuRobot struct {
	webHookURL     string
	requestBuilder gtkhttp.RequestBuilder // 请求构建器
	httpClient     *http.Client
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
		webHookURL:     webHookURL,
		requestBuilder: gtkhttp.NewRequestBuilder(),
		httpClient:     &http.Client{},
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
			gtkhttp.SetBody(data),
			gtkhttp.SetContentType("application/json; charset=utf-8"),
		}
		req *http.Request
	)

	if req, err = fr.requestBuilder.Build(ctx, http.MethodPost, fr.webHookURL, setters...); err != nil {
		return
	}
	// 发送请求
	var resp *http.Response
	if resp, err = fr.httpClient.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	if gtkhttp.IsFailureStatusCode(resp) {
		err = &gtkhttp.RequestError{
			HTTPStatusCode: resp.StatusCode,
			Err:            errors.New("Request Failed"),
		}
		return
	}
	return
}
