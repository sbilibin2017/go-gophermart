package app

import (
	"fmt"
	"testing"

	"github.com/sbilibin2017/go-gophermart/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	db, host, port, cleanup := storage.SetupPostgresContainer(t)
	defer cleanup()
	defer db.Close()

	dsn := fmt.Sprintf("postgres://testuser:testpassword@%s:%s/testdb?sslmode=disable", host, port)

	cfg := &Config{
		RunAddress:  "localhost:8082",
		DatabaseURI: dsn,
	}

	server, err := NewServer(cfg)

	require.NoError(t, err)
	require.NotNil(t, server)
	require.Equal(t, cfg.RunAddress, server.Addr)
}

func TestNewServer_DBConnectionError(t *testing.T) {
	cfg := &Config{
		RunAddress:  "localhost:9090",
		DatabaseURI: "postgres://invaliduser:invalidpass@localhost:9999/invalid_db",
	}

	server, err := NewServer(cfg)

	require.Error(t, err)
	require.Nil(t, server)
}
