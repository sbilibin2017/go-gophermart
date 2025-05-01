package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func InitLoggerWithInfoLevel() {
	initWithLevel(zapcore.InfoLevel)
}

func initWithLevel(level zapcore.Level) error {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(level)
	logger, err := zapConfig.Build()
	if err != nil {
		return err
	}
	Logger = logger.Sugar()
	return nil
}
