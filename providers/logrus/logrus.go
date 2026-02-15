package logrus

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/LandcLi/landc-logface/internal/logger"
)

// LogrusLogger logrus日志库适配器
type LogrusLogger struct {
	logger         *logrus.Logger
	level          logger.LogLevel
	fields         []logger.Field
	ctx            context.Context
	name           string
	maxMessageSize int // 单条日志最大大小（KB）
}

// NewLogrusLogger 创建logrus日志实例
func NewLogrusLogger(name string, opts ...logger.Option) *LogrusLogger {
	options := &logger.LoggerOptions{
		Level:          logger.InfoLevel,
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

	// 创建logrus实例
	logrusLogger := logrus.New()

	// 设置日志级别
	logrusLevel := logrus.InfoLevel
	switch options.Level {
	case logger.DebugLevel:
		logrusLevel = logrus.DebugLevel
	case logger.InfoLevel:
		logrusLevel = logrus.InfoLevel
	case logger.WarnLevel:
		logrusLevel = logrus.WarnLevel
	case logger.ErrorLevel:
		logrusLevel = logrus.ErrorLevel
	case logger.FatalLevel, logger.PanicLevel:
		logrusLevel = logrus.FatalLevel
	}
	logrusLogger.SetLevel(logrusLevel)

	// 设置日志格式
	if options.Format == "json" {
		logrusLogger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		logrusLogger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		})
	}

	// 设置输出
	if options.OutputPath == "stdout" || options.OutputPath == "" {
		logrusLogger.SetOutput(os.Stdout)
	} else {
		logrusLogger.SetOutput(&lumberjack.Logger{
			Filename:   options.OutputPath,
			MaxSize:    int(options.MaxLogSize),
			MaxBackups: options.MaxLogFiles,
			MaxAge:     int(options.MaxLogAge / (24 * time.Hour)),
			Compress:   options.CompressLogs,
		})
	}

	return &LogrusLogger{
		logger:         logrusLogger,
		level:          options.Level,
		name:           name,
		maxMessageSize: options.MaxMessageSize,
	}
}

// limitMessageSize 限制日志消息大小
func (l *LogrusLogger) limitMessageSize(msg string) string {
	if l.maxMessageSize > 0 {
		maxSize := l.maxMessageSize * 1024 // 转换为字节
		if len(msg) > maxSize {
			return msg[:maxSize-3] + "..."
		}
	}
	return msg
}

// Debug 输出调试级日志
func (l *LogrusLogger) Debug(msg string, fields ...logger.Field) {
	msg = l.limitMessageSize(msg)
	entry := l.logger.WithFields(l.convertFields(append(l.fields, fields...)))
	entry.Debug(msg)
}

// Debugf 输出格式化的调试级日志
func (l *LogrusLogger) Debugf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	msg = l.limitMessageSize(msg)
	l.logger.Debugf(msg)
}

// Info 输出信息级日志
func (l *LogrusLogger) Info(msg string, fields ...logger.Field) {
	msg = l.limitMessageSize(msg)
	entry := l.logger.WithFields(l.convertFields(append(l.fields, fields...)))
	entry.Info(msg)
}

// Infof 输出格式化的信息级日志
func (l *LogrusLogger) Infof(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	msg = l.limitMessageSize(msg)
	l.logger.Infof(msg)
}

// Warn 输出警告级日志
func (l *LogrusLogger) Warn(msg string, fields ...logger.Field) {
	msg = l.limitMessageSize(msg)
	entry := l.logger.WithFields(l.convertFields(append(l.fields, fields...)))
	entry.Warn(msg)
}

// Warnf 输出格式化的警告级日志
func (l *LogrusLogger) Warnf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	msg = l.limitMessageSize(msg)
	l.logger.Warnf(msg)
}

// Error 输出错误级日志
func (l *LogrusLogger) Error(msg string, fields ...logger.Field) {
	msg = l.limitMessageSize(msg)
	entry := l.logger.WithFields(l.convertFields(append(l.fields, fields...)))
	entry.Error(msg)
}

// Errorf 输出格式化的错误级日志
func (l *LogrusLogger) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	msg = l.limitMessageSize(msg)
	l.logger.Errorf(msg)
}

// Fatal 输出致命级日志并退出程序
func (l *LogrusLogger) Fatal(msg string, fields ...logger.Field) {
	msg = l.limitMessageSize(msg)
	entry := l.logger.WithFields(l.convertFields(append(l.fields, fields...)))
	entry.Fatal(msg)
}

// Fatalf 输出格式化的致命级日志并退出程序
func (l *LogrusLogger) Fatalf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	msg = l.limitMessageSize(msg)
	l.logger.Fatalf(msg)
}

// Panic 输出恐慌级日志并触发panic
func (l *LogrusLogger) Panic(msg string, fields ...logger.Field) {
	msg = l.limitMessageSize(msg)
	entry := l.logger.WithFields(l.convertFields(append(l.fields, fields...)))
	entry.Panic(msg)
}

// Panicf 输出格式化的恐慌级日志并触发panic
func (l *LogrusLogger) Panicf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	msg = l.limitMessageSize(msg)
	l.logger.Panicf(msg)
}

// WithFields 添加字段到日志
func (l *LogrusLogger) WithFields(fields ...logger.Field) logger.Logger {
	return &LogrusLogger{
		logger:         l.logger,
		level:          l.level,
		fields:         append(l.fields, fields...),
		ctx:            l.ctx,
		name:           l.name,
		maxMessageSize: l.maxMessageSize,
	}
}

