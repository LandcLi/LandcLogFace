//go:build logrus_provider

package tests

import (
	"testing"

	"github.com/LandcLi/landc-logface/lclogface"
	_ "github.com/LandcLi/landc-logface/providers/logrus"
)

// TestLogrusProvider 测试Logrus提供者
func TestLogrusProvider(t *testing.T) {
	logger := lclogface.GetLoggerWithProvider("test", "logrus",
		lclogface.WithLevel(lclogface.DebugLevel),
		lclogface.WithFormat("text"),
	)
	if logger == nil {
		t.Fatal("创建Logrus日志失败")
	}

	logger.Info("Logrus日志测试")
	logger.Debug("Logrus调试日志")
	logger.Warn("Logrus警告日志")
	logger.Error("Logrus错误日志")
}

// TestLogrusWithOptions 测试Logrus带选项
func TestLogrusWithOptions(t *testing.T) {
	logger := lclogface.GetLoggerWithProvider("test", "logrus",
		lclogface.WithLevel(lclogface.InfoLevel),
		lclogface.WithFormat("json"),
		lclogface.WithOutputPath("stdout"),
		lclogface.WithMaxMessageSize(10),
	)
	if logger == nil {
		t.Fatal("创建Logrus日志失败")
	}

	logger.Info("Logrus带选项测试")
}
