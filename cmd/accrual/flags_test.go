package main

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
}

func TestFlagsWithCommandLineArgs(t *testing.T) {
	resetFlags()
	os.Args = []string{"cmd", "-a", "localhost:8080", "-d", "postgres://user:pass@localhost:5432/db"}
	config := flags()
	assert.Equal(t, "localhost:8080", config.RunAddress)
	assert.Equal(t, "postgres://user:pass@localhost:5432/db", config.DatabaseURI)
}

func TestFlagsWithEnvironmentVariables(t *testing.T) {
	resetFlags()
	os.Setenv("RUN_ADDRESS", "envAddress:8080")
	os.Setenv("DATABASE_URI", "envDbUri")
	os.Args = []string{"cmd"}
	config := flags()
	assert.Equal(t, "envAddress:8080", config.RunAddress)
	assert.Equal(t, "envDbUri", config.DatabaseURI)
	os.Unsetenv("RUN_ADDRESS")
	os.Unsetenv("DATABASE_URI")
}

func TestFlagsWithBothEnvironmentVariablesAndFlags(t *testing.T) {
	resetFlags()
	os.Setenv("RUN_ADDRESS", "envAddressOverride:8080")
	os.Setenv("DATABASE_URI", "envDbUriOverride")
	os.Args = []string{"cmd", "-a", "flagAddress:8080", "-d", "flagDbUri"}
	config := flags()
	assert.Equal(t, "envAddressOverride:8080", config.RunAddress)
	assert.Equal(t, "envDbUriOverride", config.DatabaseURI)
	os.Unsetenv("RUN_ADDRESS")
	os.Unsetenv("DATABASE_URI")
}

func TestFlagsWithNoEnvironmentVariables(t *testing.T) {
	resetFlags()
	os.Clearenv()
	os.Args = []string{"cmd", "-a", "localhost:8080", "-d", "postgres://user:pass@localhost:5432/db"}
	config := flags()
	assert.Equal(t, "localhost:8080", config.RunAddress)
	assert.Equal(t, "postgres://user:pass@localhost:5432/db", config.DatabaseURI)
}

func TestFlagsWithNoFlagsOrEnvironmentVariables(t *testing.T) {
	resetFlags()
	os.Clearenv()
	os.Args = []string{"cmd"}
	config := flags()
	assert.Equal(t, "", config.RunAddress)
	assert.Equal(t, "", config.DatabaseURI)
}
