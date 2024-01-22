/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-19 22:29:06
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-22 12:47:11
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gfweixin_test

import (
	"context"
	"github.com/liusuxian/go-toolkit/gf/gfweixin"
	"testing"
)

func TestAuthCode2Session(t *testing.T) {
	weChatService := gfweixin.NewWeChatService("wx65064684d6c0f73f", "e20bf5f51062ab55ed2b1cec8e540502", nil)
	ctx := context.Background()
	var (
		resMap map[string]any
		err    error
	)
	if resMap, err = weChatService.AuthCode2Session(ctx, "0d3AYx000aRfrR17go400dpCcR2AYx0c"); err != nil {
		t.Logf("AuthCode2Session err: %+v\n", err)
		return
	}
	t.Logf("resMap: %+v\n", resMap)
}

func TestGetStableAccessToken(t *testing.T) {
	weChatService := gfweixin.NewWeChatService("wx65064684d6c0f73f", "e20bf5f51062ab55ed2b1cec8e540502", nil)
	ctx := context.Background()
	var (
		accessToken string
		err         error
	)
	if accessToken, err = weChatService.GetStableAccessToken(ctx); err != nil {
		t.Logf("GetStableAccessToken err: %+v\n", err)
		return
	}
	t.Logf("accessToken: %v\n", accessToken)
	if accessToken, err = weChatService.GetStableAccessToken(ctx, true); err != nil {
		t.Logf("GetStableAccessToken err: %+v\n", err)
		return
	}
	t.Logf("accessToken: %v\n", accessToken)
}

func TestGetPhoneNumber(t *testing.T) {
	weChatService := gfweixin.NewWeChatService("wx65064684d6c0f73f", "e20bf5f51062ab55ed2b1cec8e540502", nil)
	ctx := context.Background()
	var (
		accessToken string
		err         error
	)
	if accessToken, err = weChatService.GetStableAccessToken(ctx); err != nil {
		t.Logf("GetStableAccessToken err: %+v\n", err)
		return
	}
	var (
		phoneNumber        string
		invalidAccessToken bool
	)
	if phoneNumber, invalidAccessToken, err = weChatService.GetPhoneNumber(ctx, "", accessToken); err != nil {
		t.Logf("GetPhoneNumber err: %+v\n", err)
		return
	}
	t.Logf("GetPhoneNumber invalidAccessToken: %v\n", invalidAccessToken)
	t.Logf("phoneNumber: %+v\n", phoneNumber)
}
