package log

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestInit_InfoLevel(t *testing.T) {
	err := Init(LevelInfo)
	assert.NoError(t, err)
	assert.NotNil(t, logger)
}

func TestInit_ErrorLevel(t *testing.T) {
	err := Init(LevelError)
	assert.NoError(t, err)
	assert.NotNil(t, logger)
}

func TestInit_DefaultLevel(t *testing.T) {
	err := Init("unknown") // should default to info
	assert.NoError(t, err)
	assert.NotNil(t, logger)
}

func TestInfoAndErrorLogs(t *testing.T) {
	var buf bytes.Buffer
	writer := zapcore.AddSync(&buf)

	encoderCfg := zap.NewProductionEncoderConfig()
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		writer,
		zapcore.InfoLevel,
	)
	testLogger := zap.New(core).Sugar()

	// Override the global logger with testLogger
	mu.Lock()
	logger = testLogger
	mu.Unlock()

	Info("info message", "key1", "val1")
	Error("error message", "key2", "val2")

	logged := buf.String()

	assert.Contains(t, logged, `"msg":"info message"`)
	assert.Contains(t, logged, `"key1":"val1"`)
	assert.Contains(t, logged, `"msg":"error message"`)
	assert.Contains(t, logged, `"key2":"val2"`)
}

func TestSync_NoPanic(t *testing.T) {
	err := Init(LevelInfo)
	assert.NoError(t, err)
	assert.NotPanics(t, func() {
		Sync()
	})
}
