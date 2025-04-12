package fixtures

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Fixture function to set up a PostgreSQL database
func SetupPostgresDB(t *testing.T) (*sqlx.DB, func(ctx context.Context, query string) error, func()) {
	// Start a PostgreSQL container using testcontainers
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15", // You can use any version of PostgreSQL
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_USER":     "user",
			"POSTGRES_DB":       "testdb",
		},
		// Use a port-based wait strategy from the wait package
		WaitingFor: wait.ForListeningPort("5432/tcp"), // Wait for PostgreSQL to be ready
	}

	// Create the container using testcontainers.NewGenericContainer
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	// Get the container's mapped port
	mappedPort, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	// Get the host machine's IP address
	host, err := container.Host(ctx)
	require.NoError(t, err)

	// Wait for PostgreSQL to be ready by checking if the port is open
	require.Eventually(t, func() bool {
		_, err := sqlx.Open("pgx", fmt.Sprintf("postgres://user:password@%s:%s/testdb?sslmode=disable", host, mappedPort.Port()))
		return err == nil
	}, 30*time.Second, 1*time.Second, "PostgreSQL didn't become ready in time")

	// Connect to the database
	dsn := fmt.Sprintf("postgres://user:password@%s:%s/testdb?sslmode=disable", host, mappedPort.Port())
	db, err := sqlx.Open("pgx", dsn)
	require.NoError(t, err)

	// Return the DB connection, a query executor function, and a cleanup function
	return db, func(ctx context.Context, query string) error {
			_, err := db.ExecContext(ctx, query)
			return err
		}, func() {
			_ = db.Close()
			_ = container.Terminate(ctx)
		}
}
