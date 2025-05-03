package logger

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
	zapConfig.DisableCaller = true
	zapConfig.Level = zap.NewAtomicLevelAt(level)
	zapConfig.Encoding = "json"
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "@timestamp"
	encoderConfig.CallerKey = "@caller"
	encoderConfig.MessageKey = "message"
	encoderConfig.LevelKey = "@level"
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConfig.EncoderConfig = encoderConfig
	lgr, _ := zapConfig.Build()
	logger = lgr.Sugar()
}

func getCallerInfo() (string, int) {
	_, file, line, _ := runtime.Caller(3)
	return file, line
}

func log(level zapcore.Level, msg string, args ...any) {
	file, line := getCallerInfo()
	args = append(args, "@caller", file+":"+strconv.Itoa(line))
	if level == zapcore.ErrorLevel {
		logger.Errorw(msg, args...)
	} else {
		logger.Infow(msg, args...)
	}
}

func Info(msg string, args ...any) {
	if logger != nil {
		log(zapcore.InfoLevel, msg, args...)
	}
}

func Error(msg string, args ...any) {
	if logger != nil {
		log(zapcore.ErrorLevel, msg, args...)
	}
}
