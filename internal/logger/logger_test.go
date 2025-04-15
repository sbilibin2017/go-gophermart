package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestInit(t *testing.T) {
	Logger = nil
	assert.Nil(t, Logger, "Logger must be nil")
	Init(zapcore.DebugLevel)
	assert.NotNil(t, Logger, "Logger must be initialized")
	assert.NotPanics(t, func() {
		Logger.Debugf("Debug log test")
		Logger.Infof("Info log test")
		Logger.Warnf("Warn log test")
	})
}
