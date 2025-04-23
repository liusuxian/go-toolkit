/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-07 19:27:37
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-07 19:27:39
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package checks

import (
	"errors"
	"testing"
)

func NoError(t *testing.T, err error, message ...string) {
	t.Helper()
	if err != nil {
		t.Error(err, message)
	}
}

func NoErrorF(t *testing.T, err error, message ...string) {
	t.Helper()
	if err != nil {
		t.Fatal(err, message)
	}
}

func HasError(t *testing.T, err error, message ...string) {
	t.Helper()
	if err == nil {
		t.Error(err, message)
	}
}

func ErrorIs(t *testing.T, err, target error, message ...string) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatal(message)
	}
}

func ErrorIsF(t *testing.T, err, target error, format string, message ...string) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf(format, message)
	}
}

func ErrorIsNot(t *testing.T, err, target error, message ...string) {
	t.Helper()
	if errors.Is(err, target) {
		t.Fatal(message)
	}
}

func ErrorIsNotf(t *testing.T, err, target error, format string, message ...string) {
	t.Helper()
	if errors.Is(err, target) {
		t.Fatalf(format, message)
	}
}
