package app

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {
	cmd := NewCommand()

	go func() {
		cmd.Execute()
	}()

	time.Sleep(2 * time.Second)
}

func TestFlags(t *testing.T) {
	cmd := NewCommand()

	cmd.Flags().Set(FlagLongRunAddress, ":9090")
	cmd.Flags().Set(FlagLongDatabaseURI, "mongodb://localhost:27017")
	cmd.Flags().Set(FlagLongAccrualSystemAddress, "http://localhost:8081")

	assert.Equal(t, ":9090", viper.GetString(FlagLongRunAddress))
	assert.Equal(t, "mongodb://localhost:27017", viper.GetString(FlagLongDatabaseURI))
	assert.Equal(t, "http://localhost:8081", viper.GetString(FlagLongAccrualSystemAddress))
}

func TestEnvVars(t *testing.T) {
	os.Setenv(EnvRunAddress, ":9090")
	os.Setenv(EnvDatabaseURI, "mongodb://localhost:27017")
	os.Setenv(EnvAccrualSystemAddress, "http://localhost:8081")

	viper.BindEnv(FlagLongRunAddress, EnvRunAddress)
	viper.BindEnv(FlagLongDatabaseURI, EnvDatabaseURI)
	viper.BindEnv(FlagLongAccrualSystemAddress, EnvAccrualSystemAddress)

	assert.Equal(t, ":9090", viper.GetString(FlagLongRunAddress))
	assert.Equal(t, "mongodb://localhost:27017", viper.GetString(FlagLongDatabaseURI))
	assert.Equal(t, "http://localhost:8081", viper.GetString(FlagLongAccrualSystemAddress))
}
