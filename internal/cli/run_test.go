package cli

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name       string
		cmd        *cobra.Command
		wantReturn int
	}{
		{
			name: "successful command execution",
			cmd: &cobra.Command{
				RunE: func(cmd *cobra.Command, args []string) error {
					return nil
				},
			},
			wantReturn: 0,
		},
		{
			name: "command execution with error",
			cmd: &cobra.Command{
				RunE: func(cmd *cobra.Command, args []string) error {
					return fmt.Errorf("some error")
				},
			},
			wantReturn: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Run(tt.cmd)
			assert.Equal(t, tt.wantReturn, got, "Run function should return the correct exit code")
		})
	}
}
