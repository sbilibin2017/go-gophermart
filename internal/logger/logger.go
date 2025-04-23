package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	Init(zapcore.InfoLevel)
}

var Logger *zap.Logger

func Init(level zapcore.Level) {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(level)
	Logger, _ = zapConfig.Build()
}
