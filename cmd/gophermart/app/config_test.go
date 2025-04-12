package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig(":8080", "postgres://localhost/db", "http://localhost:8081")
	assert.Equal(t, ":8080", config.GetRunAddress())
	assert.Equal(t, "postgres://localhost/db", config.GetDatabaseURI())
	assert.Equal(t, "http://localhost:8081", config.GetAccrualSystemAddress())
	assert.Equal(t, "test", config.GetJWTSecretKey())
	assert.Equal(t, 365*24*time.Hour, config.GetJWTExpireTime())
}

func TestSetters(t *testing.T) {
	config := NewConfig(":8080", "postgres://localhost/db", "http://localhost:8081")
	config.SetRunAddress(":9090")
	config.SetDatabaseURI("postgres://localhost/otherdb")
	config.SetAccrualSystemAddress("http://localhost:9090")
	assert.Equal(t, ":9090", config.GetRunAddress())
	assert.Equal(t, "postgres://localhost/otherdb", config.GetDatabaseURI())
	assert.Equal(t, "http://localhost:9090", config.GetAccrualSystemAddress())
}

func TestJWTSettings(t *testing.T) {
	config := NewConfig(":8080", "postgres://localhost/db", "http://localhost:8081")
	assert.Equal(t, "test", config.GetJWTSecretKey())
	assert.Equal(t, 365*24*time.Hour, config.GetJWTExpireTime())
}
