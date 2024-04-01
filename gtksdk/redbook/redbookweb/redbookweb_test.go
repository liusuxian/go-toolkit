/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-27 22:00:33
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-03-31 20:41:25
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package redbookweb_test

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/liusuxian/go-toolkit/gtkconv"
	"github.com/liusuxian/go-toolkit/gtkenv"
	"github.com/liusuxian/go-toolkit/gtksdk/redbook/redbookweb"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func TestSendCode(t *testing.T) {
	var (
		assert = assert.New(t)
		ctx    = context.Background()
		resp   *redbookweb.SendCodeResponse
		err    error
	)
	err = godotenv.Load(".env")
	assert.NoError(err)
	phone := gtkenv.Get("phone")
	t.Logf("phone: %s\n", phone)
	c := redbookweb.NewClient()
	resp, err = c.SendCode(ctx, redbookweb.SendCodeRequest{
		Zone:  "86",
		Phone: phone,
	})
	t.Logf("err: %v\n\n", err)
	t.Logf("resp: %v\n\n", resp)
}

func TestLoginWithVerifyCodeAndCustomerLoginAndLoginAndUserInfo(t *testing.T) {
	var (
		assert   = assert.New(t)
		ctx      = context.Background()
		resp1    *redbookweb.LoginWithVerifyCodeResponse
		cookies1 []*http.Cookie
		resp2    *redbookweb.CustomerLoginResponse
		cookies2 []*http.Cookie
		resp3    *redbookweb.LoginResponse
		cookies3 []*http.Cookie
		resp4    *redbookweb.UserInfoResponse
		cookies4 []*http.Cookie
		err      error
	)
	err = godotenv.Load(".env")
	assert.NoError(err)
	phone := gtkenv.Get("phone")
	t.Logf("phone: %s\n", phone)

	c := redbookweb.NewClient()
	if resp1, cookies1, err = c.LoginWithVerifyCode(ctx, redbookweb.LoginWithVerifyCodeRequest{
		Zone:       "86",
		Mobile:     phone,
		VerifyCode: "198129", // 需要收到验证码以后手动输入
	}); err != nil {
		t.Logf("err1: %v\n\n", err)
		return
	}
	t.Logf("resp1: %v\n\n", resp1)
	cookieList := make([]string, 0, len(cookies1))
	for _, cookie := range cookies1 {
		cookieList = append(cookieList, fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
	}
	t.Logf("cookies1: %v\n\n", strings.Join(cookieList, "; "))

	if resp2, cookies2, err = c.CustomerLogin(ctx, redbookweb.CustomerLoginRequest{
		Ticket: gtkconv.ToString(resp1.Data),
	}, cookies1); err != nil {
		t.Logf("err2: %v\n\n", err)
		return
	}
	t.Logf("resp2: %v\n\n", resp2)
	cookies2 = append(cookies2, cookies1...)
	cookieList = make([]string, 0, len(cookies2))
	for _, cookie := range cookies2 {
		cookieList = append(cookieList, fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
	}
	t.Logf("cookies2: %v\n\n", strings.Join(cookieList, "; "))

	if resp3, cookies3, err = c.Login(ctx, cookies2); err != nil {
		t.Logf("err3: %v\n\n", err)
		return
	}
	t.Logf("resp3: %v\n\n", resp3)
	cookies3 = append(cookies3, cookies2...)
	cookieList = make([]string, 0, len(cookies3))
	for _, cookie := range cookies3 {
		cookieList = append(cookieList, fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
	}
	t.Logf("cookies3: %v\n\n", strings.Join(cookieList, "; "))

	if resp4, cookies4, err = c.UserInfo(ctx, cookies3); err != nil {
		t.Logf("err4: %v\n\n", err)
		return
	}
	t.Logf("resp4: %v\n\n", resp4)
	cookies4 = append(cookies4, cookies3...)
	cookieList = make([]string, 0, len(cookies4))
	for _, cookie := range cookies4 {
		cookieList = append(cookieList, fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
	}
	t.Logf("cookies4: %v\n\n", strings.Join(cookieList, "; "))
}
