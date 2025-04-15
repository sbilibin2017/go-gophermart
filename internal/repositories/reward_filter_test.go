package repositories

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func init() {
	logger.Init(zapcore.InfoLevel)
}

func TestRewardFilterRepository_Filter(t *testing.T) {
	tests := []struct {
		name           string
		filter         *RewardFilterDB
		mockExpect     func(mock sqlmock.Sqlmock, filter *RewardFilterDB)
		expectedReward *RewardFilteredDB
		expectedError  bool
	}{
		{
			name: "Success",
			filter: &RewardFilterDB{
				Description: "bonus",
			},
			mockExpect: func(mock sqlmock.Sqlmock, filter *RewardFilterDB) {
				mock.ExpectQuery(`SELECT match, reward, reward_type FROM rewards WHERE match ILIKE \$1`).
					WithArgs("%bonus%").
					WillReturnRows(sqlmock.NewRows([]string{"match", "reward", "reward_type"}).
						AddRow("bonus", 1000, "cash"))
			},
			expectedReward: &RewardFilteredDB{
				Match:      "bonus",
				Reward:     1000,
				RewardType: "cash",
			},
			expectedError: false,
		},
		{
			name: "Error - No Matching Reward",
			filter: &RewardFilterDB{
				Description: "nonexistent",
			},
			mockExpect: func(mock sqlmock.Sqlmock, filter *RewardFilterDB) {
				mock.ExpectQuery(`SELECT match, reward, reward_type FROM rewards WHERE match ILIKE \$1`).
					WithArgs("%nonexistent%").
					WillReturnError(fmt.Errorf("no rows found"))
			},
			expectedReward: nil,
			expectedError:  true,
		},
		{
			name: "Error - Database Error",
			filter: &RewardFilterDB{
				Description: "bonus",
			},
			mockExpect: func(mock sqlmock.Sqlmock, filter *RewardFilterDB) {
				mock.ExpectQuery(`SELECT match, reward, reward_type FROM rewards WHERE match ILIKE \$1`).
					WithArgs("%bonus%").
					WillReturnError(fmt.Errorf("db error"))
			},
			expectedReward: nil,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to open mock database: %v", err)
			}
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "pgx")
			repo := NewRewardFilterRepository(sqlxDB)

			tt.mockExpect(mock, tt.filter)

			reward, err := repo.Filter(context.Background(), tt.filter)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedReward, reward)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
