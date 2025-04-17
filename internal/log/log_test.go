package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestLoggerInitialization(t *testing.T) {
	Init()
	assert.NotNil(t, Logger, "Logger should not be nil")
	currentLevel := Logger.Desugar().Core().Enabled(zapcore.InfoLevel)
	assert.True(t, currentLevel, "Expected log level to be InfoLevel")
}
