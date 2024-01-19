/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-19 22:57:39
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-19 23:15:54
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package tools_test

import (
	"go-toolkit/tools"
	"testing"
)

func TestJsonMarshal(t *testing.T) {
	t.Log(tools.JsonMarshal(map[string]any{
		"a": 1,
		"b": 1.2,
		"c": map[string]any{
			"d": 3,
		},
	}))
}

func TestMustJsonMarshal(t *testing.T) {
	t.Log(tools.MustJsonMarshal(map[string]any{
		"a": 1,
		"b": 1.2,
		"c": map[string]any{
			"d": 3,
		},
	}))
}
