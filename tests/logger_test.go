package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/LandcLi/LandcLogFace"
)

// TestLoggerInterface 测试Logger接口的基本功能
func TestLoggerInterface(t *testing.T) {
	// 测试控制台日志
	logger := LandcLogFace.GetLoggerWithProvider("test", "console")

	// 测试日志级别设置
	logger.SetLevel(LandcLogFace.DebugLevel)
	if logger.GetLevel() != LandcLogFace.DebugLevel {
		t.Errorf("Expected level DebugLevel, got %v", logger.GetLevel())
	}

	// 测试日志输出（这里只是测试方法调用，不测试输出内容）
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warn message")
	logger.Error("Error message")

	// 测试格式化日志
	logger.Debugf("Debug: %s", "test")
	logger.Infof("Info: %s", "test")
	logger.Warnf("Warn: %s", "test")
	logger.Errorf("Error: %s", "test")

	// 测试字段
	logger.Info("With fields",
		LandcLogFace.Field{Key: "key1", Value: "value1"},
		LandcLogFace.Field{Key: "key2", Value: 123},
	)

	// 测试链式调用
	logger.WithField("chain", "value").Info("Chained logger")

	// 测试上下文
	ctx := context.Background()
	logger.WithContext(ctx).Info("With context")

	// 测试错误
	err := errors.New("test error")
	logger.WithError(err).Error("With error")

	// 测试时间
	now := time.Now()
	logger.WithTime(now).Info("With time")

	// 测试日志级别检查
	if !logger.IsDebugEnabled() {
		t.Error("Expected DebugLevel to be enabled")
	}
	if !logger.IsInfoEnabled() {
		t.Error("Expected InfoLevel to be enabled")
	}
	if !logger.IsWarnEnabled() {
		t.Error("Expected WarnLevel to be enabled")
	}
	if !logger.IsErrorEnabled() {
		t.Error("Expected ErrorLevel to be enabled")
	}

	// 测试Sync
	if err := logger.Sync(); err != nil {
		t.Errorf("Sync failed: %v", err)
	}
}

// TestGlobalLogger 测试全局日志
func TestGlobalLogger(t *testing.T) {
	// 测试获取全局日志实例
	globalLogger := LandcLogFace.GetLogger()
	if globalLogger == nil {
		t.Error("Failed to get global logger")
	}

	// 测试全局日志函数
	LandcLogFace.Debug("Global debug")
	LandcLogFace.Info("Global info")
	LandcLogFace.Warn("Global warn")
	LandcLogFace.Error("Global error")

	// 测试全局格式化函数
	LandcLogFace.Debugf("Global debug: %s", "test")
	LandcLogFace.Infof("Global info: %s", "test")
	LandcLogFace.Warnf("Global warn: %s", "test")
	LandcLogFace.Errorf("Global error: %s", "test")

	// 测试全局函数带字段
	LandcLogFace.Info("Global with fields", LandcLogFace.Field{Key: "key", Value: "value"})
}

// TestZapLogger 测试zap日志适配器
func TestZapLogger(t *testing.T) {
	logger := LandcLogFace.GetLoggerWithProvider("test-zap", "zap")
	logger.SetLevel(LandcLogFace.DebugLevel)

	if logger == nil {
		t.Error("Failed to create zap logger")
	}

	// 测试zap日志功能
	logger.Debug("Zap debug")
	logger.Info("Zap info")
	logger.Warn("Zap warn")
	logger.Error("Zap error")

	// 测试字段
	logger.Info("Zap with fields",
		LandcLogFace.Field{Key: "zap", Value: "test"},
	)

	// 测试Sync
	if err := logger.Sync(); err != nil {
		t.Errorf("Zap sync failed: %v", err)
	}
}

// TestLogrusLogger 测试logrus日志适配器
func TestLogrusLogger(t *testing.T) {
	logger := LandcLogFace.GetLoggerWithProvider("test-logrus", "logrus")
	logger.SetLevel(LandcLogFace.DebugLevel)

	if logger == nil {
		t.Error("Failed to create logrus logger")
	}

	// 测试logrus日志功能
	logger.Debug("Logrus debug")
	logger.Info("Logrus info")
	logger.Warn("Logrus warn")
	logger.Error("Logrus error")

	// 测试字段
	logger.Info("Logrus with fields",
		LandcLogFace.Field{Key: "logrus", Value: "test"},
	)

	// 测试Sync
	if err := logger.Sync(); err != nil {
		t.Errorf("Logrus sync failed: %v", err)
	}
}

// TestStdLogger 测试标准库日志适配器
func TestStdLogger(t *testing.T) {
	logger := LandcLogFace.GetLoggerWithProvider("test-std", "std")
	logger.SetLevel(LandcLogFace.DebugLevel)

	if logger == nil {
		t.Error("Failed to create std logger")
	}

	// 测试std日志功能
	logger.Debug("Std debug")
	logger.Info("Std info")
	logger.Warn("Std warn")
	logger.Error("Std error")

	// 测试字段
	logger.Info("Std with fields",
		LandcLogFace.Field{Key: "std", Value: "test"},
	)

	// 测试Sync
	if err := logger.Sync(); err != nil {
		t.Errorf("Std sync failed: %v", err)
	}
}

