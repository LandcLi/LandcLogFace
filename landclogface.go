package LandcLogFace

import (
	"time"

	"github.com/LandcLi/LandcLogFace/pkg/adapters"
	"github.com/LandcLi/LandcLogFace/pkg/logger"

	"github.com/gin-gonic/gin"
)

// 导入核心包

// 导出核心类型和函数

// LogLevel 定义日志级别
type LogLevel = logger.LogLevel

// Field 定义日志字段
type Field = logger.Field

// Logger 日志门面接口
type Logger = logger.Logger

// Option 日志配置选项
type Option = logger.Option

// LoggerOptions 日志配置选项
type LoggerOptions = logger.LoggerOptions

// LogConfig 统一的日志配置类
type LogConfig = logger.LogConfig

// 导出日志级别常量
const (
	DebugLevel LogLevel = logger.DebugLevel
	InfoLevel  LogLevel = logger.InfoLevel
	WarnLevel  LogLevel = logger.WarnLevel
	ErrorLevel LogLevel = logger.ErrorLevel
	FatalLevel LogLevel = logger.FatalLevel
	PanicLevel LogLevel = logger.PanicLevel
)

// 导出核心函数

// GetLogFactory 获取全局日志工厂实例
func GetLogFactory() *logger.LogFactory {
	return logger.GetLogFactory()
}

// GetLogger 获取全局日志实例
func GetLogger() Logger {
	return logger.GetLogger()
}

// GetLoggerWithName 获取指定名称的日志实例
func GetLoggerWithName(name string) Logger {
	return logger.GetLoggerWithName(name)
}

// GetLoggerWithProvider 获取指定提供者的日志实例
func GetLoggerWithProvider(name string, provider string) Logger {
	return logger.GetLoggerWithProvider(name, provider)
}

// GetLoggerWithConfig 根据配置获取日志实例
func GetLoggerWithConfig(name string, config map[string]interface{}) Logger {
	return logger.GetLoggerWithConfig(name, config)
}

// GetLoggerWithLogConfig 根据LogConfig获取日志实例
func GetLoggerWithLogConfig(config *LogConfig) Logger {
	return logger.GetLoggerWithLogConfig(config)
}

// SetGlobalLogger 设置全局日志实例
func SetGlobalLogger(log Logger) {
	logger.SetGlobalLogger(log)
}

// NewLogConfig 创建默认的日志配置
func NewLogConfig() *LogConfig {
	return logger.NewLogConfig()
}

// 导出选项函数

// WithLevel 设置日志级别
func WithLevel(level LogLevel) Option {
	return logger.WithLevel(level)
}

// WithFormat 设置日志格式
func WithFormat(format string) Option {
	return logger.WithFormat(format)
}

// WithOutputPath 设置日志输出路径
func WithOutputPath(path string) Option {
	return logger.WithOutputPath(path)
}

// WithConfig 设置额外配置
func WithConfig(config map[string]interface{}) Option {
	return logger.WithConfig(config)
}

// WithMaxLogSize 设置单个日志文件最大大小（MB）
func WithMaxLogSize(size int64) Option {
	return logger.WithMaxLogSize(size)
}

// WithMaxLogAge 设置日志文件最大保留时间
func WithMaxLogAge(age time.Duration) Option {
	return logger.WithMaxLogAge(age)
}

// WithMaxLogFiles 设置最大保留日志文件数量
func WithMaxLogFiles(files int) Option {
	return logger.WithMaxLogFiles(files)
}

// WithCompressLogs 设置是否压缩旧日志
func WithCompressLogs(compress bool) Option {
	return logger.WithCompressLogs(compress)
}

// WithMaxMessageSize 设置单条日志最大大小（KB）
func WithMaxMessageSize(size int) Option {
	return logger.WithMaxMessageSize(size)
}

// 导出全局日志函数

// Debug 全局调试级日志
func Debug(msg string, fields ...Field) {
	logger.Debug(msg, fields...)
}

// Debugf 全局格式化调试级日志
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

// Info 全局信息级日志
func Info(msg string, fields ...Field) {
	logger.Info(msg, fields...)
}

// Infof 全局格式化信息级日志
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Warn 全局警告级日志
func Warn(msg string, fields ...Field) {
	logger.Warn(msg, fields...)
}

// Warnf 全局格式化警告级日志
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

// Error 全局错误级日志
func Error(msg string, fields ...Field) {
	logger.Error(msg, fields...)
}

// Errorf 全局格式化错误级日志
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

// Fatal 全局致命级日志
func Fatal(msg string, fields ...Field) {
	logger.Fatal(msg, fields...)
}

// Fatalf 全局格式化致命级日志
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

// Panic 全局恐慌级日志
func Panic(msg string, fields ...Field) {
	logger.Panic(msg, fields...)
}

// Panicf 全局格式化恐慌级日志
func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

// 导出适配器函数

// NewGinLogger 创建一个新的gin日志适配器
func NewGinLogger(log Logger) *adapters.GinLogger {
	return adapters.NewGinLogger(log)
}

// NewGFLogger 创建一个新的goframe日志适配器
func NewGFLogger(log Logger) *adapters.GFLogger {
	return adapters.NewGFLogger(log)
}

// UseWithGin 将日志适配器应用到gin引擎
func UseWithGin(r *gin.Engine, log Logger) {
	adapters.UseWithGin(r, log)
}

// UseWithGF 将日志适配器应用到goframe
func UseWithGF(log Logger) *adapters.GFLogger {
	return adapters.UseWithGF(log)
}
