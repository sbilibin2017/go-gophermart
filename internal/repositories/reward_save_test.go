package repositories

import (
	"context"
	"testing"
	"time"

	_ "github.com/glebarez/sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRewardSaveTest(t *testing.T) *RewardSaveRepository {
	db, err := sqlx.Open("sqlite", ":memory:")
	require.NoError(t, err)
	_, err = db.Exec(`
		CREATE TABLE rewards (
			match TEXT PRIMARY KEY,
			reward INTEGER,
			reward_type TEXT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		);
	`)
	require.NoError(t, err)
	return NewRewardSaveRepository(db, nil)
}

func TestRewardSaveRepository_Save(t *testing.T) {
	repo := setupRewardSaveTest(t)

	t.Run("Save new reward", func(t *testing.T) {
		reward := map[string]any{
			"match":       "user123",
			"reward":      int64(50),
			"reward_type": "%",
		}
		err := repo.Save(context.Background(), reward)
		require.NoError(t, err)
		var match string
		var rewardVal int64
		var rewardType string
		var createdAt, updatedAt time.Time
		err = repo.db.QueryRowx("SELECT match, reward, reward_type, created_at, updated_at FROM rewards WHERE match = ?", "user123").
			Scan(&match, &rewardVal, &rewardType, &createdAt, &updatedAt)
		require.NoError(t, err)
		assert.Equal(t, "user123", match)
		assert.Equal(t, int64(50), rewardVal)
		assert.Equal(t, "%", rewardType)
		assert.False(t, createdAt.IsZero())
		assert.False(t, updatedAt.IsZero())
	})

	t.Run("Update existing reward", func(t *testing.T) {
		reward := map[string]any{
			"match":       "user123",
			"reward":      int64(75),
			"reward_type": "pt",
		}
		err := repo.Save(context.Background(), reward)
		require.NoError(t, err)
		var rewardVal int64
		var rewardType string
		err = repo.db.QueryRowx("SELECT reward, reward_type FROM rewards WHERE match = ?", "user123").
			Scan(&rewardVal, &rewardType)
		require.NoError(t, err)
		assert.Equal(t, int64(75), rewardVal)
		assert.Equal(t, "pt", rewardType)
	})

	t.Run("Update reward with partial fields", func(t *testing.T) {
		initial := map[string]any{
			"match":       "user456",
			"reward":      int64(30),
			"reward_type": "%",
		}
		err := repo.Save(context.Background(), initial)
		require.NoError(t, err)
		partial := map[string]any{
			"match":       "user456",
			"reward":      int64(100),
			"reward_type": "%",
		}
		err = repo.Save(context.Background(), partial)
		require.NoError(t, err)
		var rewardVal int64
		var rewardType string
		err = repo.db.QueryRowx("SELECT reward, reward_type FROM rewards WHERE match = ?", "user456").
			Scan(&rewardVal, &rewardType)
		require.NoError(t, err)
		assert.Equal(t, int64(100), rewardVal)
		assert.Equal(t, "%", rewardType)
	})
}
