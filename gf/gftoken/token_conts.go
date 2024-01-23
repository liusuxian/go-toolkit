/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-20 15:38:07
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-24 01:10:38
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gftoken

import "fmt"

const (
	FAIL  = -1
	ERROR = -99
)

const (
	CacheModeCache   = 1
	CacheModeRedis   = 2
	CacheModeFile    = 3
	CacheModeFileDat = "gftoken.dat"

	MiddlewareTypeGroup  = 1
	MiddlewareTypeBind   = 2
	MiddlewareTypeGlobal = 3

	DefaultTimeout        = 10 * 24 * 60 * 60 * 1000
	DefaultCacheKey       = "gftoken:"
	DefaultTokenDelimiter = "_"
	DefaultEncryptKey     = "12345678912345678912345678912345"
	DefaultAuthFailMsg    = "请求错误或登录超时"

	TraceId = "d5dfce77cdff812161134e55de3c5207"

	KeyUserKey     = "userKey"
	KeyRefreshTime = "refreshTime"
	KeyCreateTime  = "createTime"
	KeyUuid        = "uuid"
	KeyData        = "data"
	KeyToken       = "token"
)

const (
	DefaultLogPrefix   = "[gftoken]" // 日志前缀
	MsgErrUserKeyEmpty = "userKey is empty"
	MsgErrReqMethod    = "request method is error! "
	MsgErrAuthHeader   = "Authorization : %s get token key fail"
	MsgErrTokenEmpty   = "token is empty"
	MsgErrTokenEncrypt = "token encrypt error"
	MsgErrTokenDecode  = "token decode error"
	MsgErrTokenLen     = "token len error"
	MsgErrAuthUuid     = "user auth uuid error"
)

func msgLog(msg string, params ...any) string {
	if len(params) == 0 {
		return DefaultLogPrefix + msg
	}
	return DefaultLogPrefix + fmt.Sprintf(msg, params...)
}
