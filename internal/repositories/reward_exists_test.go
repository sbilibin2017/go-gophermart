package repositories

import (
	"context"
	"testing"
	"time"

	_ "github.com/glebarez/sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setupRewardExistsTest(t *testing.T) *RewardExistsRepository {
	db, err := sqlx.Open("sqlite", ":memory:")
	assert.NoError(t, err)
	_, err = db.Exec(`
		CREATE TABLE rewards (
			match TEXT PRIMARY KEY,
			reward INTEGER,
			reward_type TEXT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		);
	`)
	assert.NoError(t, err)
	_, err = db.Exec(`
		INSERT INTO rewards (match, reward, reward_type, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?);
	`, "match123", 500, "cashback", time.Now(), time.Now())
	assert.NoError(t, err)
	return NewRewardExistsRepository(db, nil)
}

func TestRewardExists(t *testing.T) {
	repo := setupRewardExistsTest(t)

	t.Run("Reward exists", func(t *testing.T) {
		exists, err := repo.Exists(context.Background(), "match123")
		assert.NoError(t, err)
		assert.True(t, exists, "Reward should exist")
	})

	t.Run("Reward does not exist", func(t *testing.T) {
		exists, err := repo.Exists(context.Background(), "nonexistent_match")
		assert.NoError(t, err)
		assert.False(t, exists, "Reward should not exist")
	})

	t.Run("Empty database", func(t *testing.T) {
		db, err := sqlx.Open("sqlite", ":memory:")
		assert.NoError(t, err)
		_, err = db.Exec(`
			CREATE TABLE rewards (
				match TEXT PRIMARY KEY,
				reward INTEGER,
				reward_type TEXT,
				created_at TIMESTAMP,
				updated_at TIMESTAMP
			);
		`)
		assert.NoError(t, err)
		repo := NewRewardExistsRepository(db, nil)
		exists, err := repo.Exists(context.Background(), "match123")
		assert.NoError(t, err)
		assert.False(t, exists, "Reward should not exist in an empty database")
	})
}
