package log

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initTestLogger(buf *bytes.Buffer) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = ""
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(buf),
		zapcore.InfoLevel,
	)
	testLogger := zap.New(core).Sugar()
	logger = testLogger
}

func TestInfoLogging(t *testing.T) {
	var buf bytes.Buffer
	initTestLogger(&buf)
	Info("This is an info message", "key1", "value1")
	logOutput := buf.String()
	assert.Contains(t, logOutput, "This is an info message")
	assert.Contains(t, logOutput, `"key1":"value1"`)
}

func TestErrorLogging(t *testing.T) {
	var buf bytes.Buffer
	initTestLogger(&buf)
	Error("This is an error message", "err", "some error")
	logOutput := buf.String()
	assert.Contains(t, logOutput, "This is an error message")
	assert.Contains(t, logOutput, `"err":"some error"`)
}
