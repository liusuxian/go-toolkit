/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2026-01-13 11:48:27
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2026-01-13 16:42:15
 * @Description:
 *
 * Copyright (c) 2026 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtksms

// SendSmsRequest 发送短信请求
type SendSmsRequest struct {
	PhoneNumbers  string `json:"phone_numbers"`  // 手机号
	SignName      string `json:"sign_name"`      // 签名
	TemplateCode  string `json:"template_code"`  // 模板代码
	TemplateParam string `json:"template_param"` // 模板参数
}
