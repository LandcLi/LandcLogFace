package zap

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/LandcLi/LandcLogFace/internal/logger"
)

// ZapLogger zap日志库适配器
type ZapLogger struct {
	logger         *zap.Logger
	level          logger.LogLevel
	fields         []logger.Field
	ctx            context.Context
	name           string
	maxMessageSize int // 单条日志最大大小（KB）
}

// NewZapLogger 创建zap日志实例
func NewZapLogger(name string, opts ...logger.Option) *ZapLogger {
	options := &logger.LoggerOptions{
		Level:          logger.InfoLevel,
		Format:         "json",
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

	// 配置zap
	zapLevel := zapcore.InfoLevel
	switch options.Level {
	case logger.DebugLevel:
		zapLevel = zapcore.DebugLevel
	case logger.InfoLevel:
		zapLevel = zapcore.InfoLevel
	case logger.WarnLevel:
		zapLevel = zapcore.WarnLevel
	case logger.ErrorLevel:
		zapLevel = zapcore.ErrorLevel
	case logger.FatalLevel, logger.PanicLevel:
		zapLevel = zapcore.FatalLevel
	}

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if options.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 配置输出
	var writer zapcore.WriteSyncer
	if options.OutputPath == "stdout" || options.OutputPath == "" {
		writer = zapcore.AddSync(os.Stdout)
	} else {
		writer = zapcore.AddSync(&lumberjack.Logger{
			Filename:   options.OutputPath,
			MaxSize:    int(options.MaxLogSize),
			MaxBackups: options.MaxLogFiles,
			MaxAge:     int(options.MaxLogAge / (24 * time.Hour)),
			Compress:   options.CompressLogs,
		})
	}

	// 创建core
	core := zapcore.NewCore(encoder, writer, zapLevel)

	// 创建logger
	zapLogger := zap.New(core, zap.AddCaller())

	return &ZapLogger{
		logger:         zapLogger,
		level:          options.Level,
		name:           name,
		maxMessageSize: options.MaxMessageSize,
	}
}

// limitMessageSize 限制日志消息大小
func (z *ZapLogger) limitMessageSize(msg string) string {
	if z.maxMessageSize > 0 {
		maxSize := z.maxMessageSize * 1024 // 转换为字节
		if len(msg) > maxSize {
			return msg[:maxSize-3] + "..."
		}
	}
	return msg
}

// Debug 输出调试级日志
func (z *ZapLogger) Debug(msg string, fields ...logger.Field) {
	msg = z.limitMessageSize(msg)
	z.logger.Debug(msg, z.convertFields(fields)...)
}

// Debugf 输出格式化的调试级日志
func (z *ZapLogger) Debugf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	msg = z.limitMessageSize(msg)
	z.logger.Debug(msg)
}

// Info 输出信息级日志
func (z *ZapLogger) Info(msg string, fields ...logger.Field) {
	msg = z.limitMessageSize(msg)
	z.logger.Info(msg, z.convertFields(fields)...)
}

// Infof 输出格式化的信息级日志
func (z *ZapLogger) Infof(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	msg = z.limitMessageSize(msg)
	z.logger.Info(msg)
}

// Warn 输出警告级日志
func (z *ZapLogger) Warn(msg string, fields ...logger.Field) {
	msg = z.limitMessageSize(msg)
	z.logger.Warn(msg, z.convertFields(fields)...)
}

// Warnf 输出格式化的警告级日志
func (z *ZapLogger) Warnf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	msg = z.limitMessageSize(msg)
	z.logger.Warn(msg)
}

// Error 输出错误级日志
func (z *ZapLogger) Error(msg string, fields ...logger.Field) {
	msg = z.limitMessageSize(msg)
	z.logger.Error(msg, z.convertFields(fields)...)
}

// Errorf 输出格式化的错误级日志
func (z *ZapLogger) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	msg = z.limitMessageSize(msg)
	z.logger.Error(msg)
}

// Fatal 输出致命级日志并退出程序
func (z *ZapLogger) Fatal(msg string, fields ...logger.Field) {
	msg = z.limitMessageSize(msg)
	z.logger.Fatal(msg, z.convertFields(fields)...)
}

