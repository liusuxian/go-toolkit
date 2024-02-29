/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-27 22:00:33
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-28 15:39:48
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package redbookweb_test

import (
	"context"
	"fmt"
	"github.com/liusuxian/go-toolkit/gtkconf"
	"github.com/liusuxian/go-toolkit/gtkconv"
	"github.com/liusuxian/go-toolkit/gtksdk/redbook/redbookweb"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

type Config struct {
	Phone string `json:"phone" dc:"phone"`
}

func TestSendCode(t *testing.T) {
	var (
		assert   = assert.New(t)
		ctx      = context.Background()
		config   = Config{}
		localCfg *gtkconf.Config
		resp     *redbookweb.SendCodeResponse
		err      error
	)
	localCfg, err = gtkconf.NewConfig("../../../test_config/redbookweb.json")
	assert.NoError(err)
	err = localCfg.StructKey("test", &config)
	assert.NoError(err)

	c := redbookweb.NewClient()
	resp, err = c.SendCode(ctx, redbookweb.SendCodeRequest{
		Zone:  "86",
		Phone: config.Phone,
	})
	t.Logf("err: %v\n\n", err)
	t.Logf("resp: %v\n\n", resp)
}

func TestLoginWithVerifyCodeAndCustomerLoginAndLoginAndUserInfo(t *testing.T) {
	var (
		assert   = assert.New(t)
		ctx      = context.Background()
		config   = Config{}
		localCfg *gtkconf.Config
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
	localCfg, err = gtkconf.NewConfig("../../../test_config/redbookweb.json")
	assert.NoError(err)
	err = localCfg.StructKey("test", &config)
	assert.NoError(err)

	c := redbookweb.NewClient()
	if resp1, cookies1, err = c.LoginWithVerifyCode(ctx, redbookweb.LoginWithVerifyCodeRequest{
		Zone:       "86",
		Mobile:     config.Phone,
		VerifyCode: "138264", // 需要收到验证码以后手动输入
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
