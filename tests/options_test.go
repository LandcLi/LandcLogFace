package tests

import (
	"testing"

	"github.com/LandcLi/LandcLogFace"
)

func TestCreateLoggerWithOptions(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		opts     []LandcLogFace.Option
	}{
		{
			name:     "console with options",
			provider: "console",
			opts: []LandcLogFace.Option{
				LandcLogFace.WithLevel(LandcLogFace.WarnLevel),
				LandcLogFace.WithFormat("json"),
			},
		},
		{
			name:     "std with options",
			provider: "std",
			opts: []LandcLogFace.Option{
				LandcLogFace.WithLevel(LandcLogFace.ErrorLevel),
				LandcLogFace.WithFormat("text"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := LandcLogFace.GetLoggerWithProvider("test", tt.provider, tt.opts...)
			if log == nil {
				t.Fatalf("创建日志实例失败")
			}

			log.Info("测试日志")
		})
	}
}

func TestGetLoggerWithProviderWithOptions(t *testing.T) {
	log := LandcLogFace.GetLoggerWithProvider("test-app", "console",
		LandcLogFace.WithLevel(LandcLogFace.DebugLevel),
		LandcLogFace.WithFormat("json"),
		LandcLogFace.WithOutputPath("stdout"),
	)

	if log == nil {
		t.Fatal("创建日志实例失败")
	}

	log.Debug("调试日志")
	log.Info("信息日志")
}
