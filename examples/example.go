package examples

import (
	"github.com/LandcLi/LandcLogFace"
	"context"
	"errors"
	"fmt"
	"time"
)

// ExampleBasicUsage 基本使用示例
func ExampleBasicUsage() {
	// 获取全局日志实例
	logger := LandcLogFace.GetLogger()

	// 输出不同级别的日志
	logger.Debug("这是一条调试日志")
	logger.Info("这是一条信息日志")
	logger.Warn("这是一条警告日志")
	logger.Error("这是一条错误日志")

	// 输出格式化日志
	logger.Infof("Hello, %s!", "world")
	logger.Errorf("Error: %d", 404)

	// 使用字段
	logger.Info("用户登录",
		LandcLogFace.Field{Key: "user_id", Value: 123},
		LandcLogFace.Field{Key: "ip", Value: "192.168.1.1"},
	)

	// 链式调用
	logger.WithField("module", "auth").
		WithField("action", "login").
		Info("用户认证")

	fmt.Println("基本使用示例完成")
}

// ExampleWithDifferentProviders 使用不同的日志提供者
func ExampleWithDifferentProviders() {
	// 使用console提供者
	consoleLogger := LandcLogFace.GetLogFactory().CreateLoggerWithProvider("app", "console")
	consoleLogger.Info("使用控制台日志")

	// 使用zap提供者
	zapLogger := LandcLogFace.GetLogFactory().CreateLoggerWithProvider("app", "zap")
	zapLogger.Info("使用zap日志")

	// 使用logrus提供者
	logrusLogger := LandcLogFace.GetLogFactory().CreateLoggerWithProvider("app", "logrus")
	logrusLogger.Info("使用logrus日志")

	// 使用std提供者
	stdLogger := LandcLogFace.GetLogFactory().CreateLoggerWithProvider("app", "std")
	stdLogger.Info("使用标准库日志")

	fmt.Println("不同提供者使用示例完成")
}

// ExampleWithConfiguration 带配置的使用示例
func ExampleWithConfiguration() {
	// 配置map
	config := map[string]interface{}{
		"provider":   "zap",
		"level":      LandcLogFace.DebugLevel,
		"format":     "json",
		"outputPath": "stdout",
	}

	// 根据配置创建日志实例
	logger := LandcLogFace.GetLogFactory().CreateLoggerWithConfig("app", config)
	logger.Debug("带配置的调试日志")
	logger.Info("带配置的信息日志")

	fmt.Println("带配置使用示例完成")
}

// ExampleWithFileRotation 日志文件轮转示例
func ExampleWithFileRotation() {
	// 使用zap日志库，配置文件轮转
	zapLogger := LandcLogFace.GetLoggerWithProvider("app", "zap")
	zapLogger.SetLevel(LandcLogFace.InfoLevel)

	zapLogger.Info("使用文件轮转的zap日志")

	// 使用logrus日志库，配置文件轮转
	logrusLogger := LandcLogFace.GetLoggerWithProvider("app", "logrus")
	logrusLogger.SetLevel(LandcLogFace.InfoLevel)

	logrusLogger.Info("使用文件轮转的logrus日志")

	fmt.Println("日志文件轮转示例完成")
}

// ExampleWithMaxMessageSize 单条日志大小限制示例
func ExampleWithMaxMessageSize() {
	// 配置单条日志最大大小
	logger := LandcLogFace.GetLoggerWithProvider("app", "zap")
	logger.SetLevel(LandcLogFace.InfoLevel)

	logger.Info("单条日志大小限制示例")

	fmt.Println("单条日志大小限制示例完成")
}

// ExampleWithLogConfig 使用统一配置类示例
func ExampleWithLogConfig() {
	// 创建并配置LogConfig
	config := LandcLogFace.NewLogConfig()
	config.WithProvider("zap").
		WithName("myapp").
		WithLevel(LandcLogFace.DebugLevel).
		WithFormat("json").
		WithOutputPath("app.log").
		WithMaxLogSize(50).             // 单个日志文件最大50MB
		WithMaxLogAge(7*24*time.Hour).  // 日志文件保留7天
		WithMaxLogFiles(5).             // 最多保留5个日志文件
		WithCompressLogs(true).         // 压缩旧日志
		WithMaxMessageSize(10).         // 单条日志最大10KB
		WithExtraConfig("key", "value") // 添加额外配置

	// 根据LogConfig创建日志实例
	logger := LandcLogFace.GetLoggerWithLogConfig(config)
	logger.Debug("使用LogConfig的调试日志")
	logger.Info("使用LogConfig的信息日志")

	// 另一种创建方式
	anotherConfig := LandcLogFace.NewLogConfig()
	anotherConfig.Provider = "logrus"
	anotherConfig.Name = "anotherapp"
	anotherConfig.Level = LandcLogFace.InfoLevel
	anotherConfig.OutputPath = "another.log"

	anotherLogger := LandcLogFace.GetLoggerWithLogConfig(anotherConfig)
	anotherLogger.Info("使用另一个LogConfig的日志")

	fmt.Println("使用统一配置类示例完成")
}

