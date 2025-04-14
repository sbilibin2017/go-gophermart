package cli

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestRun_Success(t *testing.T) {
	cmd := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	exitCode := Run(cmd)
	assert.Equal(t, 0, exitCode, "expected exit code to be 0 on success")
}

func TestRun_Error(t *testing.T) {
	cmd := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("some error")
		},
	}
	exitCode := Run(cmd)
	assert.Equal(t, 1, exitCode, "expected exit code to be 1 on error")
}
