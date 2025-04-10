package cli

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func mockCommandSuccess() *cobra.Command {
	return &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
}

func mockCommandError() *cobra.Command {
	return &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("execution failed")
		},
	}
}

func TestRun_Success(t *testing.T) {
	cmd := mockCommandSuccess()
	code := Run(cmd)
	assert.Equal(t, 0, code, "Run should return 0 on success")
}

func TestRun_Error(t *testing.T) {
	cmd := mockCommandError()
	code := Run(cmd)
	assert.Equal(t, 1, code, "Run should return 1 on error")
}
