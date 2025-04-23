package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
	Init()
	assert.NotNil(t, Logger, "Logger should be initialized")
}
