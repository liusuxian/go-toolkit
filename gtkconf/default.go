/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-12-21 01:03:01
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-21 02:34:10
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkconf

import (
	"errors"
	"fmt"
	"github.com/liusuxian/go-toolkit/gtkconv"
	"github.com/spf13/viper"
	"time"
)

// 默认配置（包含启动配置文件）
var defaultConfig *Config

func init() {
	cfg, err := NewConfig()
	if err != nil {
		var notFoundErr viper.ConfigFileNotFoundError
		if !errors.As(err, &notFoundErr) {
			panic(fmt.Errorf("read default config error: %w", err))
		}
		cfg = &Config{v: viper.New()}
	}
	defaultConfig = cfg
}

// Get 获取 value
func Get(key string) (val any) {
	return defaultConfig.v.Get(key)
}

// GetBool 获取 bool
func GetBool(key string) (val bool) {
	return gtkconv.ToBool(defaultConfig.v.Get(key))
}

// GetDuration 获取 Duration
func GetDuration(key string) (val time.Duration) {
	return gtkconv.ToDuration(defaultConfig.v.Get(key))
}

// GetFloat32 获取 float32
func GetFloat32(key string) (val float32) {
	return gtkconv.ToFloat32(defaultConfig.v.Get(key))
}

// GetFloat64 获取 float64
func GetFloat64(key string) (val float64) {
	return gtkconv.ToFloat64(defaultConfig.v.Get(key))
}

// GetInt 获取 int
func GetInt(key string) (val int) {
	return gtkconv.ToInt(defaultConfig.v.Get(key))
}

// GetInt8 获取 int8
func GetInt8(key string) (val int8) {
	return gtkconv.ToInt8(defaultConfig.v.Get(key))
}

// GetInt16 获取 int16
func GetInt16(key string) (val int16) {
	return gtkconv.ToInt16(defaultConfig.v.Get(key))
}

// GetInt32 获取 int32
func GetInt32(key string) (val int32) {
	return gtkconv.ToInt32(defaultConfig.v.Get(key))
}

// GetInt64 获取 int64
func GetInt64(key string) (val int64) {
	return gtkconv.ToInt64(defaultConfig.v.Get(key))
}

// GetAnySlice 获取 []any
func GetAnySlice(key string) (vals []any) {
	return gtkconv.ToSlice(defaultConfig.v.Get(key))
}

// GetBoolSlice 获取 []bool
func GetBoolSlice(key string) (vals []bool) {
	return gtkconv.ToBoolSlice(defaultConfig.v.Get(key))
}

// GetStringSlice 获取 []string
func GetStringSlice(key string) (vals []string) {
	return gtkconv.ToStringSlice(defaultConfig.v.Get(key))
}

// GetIntSlice 获取 []int
func GetIntSlice(key string) (vals []int) {
	return gtkconv.ToIntSlice(defaultConfig.v.Get(key))
}

// GetDurationSlice 获取 []time.Duration
func GetDurationSlice(key string) (vals []time.Duration) {
	return gtkconv.ToDurationSlice(defaultConfig.v.Get(key))
}

// GetSizeInBytes 获取某个配置项对应的值所占用的内存大小（以字节为单位）
func GetSizeInBytes(key string) (val uint) {
	sizeStr := gtkconv.ToString(defaultConfig.v.Get(key))
	return parseSizeInBytes(sizeStr)
}

// GetString 获取 string
func GetString(key string) (val string) {
	return gtkconv.ToString(defaultConfig.v.Get(key))
}

// GetStringMap 获取 map[string]any
func GetStringMap(key string) (val map[string]any) {
	return gtkconv.ToStringMap(defaultConfig.v.Get(key))
}

// GetStringMapBool 获取 map[string]bool
func GetStringMapBool(key string) (val map[string]bool) {
	return gtkconv.ToStringMapBool(defaultConfig.v.Get(key))
}

// GetStringMapInt 获取 map[string]int
func GetStringMapInt(key string) (val map[string]int) {
	return gtkconv.ToStringMapInt(defaultConfig.v.Get(key))
}

// GetStringMapInt64 获取 map[string]int64
func GetStringMapInt64(key string) (val map[string]int64) {
	return gtkconv.ToStringMapInt64(defaultConfig.v.Get(key))
}

