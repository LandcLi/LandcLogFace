package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ZapLogger zap日志库适配器
type ZapLogger struct {
	logger         *zap.Logger
	level          LogLevel
	fields         []Field
	ctx            context.Context
	name           string
	maxMessageSize int // 单条日志最大大小（KB）
}

// NewZapLogger 创建zap日志实例
func NewZapLogger(name string, opts ...Option) *ZapLogger {
	options := &LoggerOptions{
		Level:          InfoLevel,
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
	case DebugLevel:
		zapLevel = zapcore.DebugLevel
	case InfoLevel:
		zapLevel = zapcore.InfoLevel
	case WarnLevel:
		zapLevel = zapcore.WarnLevel
	case ErrorLevel:
		zapLevel = zapcore.ErrorLevel
	case FatalLevel:
		zapLevel = zapcore.FatalLevel
	case PanicLevel:
		zapLevel = zapcore.PanicLevel
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

	// 配置输出
	var core zapcore.Core
	if options.OutputPath == "stdout" {
		// 输出到标准输出
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			zapLevel,
		)
	} else {
		// 输出到文件，使用lumberjack进行轮转
		lumberjackLogger := &lumberjack.Logger{
			Filename:   options.OutputPath,
			MaxSize:    int(options.MaxLogSize),             // MB
			MaxAge:     int(options.MaxLogAge.Hours() / 24), // 天
			MaxBackups: options.MaxLogFiles,
			Compress:   options.CompressLogs,
		}

		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(lumberjackLogger),
			zapLevel,
		)
	}

	// 构建logger
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	// 添加名称字段
	logger = logger.Named(name)

	return &ZapLogger{
		logger:         logger,
		level:          options.Level,
		fields:         make([]Field, 0),
		ctx:            context.Background(),
		name:           name,
		maxMessageSize: options.MaxMessageSize,
	}
}

// SetLevel 设置日志级别
func (z *ZapLogger) SetLevel(level LogLevel) {
	z.level = level
	// 更新zap的日志级别
	// 注意：zap的AtomicLevel需要通过Core来设置，这里简化处理
}

// GetLevel 获取当前日志级别
func (z *ZapLogger) GetLevel() LogLevel {
	return z.level
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

// toZapFields 将自定义字段转换为zap字段
func (z *ZapLogger) toZapFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(z.fields)+len(fields))

	// 添加已有的字段
	for _, field := range z.fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}

	// 添加新的字段
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}

	return zapFields
}

// Debug 输出调试级日志
func (z *ZapLogger) Debug(msg string, fields ...Field) {
	if z.level <= DebugLevel {
		z.logger.Debug(z.limitMessageSize(msg), z.toZapFields(fields)...)
	}
}

// Debugf 输出格式化的调试级日志
func (z *ZapLogger) Debugf(format string, args ...interface{}) {
	if z.level <= DebugLevel {
		msg := fmt.Sprintf(format, args...)
		z.logger.Sugar().Debug(z.limitMessageSize(msg))
	}
}

// Info 输出信息级日志
func (z *ZapLogger) Info(msg string, fields ...Field) {
	if z.level <= InfoLevel {
		z.logger.Info(z.limitMessageSize(msg), z.toZapFields(fields)...)
	}
}

// Infof 输出格式化的信息级日志
func (z *ZapLogger) Infof(format string, args ...interface{}) {
	if z.level <= InfoLevel {
		msg := fmt.Sprintf(format, args...)
		z.logger.Sugar().Info(z.limitMessageSize(msg))
	}
}

// Warn 输出警告级日志
func (z *ZapLogger) Warn(msg string, fields ...Field) {
	if z.level <= WarnLevel {
		z.logger.Warn(z.limitMessageSize(msg), z.toZapFields(fields)...)
	}
}

// Warnf 输出格式化的警告级日志
func (z *ZapLogger) Warnf(format string, args ...interface{}) {
	if z.level <= WarnLevel {
		msg := fmt.Sprintf(format, args...)
		z.logger.Sugar().Warn(z.limitMessageSize(msg))
	}
}

// Error 输出错误级日志
func (z *ZapLogger) Error(msg string, fields ...Field) {
	if z.level <= ErrorLevel {
		z.logger.Error(z.limitMessageSize(msg), z.toZapFields(fields)...)
	}
}

