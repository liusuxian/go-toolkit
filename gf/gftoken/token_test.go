/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-20 15:38:07
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-24 01:19:14
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gftoken_test

import (
	"context"
	"github.com/liusuxian/go-toolkit/gf/gftoken"
	"testing"
)

func TestAuthPathGlobal(t *testing.T) {
	ctx := context.Background()

	t.Log("Global auth path test ")
	// 启动gtoken
	gfToken := &gftoken.Token{
		//Timeout:         10 * 1000,
		AuthPaths:        []string{"/user", "/system"},             // 这里是按照前缀拦截，拦截/user /user/list /user/add ...
		AuthExcludePaths: []string{"/user/info", "/system/user/*"}, // 不拦截路径  /user/info,/system/user/info,/system/user,
		MiddlewareType:   gftoken.MiddlewareTypeGlobal,             // 开启全局拦截
	}

	authPath(gfToken, t)
	flag := gfToken.AuthPath(ctx, "/test")
	if flag {
		t.Error("error:", "/test auth path error")
	}

}

func TestBindAuthPath(t *testing.T) {
	t.Log("Bind auth path test ")
	// 启动gtoken
	gfToken := &gftoken.Token{
		//Timeout:         10 * 1000,
		AuthPaths:        []string{"/user", "/system"},             // 这里是按照前缀拦截，拦截/user /user/list /user/add ...
		AuthExcludePaths: []string{"/user/info", "/system/user/*"}, // 不拦截路径  /user/info,/system/user/info,/system/user,
		MiddlewareType:   gftoken.MiddlewareTypeBind,               // 开启局部拦截
	}

	authPath(gfToken, t)
}

func TestGroupAuthPath(t *testing.T) {
	ctx := context.Background()

	t.Log("Group auth path test ")
	// 启动gtoken
	gfToken := &gftoken.Token{
		//Timeout:         10 * 1000,
		AuthExcludePaths: []string{"/login", "/user/info", "/system/user/*"}, // 不拦截路径  /user/info,/system/user/info,/system/user,
		MiddlewareType:   gftoken.MiddlewareTypeGroup,                        // 开启组拦截
	}

	flag := gfToken.AuthPath(ctx, "/login")
	if flag {
		t.Error("error:", "/login auth path error")
	}

	flag = gfToken.AuthPath(ctx, "/user/info")
	if flag {
		t.Error("error:", "/user/info auth path error")
	}

	flag = gfToken.AuthPath(ctx, "/system/user/info")
	if flag {
		t.Error("error:", "/system/user/info auth path error")
	}

	flag = gfToken.AuthPath(ctx, "/system/test")
	if !flag {
		t.Error("error:", "/system/test auth path error")
	}
}

func TestAuthPathNoExclude(t *testing.T) {
	ctx := context.Background()

	t.Log("auth no exclude path test ")
	// 启动gtoken
	gfToken := &gftoken.Token{
		//Timeout:         10 * 1000,
		AuthPaths:      []string{"/user", "/system"}, // 这里是按照前缀拦截，拦截/user /user/list /user/add ...
		MiddlewareType: gftoken.MiddlewareTypeGlobal, // 关闭全局拦截
	}

	authFlag := gfToken.AuthPath
	if authFlag(ctx, "/test") {
		t.Error(ctx, "error:", "/test auth path error")
	}
	if !authFlag(ctx, "/system/dept") {
		t.Error(ctx, "error:", "/system/dept auth path error")
	}

	if !authFlag(ctx, "/user/info") {
		t.Error(ctx, "error:", "/user/info auth path error")
	}

	if !authFlag(ctx, "/system/user") {
		t.Error(ctx, "error:", "/system/user auth path error")
	}
}

func TestAuthPathExclude(t *testing.T) {
	ctx := context.Background()

	t.Log("auth path test ")
	// 启动gtoken
	gfToken := &gftoken.Token{
		//Timeout:         10 * 1000,
		AuthPaths:        []string{"/*"},                           // 这里是按照前缀拦截，拦截/user /user/list /user/add ...
		AuthExcludePaths: []string{"/user/info", "/system/user/*"}, // 不拦截路径  /user/info,/system/user/info,/system/user,
		MiddlewareType:   gftoken.MiddlewareTypeGlobal,             // 开启全局拦截
	}

	authFlag := gfToken.AuthPath
	if !authFlag(ctx, "/test") {
		t.Error("error:", "/test auth path error")
	}
	if !authFlag(ctx, "//system/dept") {
		t.Error("error:", "/system/dept auth path error")
	}

	if authFlag(ctx, "/user/info") {
		t.Error("error:", "/user/info auth path error")
	}

	if authFlag(ctx, "/system/user") {
		t.Error("error:", "/system/user auth path error")
	}

	if authFlag(ctx, "/system/user/info") {
		t.Error("error:", "/system/user/info auth path error")
	}

}

func authPath(gfToken *gftoken.Token, t *testing.T) {
	ctx := context.Background()

	flag := gfToken.AuthPath(ctx, "/user/info")
	if flag {
		t.Error("error:", "/user/info auth path error")
	}

	flag = gfToken.AuthPath(ctx, "/system/user")
	if flag {
		t.Error("error:", "/system/user auth path error")
	}

	flag = gfToken.AuthPath(ctx, "/system/user/info")
	if flag {
		t.Error("error:", "/system/user/info auth path error")
	}

	flag = gfToken.AuthPath(ctx, "/system/dept")
	if !flag {
		t.Error("error:", "/system/dept auth path error")
	}

	flag = gfToken.AuthPath(ctx, "/user/list")
	if !flag {
		t.Error("error:", "/user/list auth path error")
	}

	flag = gfToken.AuthPath(ctx, "/user/add")
	if !flag {
		t.Error("error:", "/user/add auth path error")
	}
}

func TestEncryptDecryptToken(t *testing.T) {
	t.Log("encrypt and decrypt token test ")
	ctx := context.Background()

	gfToken := gftoken.Token{}
	gfToken.InitConfig()

	userKey := "123123"
	token := gfToken.EncryptToken(ctx, userKey, "")
	if !token.Success() {
		t.Error(token.Json())
	}
	t.Log(token.DataString())

	token2 := gfToken.DecryptToken(ctx, token.GetString("token"))
	if !token2.Success() {
		t.Error(token2.Json())
	}
	t.Log(token2.DataString())
	if userKey != token2.GetString("userKey") {
		t.Error("error:", "token decrypt userKey error")
	}
	if token.GetString("uuid") != token2.GetString("uuid") {
		t.Error("error:", "token decrypt uuid error")
	}

}

func BenchmarkEncryptDecryptToken(b *testing.B) {
	b.Log("encrypt and decrypt token test ")

	ctx := context.Background()
	gfToken := gftoken.Token{}
	gfToken.InitConfig()

	userKey := "123123"
	token := gfToken.EncryptToken(ctx, userKey, "")
	if !token.Success() {
		b.Error(token.Json())
	}
	b.Log(token.DataString())

	for i := 0; i < b.N; i++ {
		token2 := gfToken.DecryptToken(ctx, token.GetString("token"))
		if !token2.Success() {
			b.Error(token2.Json())
		}
		b.Log(token2.DataString())
		if userKey != token2.GetString("userKey") {
			b.Error("error:", "token decrypt userKey error")
		}
		if token.GetString("uuid") != token2.GetString("uuid") {
			b.Error("error:", "token decrypt uuid error")
		}
	}
}