// ExampleWithContext 使用上下文
func ExampleWithContext() {
	logger := LandcLogFace.GetLogger()

	// 创建上下文
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "123456")

	// 添加上下文到日志
	ctxLogger := logger.WithContext(ctx)
	ctxLogger.Info("带上下文的日志")

	fmt.Println("上下文使用示例完成")
}

// ExampleWithError 使用错误
func ExampleWithError() {
	logger := LandcLogFace.GetLogger()

	// 创建错误
	err := errors.New("数据库连接失败")

	// 添加错误到日志
	errLogger := logger.WithError(err)
	errLogger.Error("操作失败")

	// 直接在日志中添加错误字段
	logger.Error("操作失败", LandcLogFace.Field{Key: "error", Value: err})

	fmt.Println("错误使用示例完成")
}

// ExampleWithFields 使用字段
func ExampleWithFields() {
	logger := LandcLogFace.GetLogger()

	// 创建带字段的日志实例
	userLogger := logger.WithFields(
		LandcLogFace.Field{Key: "service", Value: "user"},
		LandcLogFace.Field{Key: "version", Value: "1.0.0"},
	)

	// 使用带字段的日志实例
	userLogger.Info("用户注册",
		LandcLogFace.Field{Key: "username", Value: "john"},
		LandcLogFace.Field{Key: "email", Value: "john@example.com"},
	)

	// 继续添加字段
	userLogger.WithField("status", "success").Info("注册成功")

	fmt.Println("字段使用示例完成")
}

// ExampleGlobalFunctions 全局函数使用示例
func ExampleGlobalFunctions() {
	// 使用全局函数输出日志
	LandcLogFace.Debug("全局调试日志")
	LandcLogFace.Info("全局信息日志")
	LandcLogFace.Warn("全局警告日志")
	LandcLogFace.Error("全局错误日志")

	// 使用全局格式化函数
	LandcLogFace.Infof("全局格式化日志: %s", "test")
	LandcLogFace.Errorf("全局错误格式化日志: %d", 500)

	// 使用全局函数带字段
	LandcLogFace.Info("全局带字段日志", LandcLogFace.Field{Key: "key", Value: "value"})

	fmt.Println("全局函数使用示例完成")
}

// ExampleLogLevelCheck 日志级别检查示例
func ExampleLogLevelCheck() {
	logger := LandcLogFace.GetLogger()

	// 检查日志级别是否启用
	if logger.IsDebugEnabled() {
		// 执行耗时操作
		data := "耗时操作的结果"
		logger.Debug("调试信息", LandcLogFace.Field{Key: "data", Value: data})
	}

	if logger.IsInfoEnabled() {
		logger.Info("信息日志已启用")
	}

	fmt.Println("日志级别检查示例完成")
}

// ExampleCustomProvider 自定义日志提供者示例
func ExampleCustomProvider() {
	// 创建自定义日志提供者
	customProvider := &CustomLoggerProvider{}

	// 注册自定义提供者
	LandcLogFace.GetLogFactory().RegisterProvider("custom", customProvider)

	// 使用自定义提供者
	customLogger := LandcLogFace.GetLogFactory().CreateLoggerWithProvider("app", "custom")
	customLogger.Info("使用自定义日志提供者")

	// 注销自定义提供者
	LandcLogFace.GetLogFactory().UnregisterProvider("custom")

	fmt.Println("自定义提供者示例完成")
}

// CustomLogger 自定义日志实现
type CustomLogger struct {
	name string
}

// NewCustomLogger 创建自定义日志实例
func NewCustomLogger(name string) *CustomLogger {
	return &CustomLogger{name: name}
}

// SetLevel 设置日志级别
func (c *CustomLogger) SetLevel(level LandcLogFace.LogLevel) {
	// 实现自定义逻辑
}

// GetLevel 获取当前日志级别
func (c *CustomLogger) GetLevel() LandcLogFace.LogLevel {
	return LandcLogFace.InfoLevel
}

