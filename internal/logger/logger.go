package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func init() {
	initWithLevel(zapcore.InfoLevel)
}

func initWithLevel(level zapcore.Level) {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(level)
	logger, _ := zapConfig.Build()
	Logger = logger.Sugar()
}
