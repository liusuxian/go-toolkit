/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-04-12 18:19:13
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-21 02:21:06
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkconf

import (
	"github.com/fsnotify/fsnotify"
	"github.com/go-viper/mapstructure/v2"
	"github.com/liusuxian/go-toolkit/gtkconv"
	"github.com/liusuxian/go-toolkit/gtkenv"
	"github.com/liusuxian/go-toolkit/internal/utils"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"time"
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

// Option 配置选项函数
type Option func(c *configBuilder)

// configBuilder 配置构建器
type configBuilder struct {
	viperOpts   []viper.Option // viper 选项
	configFile  string         // 配置文件路径（完整路径，包含文件名和扩展名）
	configName  string         // 配置文件名（不含扩展名）
	configType  string         // 配置文件类型（如：yaml, json, toml）
	configPaths []string       // 配置文件搜索路径
}

// WithConfigFile 设置配置文件路径（完整路径，包含文件名和扩展名）
func WithConfigFile(path string) (opt Option) {
	return func(c *configBuilder) {
		c.configFile = path
	}
}

// WithConfigName 设置配置文件名（不含扩展名）
func WithConfigName(name string) (opt Option) {
	return func(c *configBuilder) {
		c.configName = name
	}
}

// WithConfigType 设置配置文件类型（如：yaml, json, toml）
func WithConfigType(configType string) (opt Option) {
	return func(c *configBuilder) {
		c.configType = configType
	}
}

// WithConfigPath 添加配置文件搜索路径（可多次调用）
func WithConfigPath(path string) (opt Option) {
	return func(c *configBuilder) {
		c.configPaths = append(c.configPaths, path)
	}
}

// WithKeyDelimiter 设置键分隔符，默认值为 "."
func WithKeyDelimiter(d string) (opt Option) {
	return func(c *configBuilder) {
		c.viperOpts = append(c.viperOpts, viper.KeyDelimiter(d))
	}
}

// Config 配置结构
type Config struct {
	v *viper.Viper
}

// NewConfig 新建Config
func NewConfig(opts ...Option) (cfg *Config, err error) {
	// 构建配置
	builder := &configBuilder{
		viperOpts:   make([]viper.Option, 0),
		configPaths: make([]string, 0),
	}
	// 应用所有选项
	for _, opt := range opts {
		opt(builder)
	}
	// 创建 viper 实例
	var v *viper.Viper
	if len(builder.viperOpts) > 0 {
		v = viper.NewWithOptions(builder.viperOpts...)
	} else {
		v = viper.New()
	}
	// 配置文件路径
	if builder.configFile != "" {
		// 使用完整路径
		v.SetConfigFile(builder.configFile)
		if builder.configType != "" {
			v.SetConfigType(builder.configType)
		} else {
			if configType := utils.ExtName(builder.configFile); configType != "" {
				v.SetConfigType(configType)
			}
		}
	} else {
		// 使用名称 + 搜索路径
		if builder.configName != "" {
			v.SetConfigName(builder.configName)
		} else {
			if gtkenv.Contains("GTK_CONFIG_NAME") {
				v.SetConfigName(gtkenv.Get("GTK_CONFIG_NAME")) // 从环境变量中获取配置文件名
			} else {
				v.SetConfigName("config") // 默认值
			}
		}
		// 配置文件类型
		if builder.configType != "" {
			v.SetConfigType(builder.configType)
		}
		// 添加搜索路径
		if len(builder.configPaths) > 0 {
			for _, path := range builder.configPaths {
				v.AddConfigPath(path)
			}
		} else {
			if gtkenv.Contains("GTK_CONFIG_FILE_PATH") {
				v.AddConfigPath(gtkenv.Get("GTK_CONFIG_FILE_PATH")) // 从环境变量中获取配置文件的搜索目录
			} else {
				// 默认搜索路径
				v.AddConfigPath("./")
				v.AddConfigPath("./config/")
				v.AddConfigPath("./manifest/config/")
			}
		}
	}
	// 加载配置文件内容
	if err = v.ReadInConfig(); err != nil {
		return
	}
	cfg = &Config{v: v}
	return
}

// NewRemoteConfig 新建远程Config
func NewRemoteConfig(provider, endpoint, path string, opts ...Option) (cfg *Config, err error) {
	// 构建配置
	builder := &configBuilder{
		viperOpts:   make([]viper.Option, 0),
		configPaths: make([]string, 0),
	}
	// 应用所有选项
	for _, opt := range opts {
		opt(builder)
	}
	// 创建 viper 实例
	var v *viper.Viper
	if len(builder.viperOpts) > 0 {
		v = viper.NewWithOptions(builder.viperOpts...)
	} else {
		v = viper.New()
	}
	// 添加远程配置提供者
	v.AddRemoteProvider(provider, endpoint, path)
	// 设置配置类型
	if builder.configType != "" {
		v.SetConfigType(builder.configType)
	} else {
		if configType := utils.ExtName(path); configType != "" {
			v.SetConfigType(configType)
		}
	}
	// 加载配置文件内容
	if err = v.ReadRemoteConfig(); err != nil {
		return
	}
	cfg = &Config{v: v}
	return
}

// NewSecureRemoteConfig 新建远程Config
func NewSecureRemoteConfig(provider, endpoint, path, secretkeyring string, opts ...Option) (cfg *Config, err error) {
	// 构建配置
	builder := &configBuilder{
		viperOpts:   make([]viper.Option, 0),
		configPaths: make([]string, 0),
	}
	// 应用所有选项
	for _, opt := range opts {
		opt(builder)
	}
	// 创建 viper 实例
	var v *viper.Viper
	if len(builder.viperOpts) > 0 {
		v = viper.NewWithOptions(builder.viperOpts...)
	} else {
		v = viper.New()
	}
	// 添加安全远程配置提供者
	v.AddSecureRemoteProvider(provider, endpoint, path, secretkeyring)
	// 设置配置类型
	if builder.configType != "" {
		v.SetConfigType(builder.configType)
	} else {
		if configType := utils.ExtName(path); configType != "" {
			v.SetConfigType(configType)
		}
	}
	// 加载配置文件内容
	if err = v.ReadRemoteConfig(); err != nil {
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
