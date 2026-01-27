package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogrusLogger logrus日志库适配器
type LogrusLogger struct {
	logger         *logrus.Logger
	level          LogLevel
	fields         []Field
	ctx            context.Context
	name           string
	maxMessageSize int // 单条日志最大大小（KB）
}

// NewLogrusLogger 创建logrus日志实例
func NewLogrusLogger(name string, opts ...Option) *LogrusLogger {
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

	// 创建logrus实例
	logger := logrus.New()

	// 设置日志级别
	logrusLevel := logrus.InfoLevel
	switch options.Level {
	case DebugLevel:
		logrusLevel = logrus.DebugLevel
	case InfoLevel:
		logrusLevel = logrus.InfoLevel
	case WarnLevel:
		logrusLevel = logrus.WarnLevel
	case ErrorLevel:
		logrusLevel = logrus.ErrorLevel
	case FatalLevel:
		logrusLevel = logrus.FatalLevel
	case PanicLevel:
		logrusLevel = logrus.PanicLevel
	}
	logger.SetLevel(logrusLevel)

	// 设置输出格式
	if options.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05.000",
		})
	}

	// 设置输出目标
	if options.OutputPath != "stdout" {
		// 使用lumberjack进行日志轮转
		lumberjackLogger := &lumberjack.Logger{
			Filename:   options.OutputPath,
			MaxSize:    int(options.MaxLogSize),             // MB
			MaxAge:     int(options.MaxLogAge.Hours() / 24), // 天
			MaxBackups: options.MaxLogFiles,
			Compress:   options.CompressLogs,
		}
		logger.SetOutput(lumberjackLogger)
	}

	return &LogrusLogger{
		logger:         logger,
		level:          options.Level,
		fields:         make([]Field, 0),
		ctx:            context.Background(),
		name:           name,
		maxMessageSize: options.MaxMessageSize,
	}
}

// SetLevel 设置日志级别
func (l *LogrusLogger) SetLevel(level LogLevel) {
	l.level = level
	// 更新logrus的日志级别
	logrusLevel := logrus.InfoLevel
	switch level {
	case DebugLevel:
		logrusLevel = logrus.DebugLevel
	case InfoLevel:
		logrusLevel = logrus.InfoLevel
	case WarnLevel:
		logrusLevel = logrus.WarnLevel
	case ErrorLevel:
		logrusLevel = logrus.ErrorLevel
	case FatalLevel:
		logrusLevel = logrus.FatalLevel
	case PanicLevel:
		logrusLevel = logrus.PanicLevel
	}
	l.logger.SetLevel(logrusLevel)
}

