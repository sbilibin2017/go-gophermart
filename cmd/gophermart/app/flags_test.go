package app_test

import (
	"flag"
	"os"
	"testing"

	"github.com/sbilibin2017/go-gophermart/cmd/gophermart/app"
	"github.com/stretchr/testify/assert"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func TestParseFlags_Defaults(t *testing.T) {
	resetFlags()
	os.Clearenv()
	os.Args = []string{"cmd.test"}
	cfg := app.ParseFlags()
	assert.Equal(t, app.DefaultRunAddress, cfg.RunAddress)
	assert.Equal(t, app.DefaultDatabaseURI, cfg.DatabaseURI)
	assert.Equal(t, app.DefaultAccrualSystemAddress, cfg.AccrualSystemAddress)
}

func TestParseFlags_EnvVars(t *testing.T) {
	resetFlags()
	os.Clearenv()
	os.Setenv(app.EnvRunAddress, "localhost:9000")
	os.Setenv(app.EnvDatabaseURI, "postgres://user:pass@localhost/db")
	os.Setenv(app.EnvAccrualSystemAddress, "http://accrual-system")
	defer func() {
		os.Unsetenv(app.EnvRunAddress)
		os.Unsetenv(app.EnvDatabaseURI)
		os.Unsetenv(app.EnvAccrualSystemAddress)
	}()
	os.Args = []string{"cmd.test"}
	cfg := app.ParseFlags()
	assert.Equal(t, "localhost:9000", cfg.RunAddress)
	assert.Equal(t, "postgres://user:pass@localhost/db", cfg.DatabaseURI)
	assert.Equal(t, "http://accrual-system", cfg.AccrualSystemAddress)
}

func TestParseFlags_CLIFlags(t *testing.T) {
	resetFlags()
	os.Clearenv()
	os.Args = []string{
		"cmd.test",
		"-" + app.FlagRunAddress, "127.0.0.1:9999",
		"-" + app.FlagDatabaseURI, "postgres://cli:cli@localhost/cli",
		"-" + app.FlagAccrualSystemAddress, "http://cli-accrual",
	}
	cfg := app.ParseFlags()
	assert.Equal(t, "127.0.0.1:9999", cfg.RunAddress)
	assert.Equal(t, "postgres://cli:cli@localhost/cli", cfg.DatabaseURI)
	assert.Equal(t, "http://cli-accrual", cfg.AccrualSystemAddress)
}

func TestParseFlags_EnvOverridesCLI(t *testing.T) {
	resetFlags()
	os.Clearenv()
	os.Setenv(app.EnvRunAddress, "env:9090")
	os.Setenv(app.EnvDatabaseURI, "postgres://env/env")
	os.Setenv(app.EnvAccrualSystemAddress, "http://env-accrual")
	defer func() {
		os.Unsetenv(app.EnvRunAddress)
		os.Unsetenv(app.EnvDatabaseURI)
		os.Unsetenv(app.EnvAccrualSystemAddress)
	}()
	os.Args = []string{
		"cmd.test",
		"-" + app.FlagRunAddress, "cli:1111",
		"-" + app.FlagDatabaseURI, "postgres://cli/cli",
		"-" + app.FlagAccrualSystemAddress, "http://cli-accrual",
	}
	cfg := app.ParseFlags()
	assert.Equal(t, "env:9090", cfg.RunAddress)
	assert.Equal(t, "postgres://env/env", cfg.DatabaseURI)
	assert.Equal(t, "http://env-accrual", cfg.AccrualSystemAddress)
}
