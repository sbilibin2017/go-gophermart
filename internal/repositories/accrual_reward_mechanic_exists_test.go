package repositories

import (
	"context"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Функция для настройки контейнера PostgreSQL
func setupPostgresContainer2(t *testing.T) *sqlx.DB {
	// Настройка контейнера PostgreSQL с использованием testcontainers
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13", // Используем образ PostgreSQL
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"), // Ждем, пока порт 5432 откроется для подключения
	}

	// Создаем контейнер
	postgresContainer, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	// Получаем порты, на которых слушает контейнер
	host, err := postgresContainer.Host(context.Background())
	require.NoError(t, err)

	port, err := postgresContainer.MappedPort(context.Background(), "5432")
	require.NoError(t, err)

	// Создаем строку подключения
	dsn := "postgres://testuser:testpass@" + host + ":" + port.Port() + "/testdb?sslmode=disable"

	// Подключаемся к базе данных
	db, err := sqlx.Connect("pgx", dsn)
	require.NoError(t, err)

	// Создаем таблицу для тестов
	_, err = db.Exec(`
		CREATE TABLE accrual_reward_mechanic (
			match TEXT PRIMARY KEY,
			reward BIGINT NOT NULL,
			reward_type TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	require.NoError(t, err)

	return db
}

// Тестирование репозитория Exists
func TestAccrualRewardMechanicExistsRepository_Exists(t *testing.T) {
	// Подключаемся к тестовой базе данных
	db := setupPostgresContainer2(t)
	defer db.Close()

	logger.InitWithInfoLevel()

	// Создаем репозиторий
	repo := NewAccrualRewardMechanicExistsRepository(db)

	// Тестируем существование записи
	t.Run("ExistsRecord", func(t *testing.T) {
		// Добавляем запись в базу данных
		_, err := db.Exec("INSERT INTO accrual_reward_mechanic (match, reward, reward_type, created_at, updated_at) VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
			"test_match", 100, "typeA")
		assert.NoError(t, err)

		// Проверяем существование записи
		exists, err := repo.Exists(context.Background(), "test_match")
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	// Тестируем, что записи нет
	t.Run("NotExistsRecord", func(t *testing.T) {
		// Проверяем, что записи не существует
		exists, err := repo.Exists(context.Background(), "nonexistent_match")
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}