// Errorf 输出格式化的错误级日志
func (z *ZapLogger) Errorf(format string, args ...interface{}) {
	if z.level <= ErrorLevel {
		msg := fmt.Sprintf(format, args...)
		z.logger.Sugar().Error(z.limitMessageSize(msg))
	}
}

// Fatal 输出致命级日志并退出程序
func (z *ZapLogger) Fatal(msg string, fields ...Field) {
	if z.level <= FatalLevel {
		z.logger.Fatal(z.limitMessageSize(msg), z.toZapFields(fields)...)
		os.Exit(1)
	}
}

// Fatalf 输出格式化的致命级日志并退出程序
func (z *ZapLogger) Fatalf(format string, args ...interface{}) {
	if z.level <= FatalLevel {
		msg := fmt.Sprintf(format, args...)
		z.logger.Sugar().Fatal(z.limitMessageSize(msg))
		os.Exit(1)
	}
}

// Panic 输出恐慌级日志并触发panic
func (z *ZapLogger) Panic(msg string, fields ...Field) {
	if z.level <= PanicLevel {
		z.logger.Panic(z.limitMessageSize(msg), z.toZapFields(fields)...)
	}
}

// Panicf 输出格式化的恐慌级日志并触发panic
func (z *ZapLogger) Panicf(format string, args ...interface{}) {
	if z.level <= PanicLevel {
		msg := fmt.Sprintf(format, args...)
		z.logger.Sugar().Panic(z.limitMessageSize(msg))
	}
}

// WithFields 添加字段到日志
func (z *ZapLogger) WithFields(fields ...Field) Logger {
	newLogger := *z
	newLogger.fields = append(newLogger.fields, fields...)
	// 创建新的zap.Logger实例
	zapFields := make([]zap.Field, 0, len(newLogger.fields))
	for _, field := range newLogger.fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}
	newLogger.logger = z.logger.With(zapFields...)
	return &newLogger
}

// WithField 添加单个字段到日志
func (z *ZapLogger) WithField(key string, value interface{}) Logger {
	return z.WithFields(Field{Key: key, Value: value})
}

// WithContext 添加上下文到日志
func (z *ZapLogger) WithContext(ctx context.Context) Logger {
	newLogger := *z
	newLogger.ctx = ctx
	return &newLogger
}

// WithError 添加错误信息到日志
func (z *ZapLogger) WithError(err error) Logger {
	return z.WithField("error", err)
}

// WithTime 添加时间到日志
func (z *ZapLogger) WithTime(t time.Time) Logger {
	return z.WithField("time", t)
}

// IsDebugEnabled 检查调试级别是否启用
func (z *ZapLogger) IsDebugEnabled() bool {
	return z.level <= DebugLevel
}

// IsInfoEnabled 检查信息级别是否启用
func (z *ZapLogger) IsInfoEnabled() bool {
	return z.level <= InfoLevel
}

// IsWarnEnabled 检查警告级别是否启用
func (z *ZapLogger) IsWarnEnabled() bool {
	return z.level <= WarnLevel
}

// IsErrorEnabled 检查错误级别是否启用
func (z *ZapLogger) IsErrorEnabled() bool {
	return z.level <= ErrorLevel
}

// IsFatalEnabled 检查致命级别是否启用
func (z *ZapLogger) IsFatalEnabled() bool {
	return z.level <= FatalLevel
}

// IsPanicEnabled 检查恐慌级别是否启用
func (z *ZapLogger) IsPanicEnabled() bool {
	return z.level <= PanicLevel
}

// Sync 刷新日志缓冲区
func (z *ZapLogger) Sync() error {
	return z.logger.Sync()
}

// ZapLoggerProvider zap日志提供者
type ZapLoggerProvider struct{}

// NewZapLoggerProvider 创建zap日志提供者
func NewZapLoggerProvider() *ZapLoggerProvider {
	return &ZapLoggerProvider{}
}

// Create 创建日志实例
func (p *ZapLoggerProvider) Create(name string) Logger {
	return NewZapLogger(name)
}

// CreateWithConfig 根据配置创建日志实例
func (p *ZapLoggerProvider) CreateWithConfig(name string, config map[string]interface{}) Logger {
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
		format = "json"
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

	return NewZapLogger(name,
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
