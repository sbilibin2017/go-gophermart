package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name         string
		f            func() error
		expectedExit int
	}{
		{
			name: "Function returns no error",
			f: func() error {
				return nil
			},
			expectedExit: 0,
		},
		{
			name: "Function returns an error",
			f: func() error {
				return assert.AnError
			},
			expectedExit: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exitCode := Run(tt.f)
			assert.Equal(t, tt.expectedExit, exitCode)
		})
	}
}
