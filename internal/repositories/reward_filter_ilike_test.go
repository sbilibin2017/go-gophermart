package repositories

import (
	"context"
	"fmt"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap/zapcore"
)

type RewardDB struct {
	Match      string    `db:"match"`
	Reward     int64     `db:"reward"`
	RewardType string    `db:"reward_type"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func setupPostgreSQLContainer(t *testing.T) (*sqlx.DB, func()) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpassword",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("listening on IPv4 address"),
	}
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	mappedPort, err := postgresContainer.MappedPort(ctx, "5432")
	assert.NoError(t, err)
	host, err := postgresContainer.Host(ctx)
	assert.NoError(t, err)
	dsn := fmt.Sprintf("postgres://testuser:testpassword@%s:%s/testdb?sslmode=disable", host, mappedPort.Port())
	var db *sqlx.DB
	for i := 0; i < 10; i++ {
		db, err = sqlx.Connect("pgx", dsn)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	assert.NoError(t, err)
	assert.NotNil(t, db)
	_, err = db.Exec(`
        CREATE TABLE rewards (
            id SERIAL PRIMARY KEY,
            match TEXT,
            reward INT,
            reward_type TEXT,
            created_at TIMESTAMP,
            updated_at TIMESTAMP
        );
    `)
	assert.NoError(t, err)
	_, err = db.Exec(`
        INSERT INTO rewards (match, reward, reward_type, created_at, updated_at)
        VALUES
        ('this is a reward', 100, 'type1', '2025-04-24 10:00:00', '2025-04-24 10:00:00'),
        ('another reward', 200, 'type2', '2025-04-24 10:05:00', '2025-04-24 10:05:00'),
        ('special reward', 300, 'type3', '2025-04-24 10:10:00', '2025-04-24 10:10:00');
    `)
	assert.NoError(t, err)
	return db, func() {
		postgresContainer.Terminate(ctx)
	}
}

func TestRewardFilterILikeRepository_FilterILike(t *testing.T) {
	db, teardown := setupPostgreSQLContainer(t)
	defer teardown()
	logger.Init(zapcore.InfoLevel)
	repo := NewRewardFilterILikeRepository(db, nil)

	t.Run("Matching description found", func(t *testing.T) {
		result, err := repo.FilterILike(context.Background(), "reward", []string{"id", "match", "reward", "reward_type"})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Contains(t, result, "match")
		assert.Equal(t, "this is a reward", result["match"])
		assert.Equal(t, int64(100), result["reward"])
		assert.Equal(t, "type1", result["reward_type"])
	})

	t.Run("No matching description", func(t *testing.T) {
		result, err := repo.FilterILike(context.Background(), "nonexistent", []string{"id", "match", "reward", "reward_type"})
		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("Empty description filter", func(t *testing.T) {
		result, err := repo.FilterILike(context.Background(), "", []string{"id", "match", "reward", "reward_type"})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 4)
		assert.Contains(t, result, "match")
		assert.Equal(t, "this is a reward", result["match"])
	})

	t.Run("No fields to select", func(t *testing.T) {
		result, err := repo.FilterILike(context.Background(), "reward", []string{})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Contains(t, result, "id")
		assert.Contains(t, result, "match")
		assert.Contains(t, result, "reward")
		assert.Contains(t, result, "reward_type")
		assert.Contains(t, result, "created_at")
		assert.Contains(t, result, "updated_at")
	})
}
