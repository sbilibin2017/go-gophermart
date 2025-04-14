package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRunAddress(t *testing.T) {
	cfg := &Config{
		RunAddress: "localhost:8080",
	}
	assert.Equal(t, "localhost:8080", cfg.GetRunAddress(), "RunAddress should match expected value")
}

func TestGetDatabaseURI(t *testing.T) {
	cfg := &Config{
		DatabaseURI: "postgres://user:pass@localhost:5432/dbname",
	}
	assert.Equal(t, "postgres://user:pass@localhost:5432/dbname", cfg.GetDatabaseURI(), "DatabaseURI should match expected value")
}
