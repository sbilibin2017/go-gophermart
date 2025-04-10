package unitofwork

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestPostgresUnitOfWork(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	defer container.Terminate(ctx)

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	fmt.Printf("Host: %s, Port: %s\n", host, port)

	dbURL := fmt.Sprintf("postgres://postgres:password@%s:%s/testdb?sslmode=disable", host, port.Port())

	db, err := sql.Open("pgx", dbURL)
	require.NoError(t, err)

	err = db.PingContext(ctx)
	require.NoError(t, err, "Failed to connect to the database")

	defer db.Close()

	uow := NewDBUnitOfWork(db)

	operation := func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT)")
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, "INSERT INTO users (name) VALUES ($1)", "Test User")
		return err
	}

	err = uow.Do(ctx, operation)
	require.NoError(t, err)

	var name string
	err = db.QueryRowContext(ctx, "SELECT name FROM users WHERE name = $1", "Test User").Scan(&name)
	require.NoError(t, err)
	assert.Equal(t, "Test User", name)
}

func TestPostgresUnitOfWorkRollback(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	defer container.Terminate(ctx)

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	fmt.Printf("Host: %s, Port: %s\n", host, port)

	dbURL := fmt.Sprintf("postgres://postgres:password@%s:%s/testdb?sslmode=disable", host, port.Port())

	db, err := sql.Open("pgx", dbURL)
	require.NoError(t, err)

	err = db.PingContext(ctx)
	require.NoError(t, err, "Failed to connect to the database")

	defer db.Close()

	_, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT)")
	require.NoError(t, err, "Failed to create table")

	uow := NewDBUnitOfWork(db)

	operation := func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, "INSERT INTO users (name) VALUES ($1)", "Test User")
		if err != nil {
			return err
		}

		return fmt.Errorf("simulated error")
	}

	err = uow.Do(ctx, operation)
	require.Error(t, err, "Expected an error to trigger rollback")

	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE name = $1", "Test User").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count, "Expected no records to be inserted as transaction was rolled back")
}
