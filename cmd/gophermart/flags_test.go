package main

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func resetFlags() {
	// Reset the flag set before each test to avoid redefinition issues
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.PanicOnError)
}

func resetOsArgs(args []string) {
	// Reset the os.Args before each test
	os.Args = args
}

func TestFlags(t *testing.T) {
	tests := []struct {
		name                   string
		envVars                map[string]string
		args                   []string
		expectedRunAddress     string
		expectedDatabaseURI    string
		expectedAccrualAddress string
	}{
		{
			name: "With environment variables and command-line args",
			envVars: map[string]string{
				"RUN_ADDRESS":            "localhost:8080",
				"DATABASE_URI":           "localhost:5432",
				"ACCRUAL_SYSTEM_ADDRESS": "localhost:9000",
			},
			args:                   []string{"cmd", "-a", "127.0.0.1:8081", "-d", "127.0.0.1:5433", "-r", "127.0.0.1:9001"},
			expectedRunAddress:     "localhost:8080",
			expectedDatabaseURI:    "localhost:5432",
			expectedAccrualAddress: "localhost:9000",
		},
		{
			name: "With empty environment variables",
			envVars: map[string]string{
				"RUN_ADDRESS":            "",
				"DATABASE_URI":           "",
				"ACCRUAL_SYSTEM_ADDRESS": "",
			},
			args:                   []string{"cmd", "-a", "127.0.0.1:8081", "-d", "127.0.0.1:5433", "-r", "127.0.0.1:9001"},
			expectedRunAddress:     "127.0.0.1:8081",
			expectedDatabaseURI:    "127.0.0.1:5433",
			expectedAccrualAddress: "127.0.0.1:9001",
		},
		{
			name: "With no command-line args and no environment variables",
			envVars: map[string]string{
				"RUN_ADDRESS":            "",
				"DATABASE_URI":           "",
				"ACCRUAL_SYSTEM_ADDRESS": "",
			},
			args:                   []string{"cmd"}, // No command-line args
			expectedRunAddress:     "",
			expectedDatabaseURI:    "",
			expectedAccrualAddress: "",
		},
		{
			name: "With missing environment variables",
			envVars: map[string]string{
				"RUN_ADDRESS":            "",
				"DATABASE_URI":           "",
				"ACCRUAL_SYSTEM_ADDRESS": "",
			},
			args:                   []string{"cmd", "-a", "127.0.0.1:8081", "-d", "127.0.0.1:5433", "-r", "127.0.0.1:9001"},
			expectedRunAddress:     "127.0.0.1:8081",
			expectedDatabaseURI:    "127.0.0.1:5433",
			expectedAccrualAddress: "127.0.0.1:9001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			// Reset flags and os.Args
			t.Cleanup(func() {
				resetFlags()
				resetOsArgs([]string{"cmd"})
			})

			// Set the command-line arguments
			os.Args = tt.args

			// Call flags function
			flags()

			// Validate the values
			assert.Equal(t, tt.expectedRunAddress, runAddress)
			assert.Equal(t, tt.expectedDatabaseURI, databaseURI)
			assert.Equal(t, tt.expectedAccrualAddress, accrualSystemAddress)
		})
	}
}
