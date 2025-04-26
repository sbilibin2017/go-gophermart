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
func setupPostgresContainer(t *testing.T) *sqlx.DB {
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

// Тестирование репозитория Save
func TestAccrualRewardMechanicSaveRepository_Save(t *testing.T) {
	// Подключаемся к тестовой базе данных
	db := setupPostgresContainer(t)
	defer db.Close()

	logger.InitWithInfoLevel()

	// Создаем репозиторий
	repo := NewAccrualRewardMechanicSaveRepository(db)

	// Тестируем сохранение новой записи
	t.Run("SaveNewRecord", func(t *testing.T) {
		err := repo.Save(context.Background(), "test_match", 100, "typeA")
		assert.NoError(t, err)

		// Проверяем, что запись была добавлена в таблицу
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM accrual_reward_mechanic WHERE match = $1", "test_match").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	// Тестируем обновление существующей записи
	t.Run("UpdateExistingRecord", func(t *testing.T) {

		// Обновляем запись
		repo.Save(context.Background(), "test_match", 200, "typeA")
		err := repo.Save(context.Background(), "test_match", 300, "typeB")
		assert.NoError(t, err)

		// Проверяем, что данные были обновлены
		var reward int64
		var rewardType string
		var updatedAt string
		err = db.QueryRow("SELECT reward, reward_type, updated_at FROM accrual_reward_mechanic WHERE match = $1", "test_match").Scan(&reward, &rewardType, &updatedAt)
		assert.NoError(t, err)
		assert.Equal(t, int64(300), reward)
		assert.Equal(t, "typeB", rewardType)
		assert.NotEmpty(t, updatedAt) // Ensure updated_at is set
	})
}
