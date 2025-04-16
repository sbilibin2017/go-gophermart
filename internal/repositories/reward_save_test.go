package repositories

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRewardSaveRepository_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewRewardSaveRepository(sqlxDB)

	ctx := context.Background()

	reward := &models.RewardDB{
		Match:      "test-match",
		Reward:     100,
		RewardType: "points",
	}

	t.Run("save new reward", func(t *testing.T) {
		mock.ExpectExec("").
			WithArgs(reward.Match, reward.Reward, reward.RewardType).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Save(ctx, reward)
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("update existing reward", func(t *testing.T) {
		mock.ExpectExec("").
			WithArgs(reward.Match, reward.Reward, reward.RewardType).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Save(ctx, reward)
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectExec("").
			WithArgs(reward.Match, reward.Reward, reward.RewardType).
			WillReturnError(fmt.Errorf("db error"))

		err := repo.Save(ctx, reward)
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
