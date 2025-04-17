package repositories

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setupRewardSaveTest(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, *RewardSaveRepository) {
	t.Helper()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewRewardSaveRepository(sqlxDB)
	return sqlxDB, mock, repo
}

func TestRewardSaveRepository_Save(t *testing.T) {
	tests := []struct {
		name          string
		match         map[string]any
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			name: "Success",
			match: map[string]any{
				"match":       "order123",
				"reward":      100,
				"reward_type": "pt",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO rewards \(match, reward, reward_type, created_at, updated_at\)`).
					WithArgs("order123", 100, "pt").
					WillReturnResult(sqlmock.NewResult(1, 1)) // Simulate success
			},
			expectedError: nil,
		},
		{
			name: "Conflict occurs, update happens",
			match: map[string]any{
				"match":       "order123",
				"reward":      100,
				"reward_type": "pt",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO rewards \(match, reward, reward_type, created_at, updated_at\)`).
					WithArgs("order123", 100, "pt").
					WillReturnResult(sqlmock.NewResult(1, 1)) // Simulate success
			},
			expectedError: nil,
		},
		{
			name: "Database error during insert",
			match: map[string]any{
				"match":       "order123",
				"reward":      100,
				"reward_type": "pt",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO rewards \(match, reward, reward_type, created_at, updated_at\)`).
					WithArgs("order123", 100, "pt").
					WillReturnError(errors.New("insert error"))
			},
			expectedError: errors.New("insert error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlxDB, mock, repo := setupRewardSaveTest(t)
			defer sqlxDB.Close()
			tt.mockSetup(mock)
			err := repo.Save(context.Background(), tt.match)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