// Fatalf 输出格式化的致命级日志并退出程序
func (z *ZapLogger) Fatalf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	msg = z.limitMessageSize(msg)
	z.logger.Fatal(msg)
}

// Panic 输出恐慌级日志并触发panic
func (z *ZapLogger) Panic(msg string, fields ...logger.Field) {
	msg = z.limitMessageSize(msg)
	z.logger.Panic(msg, z.convertFields(fields)...)
}

// Panicf 输出格式化的恐慌级日志并触发panic
func (z *ZapLogger) Panicf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	msg = z.limitMessageSize(msg)
	z.logger.Panic(msg)
}

// WithFields 添加字段到日志
func (z *ZapLogger) WithFields(fields ...logger.Field) logger.Logger {
	return &ZapLogger{
		logger:         z.logger.With(z.convertFields(fields)...),
		level:          z.level,
		fields:         append(z.fields, fields...),
		ctx:            z.ctx,
		name:           z.name,
		maxMessageSize: z.maxMessageSize,
	}
}

// WithField 添加单个字段到日志
func (z *ZapLogger) WithField(key string, value interface{}) logger.Logger {
	return z.WithFields(logger.Field{Key: key, Value: value})
}

// WithContext 添加上下文到日志
func (z *ZapLogger) WithContext(ctx context.Context) logger.Logger {
	return &ZapLogger{
		logger:         z.logger,
		level:          z.level,
		fields:         z.fields,
		ctx:            ctx,
		name:           z.name,
		maxMessageSize: z.maxMessageSize,
	}
}

// WithError 添加错误信息到日志
func (z *ZapLogger) WithError(err error) logger.Logger {
	return z.WithField("error", err)
}

// WithTime 添加时间到日志
func (z *ZapLogger) WithTime(t time.Time) logger.Logger {
	return z.WithField("time", t)
}

// SetLevel 设置日志级别
func (z *ZapLogger) SetLevel(level logger.LogLevel) {
	z.level = level
}

// GetLevel 获取日志级别
func (z *ZapLogger) GetLevel() logger.LogLevel {
	return z.level
}

// IsDebugEnabled 检查调试级别是否启用
func (z *ZapLogger) IsDebugEnabled() bool {
	return z.level <= logger.DebugLevel
}

// IsInfoEnabled 检查信息级别是否启用
func (z *ZapLogger) IsInfoEnabled() bool {
	return z.level <= logger.InfoLevel
}

// IsWarnEnabled 检查警告级别是否启用
func (z *ZapLogger) IsWarnEnabled() bool {
	return z.level <= logger.WarnLevel
}

// IsErrorEnabled 检查错误级别是否启用
func (z *ZapLogger) IsErrorEnabled() bool {
	return z.level <= logger.ErrorLevel
}

// IsFatalEnabled 检查致命级别是否启用
func (z *ZapLogger) IsFatalEnabled() bool {
	return z.level <= logger.FatalLevel
}

// IsPanicEnabled 检查恐慌级别是否启用
func (z *ZapLogger) IsPanicEnabled() bool {
	return z.level <= logger.PanicLevel
}

// Sync 刷新日志缓冲区
func (z *ZapLogger) Sync() error {
	return z.logger.Sync()
}

// convertFields 转换字段
func (z *ZapLogger) convertFields(fields []logger.Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	return zapFields
}

// ZapLoggerProvider zap日志提供者
type ZapLoggerProvider struct{}

// NewZapLoggerProvider 创建zap日志提供者
func NewZapLoggerProvider() *ZapLoggerProvider {
	return &ZapLoggerProvider{}
}

// Create 创建日志实例
func (p *ZapLoggerProvider) Create(name string, opts ...logger.Option) logger.Logger {
	return NewZapLogger(name, opts...)
}

// CreateWithConfig 根据配置创建日志实例
func (p *ZapLoggerProvider) CreateWithConfig(name string, config map[string]interface{}) logger.Logger {
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

	return NewZapLogger(name, opts...)
}

// RegisterProvider 注册zap日志提供者
func RegisterProvider() {
	logger.GetLogFactory().RegisterProvider("zap", NewZapLoggerProvider())
}

// init 自动注册zap日志提供者
func init() {
	RegisterProvider()
}
