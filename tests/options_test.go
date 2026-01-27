package tests

import (
	"testing"

	"github.com/LandcLi/LandcLogFace/internal/logger"
)

func TestCreateLoggerWithOptions(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		opts     []logger.Option
	}{
		{
			name:     "zap with options",
			provider: "zap",
			opts: []logger.Option{
				logger.WithLevel(logger.DebugLevel),
				logger.WithFormat("json"),
			},
		},
		{
			name:     "logrus with options",
			provider: "logrus",
			opts: []logger.Option{
				logger.WithLevel(logger.InfoLevel),
				logger.WithFormat("text"),
			},
		},
		{
			name:     "console with options",
			provider: "console",
			opts: []logger.Option{
				logger.WithLevel(logger.WarnLevel),
				logger.WithFormat("json"),
			},
		},
		{
			name:     "std with options",
			provider: "std",
			opts: []logger.Option{
				logger.WithLevel(logger.ErrorLevel),
				logger.WithFormat("text"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := logger.GetLogFactory().CreateLoggerWithProvider("test", tt.provider, tt.opts...)
			if log == nil {
				t.Fatalf("创建日志实例失败")
			}

			log.Info("测试日志")
		})
	}
}

func TestGetLoggerWithProviderWithOptions(t *testing.T) {
	log := logger.GetLoggerWithProvider("test-app", "zap",
		logger.WithLevel(logger.DebugLevel),
		logger.WithFormat("json"),
		logger.WithOutputPath("stdout"),
	)

	if log == nil {
		t.Fatal("创建日志实例失败")
	}

	log.Debug("调试日志")
	log.Info("信息日志")
}