// Debug 输出调试级日志
func (c *CustomLogger) Debug(msg string, fields ...LandcLogFace.Field) {
	fmt.Printf("[CUSTOM] [DEBUG] [%s] %s\n", c.name, msg)
}

// Debugf 输出格式化的调试级日志
func (c *CustomLogger) Debugf(format string, args ...interface{}) {
	fmt.Printf("[CUSTOM] [DEBUG] [%s] "+format+"\n", append([]interface{}{c.name}, args...)...)
}

// Info 输出信息级日志
func (c *CustomLogger) Info(msg string, fields ...LandcLogFace.Field) {
	fmt.Printf("[CUSTOM] [INFO] [%s] %s\n", c.name, msg)
}

// Infof 输出格式化的信息级日志
func (c *CustomLogger) Infof(format string, args ...interface{}) {
	fmt.Printf("[CUSTOM] [INFO] [%s] "+format+"\n", append([]interface{}{c.name}, args...)...)
}

// Warn 输出警告级日志
func (c *CustomLogger) Warn(msg string, fields ...LandcLogFace.Field) {
	fmt.Printf("[CUSTOM] [WARN] [%s] %s\n", c.name, msg)
}

// Warnf 输出格式化的警告级日志
func (c *CustomLogger) Warnf(format string, args ...interface{}) {
	fmt.Printf("[CUSTOM] [WARN] [%s] "+format+"\n", append([]interface{}{c.name}, args...)...)
}

// Error 输出错误级日志
func (c *CustomLogger) Error(msg string, fields ...LandcLogFace.Field) {
	fmt.Printf("[CUSTOM] [ERROR] [%s] %s\n", c.name, msg)
}

// Errorf 输出格式化的错误级日志
func (c *CustomLogger) Errorf(format string, args ...interface{}) {
	fmt.Printf("[CUSTOM] [ERROR] [%s] "+format+"\n", append([]interface{}{c.name}, args...)...)
}

// Fatal 输出致命级日志并退出程序
func (c *CustomLogger) Fatal(msg string, fields ...LandcLogFace.Field) {
	fmt.Printf("[CUSTOM] [FATAL] [%s] %s\n", c.name, msg)
}

// Fatalf 输出格式化的致命级日志并退出程序
func (c *CustomLogger) Fatalf(format string, args ...interface{}) {
	fmt.Printf("[CUSTOM] [FATAL] [%s] "+format+"\n", append([]interface{}{c.name}, args...)...)
}

// Panic 输出恐慌级日志并触发panic
func (c *CustomLogger) Panic(msg string, fields ...LandcLogFace.Field) {
	fmt.Printf("[CUSTOM] [PANIC] [%s] %s\n", c.name, msg)
}

// Panicf 输出格式化的恐慌级日志并触发panic
func (c *CustomLogger) Panicf(format string, args ...interface{}) {
	fmt.Printf("[CUSTOM] [PANIC] [%s] "+format+"\n", append([]interface{}{c.name}, args...)...)
}

// WithFields 添加字段到日志
func (c *CustomLogger) WithFields(fields ...LandcLogFace.Field) LandcLogFace.Logger {
	return c
}

// WithField 添加单个字段到日志
func (c *CustomLogger) WithField(key string, value interface{}) LandcLogFace.Logger {
	return c
}

// WithContext 添加上下文到日志
func (c *CustomLogger) WithContext(ctx context.Context) LandcLogFace.Logger {
	return c
}

// WithError 添加错误信息到日志
func (c *CustomLogger) WithError(err error) LandcLogFace.Logger {
	return c
}

// WithTime 添加时间到日志
func (c *CustomLogger) WithTime(t time.Time) LandcLogFace.Logger {
	return c
}

// IsDebugEnabled 检查调试级别是否启用
func (c *CustomLogger) IsDebugEnabled() bool {
	return true
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

// Sync 刷新日志缓冲区
func (c *CustomLogger) Sync() error {
	return nil
}

// CustomLoggerProvider 自定义日志提供者
type CustomLoggerProvider struct{}

// Create 创建日志实例
func (p *CustomLoggerProvider) Create(name string) LandcLogFace.Logger {
	return NewCustomLogger(name)
}

// CreateWithConfig 根据配置创建日志实例
func (p *CustomLoggerProvider) CreateWithConfig(name string, config map[string]interface{}) LandcLogFace.Logger {
	return NewCustomLogger(name)
}
