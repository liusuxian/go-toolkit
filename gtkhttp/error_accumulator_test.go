/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-07 18:57:40
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-23 18:48:48
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp_test

import (
	"bytes"
	"errors"
	utils "github.com/liusuxian/go-toolkit/gtkhttp"
	"github.com/liusuxian/go-toolkit/gtkhttp/test"
	"testing"
)

func TestErrorAccumulatorBytes(t *testing.T) {
	accumulator := &utils.DefaultErrorAccumulator{
		Buffer: &bytes.Buffer{},
	}

	if errBytes := accumulator.Bytes(); len(errBytes) != 0 {
		t.Fatalf("did not return nil with empty bytes: %s", string(errBytes))
	}

	if err := accumulator.Write([]byte("{}")); err != nil {
		t.Fatalf("%+v", err)
	}

	if errBytes := accumulator.Bytes(); len(errBytes) == 0 {
		t.Fatalf("did not return error bytes when has error: %s", string(errBytes))
	}
}

func TestErrorByteWriteErrors(t *testing.T) {
	accumulator := &utils.DefaultErrorAccumulator{
		Buffer: &test.FailingErrorBuffer{},
	}

	if err := accumulator.Write([]byte("{")); !errors.Is(err, test.ErrTestErrorAccumulatorWriteFailed) {
		t.Fatalf("did not return error when write failed: %v", err)
	}
}
