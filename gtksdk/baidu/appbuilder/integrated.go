/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-26 11:56:58
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-28 14:26:36
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package appbuilder

import (
	"context"
	"fmt"
	"github.com/liusuxian/go-toolkit/gtkhttp"
	"net/http"
)

const (
	integratedURL = "/rpc/2.0/cloud_hub/v1/ai_engine/agi_platform/v1/instance/integrated"
)

// Integrated 集成API
func (c *Client) Integrated(ctx context.Context, request IntegratedRequest) (response *IntegratedResponse, err error) {
	var (
		setters = []gtkhttp.RequestOption{
			gtkhttp.SetBody(map[string]any{
				"query":           request.Query,
				"response_mode":   "blocking",
				"conversation_id": request.ConversationId,
			}),
			gtkhttp.SetContentType("application/json; charset=utf-8"),
			gtkhttp.SetKeyValue("Accept", "application/json; charset=utf-8"),
			gtkhttp.SetKeyValue("X-Appbuilder-Authorization", fmt.Sprintf("Bearer %s", c.config.AppToken)),
		}
		req *http.Request
	)
	if req, err = c.requestBuilder.Build(ctx, http.MethodPost, c.fullURL(integratedURL), setters...); err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}

// IntegratedStream 集成API，流式返回
func (c *Client) IntegratedStream(ctx context.Context, request IntegratedRequest) (stream *IntegratedResponseStream, err error) {
	var (
		setters = []gtkhttp.RequestOption{
			gtkhttp.SetBody(map[string]any{
				"query":           request.Query,
				"response_mode":   "streaming",
				"conversation_id": request.ConversationId,
			}),
			gtkhttp.SetContentType("application/json; charset=utf-8"),
			gtkhttp.SetKeyValue("Accept", "text/event-stream"),
			gtkhttp.SetKeyValue("Cache-Control", "no-cache"),
			gtkhttp.SetKeyValue("Connection", "keep-alive"),
			gtkhttp.SetKeyValue("X-Appbuilder-Authorization", fmt.Sprintf("Bearer %s", c.config.AppToken)),
		}
		req *http.Request
	)
	if req, err = c.requestBuilder.Build(ctx, http.MethodPost, c.fullURL(integratedURL), setters...); err != nil {
		return
	}

	var resp *streamReader[IntegratedResponseResult]
	if resp, err = sendRequestStream[IntegratedResponseResult](c, req); err != nil {
		return
	}

	stream = &IntegratedResponseStream{
		streamReader: resp,
	}
	return
}
