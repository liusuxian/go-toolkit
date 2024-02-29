/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-27 22:09:13
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-27 23:04:38
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package redbookweb

// SendCodeRequest 发送验证码请求数据
type SendCodeRequest struct {
	Zone  string `json:"zone"`  // 地区
	Phone string `json:"phone"` // 手机号
}

// SendCodeResponse 发送验证码响应数据
type SendCodeResponse struct {
	AlertMsg   string `json:"alertMsg"`
	Data       any    `json:"data"`
	ErrorCode  int    `json:"errorCode"`
	ErrorMsg   string `json:"errorMsg"`
	StatusCode int    `json:"statusCode"`
	Success    bool   `json:"success"`
	TraceLogId string `json:"traceLogId"`
}

// LoginWithVerifyCodeRequest 使用验证码登录请求数据
type LoginWithVerifyCodeRequest struct {
	Zone       string `json:"zone"`       // 地区
	Mobile     string `json:"mobile"`     // 手机号
	VerifyCode string `json:"verifyCode"` // 验证码
}

// LoginWithVerifyCodeResponse 使用验证码登录响应数据
type LoginWithVerifyCodeResponse struct {
	AlertMsg   string `json:"alertMsg"`
	Data       any    `json:"data"`
	ErrorCode  int    `json:"errorCode"`
	ErrorMsg   string `json:"errorMsg"`
	StatusCode int    `json:"statusCode"`
	Success    bool   `json:"success"`
	TraceLogId string `json:"traceLogId"`
}
