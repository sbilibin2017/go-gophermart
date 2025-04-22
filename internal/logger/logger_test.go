package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestLoggerInitialization(t *testing.T) {
	Init(zapcore.InfoLevel)
	require.NotNil(t, Logger)
	assert.Equal(t, true, Logger.Core().Enabled(zapcore.InfoLevel))
	assert.NotPanics(t, func() {
		Logger.Info("Test log")
	})
}
