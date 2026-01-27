package logger

import (
	"time"
)

// LogConfig 统一的日志配置类
type LogConfig struct {
	// 基础配置
	Provider     string        `json:"provider" yaml:"provider"`     // 日志提供者名称
	Name         string        `json:"name" yaml:"name"`             // 日志名称
	Level        LogLevel      `json:"level" yaml:"level"`           // 日志级别
	Format       string        `json:"format" yaml:"format"`         // 日志格式（text/json）
	OutputPath   string        `json:"outputPath" yaml:"outputPath"` // 日志输出路径

	// 日志文件轮转配置
	MaxLogSize    int64         `json:"maxLogSize" yaml:"maxLogSize"`       // 单个日志文件最大大小（MB）
	MaxLogAge     time.Duration `json:"maxLogAge" yaml:"maxLogAge"`         // 日志文件最大保留时间
	MaxLogFiles   int           `json:"maxLogFiles" yaml:"maxLogFiles"`     // 最大保留日志文件数量
	CompressLogs  bool          `json:"compressLogs" yaml:"compressLogs"`   // 是否压缩旧日志
	MaxMessageSize int          `json:"maxMessageSize" yaml:"maxMessageSize"` // 单条日志最大大小（KB）

	// 额外配置
	ExtraConfig   map[string]interface{} `json:"extraConfig" yaml:"extraConfig"` // 额外的提供者特定配置
}

// NewLogConfig 创建默认的日志配置
func NewLogConfig() *LogConfig {
	return &LogConfig{
		Provider:     "console",
		Name:         "app",
		Level:        InfoLevel,
		Format:       "text",
		OutputPath:   "stdout",
		MaxLogSize:   100, // 默认100MB
		MaxLogAge:    7 * 24 * time.Hour, // 默认7天
		MaxLogFiles:  10, // 默认10个文件
		CompressLogs: false, // 默认不压缩
		MaxMessageSize: 0, // 默认不限制
		ExtraConfig:  make(map[string]interface{}),
	}
}

// WithProvider 设置日志提供者
func (c *LogConfig) WithProvider(provider string) *LogConfig {
	c.Provider = provider
	return c
}

// WithName 设置日志名称
func (c *LogConfig) WithName(name string) *LogConfig {
	c.Name = name
	return c
}

// WithLevel 设置日志级别
func (c *LogConfig) WithLevel(level LogLevel) *LogConfig {
	c.Level = level
	return c
}

// WithFormat 设置日志格式
func (c *LogConfig) WithFormat(format string) *LogConfig {
	c.Format = format
	return c
}

// WithOutputPath 设置日志输出路径
func (c *LogConfig) WithOutputPath(path string) *LogConfig {
	c.OutputPath = path
	return c
}

// WithMaxLogSize 设置单个日志文件最大大小（MB）
func (c *LogConfig) WithMaxLogSize(size int64) *LogConfig {
	c.MaxLogSize = size
	return c
}

// WithMaxLogAge 设置日志文件最大保留时间
func (c *LogConfig) WithMaxLogAge(age time.Duration) *LogConfig {
	c.MaxLogAge = age
	return c
}

// WithMaxLogFiles 设置最大保留日志文件数量
func (c *LogConfig) WithMaxLogFiles(files int) *LogConfig {
	c.MaxLogFiles = files
	return c
}

// WithCompressLogs 设置是否压缩旧日志
func (c *LogConfig) WithCompressLogs(compress bool) *LogConfig {
	c.CompressLogs = compress
	return c
}

// WithMaxMessageSize 设置单条日志最大大小（KB）
func (c *LogConfig) WithMaxMessageSize(size int) *LogConfig {
	c.MaxMessageSize = size
	return c
}

// WithExtraConfig 设置额外配置
func (c *LogConfig) WithExtraConfig(key string, value interface{}) *LogConfig {
	if c.ExtraConfig == nil {
		c.ExtraConfig = make(map[string]interface{})
	}
	c.ExtraConfig[key] = value
	return c
}

// WithExtraConfigs 设置多个额外配置
func (c *LogConfig) WithExtraConfigs(configs map[string]interface{}) *LogConfig {
	if c.ExtraConfig == nil {
		c.ExtraConfig = make(map[string]interface{})
	}
	for k, v := range configs {
		c.ExtraConfig[k] = v
	}
	return c
}

// ToOptions 将配置转换为选项函数
func (c *LogConfig) ToOptions() []Option {
	options := []Option{
		WithLevel(c.Level),
		WithFormat(c.Format),
		WithOutputPath(c.OutputPath),
		WithMaxLogSize(c.MaxLogSize),
		WithMaxLogAge(c.MaxLogAge),
		WithMaxLogFiles(c.MaxLogFiles),
		WithCompressLogs(c.CompressLogs),
		WithMaxMessageSize(c.MaxMessageSize),
		WithConfig(c.ExtraConfig),
	}
	return options
}

// Validate 验证配置的有效性
func (c *LogConfig) Validate() bool {
	// 验证提供者
	if c.Provider == "" {
		c.Provider = "console"
	}

	// 验证名称
	if c.Name == "" {
		c.Name = "app"
	}

	// 验证输出路径
	if c.OutputPath == "" {
		c.OutputPath = "stdout"
	}

	// 验证格式
	if c.Format == "" {
		c.Format = "text"
	} else if c.Format != "text" && c.Format != "json" {
		c.Format = "text"
	}

	// 验证文件轮转配置
	if c.MaxLogSize <= 0 {
		c.MaxLogSize = 100
	}

	if c.MaxLogAge <= 0 {
		c.MaxLogAge = 7 * 24 * time.Hour
	}

	if c.MaxLogFiles <= 0 {
		c.MaxLogFiles = 10
	}

	return true
}
