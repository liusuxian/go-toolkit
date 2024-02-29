/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-27 22:07:07
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-27 22:07:12
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package appbuilder

// IntegratedRequest 请求数据
type IntegratedRequest struct {
	Query          string `json:"query" dc:"用户的请求query"`                                                              // 用户的请求query
	ConversationId string `json:"conversation_id" dc:"对话ID，仅对话型应用生效。在对话型应用中，空:表示表新建会话，非空:表示在对应的会话中继续进行对话，服务内部维护对话历史"` // 对话ID，仅对话型应用生效。在对话型应用中，空:表示表新建会话，非空:表示在对应的会话中继续进行对话，服务内部维护对话历史
}

// IntegratedResponse 响应数据
type IntegratedResponse struct {
	Code    int                       `json:"code" dc:"错误码。非0为错误"` // 错误码。非0为错误
	Message string                    `json:"message" dc:"报错信息"`   // 报错信息
	Result  *IntegratedResponseResult `json:"result" dc:"返回结果"`    // 返回结果
}

// IntegratedResponseResult 返回结果
type IntegratedResponseResult struct {
	Answer         string `json:"answer" dc:"应用响应结果"`                                                   // 应用响应结果
	ConversationId string `json:"conversationId" dc:"对话ID，仅对话式应用生效。如果是对话请求中没有conversation_id，则会自动生成一个"` // 对话ID，仅对话式应用生效。如果是对话请求中没有conversation_id，则会自动生成一个
}

// IntegratedResponseStream 流式响应数据
type IntegratedResponseStream struct {
	*streamReader[IntegratedResponseResult]
}
