package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var level = zapcore.InfoLevel
var Logger *zap.SugaredLogger

func Init() {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(level)
	logInstance, _ := zapConfig.Build()
	Logger = logInstance.Sugar()
}


