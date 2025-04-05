package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	runAddress := "http://localhost:8080"
	databaseURI := "mongodb://localhost:27017"
	accrualSystemAddress := "http://localhost:9090"
	config := NewConfig(runAddress, databaseURI, accrualSystemAddress)
	assert.Equal(t, runAddress, config.GetRunAddress(), "Run address should be set correctly")
	assert.Equal(t, databaseURI, config.GetDatabaseURI(), "Database URI should be set correctly")
	assert.Equal(t, accrualSystemAddress, config.GetAccrualSystemAddress(), "Accrual system address should be set correctly")
}

func TestSetRunAddress(t *testing.T) {
	config := NewConfig("http://localhost:8080", "mongodb://localhost:27017", "http://localhost:9090")
	newAddress := "http://localhost:8081"
	config.SetRunAddress(newAddress)
	assert.Equal(t, newAddress, config.GetRunAddress(), "Run address should be updated correctly")
}

func TestSetDatabaseURI(t *testing.T) {
	config := NewConfig("http://localhost:8080", "mongodb://localhost:27017", "http://localhost:9090")
	newURI := "mongodb://localhost:28017"
	config.SetDatabaseURI(newURI)
	assert.Equal(t, newURI, config.GetDatabaseURI(), "Database URI should be updated correctly")
}

func TestSetAccrualSystemAddress(t *testing.T) {
	config := NewConfig("http://localhost:8080", "mongodb://localhost:27017", "http://localhost:9090")
	newAddress := "http://localhost:9091"
	config.SetAccrualSystemAddress(newAddress)
	assert.Equal(t, newAddress, config.GetAccrualSystemAddress(), "Accrual system address should be updated correctly")
}
