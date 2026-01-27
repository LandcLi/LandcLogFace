package tests

import (
	"context"
	"testing"
	"time"

	"github.com/LandcLi/LandcLogFace/internal/logger"
)

// CustomLogger 自定义日志实现
type CustomLogger struct {
	name string
}

// NewCustomLogger 创建自定义日志实例
func NewCustomLogger(name string) *CustomLogger {
	return &CustomLogger{name: name}
}

// Debug 输出调试级日志
func (c *CustomLogger) Debug(msg string, fields ...logger.Field) {}

// Debugf 输出格式化的调试级日志
func (c *CustomLogger) Debugf(format string, args ...interface{}) {}

// Info 输出信息级日志
func (c *CustomLogger) Info(msg string, fields ...logger.Field) {}

// Infof 输出格式化的信息级日志
func (c *CustomLogger) Infof(format string, args ...interface{}) {}

// Warn 输出警告级日志
func (c *CustomLogger) Warn(msg string, fields ...logger.Field) {}

// Warnf 输出格式化的警告级日志
func (c *CustomLogger) Warnf(format string, args ...interface{}) {}

// Error 输出错误级日志
func (c *CustomLogger) Error(msg string, fields ...logger.Field) {}

// Errorf 输出格式化的错误级日志
func (c *CustomLogger) Errorf(format string, args ...interface{}) {}

// Fatal 输出致命级日志并退出程序
func (c *CustomLogger) Fatal(msg string, fields ...logger.Field) {}

// Fatalf 输出格式化的致命级日志并退出程序
func (c *CustomLogger) Fatalf(format string, args ...interface{}) {}

// Panic 输出恐慌级日志并触发panic
func (c *CustomLogger) Panic(msg string, fields ...logger.Field) {}

// Panicf 输出格式化的恐慌级日志并触发panic
func (c *CustomLogger) Panicf(format string, args ...interface{}) {}

// WithFields 添加字段到日志
func (c *CustomLogger) WithFields(fields ...logger.Field) logger.Logger {
	return c
}

// WithField 添加单个字段到日志
func (c *CustomLogger) WithField(key string, value interface{}) logger.Logger {
	return c
}

// WithContext 添加上下文到日志
func (c *CustomLogger) WithContext(ctx context.Context) logger.Logger {
	return c
}

// WithError 添加错误信息到日志
func (c *CustomLogger) WithError(err error) logger.Logger {
	return c
}

// WithTime 添加时间到日志
func (c *CustomLogger) WithTime(t time.Time) logger.Logger {
	return c
}

// SetLevel 设置日志级别
func (c *CustomLogger) SetLevel(level logger.LogLevel) {}

// GetLevel 获取日志级别
func (c *CustomLogger) GetLevel() logger.LogLevel {
	return logger.InfoLevel
}

// IsDebugEnabled 检查调试级别是否启用
func (c *CustomLogger) IsDebugEnabled() bool {
	return false
}

// IsInfoEnabled 检查信息级别是否启用
func (c *CustomLogger) IsInfoEnabled() bool {
	return true
}

// IsWarnEnabled 检查警告级别是否启用
func (c *CustomLogger) IsWarnEnabled() bool {
	return true
}

// IsErrorEnabled 检查错误级别是否启用
func (c *CustomLogger) IsErrorEnabled() bool {
	return true
}

// IsFatalEnabled 检查致命级别是否启用
func (c *CustomLogger) IsFatalEnabled() bool {
	return true
}

// IsPanicEnabled 检查恐慌级别是否启用
func (c *CustomLogger) IsPanicEnabled() bool {
	return true
}

// Sync 同步日志
func (c *CustomLogger) Sync() error {
	return nil
}

// CustomLoggerProvider 自定义日志提供者
type CustomLoggerProvider struct{}

// Create 创建日志实例
func (p *CustomLoggerProvider) Create(name string, opts ...logger.Option) logger.Logger {
	return NewCustomLogger(name)
}

// CreateWithConfig 根据配置创建日志实例
func (p *CustomLoggerProvider) CreateWithConfig(name string, config map[string]interface{}) logger.Logger {
	return NewCustomLogger(name)
}

func TestCustomLoggerProvider(t *testing.T) {
	// 注册自定义提供者
	logger.GetLogFactory().RegisterProvider("custom", &CustomLoggerProvider{})

	// 使用自定义提供者
	customLogger := logger.GetLogFactory().CreateLoggerWithProvider("app", "custom")
	if customLogger == nil {
		t.Fatal("自定义日志提供者创建失败")
	}

	customLogger.Info("使用自定义日志提供者")

	// 注销自定义提供者
	logger.GetLogFactory().UnregisterProvider("custom")
}
