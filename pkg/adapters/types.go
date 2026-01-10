package adapters

import (
	"github.com/LandcLi/LandcLogFace/pkg/logger"
)

// Logger 是日志门面接口类型别名
type Logger = logger.Logger

// LogLevel 定义日志级别类型别名
type LogLevel = logger.LogLevel

// Field 定义日志字段类型别名
type Field = logger.Field

// DebugLevel 调试级别
const DebugLevel = logger.DebugLevel

// InfoLevel 信息级别
const InfoLevel = logger.InfoLevel

// WarnLevel 警告级别
const WarnLevel = logger.WarnLevel

// ErrorLevel 错误级别
const ErrorLevel = logger.ErrorLevel

// FatalLevel 致命级别
const FatalLevel = logger.FatalLevel

// PanicLevel 恐慌级别
const PanicLevel = logger.PanicLevel
