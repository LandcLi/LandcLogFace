package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// ConsoleLogger 默认的控制台日志适配器
type ConsoleLogger struct {
	level          LogLevel
	fields         []Field
	ctx            context.Context
	logger         *log.Logger
	name           string
	format         string // 日志格式（text/json）
	maxMessageSize int    // 单条日志最大大小（KB）
}

// NewConsoleLogger 创建控制台日志实例
func NewConsoleLogger(name string, opts ...Option) *ConsoleLogger {
	options := &LoggerOptions{
		Level:          InfoLevel,
		Format:         "text",
		OutputPath:     "stdout",
		MaxLogSize:     100,                // 默认100MB
		MaxLogAge:      7 * 24 * time.Hour, // 默认7天
		MaxLogFiles:    10,                 // 默认10个文件
		CompressLogs:   false,              // 默认不压缩
		MaxMessageSize: 0,                  // 默认不限制
		Config:         make(map[string]interface{}),
	}

	for _, opt := range opts {
		opt(options)
	}

	var output interface {
		Write(p []byte) (n int, err error)
	}
	if options.OutputPath == "stdout" {
		output = os.Stdout
	} else {
		// 使用lumberjack进行日志轮转
		output = &lumberjack.Logger{
			Filename:   options.OutputPath,
			MaxSize:    int(options.MaxLogSize),             // MB
			MaxAge:     int(options.MaxLogAge.Hours() / 24), // 天
			MaxBackups: options.MaxLogFiles,
			Compress:   options.CompressLogs,
		}
	}

	return &ConsoleLogger{
		level:          options.Level,
		fields:         make([]Field, 0),
		ctx:            context.Background(),
		logger:         log.New(output, "", 0),
		name:           name,
		format:         options.Format,
		maxMessageSize: options.MaxMessageSize,
	}
}

// SetLevel 设置日志级别
func (c *ConsoleLogger) SetLevel(level LogLevel) {
	c.level = level
}

// GetLevel 获取当前日志级别
func (c *ConsoleLogger) GetLevel() LogLevel {
	return c.level
}

// limitMessageSize 限制日志消息大小
func (c *ConsoleLogger) limitMessageSize(msg string) string {
	if c.maxMessageSize > 0 {
		maxSize := c.maxMessageSize * 1024 // 转换为字节
		if len(msg) > maxSize {
			return msg[:maxSize-3] + "..."
		}
	}
	return msg
}

// formatMessage 格式化日志消息
func (c *ConsoleLogger) formatMessage(level LogLevel, msg string, fields []Field) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	allFields := append(c.fields, fields...)

	if c.format == "json" {
		// 构建JSON格式的日志
		jsonFields := make(map[string]interface{})
		jsonFields["time"] = timestamp
		jsonFields["level"] = level.String()
		jsonFields["logger"] = c.name
		jsonFields["msg"] = msg

		// 添加所有字段
		for _, field := range allFields {
			jsonFields[field.Key] = field.Value
		}

		// 转换为JSON字符串
		jsonBytes, err := json.Marshal(jsonFields)
		if err != nil {
			// 如果JSON转换失败，回退到文本格式
			fieldStr := ""
			for _, field := range allFields {
				fieldStr += fmt.Sprintf(" %s=%v", field.Key, field.Value)
			}
			formattedMsg := fmt.Sprintf("%s [%s] [%s] %s%s", timestamp, level.String(), c.name, msg, fieldStr)
			return c.limitMessageSize(formattedMsg)
		}

		formattedMsg := string(jsonBytes)
		return c.limitMessageSize(formattedMsg)
	} else {
		// 文本格式
		fieldStr := ""
		for _, field := range allFields {
			fieldStr += fmt.Sprintf(" %s=%v", field.Key, field.Value)
		}

		formattedMsg := fmt.Sprintf("%s [%s] [%s] %s%s", timestamp, level.String(), c.name, msg, fieldStr)
		return c.limitMessageSize(formattedMsg)
	}
}

