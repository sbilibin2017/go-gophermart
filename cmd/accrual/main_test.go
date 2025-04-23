package main

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlags(t *testing.T) {
	resetFlags := func() {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	}

	t.Run("Flags are parsed correctly", func(t *testing.T) {
		resetFlags()
		os.Args = []string{"cmd", "-a", "localhost:8080", "-d", "postgres://user:pass@localhost/db"}
		config := flags()
		assert.Equal(t, "localhost:8080", config.RunAddress)
		assert.Equal(t, "postgres://user:pass@localhost/db", config.DatabaseURI)
	})

	t.Run("Environment variables override flags", func(t *testing.T) {
		resetFlags()
		os.Setenv("RUN_ADDRESS", "env-address:8080")
		os.Setenv("DATABASE_URI", "env-postgres://user:pass@localhost/db")
		os.Args = []string{"cmd", "-a", "localhost:8080", "-d", "postgres://user:pass@localhost/db"}
		config := flags()
		assert.Equal(t, "env-address:8080", config.RunAddress)
		assert.Equal(t, "env-postgres://user:pass@localhost/db", config.DatabaseURI)
		os.Unsetenv("RUN_ADDRESS")
		os.Unsetenv("DATABASE_URI")
	})

	t.Run("Empty flags and no environment variables", func(t *testing.T) {
		resetFlags()
		os.Args = []string{"cmd"}
		config := flags()
		assert.Empty(t, config.RunAddress)
		assert.Empty(t, config.DatabaseURI)
	})
}
