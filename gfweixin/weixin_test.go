/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-01-19 22:29:06
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-20 00:36:27
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gfweixin_test

import (
	"context"
	"github.com/liusuxian/go-toolkit/gfweixin"
	"testing"
)

func TestGetStableAccessToken(t *testing.T) {
	ctx := context.Background()
	var (
		resMap map[string]any
		err    error
	)
	if resMap, err = gfweixin.GetStableAccessToken(ctx, "wx65064684d6c0f73f", "e20bf5f51062ab55ed2b1cec8e540502"); err != nil {
		t.Logf("err: %+v\n", err)
	}
	t.Logf("resMap111: %+v\n", resMap)
	if resMap, err = gfweixin.GetStableAccessToken(ctx, "wx65064684d6c0f73f", "e20bf5f51062ab55ed2b1cec8e540502", true); err != nil {
		t.Logf("err: %+v\n", err)
	}
	t.Logf("resMap222: %+v\n", resMap)
}
