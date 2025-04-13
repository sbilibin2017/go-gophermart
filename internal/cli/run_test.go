package cli

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	errFunc := func() error {
		return errors.New("some error")
	}

	exitCode := Run(errFunc)
	assert.Equal(t, 1, exitCode, "Expected exit code 1 when an error occurs")

	successFunc := func() error {
		return nil
	}

	exitCode = Run(successFunc)
	assert.Equal(t, 0, exitCode, "Expected exit code 0 when there is no error")
}