// WithField 添加单个字段到日志
func (l *LogrusLogger) WithField(key string, value interface{}) logger.Logger {
	return l.WithFields(logger.Field{Key: key, Value: value})
}

// WithContext 添加上下文到日志
func (l *LogrusLogger) WithContext(ctx context.Context) logger.Logger {
	return &LogrusLogger{
		logger:         l.logger,
		level:          l.level,
		fields:         l.fields,
		ctx:            ctx,
		name:           l.name,
		maxMessageSize: l.maxMessageSize,
	}
}

// WithError 添加错误信息到日志
func (l *LogrusLogger) WithError(err error) logger.Logger {
	return l.WithField("error", err)
}

// WithTime 添加时间到日志
func (l *LogrusLogger) WithTime(t time.Time) logger.Logger {
	return l.WithField("time", t)
}

// SetLevel 设置日志级别
func (l *LogrusLogger) SetLevel(level logger.LogLevel) {
	l.level = level
	logrusLevel := logrus.InfoLevel
	switch level {
	case logger.DebugLevel:
		logrusLevel = logrus.DebugLevel
	case logger.InfoLevel:
		logrusLevel = logrus.InfoLevel
	case logger.WarnLevel:
		logrusLevel = logrus.WarnLevel
	case logger.ErrorLevel:
		logrusLevel = logrus.ErrorLevel
	case logger.FatalLevel, logger.PanicLevel:
		logrusLevel = logrus.FatalLevel
	}
	l.logger.SetLevel(logrusLevel)
}

// GetLevel 获取日志级别
func (l *LogrusLogger) GetLevel() logger.LogLevel {
	return l.level
}

// IsDebugEnabled 检查调试级别是否启用
func (l *LogrusLogger) IsDebugEnabled() bool {
	return l.level <= logger.DebugLevel
}

// IsInfoEnabled 检查信息级别是否启用
func (l *LogrusLogger) IsInfoEnabled() bool {
	return l.level <= logger.InfoLevel
}

// IsWarnEnabled 检查警告级别是否启用
func (l *LogrusLogger) IsWarnEnabled() bool {
	return l.level <= logger.WarnLevel
}

// IsErrorEnabled 检查错误级别是否启用
func (l *LogrusLogger) IsErrorEnabled() bool {
	return l.level <= logger.ErrorLevel
}

// IsFatalEnabled 检查致命级别是否启用
func (l *LogrusLogger) IsFatalEnabled() bool {
	return l.level <= logger.FatalLevel
}

// IsPanicEnabled 检查恐慌级别是否启用
func (l *LogrusLogger) IsPanicEnabled() bool {
	return l.level <= logger.PanicLevel
}

// Sync 刷新日志缓冲区
func (l *LogrusLogger) Sync() error {
	return nil
}

// convertFields 转换字段
func (l *LogrusLogger) convertFields(fields []logger.Field) logrus.Fields {
	logrusFields := make(logrus.Fields)
	for _, field := range fields {
		logrusFields[field.Key] = field.Value
	}
	return logrusFields
}

// LogrusLoggerProvider logrus日志提供者
type LogrusLoggerProvider struct{}

// NewLogrusLoggerProvider 创建logrus日志提供者
func NewLogrusLoggerProvider() *LogrusLoggerProvider {
	return &LogrusLoggerProvider{}
}

// Create 创建日志实例
func (p *LogrusLoggerProvider) Create(name string, opts ...logger.Option) logger.Logger {
	return NewLogrusLogger(name, opts...)
}

// CreateWithConfig 根据配置创建日志实例
func (p *LogrusLoggerProvider) CreateWithConfig(name string, config map[string]interface{}) logger.Logger {
	opts := make([]logger.Option, 0)

	if level, ok := config["level"].(logger.LogLevel); ok {
		opts = append(opts, logger.WithLevel(level))
	}
	if format, ok := config["format"].(string); ok {
		opts = append(opts, logger.WithFormat(format))
	}
	if outputPath, ok := config["outputPath"].(string); ok {
		opts = append(opts, logger.WithOutputPath(outputPath))
	}
	if maxLogSize, ok := config["maxLogSize"].(int64); ok {
		opts = append(opts, logger.WithMaxLogSize(maxLogSize))
	}
	if maxLogAge, ok := config["maxLogAge"].(time.Duration); ok {
		opts = append(opts, logger.WithMaxLogAge(maxLogAge))
	}
	if maxLogFiles, ok := config["maxLogFiles"].(int); ok {
		opts = append(opts, logger.WithMaxLogFiles(maxLogFiles))
	}
	if compressLogs, ok := config["compressLogs"].(bool); ok {
		opts = append(opts, logger.WithCompressLogs(compressLogs))
	}
	if maxMessageSize, ok := config["maxMessageSize"].(int); ok {
		opts = append(opts, logger.WithMaxMessageSize(maxMessageSize))
	}

	return NewLogrusLogger(name, opts...)
}

// RegisterProvider 注册logrus日志提供者
func RegisterProvider() {
	logger.GetLogFactory().RegisterProvider("logrus", NewLogrusLoggerProvider())
}

// init 自动注册logrus日志提供者
func init() {
	RegisterProvider()
}
