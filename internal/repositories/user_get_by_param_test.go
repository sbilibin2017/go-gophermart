package repositories

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

func TestGetByParam(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpassword",
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
	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)
	dsn := "postgres://testuser:testpassword@%s:%s/testdb?sslmode=disable"
	dsn = fmt.Sprintf(dsn, host, port.Port())
	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			login TEXT PRIMARY KEY,
			password TEXT NOT NULL
		);
	`)
	require.NoError(t, err)
	_, err = db.Exec(`
		INSERT INTO users (login, password) VALUES ('testuser', 'testpassword');
	`)
	require.NoError(t, err)
	userRepo := NewUserGetByParamRepository(db)
	saveCtx := context.Background()
	param := map[string]any{
		"login": "testuser",
	}
	user, err := userRepo.GetByParam(saveCtx, param)
	require.NoError(t, err)
	assert.NotNil(t, user, "user should not be nil")
	assert.Equal(t, "testuser", user["login"], "login should be 'testuser'")
	assert.Equal(t, "testpassword", user["password"], "password should be 'testpassword'")

	param = map[string]any{
		"login": "nonexistentuser",
	}
	user, err = userRepo.GetByParam(saveCtx, param)
	require.NoError(t, err)
	assert.Nil(t, user, "user should be nil for nonexistent login")
}
