package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var err error

func Init(level zapcore.Level) error {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(level)
	Logger, err = zapConfig.Build()
	if err != nil {
		return err
	}
	return nil
}
