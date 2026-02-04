/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-08-26 10:51:39
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2026-02-04 15:26:44
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtktime_test

import (
	"github.com/liusuxian/go-toolkit/gtktime"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRemainingTimeUntilTomorrow(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(time.Hour, gtktime.RemainingTimeUntilTomorrow(time.Date(2022, 12, 30, 23, 0, 0, 0, time.UTC)))
	assert.Equal(time.Hour, gtktime.RemainingTimeUntilTomorrow(time.Date(2022, 12, 31, 23, 0, 0, 0, time.UTC)))
	assert.Equal(time.Hour, gtktime.RemainingTimeUntilTomorrow(time.Date(2023, 1, 1, 23, 0, 0, 0, time.UTC)))
}

func TestNextAligned(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(time.Date(2022, 12, 30, 0, 10, 0, 0, time.UTC), gtktime.NextAligned(time.Date(2022, 12, 30, 0, 7, 0, 0, time.UTC), 10*time.Minute))
	assert.Equal(time.Date(2022, 12, 30, 0, 20, 0, 0, time.UTC), gtktime.NextAligned(time.Date(2022, 12, 30, 0, 13, 0, 0, time.UTC), 10*time.Minute))
	assert.Equal(time.Date(2022, 12, 30, 0, 30, 0, 0, time.UTC), gtktime.NextAligned(time.Date(2022, 12, 30, 0, 21, 0, 0, time.UTC), 10*time.Minute))
	assert.Equal(time.Date(2022, 12, 30, 0, 40, 0, 0, time.UTC), gtktime.NextAligned(time.Date(2022, 12, 30, 0, 39, 0, 0, time.UTC), 10*time.Minute))
	assert.Equal(time.Date(2022, 12, 30, 0, 50, 0, 0, time.UTC), gtktime.NextAligned(time.Date(2022, 12, 30, 0, 40, 0, 0, time.UTC), 10*time.Minute))
	assert.Equal(time.Date(2022, 12, 30, 1, 0, 0, 0, time.UTC), gtktime.NextAligned(time.Date(2022, 12, 30, 0, 50, 0, 0, time.UTC), 10*time.Minute))
}
