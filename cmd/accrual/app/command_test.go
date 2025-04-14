package app

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCommand_DefaultFlags(t *testing.T) {
	cmd := NewCommand()
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	require.NoError(t, err, "command should execute without error")
	assert.Equal(t, flagDefaultRunAddress, viper.GetString(flagRunAddress), "run address should match default")
	assert.Equal(t, flagDefaultDatabaseURI, viper.GetString(flagDatabaseURI), "database URI should match default")
	var cfg Config
	err = viper.Unmarshal(&cfg)
	require.NoError(t, err, "should unmarshal config without error")
	assert.Equal(t, flagDefaultRunAddress, cfg.GetRunAddress(), "Config.GetRunAddress should match default")
	assert.Equal(t, flagDefaultDatabaseURI, cfg.GetDatabaseURI(), "Config.GetDatabaseURI should match default")
}
