package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccrualConfig(t *testing.T) {
	runAddr := "localhost:8080"
	dbURI := "postgres://user:pass@localhost:5432/db"

	cfg := NewAccrualConfig(runAddr, dbURI)

	assert.NotNil(t, cfg)
	assert.Equal(t, runAddr, cfg.GetRunAddress())
	assert.Equal(t, dbURI, cfg.GetDatabaseURI())
}