// GetLevel 获取当前日志级别
func (l *LogrusLogger) GetLevel() LogLevel {
	return l.level
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

// toLogrusFields 将自定义字段转换为logrus字段
func (l *LogrusLogger) toLogrusFields(fields []Field) logrus.Fields {
	logrusFields := make(logrus.Fields)

	// 添加已有的字段
	for _, field := range l.fields {
		logrusFields[field.Key] = field.Value
	}

	// 添加新的字段
	for _, field := range fields {
		logrusFields[field.Key] = field.Value
	}

	return logrusFields
}

// Debug 输出调试级日志
func (l *LogrusLogger) Debug(msg string, fields ...Field) {
	if l.level <= DebugLevel {
		l.logger.WithFields(l.toLogrusFields(fields)).Debug(l.limitMessageSize(msg))
	}
}

// Debugf 输出格式化的调试级日志
func (l *LogrusLogger) Debugf(format string, args ...interface{}) {
	if l.level <= DebugLevel {
		msg := fmt.Sprintf(format, args...)
		l.logger.Debug(l.limitMessageSize(msg))
	}
}

// Info 输出信息级日志
func (l *LogrusLogger) Info(msg string, fields ...Field) {
	if l.level <= InfoLevel {
		l.logger.WithFields(l.toLogrusFields(fields)).Info(l.limitMessageSize(msg))
	}
}

// Infof 输出格式化的信息级日志
func (l *LogrusLogger) Infof(format string, args ...interface{}) {
	if l.level <= InfoLevel {
		msg := fmt.Sprintf(format, args...)
		l.logger.Info(l.limitMessageSize(msg))
	}
}

// Warn 输出警告级日志
func (l *LogrusLogger) Warn(msg string, fields ...Field) {
	if l.level <= WarnLevel {
		l.logger.WithFields(l.toLogrusFields(fields)).Warn(l.limitMessageSize(msg))
	}
}

// Warnf 输出格式化的警告级日志
func (l *LogrusLogger) Warnf(format string, args ...interface{}) {
	if l.level <= WarnLevel {
		msg := fmt.Sprintf(format, args...)
		l.logger.Warn(l.limitMessageSize(msg))
	}
}

// Error 输出错误级日志
func (l *LogrusLogger) Error(msg string, fields ...Field) {
	if l.level <= ErrorLevel {
		l.logger.WithFields(l.toLogrusFields(fields)).Error(l.limitMessageSize(msg))
	}
}

// Errorf 输出格式化的错误级日志
func (l *LogrusLogger) Errorf(format string, args ...interface{}) {
	if l.level <= ErrorLevel {
		msg := fmt.Sprintf(format, args...)
		l.logger.Error(l.limitMessageSize(msg))
	}
}

// Fatal 输出致命级日志并退出程序
func (l *LogrusLogger) Fatal(msg string, fields ...Field) {
	if l.level <= FatalLevel {
		l.logger.WithFields(l.toLogrusFields(fields)).Fatal(l.limitMessageSize(msg))
		os.Exit(1)
	}
}

// Fatalf 输出格式化的致命级日志并退出程序
func (l *LogrusLogger) Fatalf(format string, args ...interface{}) {
	if l.level <= FatalLevel {
		msg := fmt.Sprintf(format, args...)
		l.logger.Fatal(l.limitMessageSize(msg))
		os.Exit(1)
	}
}

// Panic 输出恐慌级日志并触发panic
func (l *LogrusLogger) Panic(msg string, fields ...Field) {
	if l.level <= PanicLevel {
		l.logger.WithFields(l.toLogrusFields(fields)).Panic(l.limitMessageSize(msg))
	}
}

// Panicf 输出格式化的恐慌级日志并触发panic
func (l *LogrusLogger) Panicf(format string, args ...interface{}) {
	if l.level <= PanicLevel {
		msg := fmt.Sprintf(format, args...)
		l.logger.Panic(l.limitMessageSize(msg))
	}
}

// WithFields 添加字段到日志
func (l *LogrusLogger) WithFields(fields ...Field) Logger {
	newLogger := *l
	newLogger.fields = append(newLogger.fields, fields...)
	return &newLogger
}

// WithField 添加单个字段到日志
func (l *LogrusLogger) WithField(key string, value interface{}) Logger {
	return l.WithFields(Field{Key: key, Value: value})
}

// WithContext 添加上下文到日志
func (l *LogrusLogger) WithContext(ctx context.Context) Logger {
	newLogger := *l
	newLogger.ctx = ctx
	return &newLogger
}

// WithError 添加错误信息到日志
func (l *LogrusLogger) WithError(err error) Logger {
	return l.WithField("error", err)
}

// WithTime 添加时间到日志
func (l *LogrusLogger) WithTime(t time.Time) Logger {
	return l.WithField("time", t)
}

// IsDebugEnabled 检查调试级别是否启用
func (l *LogrusLogger) IsDebugEnabled() bool {
	return l.level <= DebugLevel
}

// IsInfoEnabled 检查信息级别是否启用
func (l *LogrusLogger) IsInfoEnabled() bool {
	return l.level <= InfoLevel
}

// IsWarnEnabled 检查警告级别是否启用
func (l *LogrusLogger) IsWarnEnabled() bool {
	return l.level <= WarnLevel
}

// IsErrorEnabled 检查错误级别是否启用
func (l *LogrusLogger) IsErrorEnabled() bool {
	return l.level <= ErrorLevel
}

// IsFatalEnabled 检查致命级别是否启用
func (l *LogrusLogger) IsFatalEnabled() bool {
	return l.level <= FatalLevel
}

// IsPanicEnabled 检查恐慌级别是否启用
func (l *LogrusLogger) IsPanicEnabled() bool {
	return l.level <= PanicLevel
}

// Sync 刷新日志缓冲区
func (l *LogrusLogger) Sync() error {
	// logrus没有Sync方法，返回nil
	return nil
}

// LogrusLoggerProvider logrus日志提供者
type LogrusLoggerProvider struct{}

// NewLogrusLoggerProvider 创建logrus日志提供者
func NewLogrusLoggerProvider() *LogrusLoggerProvider {
	return &LogrusLoggerProvider{}
}

// Create 创建日志实例
func (p *LogrusLoggerProvider) Create(name string) Logger {
	return NewLogrusLogger(name)
}

// CreateWithConfig 根据配置创建日志实例
func (p *LogrusLoggerProvider) CreateWithConfig(name string, config map[string]interface{}) Logger {
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

	return NewLogrusLogger(name,
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
