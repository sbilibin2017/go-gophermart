package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()

	assert.NotNil(t, cfg)
	assert.Equal(t, "test", cfg.JWTSecretKey)
	assert.Equal(t, 365*24*time.Hour, cfg.JWTExp)
}

func TestGetRunAddress(t *testing.T) {
	cfg := &Config{RunAddress: "localhost:8080"}
	assert.Equal(t, "localhost:8080", cfg.GetRunAddress())
}

func TestGetJWTSecretKey(t *testing.T) {
	cfg := &Config{JWTSecretKey: "my_secret"}
	assert.Equal(t, "my_secret", cfg.GetJWTSecretKey())
}

func TestGetJWTExp(t *testing.T) {
	expected := 7 * 24 * time.Hour
	cfg := &Config{JWTExp: expected}
	assert.Equal(t, expected, cfg.GetJWTExp())
}
