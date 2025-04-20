package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestNewLogger(t *testing.T) {
	log := NewLogger(zapcore.DebugLevel)
	require.NotNil(t, log)
	log.Debug("this is a debug message")
	log.Info("this is an info message")
}
