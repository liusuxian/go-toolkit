/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-07 19:03:13
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-07 19:07:00
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package test

import "errors"

var (
	ErrTestErrorAccumulatorWriteFailed = errors.New("test error accumulator failed")
)

type FailingErrorBuffer struct{}

func (b *FailingErrorBuffer) Write(_ []byte) (n int, err error) {
	return 0, ErrTestErrorAccumulatorWriteFailed
}

func (b *FailingErrorBuffer) Len() int {
	return 0
}

func (b *FailingErrorBuffer) Bytes() []byte {
	return []byte{}
}
