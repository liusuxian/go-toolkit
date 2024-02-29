/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-28 00:29:00
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-28 16:57:02
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package redbookweb

// CustomerLoginRequest 客户登录请求数据
type CustomerLoginRequest struct {
	Ticket string `json:"ticket"` // 票根
}

// CustomerLoginResponse 客户登录响应数据
type CustomerLoginResponse struct {
	Code    int               `json:"code"`
	Success bool              `json:"success"`
	Data    CustomerLoginData `json:"data"`
}

// CustomerLoginData
type CustomerLoginData struct {
	UserID           string `json:"userId"`
	NickName         string `json:"nickName"`
	Avatar           string `json:"avatar"`
	BAccountNo       string `json:"bAccountNo"`
	AccountType      string `json:"accountType"`
	PrimaryAccountNo string `json:"primaryAccountNo"`
	Email            string `json:"email"`
	AuthToken        string `json:"authToken"`
	BAccountName     string `json:"bAccountName"`
	BUserID          string `json:"bUserId"`
	Resources        any    `json:"resources"`
	CurrentResource  string `json:"currentResource"`
	Permissions      []any  `json:"permissions"`
}

// LoginRequest 登录请求数据
type LoginRequest struct {
}

// LoginResponse 登录响应数据
type LoginResponse struct {
	Result  int  `json:"result"`
	Success bool `json:"success"`
}

// UserInfoRequest 获取用户信息请求数据
type UserInfoRequest struct {
}

// UserInfoResponse 获取用户信息响应数据
type UserInfoResponse struct {
	Result  int          `json:"result"`
	Success bool         `json:"success"`
	Data    UserInfoData `json:"data"`
}

// UserInfoData
type UserInfoData struct {
	UserID      string   `json:"userId"`
	UserName    string   `json:"userName"`
	UserAvatar  string   `json:"userAvatar"`
	RedID       string   `json:"redId"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}
