package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerInitialization(t *testing.T) {
	assert.NotNil(t, Logger, "Logger должен быть инициализирован")
}
