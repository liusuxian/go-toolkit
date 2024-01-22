/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-22 23:27:03
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-22 23:30:12
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkjson_test

import (
	"github.com/liusuxian/go-toolkit/gtk/gtkjson"
	"testing"
)

func TestJsonMarshal(t *testing.T) {
	t.Log(gtkjson.JsonMarshal(map[string]any{
		"a": 1,
		"b": 1.2,
		"c": map[string]any{
			"d": 3,
		},
	}))
}

func TestMustJsonMarshal(t *testing.T) {
	t.Log(gtkjson.MustJsonMarshal(map[string]any{
		"a": 1,
		"b": 1.2,
		"c": map[string]any{
			"d": 3,
		},
	}))
}
