package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLoggerWithInfoLevel() {
	initWithLevel(zapcore.InfoLevel)
}

func initWithLevel(level zapcore.Level) {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(level)
	Logger, _ = zapConfig.Build()
}
