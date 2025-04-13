package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_GetDatabaseURI(t *testing.T) {
	cfg := &Config{
		DatabaseURI: "postgres://user:pass@localhost:5432/dbname",
	}

	assert.Equal(t, "postgres://user:pass@localhost:5432/dbname", cfg.GetDatabaseURI())
}

func TestConfig_GetRunAddress(t *testing.T) {
	cfg := &Config{
		RunAddress: ":8080",
	}

	assert.Equal(t, ":8080", cfg.GetRunAddress())
}