// Debug 输出调试级日志
func (c *ConsoleLogger) Debug(msg string, fields ...Field) {
	if c.level <= DebugLevel {
		c.logger.Println(c.formatMessage(DebugLevel, msg, fields))
	}
}

// Debugf 输出格式化的调试级日志
func (c *ConsoleLogger) Debugf(format string, args ...interface{}) {
	if c.level <= DebugLevel {
		msg := fmt.Sprintf(format, args...)
		c.logger.Println(c.formatMessage(DebugLevel, msg, nil))
	}
}

// Info 输出信息级日志
func (c *ConsoleLogger) Info(msg string, fields ...Field) {
	if c.level <= InfoLevel {
		c.logger.Println(c.formatMessage(InfoLevel, msg, fields))
	}
}

// Infof 输出格式化的信息级日志
func (c *ConsoleLogger) Infof(format string, args ...interface{}) {
	if c.level <= InfoLevel {
		msg := fmt.Sprintf(format, args...)
		c.logger.Println(c.formatMessage(InfoLevel, msg, nil))
	}
}

// Warn 输出警告级日志
func (c *ConsoleLogger) Warn(msg string, fields ...Field) {
	if c.level <= WarnLevel {
		c.logger.Println(c.formatMessage(WarnLevel, msg, fields))
	}
}

// Warnf 输出格式化的警告级日志
func (c *ConsoleLogger) Warnf(format string, args ...interface{}) {
	if c.level <= WarnLevel {
		msg := fmt.Sprintf(format, args...)
		c.logger.Println(c.formatMessage(WarnLevel, msg, nil))
	}
}

// Error 输出错误级日志
func (c *ConsoleLogger) Error(msg string, fields ...Field) {
	if c.level <= ErrorLevel {
		c.logger.Println(c.formatMessage(ErrorLevel, msg, fields))
	}
}

// Errorf 输出格式化的错误级日志
func (c *ConsoleLogger) Errorf(format string, args ...interface{}) {
	if c.level <= ErrorLevel {
		msg := fmt.Sprintf(format, args...)
		c.logger.Println(c.formatMessage(ErrorLevel, msg, nil))
	}
}

// Fatal 输出致命级日志并退出程序
func (c *ConsoleLogger) Fatal(msg string, fields ...Field) {
	if c.level <= FatalLevel {
		c.logger.Println(c.formatMessage(FatalLevel, msg, fields))
		os.Exit(1)
	}
}

// Fatalf 输出格式化的致命级日志并退出程序
func (c *ConsoleLogger) Fatalf(format string, args ...interface{}) {
	if c.level <= FatalLevel {
		msg := fmt.Sprintf(format, args...)
		c.logger.Println(c.formatMessage(FatalLevel, msg, nil))
		os.Exit(1)
	}
}

// Panic 输出恐慌级日志并触发panic
func (c *ConsoleLogger) Panic(msg string, fields ...Field) {
	if c.level <= PanicLevel {
		msg := c.formatMessage(PanicLevel, msg, fields)
		c.logger.Println(msg)
		panic(msg)
	}
}

// Panicf 输出格式化的恐慌级日志并触发panic
func (c *ConsoleLogger) Panicf(format string, args ...interface{}) {
	if c.level <= PanicLevel {
		msg := fmt.Sprintf(format, args...)
		fullMsg := c.formatMessage(PanicLevel, msg, nil)
		c.logger.Println(fullMsg)
		panic(fullMsg)
	}
}

// WithFields 添加字段到日志
func (c *ConsoleLogger) WithFields(fields ...Field) Logger {
	newLogger := *c
	newLogger.fields = append(newLogger.fields, fields...)
	return &newLogger
}

// WithField 添加单个字段到日志
func (c *ConsoleLogger) WithField(key string, value interface{}) Logger {
	return c.WithFields(Field{Key: key, Value: value})
}

