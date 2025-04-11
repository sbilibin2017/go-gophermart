package app

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func TestEnvironmentVariablesHaveHighestPriority(t *testing.T) {
	resetFlags()

	os.Setenv(EnvRunAddress, ":9090")
	os.Setenv(EnvDatabaseURI, "postgres://env/db")
	os.Setenv(EnvAccrualSystemAddress, "http://env-system")

	os.Args = []string{
		"gophermart",
		"-a", ":8080",
		"-d", "postgres://flags/db",
		"-r", "http://flags-system",
	}

	config := ParseFlags()

	assert.Equal(t, ":9090", config.RunAddress)
	assert.Equal(t, "postgres://env/db", config.DatabaseURI)
	assert.Equal(t, "http://env-system", config.AccrualSystemAddress)
}

func TestFlagsOverrideDefaults(t *testing.T) {
	resetFlags()

	os.Unsetenv(EnvRunAddress)
	os.Unsetenv(EnvDatabaseURI)
	os.Unsetenv(EnvAccrualSystemAddress)

	os.Args = []string{
		"gophermart",
		"-a", ":9090",
		"-d", "postgres://flags/db",
		"-r", "http://flags-system",
	}

	config := ParseFlags()

	assert.Equal(t, ":9090", config.RunAddress)
	assert.Equal(t, "postgres://flags/db", config.DatabaseURI)
	assert.Equal(t, "http://flags-system", config.AccrualSystemAddress)
}

func TestDefaultsUsedWhenNoEnvVarsOrFlags(t *testing.T) {
	resetFlags()

	os.Unsetenv(EnvRunAddress)
	os.Unsetenv(EnvDatabaseURI)
	os.Unsetenv(EnvAccrualSystemAddress)
	os.Args = []string{"gophermart"}

	config := ParseFlags()

	assert.Equal(t, ":8080", config.RunAddress)
	assert.Equal(t, "", config.DatabaseURI)
	assert.Equal(t, "", config.AccrualSystemAddress)
}

func TestEnvsOverrideFlagsVariables(t *testing.T) {
	resetFlags()

	os.Setenv(EnvRunAddress, ":7070")
	os.Setenv(EnvDatabaseURI, "postgres://env/db")
	os.Setenv(EnvAccrualSystemAddress, "http://env-system")

	os.Args = []string{
		"gophermart",
		"-a", ":9090",
		"-d", "postgres://flags/db",
		"-r", "http://flags-system",
	}

	config := ParseFlags()

	assert.Equal(t, ":7070", config.RunAddress)
	assert.Equal(t, "postgres://env/db", config.DatabaseURI)
	assert.Equal(t, "http://env-system", config.AccrualSystemAddress)
}
