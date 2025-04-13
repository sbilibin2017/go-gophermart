package app

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCommand_Defaults(t *testing.T) {
	viper.Reset()
	cmd := NewCommand()
	flags := cmd.Flags()
	addrFlag := flags.Lookup(FlagRunAddress)
	dbFlag := flags.Lookup(FlagDatabaseURI)
	assert.NotNil(t, addrFlag)
	assert.NotNil(t, dbFlag)
	assert.Equal(t, DefaultRunAddress, addrFlag.DefValue)
	assert.Equal(t, DefaultDatabaseURI, dbFlag.DefValue)
}

func TestNewCommand_RunE_WithFlags(t *testing.T) {
	viper.Reset()
	viper.Set(FlagRunAddress, "127.0.0.1:9999")
	viper.Set(FlagDatabaseURI, "postgres://test:test@localhost:5432/db")
	cmd := NewCommand()
	called := false
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		called = true
		addr := viper.GetString(FlagRunAddress)
		db := viper.GetString(FlagDatabaseURI)
		assert.Equal(t, "127.0.0.1:9999", addr)
		assert.Equal(t, "postgres://test:test@localhost:5432/db", db)
		return nil
	}
	err := cmd.Execute()
	assert.NoError(t, err)
	assert.True(t, called)
}

func TestRun_NewServerError(t *testing.T) {
	cfg := &Config{
		RunAddress:  "localhost:9999",
		DatabaseURI: "invalid_dsn",
	}
	err := run(context.Background(), cfg)
	require.Error(t, err)
}

func TestNewCommand_RunE_Override(t *testing.T) {
	originalRunE := runE
	defer func() {
		runE = originalRunE
	}()
	viper.Reset()
	viper.Set(FlagRunAddress, "127.0.0.1:8081")
	viper.Set(FlagDatabaseURI, "postgres://user:pass@localhost:5432/testdb")
	called := false
	runE = func(cmd *cobra.Command, args []string) error {
		called = true
		cfg := &Config{
			RunAddress:  viper.GetString(FlagRunAddress),
			DatabaseURI: viper.GetString(FlagDatabaseURI),
		}
		assert.Equal(t, "127.0.0.1:8081", cfg.RunAddress)
		assert.Equal(t, "postgres://user:pass@localhost:5432/testdb", cfg.DatabaseURI)
		return nil
	}
	cmd := NewCommand()
	err := cmd.Execute()
	require.NoError(t, err)
	assert.True(t, called)
}

func TestRunE_OriginalImplementation(t *testing.T) {
	viper.Reset()
	viper.Set(FlagRunAddress, "localhost:8082")
	viper.Set(FlagDatabaseURI, "postgres://user:pass@localhost:5432/testdb")
	originalRun := run
	originalServerRun := server.Run
	defer func() {
		run = originalRun
		runServer = originalServerRun
	}()
	runServer = func(ctx context.Context, srv *http.Server) error {
		assert.Equal(t, "localhost:8082", srv.Addr)
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	}
	run = func(ctx context.Context, config *Config) error {
		assert.Equal(t, "localhost:8082", config.RunAddress)
		assert.Equal(t, "postgres://user:pass@localhost:5432/testdb", config.DatabaseURI)
		mockServer := &http.Server{Addr: config.RunAddress}
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		return server.Run(ctx, mockServer)
	}
	cmd := NewCommand()
	err := cmd.Execute()
	require.NoError(t, err)
}

func TestRun_CallsRunServer(t *testing.T) {
	originalRunServer := runServer
	originalNewServer := newServer
	defer func() {
		runServer = originalRunServer
		newServer = originalNewServer
	}()
	cfg := &Config{
		RunAddress:  "localhost:9999",
		DatabaseURI: "postgres://test:test@localhost:5432/db",
	}
	called := false
	runServer = func(ctx context.Context, srv *http.Server) error {
		called = true
		assert.Equal(t, "localhost:9999", srv.Addr)
		return nil
	}
	newServer = func(cfg *Config) (*http.Server, error) {
		return &http.Server{Addr: cfg.RunAddress}, nil
	}
	err := run(context.Background(), cfg)
	require.NoError(t, err)
	assert.True(t, called)
}
