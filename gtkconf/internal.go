/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-12-21 01:18:55
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-21 01:20:32
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkconf

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/liusuxian/go-toolkit/gtkconv"
	"strings"
	"unicode"
)

// parseSizeInBytes 将像1GB或12MB这样的字符串转换为无符号整数字节数
func parseSizeInBytes(sizeStr string) (s uint) {
	sizeStr = strings.TrimSpace(sizeStr)
	lastChar := len(sizeStr) - 1
	multiplier := uint(1)

	if lastChar > 0 {
		if sizeStr[lastChar] == 'b' || sizeStr[lastChar] == 'B' {
			if lastChar > 1 {
				switch unicode.ToLower(rune(sizeStr[lastChar-1])) {
				case 'k':
					multiplier = 1 << 10
					sizeStr = strings.TrimSpace(sizeStr[:lastChar-1])
				case 'm':
					multiplier = 1 << 20
					sizeStr = strings.TrimSpace(sizeStr[:lastChar-1])
				case 'g':
					multiplier = 1 << 30
					sizeStr = strings.TrimSpace(sizeStr[:lastChar-1])
				default:
					multiplier = 1
					sizeStr = strings.TrimSpace(sizeStr[:lastChar])
				}
			}
		}
	}

	size := max(gtkconv.ToInt(sizeStr), 0)
	return safeMul(uint(size), multiplier)
}

func safeMul(a, b uint) (s uint) {
	c := a * b
	if a > 1 && b > 1 && c/b != a {
		return 0
	}
	return c
}

// defaultDecoderConfig 默认的解码配置
func defaultDecoderConfig(output any) (opt DecoderConfigOption) {
	return func(dc *DecoderConfig) {
		dc.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			mapstructure.RecursiveStructToMapHookFunc(),
			mapstructure.StringToIPHookFunc(),
			mapstructure.StringToIPNetHookFunc(),
			mapstructure.StringToNetIPAddrHookFunc(),
			mapstructure.StringToNetIPAddrPortHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToTimeHookFunc("2006-01-02 15:04:05"),
			mapstructure.TextUnmarshallerHookFunc(),
		)
		dc.WeaklyTypedInput = true
		dc.Metadata = nil
		dc.Result = output
		dc.TagName = "json"
	}
}
