/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-26 11:56:58
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-23 19:35:37
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
func (c *Client) Integrated(ctx context.Context, request IntegratedRequest) (response IntegratedResponse, err error) {
	var (
		setters = []gtkhttp.RequestOption{
			gtkhttp.WithBody(map[string]any{
				"query":           request.Query,
				"response_mode":   "blocking",
				"conversation_id": request.ConversationId,
			}),
			gtkhttp.WithKeyValue("X-Appbuilder-Authorization", fmt.Sprintf("Bearer %s", c.appToken)),
		}
		req *http.Request
	)
	if req, err = c.httpClient.NewRequest(ctx, http.MethodPost, c.httpClient.FullURL(integratedURL), setters...); err != nil {
		return
	}

	err = c.httpClient.SendRequest(req, &response)
	return
}

// IntegratedStream 集成API，流式返回
func (c *Client) IntegratedStream(ctx context.Context, request IntegratedRequest) (stream *IntegratedResponseStream, err error) {
	var (
		setters = []gtkhttp.RequestOption{
			gtkhttp.WithBody(map[string]any{
				"query":           request.Query,
				"response_mode":   "streaming",
				"conversation_id": request.ConversationId,
			}),
			gtkhttp.WithKeyValue("X-Appbuilder-Authorization", fmt.Sprintf("Bearer %s", c.appToken)),
		}
		req *http.Request
	)
	if req, err = c.httpClient.NewRequest(ctx, http.MethodPost, c.httpClient.FullURL(integratedURL), setters...); err != nil {
		return
	}

	var resp *gtkhttp.StreamReader[IntegratedResponseResult]
	if resp, err = gtkhttp.SendRequestStream[IntegratedResponseResult](c.httpClient, req); err != nil {
		return
	}

	stream = &IntegratedResponseStream{
		StreamReader: resp,
	}
	return
}
