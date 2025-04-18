package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccrualConfig(t *testing.T) {
	tests := []struct {
		name         string
		runAddress   string
		databaseURI  string
		expectedAddr string
		expectedURI  string
	}{
		{
			name:         "valid input",
			runAddress:   "localhost:8080",
			databaseURI:  "postgres://user:password@localhost/db",
			expectedAddr: "localhost:8080",
			expectedURI:  "postgres://user:password@localhost/db",
		},
		{
			name:         "empty run address",
			runAddress:   "",
			databaseURI:  "postgres://user:password@localhost/db",
			expectedAddr: "",
			expectedURI:  "postgres://user:password@localhost/db",
		},
		{
			name:         "empty database URI",
			runAddress:   "localhost:8080",
			databaseURI:  "",
			expectedAddr: "localhost:8080",
			expectedURI:  "",
		},
		{
			name:         "both empty",
			runAddress:   "",
			databaseURI:  "",
			expectedAddr: "",
			expectedURI:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := NewAccrualConfig(tt.runAddress, tt.databaseURI)
			assert.Equal(t, tt.expectedAddr, cfg.RunAddress)
			assert.Equal(t, tt.expectedURI, cfg.DatabaseURI)
		})
	}
}
