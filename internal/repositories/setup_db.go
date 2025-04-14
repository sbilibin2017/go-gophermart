package repositories

import (
	"context"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupDB(t *testing.T, setupQuery string) (*sqlx.DB, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_USER":     "user",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = postgresC.Terminate(ctx)
	})

	host, err := postgresC.Host(ctx)
	require.NoError(t, err)

	port, err := postgresC.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dsn := fmt.Sprintf("postgres://user:password@%s:%s/testdb?sslmode=disable", host, port.Port())
	db, err := sqlx.Connect("pgx", dsn)
	require.NoError(t, err)

	_, err = db.Exec(setupQuery)
	require.NoError(t, err)

	return db, func() {
		db.Close()
	}
}
