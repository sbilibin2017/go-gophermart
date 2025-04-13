package storage

import (
	"context"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupPostgresContainer(t *testing.T) (*sqlx.DB, string, string, func()) {
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

	// Запуск контейнера PostgreSQL
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	// Получаем хост и порт контейнера
	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	// Проверяем, что хост и порт правильно получены
	require.NotEmpty(t, host, "Host should not be empty")
	require.NotEmpty(t, port.Port(), "Port should not be empty")

	// Строим DSN для подключения
	dsn := "postgres://testuser:testpassword@" + host + ":" + port.Port() + "/testdb?sslmode=disable"
	db, err := sqlx.Connect("pgx", dsn) // Здесь используется правильный DSN и драйвер
	require.NoError(t, err)

	// Пинг базы данных для проверки соединения
	err = db.PingContext(ctx)
	require.NoError(t, err)

	// Функция для очистки контейнера
	cleanup := func() {
		err := container.Terminate(ctx)
		require.NoError(t, err)
	}

	// Возвращаем подключение, хост и порт для использования в тестах
	return db, host, port.Port(), cleanup
}
