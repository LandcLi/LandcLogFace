//go:build zap_provider

package tests

import (
	"testing"

	"github.com/LandcLi/landc-logface/lclogface"
	_ "github.com/LandcLi/landc-logface/providers/zap"
)

// TestZapProvider 测试Zap提供者
func TestZapProvider(t *testing.T) {
	logger := lclogface.GetLoggerWithProvider("test", "zap",
		lclogface.WithLevel(lclogface.DebugLevel),
		lclogface.WithFormat("json"),
	)
	if logger == nil {
		t.Fatal("创建Zap日志失败")
	}

	logger.Info("Zap日志测试")
	logger.Debug("Zap调试日志")
	logger.Warn("Zap警告日志")
	logger.Error("Zap错误日志")
}

// TestZapWithOptions 测试Zap带选项
func TestZapWithOptions(t *testing.T) {
	logger := lclogface.GetLoggerWithProvider("test", "zap",
		lclogface.WithLevel(lclogface.InfoLevel),
		lclogface.WithFormat("json"),
		lclogface.WithOutputPath("stdout"),
		lclogface.WithMaxMessageSize(10),
	)
	if logger == nil {
		t.Fatal("创建Zap日志失败")
	}

	logger.Info("Zap带选项测试")
}
