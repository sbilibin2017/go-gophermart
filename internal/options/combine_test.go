package options

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombine(t *testing.T) {
	envKey := "TEST_ENV"
	flagKey := "test_flag"
	defaultValue := "default"

	defer os.Unsetenv(envKey)

	t.Run("returns value from env if present", func(t *testing.T) {
		os.Setenv(envKey, "from_env")
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		flag.String(flagKey, "from_flag", "usage")
		flag.CommandLine.Parse([]string{})
		result := Combine(envKey, flagKey, defaultValue)
		assert.Equal(t, "from_env", result)
	})

	t.Run("returns value from flag if env not present", func(t *testing.T) {
		os.Unsetenv(envKey)
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		flag.String(flagKey, "", "usage")
		flag.CommandLine.Parse([]string{"-" + flagKey + "=from_flag"})
		result := Combine(envKey, flagKey, defaultValue)
		assert.Equal(t, "from_flag", result)
	})

	t.Run("returns default if neither env nor flag present", func(t *testing.T) {
		os.Unsetenv(envKey)
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		result := Combine(envKey, flagKey, defaultValue)
		assert.Equal(t, defaultValue, result)
	})
}
