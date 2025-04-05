package app

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {
	os.Setenv(EnvRunAddress, "localhost:8081")
	os.Setenv(EnvDatabaseURI, "mongodb://localhost")
	os.Setenv(EnvAccrualSystemAddress, "http://accrual-system")
	defer func() {
		os.Unsetenv(EnvRunAddress)
		os.Unsetenv(EnvDatabaseURI)
		os.Unsetenv(EnvAccrualSystemAddress)
	}()

	viper.Reset()
	cmd := NewCommand()

	assert.Equal(t, "Gophermart", cmd.Use)
	assert.Equal(t, "Gophermart Service", cmd.Short)

	assert.Equal(t, "localhost:8081", viper.GetString(EnvRunAddress))
	assert.Equal(t, "mongodb://localhost", viper.GetString(EnvDatabaseURI))
	assert.Equal(t, "http://accrual-system", viper.GetString(EnvAccrualSystemAddress))
}
