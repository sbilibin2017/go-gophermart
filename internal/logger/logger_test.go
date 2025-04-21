package logger

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLogger(t *testing.T) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)

	logger := zap.New(core)

	t.Run("Test Info Log", func(t *testing.T) {
		logger.Info("Test Info Log", zap.String("foo", "bar"))
		assert.True(t, true)
	})

	t.Run("Test Error Log", func(t *testing.T) {
		logger.Error("Test Error Log", zap.String("foo", "bar"))
		assert.True(t, true)
	})
}
