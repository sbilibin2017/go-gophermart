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
	defer func() {
		os.Unsetenv(EnvRunAddress)
		os.Unsetenv(EnvDatabaseURI)
	}()
	viper.Reset()
	cmd := NewCommand()
	assert.Equal(t, "Accrual", cmd.Use)
	assert.Equal(t, "Accrual Service", cmd.Short)
	assert.Equal(t, "localhost:8081", viper.GetString(EnvRunAddress))
	assert.Equal(t, "mongodb://localhost", viper.GetString(EnvDatabaseURI))
}
