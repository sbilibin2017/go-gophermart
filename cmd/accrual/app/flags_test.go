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

func TestParseFlags_Defaults(t *testing.T) {
	// Очистка переменных окружения и флагов
	os.Clearenv()
	resetFlags()

	cfg := ParseFlags()

	assert.Equal(t, "localhost:8080", cfg.RunAddress)
	assert.Equal(t, "", cfg.DatabaseURI)
}

func TestParseFlags_EnvVariables(t *testing.T) {
	os.Setenv("RUN_ADDRESS", "127.0.0.1:3000")
	os.Setenv("DATABASE_URI", "postgres://env")

	resetFlags()
	cfg := ParseFlags()

	assert.Equal(t, "127.0.0.1:3000", cfg.RunAddress)
	assert.Equal(t, "postgres://env", cfg.DatabaseURI)

	os.Unsetenv("RUN_ADDRESS")
	os.Unsetenv("DATABASE_URI")
}

func TestParseFlags_CommandLineFlags(t *testing.T) {
	// Симулируем передачу аргументов командной строки
	os.Args = []string{
		"cmd",
		"-a=192.168.0.1:9000",
		"-d=postgres://flag",
	}
	resetFlags()
	cfg := ParseFlags()

	assert.Equal(t, "192.168.0.1:9000", cfg.RunAddress)
	assert.Equal(t, "postgres://flag", cfg.DatabaseURI)
}

func TestParseFlags_FlagsOverrideEnv(t *testing.T) {
	os.Setenv("RUN_ADDRESS", "127.0.0.1:3000")
	os.Setenv("DATABASE_URI", "postgres://env")

	os.Args = []string{
		"cmd",
		"-a=override:1234",
		"-d=postgres://flag",
	}
	resetFlags()
	cfg := ParseFlags()

	// ENV wins over default, but FLAGS override both
	assert.Equal(t, "127.0.0.1:3000", cfg.RunAddress)
	assert.Equal(t, "postgres://env", cfg.DatabaseURI)

	os.Unsetenv("RUN_ADDRESS")
	os.Unsetenv("DATABASE_URI")
}
