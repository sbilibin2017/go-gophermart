package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func init() {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	logInstance, _ := zapConfig.Build()
	Logger = logInstance.Sugar()
}
