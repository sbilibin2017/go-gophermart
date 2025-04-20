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

func TestNoFlagsOrEnvVarsSet(t *testing.T) {
	resetFlags()

	os.Unsetenv("RUN_ADDRESS")
	os.Unsetenv("DATABASE_URI")
	runAddress = ""
	databaseURI = ""

	flags()

	assert.Equal(t, "", runAddress, "runAddress should be empty")
	assert.Equal(t, "", databaseURI, "databaseURI should be empty")
}

func TestFlagsSet(t *testing.T) {
	resetFlags()

	os.Args = []string{"cmd", "-a", "localhost", "-d", "localhost/db"}

	runAddress = ""
	databaseURI = ""
	flags()

	assert.Equal(t, "localhost", runAddress, "runAddress should be 'localhost'")
	assert.Equal(t, "localhost/db", databaseURI, "databaseURI should be 'localhost/db'")
}

func TestEnvVarsSet(t *testing.T) {
	resetFlags()

	os.Setenv("RUN_ADDRESS", "env-address")
	os.Setenv("DATABASE_URI", "env-db-uri")

	os.Args = []string{"cmd"}
	runAddress = ""
	databaseURI = ""
	flags()

	assert.Equal(t, "env-address", runAddress, "runAddress should be 'env-address' from environment variable")
	assert.Equal(t, "env-db-uri", databaseURI, "databaseURI should be 'env-db-uri' from environment variable")
}

func TestFlagsOverrideEnvVars(t *testing.T) {
	resetFlags()

	os.Setenv("RUN_ADDRESS", "env-address")
	os.Setenv("DATABASE_URI", "env-db-uri")
	os.Args = []string{"cmd", "-a", "flag-address", "-d", "flag-db-uri"}

	runAddress = ""
	databaseURI = ""
	flags()

	assert.Equal(t, "env-address", runAddress, "runAddress should be 'env-address' from environment variable")
	assert.Equal(t, "env-db-uri", databaseURI, "databaseURI should be 'env-db-uri' from environment variable")
}