// TestConsoleLogger 测试控制台日志适配器
func TestConsoleLogger(t *testing.T) {
	logger := LandcLogFace.GetLoggerWithProvider("test-console", "console")
	logger.SetLevel(LandcLogFace.DebugLevel)

	if logger == nil {
		t.Error("Failed to create console logger")
	}

	// 测试console日志功能
	logger.Debug("Console debug")
	logger.Info("Console info")
	logger.Warn("Console warn")
	logger.Error("Console error")

	// 测试字段
	logger.Info("Console with fields",
		LandcLogFace.Field{Key: "console", Value: "test"},
	)

	// 测试Sync
	if err := logger.Sync(); err != nil {
		t.Errorf("Console sync failed: %v", err)
	}
}

// TestLoggerOptions 测试日志选项
func TestLoggerOptions(t *testing.T) {
	// 测试选项函数
	levelOpt := LandcLogFace.WithLevel(LandcLogFace.DebugLevel)
	formatOpt := LandcLogFace.WithFormat("json")
	outputOpt := LandcLogFace.WithOutputPath("stdout")
	configOpt := LandcLogFace.WithConfig(map[string]interface{}{"key": "value"})

	options := &LandcLogFace.LoggerOptions{}

	// 应用选项
	levelOpt(options)
	formatOpt(options)
	outputOpt(options)
	configOpt(options)

	// 验证选项
	if options.Level != LandcLogFace.DebugLevel {
		t.Errorf("Expected level DebugLevel, got %v", options.Level)
	}

	if options.Format != "json" {
		t.Errorf("Expected format 'json', got '%s'", options.Format)
	}

	if options.OutputPath != "stdout" {
		t.Errorf("Expected output path 'stdout', got '%s'", options.OutputPath)
	}

	if options.Config["key"] != "value" {
		t.Errorf("Expected config key 'key' to be 'value', got '%v'", options.Config["key"])
	}
}

// TestLogLevelString 测试日志级别字符串表示
func TestLogLevelString(t *testing.T) {
	testCases := []struct {
		level    LandcLogFace.LogLevel
		expected string
	}{
		{LandcLogFace.DebugLevel, "DEBUG"},
		{LandcLogFace.InfoLevel, "INFO"},
		{LandcLogFace.WarnLevel, "WARN"},
		{LandcLogFace.ErrorLevel, "ERROR"},
		{LandcLogFace.FatalLevel, "FATAL"},
		{LandcLogFace.PanicLevel, "PANIC"},
		{LandcLogFace.LogLevel(999), "UNKNOWN"},
	}

	for _, tc := range testCases {
		if tc.level.String() != tc.expected {
			t.Errorf("Expected level %v to be '%s', got '%s'", tc.level, tc.expected, tc.level.String())
		}
	}
}

// TestWithMethods 测试With系列方法
func TestWithMethods(t *testing.T) {
	logger := LandcLogFace.GetLoggerWithProvider("test-with", "console")

	// 测试WithField
	logger1 := logger.WithField("key1", "value1")
	if logger1 == nil {
		t.Error("Failed to create logger with field")
	}

	// 测试WithFields
	logger2 := logger.WithFields(
		LandcLogFace.Field{Key: "key1", Value: "value1"},
		LandcLogFace.Field{Key: "key2", Value: "value2"},
	)
	if logger2 == nil {
		t.Error("Failed to create logger with fields")
	}

	// 测试WithContext
	ctx := context.Background()
	logger3 := logger.WithContext(ctx)
	if logger3 == nil {
		t.Error("Failed to create logger with context")
	}

	// 测试WithError
	err := errors.New("test error")
	logger4 := logger.WithError(err)
	if logger4 == nil {
		t.Error("Failed to create logger with error")
	}

	// 测试WithTime
	now := time.Now()
	logger5 := logger.WithTime(now)
	if logger5 == nil {
		t.Error("Failed to create logger with time")
	}
}

// TestIsEnabledMethods 测试IsEnabled系列方法
func TestIsEnabledMethods(t *testing.T) {
	// 测试DebugLevel
	debugLogger := LandcLogFace.GetLoggerWithProvider("test-debug", "console")
	debugLogger.SetLevel(LandcLogFace.DebugLevel)
	if !debugLogger.IsDebugEnabled() {
		t.Error("Debug level should be enabled")
	}
	if !debugLogger.IsInfoEnabled() {
		t.Error("Info level should be enabled")
	}
	if !debugLogger.IsWarnEnabled() {
		t.Error("Warn level should be enabled")
	}
	if !debugLogger.IsErrorEnabled() {
		t.Error("Error level should be enabled")
	}

	// 测试InfoLevel
	infoLogger := LandcLogFace.GetLoggerWithProvider("test-info", "console")
	infoLogger.SetLevel(LandcLogFace.InfoLevel)
	if infoLogger.IsDebugEnabled() {
		t.Error("Debug level should not be enabled")
	}
	if !infoLogger.IsInfoEnabled() {
		t.Error("Info level should be enabled")
	}

	// 测试ErrorLevel
	errorLogger := LandcLogFace.GetLoggerWithProvider("test-error", "console")
	errorLogger.SetLevel(LandcLogFace.ErrorLevel)
	if errorLogger.IsDebugEnabled() {
		t.Error("Debug level should not be enabled")
	}
	if errorLogger.IsInfoEnabled() {
		t.Error("Info level should not be enabled")
	}
	if errorLogger.IsWarnEnabled() {
		t.Error("Warn level should not be enabled")
	}
	if !errorLogger.IsErrorEnabled() {
		t.Error("Error level should be enabled")
	}
}
