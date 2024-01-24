/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-04-12 18:19:13
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-24 19:56:12
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkconf

import (
	"github.com/fsnotify/fsnotify"
	"github.com/liusuxian/go-toolkit/gtkconv"
	"github.com/liusuxian/go-toolkit/gtkenv"
	"github.com/liusuxian/go-toolkit/gtkfile"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"strings"
	"time"
	"unicode"
)

// DecoderConfig 解码配置
type DecoderConfig = mapstructure.DecoderConfig

// DecoderConfigOption 解码配置选项
type DecoderConfigOption func(dc *DecoderConfig)

// Event 事件
type Event = fsnotify.Event

// 操作
const (
	Create = fsnotify.Create
	Write  = fsnotify.Write
	Remove = fsnotify.Remove
	Rename = fsnotify.Rename
	Chmod  = fsnotify.Chmod
)

// Config 配置结构
type Config struct {
	v *viper.Viper
}

// NewConfig 新建Config
func NewConfig(path string) (cfg *Config, err error) {
	v := viper.New()
	v.SetConfigFile(path)
	configType := gtkfile.ExtName(path)
	v.SetConfigType(configType)
	// 加载配置文件内容
	if err = v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = errors.Wrapf(err, "no such config file, path: %s, configType: %s", path, configType)
		} else {
			err = errors.Wrapf(err, "read config error, path: %s, configType: %s", path, configType)
		}
		return
	}
	cfg = &Config{v: v}
	return
}

// NewRemoteConfig 新建远程Config
func NewRemoteConfig(provider, endpoint, path, configType string) (cfg *Config, err error) {
	v := viper.New()
	v.AddRemoteProvider(provider, endpoint, path)
	v.SetConfigType(configType)
	// 加载配置文件内容
	if err = v.ReadRemoteConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = errors.Wrapf(err, "no such config file, provider: %s, endpoint: %s, path: %s, configType: %s", provider, endpoint, path, configType)
		} else {
			err = errors.Wrapf(err, "read config error, provider: %s, endpoint: %s, path: %s, configType: %s", provider, endpoint, path, configType)
		}
		return
	}
	cfg = &Config{v: v}
	return
}

// NewSecureRemoteConfig 新建远程Config
func NewSecureRemoteConfig(provider, endpoint, path, secretkeyring, configType string) (cfg *Config, err error) {
	v := viper.New()
	v.AddSecureRemoteProvider(provider, endpoint, path, secretkeyring)
	v.SetConfigType(configType)
	// 加载配置文件内容
	if err = v.ReadRemoteConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = errors.Wrapf(err, "no such config file, provider: %s, endpoint: %s, path: %s, configType: %s", provider, endpoint, path, configType)
		} else {
			err = errors.Wrapf(err, "read config error, provider: %s, endpoint: %s, path: %s, configType: %s", provider, endpoint, path, configType)
		}
		return
	}
	cfg = &Config{v: v}
	return
}

// Get 获取 value
func (c *Config) Get(key string) (val any) {
	return c.v.Get(key)
}

// GetBool 获取 bool
func (c *Config) GetBool(key string) (val bool) {
	return gtkconv.ToBool(c.v.Get(key))
}

// GetDuration 获取 Duration
func (c *Config) GetDuration(key string) (val time.Duration) {
	return gtkconv.ToDuration(c.v.Get(key))
}

// GetFloat32 获取 float32
func (c *Config) GetFloat32(key string) (val float32) {
	return gtkconv.ToFloat32(c.v.Get(key))
}

// GetFloat64 获取 float64
func (c *Config) GetFloat64(key string) (val float64) {
	return gtkconv.ToFloat64(c.v.Get(key))
}

// GetInt 获取 int
func (c *Config) GetInt(key string) (val int) {
	return gtkconv.ToInt(c.v.Get(key))
}

// GetInt8 获取 int8
func (c *Config) GetInt8(key string) (val int8) {
	return gtkconv.ToInt8(c.v.Get(key))
}

// GetInt16 获取 int16
func (c *Config) GetInt16(key string) (val int16) {
	return gtkconv.ToInt16(c.v.Get(key))
}

// GetInt32 获取 int32
func (c *Config) GetInt32(key string) (val int32) {
	return gtkconv.ToInt32(c.v.Get(key))
}

// GetInt64 获取 int64
func (c *Config) GetInt64(key string) (val int64) {
	return gtkconv.ToInt64(c.v.Get(key))
}

