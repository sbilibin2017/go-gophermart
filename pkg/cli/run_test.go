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
	code := Run(cmd)
	assert.Equal(t, 0, code)
}

func TestRun_Error(t *testing.T) {
	cmd := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("some error")
		},
	}
	code := Run(cmd)
	assert.Equal(t, 1, code)
}
