package unitofwork

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func setupPostgresContainer(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpassword",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432").WithPollInterval(2 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dsn := "postgres://testuser:testpassword@" + host + ":" + port.Port() + "/testdb?sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)

	err = db.PingContext(ctx)
	require.NoError(t, err)

	cleanup := func() {
		err := container.Terminate(ctx)
		require.NoError(t, err)
	}

	return db, cleanup
}

func TestUnitOfWork_Do_Success(t *testing.T) {
	db, cleanup := setupPostgresContainer(t)
	defer cleanup()

	uow := NewUnitOfWork(db)

	_, err := db.Exec("CREATE TABLE users (id SERIAL PRIMARY KEY, username VARCHAR(255) NOT NULL);")
	require.NoError(t, err)

	operation := func(tx *sql.Tx) error {
		_, err := tx.Exec("INSERT INTO users (username) VALUES ('testuser')")
		return err
	}

	err = uow.Do(context.Background(), operation)
	require.NoError(t, err)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'testuser'").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestUnitOfWork_Do_NilDB(t *testing.T) {
	uow := NewUnitOfWork(nil)

	operation := func(tx *sql.Tx) error {
		return nil
	}

	err := uow.Do(context.Background(), operation)
	assert.NoError(t, err)
}

func TestUnitOfWork_Do_OperationError_Rollback(t *testing.T) {
	db, cleanup := setupPostgresContainer(t)
	defer cleanup()

	_, err := db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, username TEXT NOT NULL);")
	require.NoError(t, err)

	uow := &UnitOfWork{
		db: db,
	}

	operation := func(tx *sql.Tx) error {
		_, err := tx.Exec("INSERT INTO users (username) VALUES (?)", "testuser")
		if err != nil {
			return err
		}
		return errors.New("operation failed")
	}

	err = uow.Do(context.Background(), operation)
	require.Error(t, err)

}