// GetAnySlice 获取 []any
func (c *Config) GetAnySlice(key string) (vals []any) {
	return gtkconv.ToSlice(c.v.Get(key))
}

// GetBoolSlice 获取 []bool
func (c *Config) GetBoolSlice(key string) (vals []bool) {
	return gtkconv.ToBoolSlice(c.v.Get(key))
}

// GetStringSlice 获取 []string
func (c *Config) GetStringSlice(key string) (vals []string) {
	return gtkconv.ToStringSlice(c.v.Get(key))
}

// GetIntSlice 获取 []int
func (c *Config) GetIntSlice(key string) (vals []int) {
	return gtkconv.ToIntSlice(c.v.Get(key))
}

// GetDurationSlice 获取 []time.Duration
func (c *Config) GetDurationSlice(key string) (vals []time.Duration) {
	return gtkconv.ToDurationSlice(c.v.Get(key))
}

// GetSizeInBytes 获取某个配置项对应的值所占用的内存大小（以字节为单位）
func (c *Config) GetSizeInBytes(key string) (val uint) {
	sizeStr := gtkconv.ToString(c.v.Get(key))
	return parseSizeInBytes(sizeStr)
}

// GetString 获取 string
func (c *Config) GetString(key string) (val string) {
	return gtkconv.ToString(c.v.Get(key))
}

// GetStringMap 获取 map[string]any
func (c *Config) GetStringMap(key string) (val map[string]any) {
	return gtkconv.ToStringMap(c.v.Get(key))
}

// GetStringMapBool 获取 map[string]bool
func (c *Config) GetStringMapBool(key string) (val map[string]bool) {
	return gtkconv.ToStringMapBool(c.v.Get(key))
}

// GetStringMapInt 获取 map[string]int
func (c *Config) GetStringMapInt(key string) (val map[string]int) {
	return gtkconv.ToStringMapInt(c.v.Get(key))
}

// GetStringMapInt64 获取 map[string]int64
func (c *Config) GetStringMapInt64(key string) (val map[string]int64) {
	return gtkconv.ToStringMapInt64(c.v.Get(key))
}

// GetStringMapString 获取 map[string]string
func (c *Config) GetStringMapString(key string) (val map[string]string) {
	return gtkconv.ToStringMapString(c.v.Get(key))
}

// GetStringMapStringSlice 获取 map[string][]string
func (c *Config) GetStringMapStringSlice(key string) (val map[string][]string) {
	return gtkconv.ToStringMapStringSlice(c.v.Get(key))
}

// GetTime 获取 Time
func (c *Config) GetTime(key string) (val time.Time) {
	return gtkconv.ToTime(c.v.Get(key))
}

// GetUint 获取 uint
func (c *Config) GetUint(key string) (val uint) {
	return gtkconv.ToUint(c.v.Get(key))
}

// GetUint8 获取 uint8
func (c *Config) GetUint8(key string) (val uint8) {
	return gtkconv.ToUint8(c.v.Get(key))
}

// GetUint16 获取 uint16
func (c *Config) GetUint16(key string) (val uint16) {
	return gtkconv.ToUint16(c.v.Get(key))
}

// GetUint32 获取 uint32
func (c *Config) GetUint32(key string) (val uint32) {
	return gtkconv.ToUint32(c.v.Get(key))
}

// GetUint64 获取 uint64
func (c *Config) GetUint64(key string) (val uint64) {
	return gtkconv.ToUint64(c.v.Get(key))
}

// InConfig 检查给定的键(或别名)是否在配置文件中
func (c *Config) InConfig(key string) (val bool) {
	return c.v.InConfig(key)
}

// IsSet 检查是否在任何数据位置设置了键。键不区分大小写
func (c *Config) IsSet(key string) (val bool) {
	return c.v.IsSet(key)
}

// OnConfigChange 设置当配置文件更改时调用的事件处理程序(只能用于本地配置文件的变更监听)
func (c *Config) OnConfigChange(run func(e Event)) {
	c.v.OnConfigChange(run)
}

// SetDefault 设置配置项的默认值，对键不区分大小写，仅当通过flag, config或ENV没有提供值时使用默认值
func (c *Config) SetDefault(key string, value any) {
	c.v.SetDefault(key, value)
}

// Sub 返回一个新的Config实例，表示这个实例的子树，对键不区分大小写
func (c *Config) Sub(key string) (conf *Config) {
	return &Config{v: c.v.Sub(key)}
}

