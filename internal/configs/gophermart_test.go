package configs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewGophermartConfig(t *testing.T) {
	tests := []struct {
		name                   string
		runAddress             string
		databaseURI            string
		accrualSystemAddress   string
		expectedRunAddress     string
		expectedDatabaseURI    string
		expectedAccrualAddress string
		expectedJWTSecretKey   []byte
		expectedJWTExp         time.Duration
	}{
		{
			name:                   "valid input",
			runAddress:             "localhost:8080",
			databaseURI:            "postgres://user:password@localhost/db",
			accrualSystemAddress:   "http://localhost:8081",
			expectedRunAddress:     "localhost:8080",
			expectedDatabaseURI:    "postgres://user:password@localhost/db",
			expectedAccrualAddress: "http://localhost:8081",
			expectedJWTSecretKey:   []byte("test"),
			expectedJWTExp:         time.Duration(365 * 24 * time.Hour),
		},
		{
			name:                   "empty run address",
			runAddress:             "",
			databaseURI:            "postgres://user:password@localhost/db",
			accrualSystemAddress:   "http://localhost:8081",
			expectedRunAddress:     "",
			expectedDatabaseURI:    "postgres://user:password@localhost/db",
			expectedAccrualAddress: "http://localhost:8081",
			expectedJWTSecretKey:   []byte("test"),
			expectedJWTExp:         time.Duration(365 * 24 * time.Hour),
		},
		{
			name:                   "empty database URI",
			runAddress:             "localhost:8080",
			databaseURI:            "",
			accrualSystemAddress:   "http://localhost:8081",
			expectedRunAddress:     "localhost:8080",
			expectedDatabaseURI:    "",
			expectedAccrualAddress: "http://localhost:8081",
			expectedJWTSecretKey:   []byte("test"),
			expectedJWTExp:         time.Duration(365 * 24 * time.Hour),
		},
		{
			name:                   "empty accrual system address",
			runAddress:             "localhost:8080",
			databaseURI:            "postgres://user:password@localhost/db",
			accrualSystemAddress:   "",
			expectedRunAddress:     "localhost:8080",
			expectedDatabaseURI:    "postgres://user:password@localhost/db",
			expectedAccrualAddress: "",
			expectedJWTSecretKey:   []byte("test"),
			expectedJWTExp:         time.Duration(365 * 24 * time.Hour),
		},
		{
			name:                   "all empty",
			runAddress:             "",
			databaseURI:            "",
			accrualSystemAddress:   "",
			expectedRunAddress:     "",
			expectedDatabaseURI:    "",
			expectedAccrualAddress: "",
			expectedJWTSecretKey:   []byte("test"),
			expectedJWTExp:         time.Duration(365 * 24 * time.Hour),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := NewGophermartConfig(tt.runAddress, tt.databaseURI, tt.accrualSystemAddress)
			assert.Equal(t, tt.expectedRunAddress, cfg.RunAddress)
			assert.Equal(t, tt.expectedDatabaseURI, cfg.DatabaseURI)
			assert.Equal(t, tt.expectedAccrualAddress, cfg.AccrualSystemAddress)
			assert.Equal(t, tt.expectedJWTSecretKey, cfg.JWTSecretKey)
			assert.Equal(t, tt.expectedJWTExp, cfg.JWTExp)
		})
	}
}
