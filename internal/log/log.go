package log

import (
	"runtime"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var level = zapcore.InfoLevel
var logger *zap.SugaredLogger

func init() {
	zapConfig := zap.NewProductionConfig()
	zapConfig.EncoderConfig.CallerKey = ""
	zapConfig.EncoderConfig.MessageKey = ""
	zapConfig.Level = zap.NewAtomicLevelAt(level)
	logInstance, _ := zapConfig.Build()
	logger = logInstance.Sugar()
}

func getCaller() string {
	_, file, line, _ := runtime.Caller(2)
	return file + ":" + strconv.Itoa(line)
}

func Info(msg string, keysAndValues ...interface{}) {
	if logger != nil {
		keysAndValues = append(keysAndValues, "caller", getCaller(), "msg", msg)
		logger.Infow(msg, keysAndValues...)
	}
}

func Error(msg string, keysAndValues ...interface{}) {
	if logger != nil {
		keysAndValues = append(keysAndValues, "caller", getCaller(), "msg", msg)
		logger.Errorw(msg, keysAndValues...)
	}
}
