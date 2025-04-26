package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func initLogger(level zapcore.Level) {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(level)
	rawLogger, _ := zapConfig.Build()
	Logger = rawLogger.Sugar()
}

func InitWithInfoLevel() {
	initLogger(zapcore.InfoLevel)
}
