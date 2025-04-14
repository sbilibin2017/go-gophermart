package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun_Default(t *testing.T) {
	code := Run()
	assert.Equal(t, 0, code, "Run() should return exit code 0 by default")
}
