package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func init() {
	Logger = NewLogger(zapcore.InfoLevel)
}

func NewLogger(level zapcore.Level) *zap.Logger {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(level)
	logger, _ := zapConfig.Build()
	return logger
}
