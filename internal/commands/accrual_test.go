package commands

import (
	"context"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func dummyAccrualRunner(ctx context.Context) error {
	return nil
}

func dummyAccrualCtxer() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	return ctx, cancel
}

func TestNewAccrualCommand(t *testing.T) {
	cmd := NewAccrualCommand(dummyAccrualCtxer, dummyAccrualRunner)
	assert.NotNil(t, cmd)
	assert.Equal(t, "Accrual", cmd.Use)
	assert.Equal(t, "Accrual Service", cmd.Short)
	runAddressFlag := cmd.PersistentFlags().Lookup(FlagAccrualRunAddress)
	assert.NotNil(t, runAddressFlag)
	assert.Equal(t, "localhost:8081", runAddressFlag.Value.String())
	databaseURIFlag := cmd.PersistentFlags().Lookup(FlagAccrualDatabaseURI)
	assert.NotNil(t, databaseURIFlag)
	assert.Equal(t, "", databaseURIFlag.Value.String())
	os.Setenv(EnvAccrualRunAddress, "env-address")
	os.Setenv(EnvAccrualDatabaseURI, "env-database-uri")
	viper.AutomaticEnv()
	assert.Equal(t, "env-address", viper.GetString(EnvAccrualRunAddress))
	assert.Equal(t, "env-database-uri", viper.GetString(EnvAccrualDatabaseURI))
	cmd.PersistentFlags().Parse([]string{
		"--run-address", "cmd-address",
		"--database-uri", "cmd-database-uri",
	})
	assert.Equal(t, "cmd-address", runAddressFlag.Value.String())
	assert.Equal(t, "cmd-database-uri", databaseURIFlag.Value.String())
	err := cmd.RunE(cmd, []string{})
	assert.NoError(t, err)
}
