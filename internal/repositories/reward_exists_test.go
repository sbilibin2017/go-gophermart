package repositories

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setupRewardExists(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, *RewardExistsRepository) {
	t.Helper()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewRewardExistsRepository(sqlxDB)
	return sqlxDB, mock, repo
}

func TestRewardExistsRepository_Exists(t *testing.T) {
	tests := []struct {
		name           string
		filter         map[string]any
		mockSetup      func(mock sqlmock.Sqlmock)
		expectedResult bool
		expectedError  error
	}{
		{
			name: "Exists returns true",
			filter: map[string]any{
				"match": "order123",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT EXISTS \(SELECT 1 FROM rewards WHERE match = \$1\)`).
					WithArgs("order123").
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name: "Exists returns false",
			filter: map[string]any{
				"match": "order123",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT EXISTS \(SELECT 1 FROM rewards WHERE match = \$1\)`).
					WithArgs("order123").
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
			},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name: "Database query error",
			filter: map[string]any{
				"match": "order123",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT EXISTS \(SELECT 1 FROM rewards WHERE match = \$1\)`).
					WithArgs("order123").
					WillReturnError(errors.New("database error"))
			},
			expectedResult: false,
			expectedError:  errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlxDB, mock, repo := setupRewardExists(t)
			defer sqlxDB.Close()
			tt.mockSetup(mock)
			result, err := repo.Exists(context.Background(), tt.filter)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
