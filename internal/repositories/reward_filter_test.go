package repositories

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRewardFilterRepository_Filter(t *testing.T) {
	setupQuery := `CREATE TABLE rewards (
		reward BIGINT NOT NULL,
		reward_type TEXT NOT NULL,
		match TEXT NOT NULL
	)`
	db, cleanup := setupDB(t, setupQuery)
	defer cleanup()

	_, err := db.Exec(`
		INSERT INTO rewards (reward, reward_type, match)
		VALUES 
			(100, 'bonus', 'Welcome Bonus'),
			(200, 'cashback', 'Weekly Cashback'),
			(300, 'bonus', 'Loyalty Bonus')
	`)
	require.NoError(t, err)

	repo := NewRewardFilterRepository(db)

	testCases := []struct {
		name     string
		filter   []*types.RewardFilter
		expected map[string]*types.RewardDB
	}{
		{
			name: "Match description with 'bonus'",
			filter: []*types.RewardFilter{
				{Description: "bonus"},
			},
			expected: map[string]*types.RewardDB{
				"Welcome Bonus": {Reward: 100, RewardType: "bonus", Match: "Welcome Bonus"},
				"Loyalty Bonus": {Reward: 300, RewardType: "bonus", Match: "Loyalty Bonus"},
			},
		},
		{
			name: "No matches",
			filter: []*types.RewardFilter{
				{Description: "nonexistent"},
			},
			expected: map[string]*types.RewardDB{}, // Ожидаем пустую карту
		},
		{
			name: "Multiple descriptions",
			filter: []*types.RewardFilter{
				{Description: "bonus"},
				{Description: "cashback"},
			},
			expected: map[string]*types.RewardDB{
				"Welcome Bonus":   {Reward: 100, RewardType: "bonus", Match: "Welcome Bonus"},
				"Loyalty Bonus":   {Reward: 300, RewardType: "bonus", Match: "Loyalty Bonus"},
				"Weekly Cashback": {Reward: 200, RewardType: "cashback", Match: "Weekly Cashback"},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rewards, err := repo.Filter(context.Background(), tc.filter)
			require.NoError(t, err)

			resultMap := make(map[string]*types.RewardDB)
			for _, r := range rewards {
				resultMap[r.Match] = r
			}

			assert.Equal(t, tc.expected, resultMap)
		})
	}
}

func TestRewardFilterRepository_Filter_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewRewardFilterRepository(sqlxDB)

	filters := []*types.RewardFilter{{Description: "bonus"}}

	mock.ExpectQuery(`SELECT reward, reward_type, match FROM rewards WHERE match ILIKE ANY`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(errors.New("mocked db error"))

	rewards, err := repo.Filter(context.Background(), filters)

	require.Error(t, err)
	assert.Nil(t, rewards)
}
