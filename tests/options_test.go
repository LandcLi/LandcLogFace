package tests

import (
	"testing"

	"github.com/LandcLi/landc-logface/lclogface"
)

func TestCreateLoggerWithOptions(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		opts     []lclogface.Option
	}{
		{
			name:     "console with options",
			provider: "console",
			opts: []lclogface.Option{
				lclogface.WithLevel(lclogface.WarnLevel),
				lclogface.WithFormat("json"),
			},
		},
		{
			name:     "std with options",
			provider: "std",
			opts: []lclogface.Option{
				lclogface.WithLevel(lclogface.ErrorLevel),
				lclogface.WithFormat("text"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := lclogface.GetLoggerWithProvider("test", tt.provider, tt.opts...)
			if log == nil {
				t.Fatalf("创建日志实例失败")
			}

			log.Info("测试日志")
		})
	}
}

func TestGetLoggerWithProviderWithOptions(t *testing.T) {
	log := lclogface.GetLoggerWithProvider("test-app", "console",
		lclogface.WithLevel(lclogface.DebugLevel),
		lclogface.WithFormat("json"),
		lclogface.WithOutputPath("stdout"),
	)

	if log == nil {
		t.Fatal("创建日志实例失败")
	}

	log.Debug("调试日志")
	log.Info("信息日志")
}
