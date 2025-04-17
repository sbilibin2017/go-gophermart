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

func TestRewardExistsRepository_Exists(t *testing.T) {
	tests := []struct {
		name              string
		filter            *dto.RewardExistsFilterDB
		mockExpectActions func(mock sqlmock.Sqlmock)
		expectedExists    bool
		expectedErr       error
	}{
		{
			name:   "Возвращает true, когда награда существует",
			filter: &dto.RewardExistsFilterDB{Match: "some_match"},
			mockExpectActions: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
				mock.ExpectQuery(`SELECT EXISTS \(SELECT 1 FROM rewards WHERE match = \$1\)`).
					WithArgs("some_match").
					WillReturnRows(rows)
			},
			expectedExists: true,
			expectedErr:    nil,
		},
		{
			name:   "Возвращает false, когда награды не существует",
			filter: &dto.RewardExistsFilterDB{Match: "nonexistent_match"},
			mockExpectActions: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
				mock.ExpectQuery(`SELECT EXISTS \(SELECT 1 FROM rewards WHERE match = \$1\)`).
					WithArgs("nonexistent_match").
					WillReturnRows(rows)
			},
			expectedExists: false,
			expectedErr:    nil,
		},
		{
			name:   "Возвращает ошибку, когда запрос не удается",
			filter: &dto.RewardExistsFilterDB{Match: "some_match"},
			mockExpectActions: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT EXISTS \(SELECT 1 FROM rewards WHERE match = \$1\)`).
					WithArgs("some_match").
					WillReturnError(assert.AnError)
			},
			expectedExists: false,
			expectedErr:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			sqlxDB := sqlx.NewDb(db, "postgres")
			repo := NewRewardExistsRepository(sqlxDB)

			tt.mockExpectActions(mock)

			exists, err := repo.Exists(context.Background(), tt.filter)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expectedExists, exists)
		})
	}
}