// WithContext 添加上下文到日志
func (c *ConsoleLogger) WithContext(ctx context.Context) Logger {
	newLogger := *c
	newLogger.ctx = ctx
	return &newLogger
}

// WithError 添加错误信息到日志
func (c *ConsoleLogger) WithError(err error) Logger {
	return c.WithField("error", err)
}

// WithTime 添加时间到日志
func (c *ConsoleLogger) WithTime(t time.Time) Logger {
	return c.WithField("time", t)
}

// IsDebugEnabled 检查调试级别是否启用
func (c *ConsoleLogger) IsDebugEnabled() bool {
	return c.level <= DebugLevel
}

// IsInfoEnabled 检查信息级别是否启用
func (c *ConsoleLogger) IsInfoEnabled() bool {
	return c.level <= InfoLevel
}

// IsWarnEnabled 检查警告级别是否启用
func (c *ConsoleLogger) IsWarnEnabled() bool {
	return c.level <= WarnLevel
}

// IsErrorEnabled 检查错误级别是否启用
func (c *ConsoleLogger) IsErrorEnabled() bool {
	return c.level <= ErrorLevel
}

// IsFatalEnabled 检查致命级别是否启用
func (c *ConsoleLogger) IsFatalEnabled() bool {
	return c.level <= FatalLevel
}

// IsPanicEnabled 检查恐慌级别是否启用
func (c *ConsoleLogger) IsPanicEnabled() bool {
	return c.level <= PanicLevel
}

// Sync 刷新日志缓冲区
func (c *ConsoleLogger) Sync() error {
	return nil
}

// ConsoleLoggerProvider 控制台日志提供者
type ConsoleLoggerProvider struct{}

// NewConsoleLoggerProvider 创建控制台日志提供者
func NewConsoleLoggerProvider() *ConsoleLoggerProvider {
	return &ConsoleLoggerProvider{}
}

// Create 创建日志实例
func (p *ConsoleLoggerProvider) Create(name string) Logger {
	return NewConsoleLogger(name)
}

// CreateWithConfig 根据配置创建日志实例
func (p *ConsoleLoggerProvider) CreateWithConfig(name string, config map[string]interface{}) Logger {
	var level LogLevel
	if lvl, ok := config["level"].(LogLevel); ok {
		level = lvl
	} else {
		level = InfoLevel
	}

	var format string
	if fmt, ok := config["format"].(string); ok {
		format = fmt
	} else {
		format = "text"
	}

	var outputPath string
	if path, ok := config["outputPath"].(string); ok {
		outputPath = path
	} else {
		outputPath = "stdout"
	}

	var maxLogSize int64
	if size, ok := config["maxLogSize"].(int64); ok {
		maxLogSize = size
	} else {
		maxLogSize = 100
	}

	var maxLogAge time.Duration
	if age, ok := config["maxLogAge"].(time.Duration); ok {
		maxLogAge = age
	} else {
		maxLogAge = 7 * 24 * time.Hour
	}

	var maxLogFiles int
	if files, ok := config["maxLogFiles"].(int); ok {
		maxLogFiles = files
	} else {
		maxLogFiles = 10
	}

	var compressLogs bool
	if compress, ok := config["compressLogs"].(bool); ok {
		compressLogs = compress
	} else {
		compressLogs = false
	}

	var maxMessageSize int
	if size, ok := config["maxMessageSize"].(int); ok {
		maxMessageSize = size
	} else {
		maxMessageSize = 0
	}

	return NewConsoleLogger(name,
		WithLevel(level),
		WithFormat(format),
		WithOutputPath(outputPath),
		WithMaxLogSize(maxLogSize),
		WithMaxLogAge(maxLogAge),
		WithMaxLogFiles(maxLogFiles),
		WithCompressLogs(compressLogs),
		WithMaxMessageSize(maxMessageSize),
		WithConfig(config),
	)
}
