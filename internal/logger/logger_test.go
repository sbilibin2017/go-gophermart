package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestInitLogger(t *testing.T) {
	Init(zapcore.InfoLevel)
	assert.NotNil(t, Logger, "Logger should be initialized")
}