// GetStringMapString 获取 map[string]string
func GetStringMapString(key string) (val map[string]string) {
	return gtkconv.ToStringMapString(defaultConfig.v.Get(key))
}

// GetStringMapStringSlice 获取 map[string][]string
func GetStringMapStringSlice(key string) (val map[string][]string) {
	return gtkconv.ToStringMapStringSlice(defaultConfig.v.Get(key))
}

// GetTime 获取 Time
func GetTime(key string) (val time.Time) {
	return gtkconv.ToTime(defaultConfig.v.Get(key))
}

// GetUint 获取 uint
func GetUint(key string) (val uint) {
	return gtkconv.ToUint(defaultConfig.v.Get(key))
}

// GetUint8 获取 uint8
func GetUint8(key string) (val uint8) {
	return gtkconv.ToUint8(defaultConfig.v.Get(key))
}

// GetUint16 获取 uint16
func GetUint16(key string) (val uint16) {
	return gtkconv.ToUint16(defaultConfig.v.Get(key))
}

// GetUint32 获取 uint32
func GetUint32(key string) (val uint32) {
	return gtkconv.ToUint32(defaultConfig.v.Get(key))
}

// GetUint64 获取 uint64
func GetUint64(key string) (val uint64) {
	return gtkconv.ToUint64(defaultConfig.v.Get(key))
}

// InConfig 检查给定的键(或别名)是否在配置文件中
func InConfig(key string) (val bool) {
	return defaultConfig.v.InConfig(key)
}

// IsSet 检查是否在任何数据位置设置了键。键不区分大小写
func IsSet(key string) (val bool) {
	return defaultConfig.v.IsSet(key)
}

// OnConfigChange 设置当配置文件更改时调用的事件处理程序(只能用于本地配置文件的变更监听)
func OnConfigChange(run func(e Event)) {
	defaultConfig.v.OnConfigChange(run)
}

// SetDefault 设置配置项的默认值，对键不区分大小写，仅当通过flag, config或ENV没有提供值时使用默认值
func SetDefault(key string, value any) {
	defaultConfig.v.SetDefault(key, value)
}

// Sub 返回一个新的Config实例，表示这个实例的子树，对键不区分大小写
func Sub(key string) (conf *Config) {
	return &Config{v: defaultConfig.v.Sub(key)}
}

// Struct 将配置解析为结构体，确保标签正确设置该结构的字段
func Struct(rawVal any, opts ...DecoderConfigOption) (err error) {
	newOpts := make([]viper.DecoderConfigOption, 0, len(opts)+1)
	newOpts = append(newOpts, viper.DecoderConfigOption(defaultDecoderConfig(rawVal)))
	for _, opt := range opts {
		newOpts = append(newOpts, viper.DecoderConfigOption(opt))
	}
	return defaultConfig.v.Unmarshal(rawVal, newOpts...)
}

// StructExact 将配置解析为结构体，如果在目标结构体中字段不存在则报错
func StructExact(rawVal any, opts ...DecoderConfigOption) (err error) {
	newOpts := make([]viper.DecoderConfigOption, 0, len(opts)+1)
	newOpts = append(newOpts, viper.DecoderConfigOption(defaultDecoderConfig(rawVal)))
	for _, opt := range opts {
		newOpts = append(newOpts, viper.DecoderConfigOption(opt))
	}
	return defaultConfig.v.UnmarshalExact(rawVal, newOpts...)
}

// StructKey 接收一个键并将其解析到结构体中，确保标签正确设置该结构的字段
func StructKey(key string, rawVal any, opts ...DecoderConfigOption) (err error) {
	newOpts := make([]viper.DecoderConfigOption, 0, len(opts)+1)
	newOpts = append(newOpts, viper.DecoderConfigOption(defaultDecoderConfig(rawVal)))
	for _, opt := range opts {
		newOpts = append(newOpts, viper.DecoderConfigOption(opt))
	}
	return defaultConfig.v.UnmarshalKey(key, rawVal, newOpts...)
}

// WatchConfig 监视配置文件的变化
func WatchConfig() {
	defaultConfig.v.WatchConfig()
}
