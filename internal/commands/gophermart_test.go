package commands

import (
	"context"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func dummyGophermartRunner(ctx context.Context) error {
	return nil
}

func dummyGophermartCtxer() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	return ctx, cancel
}

func TestNewGophermartCommand(t *testing.T) {
	cmd := NewGophermartCommand(dummyGophermartCtxer, dummyGophermartRunner)
	assert.NotNil(t, cmd)
	assert.Equal(t, "Gophermart", cmd.Use)
	assert.Equal(t, "Gophermart Service", cmd.Short)
	runAddressFlag := cmd.PersistentFlags().Lookup(FlagGophermartRunAddress)
	assert.NotNil(t, runAddressFlag)
	assert.Equal(t, "localhost:8080", runAddressFlag.Value.String())
	databaseURIFlag := cmd.PersistentFlags().Lookup(FlagGophermartDatabaseURI)
	assert.NotNil(t, databaseURIFlag)
	assert.Equal(t, "", databaseURIFlag.Value.String())
	accrualSystemFlag := cmd.PersistentFlags().Lookup(FlagGophermartAccrualSystemAddress)
	assert.NotNil(t, accrualSystemFlag)
	assert.Equal(t, "", accrualSystemFlag.Value.String())
	os.Setenv(EnvGophermartRunAddress, "env-address")
	os.Setenv(EnvGophermartDatabaseURI, "env-database-uri")
	os.Setenv(EnvGophermartAccrualSystemAddress, "env-accrual-system-address")
	viper.AutomaticEnv()
	assert.Equal(t, "env-address", viper.GetString(EnvGophermartRunAddress))
	assert.Equal(t, "env-database-uri", viper.GetString(EnvGophermartDatabaseURI))
	assert.Equal(t, "env-accrual-system-address", viper.GetString(EnvGophermartAccrualSystemAddress))
	cmd.PersistentFlags().Parse([]string{
		"--run-address", "cmd-address",
		"--database-uri", "cmd-database-uri",
		"--accrual-system-address", "cmd-accrual-system-address",
	})
	assert.Equal(t, "cmd-address", runAddressFlag.Value.String())
	assert.Equal(t, "cmd-database-uri", databaseURIFlag.Value.String())
	assert.Equal(t, "cmd-accrual-system-address", accrualSystemFlag.Value.String())
	err := cmd.RunE(cmd, []string{})
	assert.NoError(t, err)
}