// Struct 将配置解析为结构体，确保标签正确设置该结构的字段
func (c *Config) Struct(rawVal any, opts ...DecoderConfigOption) (err error) {
	newOpts := make([]viper.DecoderConfigOption, 0, len(opts)+1)
	newOpts = append(newOpts, viper.DecoderConfigOption(defaultDecoderConfig(rawVal)))
	for _, opt := range opts {
		newOpts = append(newOpts, viper.DecoderConfigOption(opt))
	}
	return c.v.Unmarshal(rawVal, newOpts...)
}

// StructExact 将配置解析为结构体，如果在目标结构体中字段不存在则报错
func (c *Config) StructExact(rawVal any, opts ...DecoderConfigOption) (err error) {
	newOpts := make([]viper.DecoderConfigOption, 0, len(opts)+1)
	newOpts = append(newOpts, viper.DecoderConfigOption(defaultDecoderConfig(rawVal)))
	for _, opt := range opts {
		newOpts = append(newOpts, viper.DecoderConfigOption(opt))
	}
	return c.v.UnmarshalExact(rawVal, newOpts...)
}

// StructKey 接收一个键并将其解析到结构体中，确保标签正确设置该结构的字段
func (c *Config) StructKey(key string, rawVal any, opts ...DecoderConfigOption) (err error) {
	newOpts := make([]viper.DecoderConfigOption, 0, len(opts)+1)
	newOpts = append(newOpts, viper.DecoderConfigOption(defaultDecoderConfig(rawVal)))
	for _, opt := range opts {
		newOpts = append(newOpts, viper.DecoderConfigOption(opt))
	}
	return c.v.UnmarshalKey(key, rawVal, newOpts...)
}

// WatchConfig 监视配置文件的变化
func (c *Config) WatchConfig() {
	c.v.WatchConfig()
}

// WatchRemoteConfig 监视远程配置文件的变化(阻塞式)
func (c *Config) WatchRemoteConfig() (err error) {
	return c.v.WatchRemoteConfig()
}

// WatchRemoteConfigOnChannel 监视远程配置文件的变化(非阻塞式)
func (c *Config) WatchRemoteConfigOnChannel() (err error) {
	return c.v.WatchRemoteConfigOnChannel()
}

// 默认配置（包含启动配置文件）
var defaultConfig *Config

func init() {
	v := viper.New()
	v.SetConfigName("config") // 设置配置文件名，不需要配置文件扩展名，配置文件的类型会自动根据扩展名自动匹配
	if gtkenv.Contains("GTK_CONFIG_NAME") {
		v.SetConfigName(gtkenv.Get("GTK_CONFIG_NAME")) // 设置配置文件名
	}
	v.AddConfigPath("./")                 // 设置配置文件的搜索目录
	v.AddConfigPath("./config/")          // 设置配置文件的搜索目录
	v.AddConfigPath("./manifest/config/") // 设置配置文件的搜索目录
	if gtkenv.Contains("GTK_CONFIG_FILE_PATH") {
		v.AddConfigPath(gtkenv.Get("GTK_CONFIG_FILE_PATH")) // 设置配置文件的搜索目录
	}
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		} else {
			panic(errors.Wrapf(err, "read default config error"))
		}
	}
	defaultConfig = &Config{v: v}
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
	newOpts := make([]viper.DecoderConfigOption, 0, len(opts))
	for _, opt := range opts {
		newOpts = append(newOpts, viper.DecoderConfigOption(opt))
	}
	return defaultConfig.v.Unmarshal(rawVal, newOpts...)
}

// StructExact 将配置解析为结构体，如果在目标结构体中字段不存在则报错
func StructExact(rawVal any, opts ...DecoderConfigOption) (err error) {
	newOpts := make([]viper.DecoderConfigOption, 0, len(opts))
	for _, opt := range opts {
		newOpts = append(newOpts, viper.DecoderConfigOption(opt))
	}
	return defaultConfig.v.UnmarshalExact(rawVal, newOpts...)
}

// StructKey 接收一个键并将其解析到结构体中
func StructKey(key string, rawVal any, opts ...DecoderConfigOption) (err error) {
	newOpts := make([]viper.DecoderConfigOption, 0, len(opts))
	for _, opt := range opts {
		newOpts = append(newOpts, viper.DecoderConfigOption(opt))
	}
	return defaultConfig.v.UnmarshalKey(key, rawVal, newOpts...)
}

// WatchConfig 监视配置文件的变化
func WatchConfig() {
	defaultConfig.v.WatchConfig()
}

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

	size := gtkconv.ToInt(sizeStr)
	if size < 0 {
		size = 0
	}

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
