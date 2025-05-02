package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestLoggerInit(t *testing.T) {
	initWithLevel(zapcore.InfoLevel)
	assert.NotNil(t, Logger, "Logger should be initialized")
}
