package repositories

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRewardSaveRepository_Save(t *testing.T) {
	tests := []struct {
		название          string
		reward            *dto.RewardDB
		mockExpectActions func(mock sqlmock.Sqlmock)
		expectedErr       error
	}{
		{
			название: "Успешно сохраняет награду",
			reward:   &dto.RewardDB{Match: "some_match", Reward: 100, RewardType: "cash"},
			mockExpectActions: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO rewards \(match, reward, reward_type, created_at, updated_at\)`).
					WithArgs("some_match", int64(100), "cash").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			название: "Возвращает ошибку при ошибке выполнения запроса",
			reward:   &dto.RewardDB{Match: "some_match", Reward: 100, RewardType: "cash"},
			mockExpectActions: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO rewards \(match, reward, reward_type, created_at, updated_at\)`).
					WithArgs("some_match", int64(100), "cash").
					WillReturnError(assert.AnError)
			},
			expectedErr: assert.AnError,
		},
		{
			название: "Успешно обновляет награду при конфликте",
			reward:   &dto.RewardDB{Match: "some_match", Reward: 150, RewardType: "points"},
			mockExpectActions: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO rewards \(match, reward, reward_type, created_at, updated_at\)`).
					WithArgs("some_match", int64(150), "points").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.название, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			sqlxDB := sqlx.NewDb(db, "postgres")
			repo := NewRewardSaveRepository(sqlxDB)

			tt.mockExpectActions(mock)

			err = repo.Save(context.Background(), tt.reward)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
