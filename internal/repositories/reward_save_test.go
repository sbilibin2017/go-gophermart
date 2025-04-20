package repositories

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func setupReewardSaveTestDB() (*sql.DB, func()) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		panic(err)
	}
	tearDown := func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}
	_, err = db.Exec(`
		CREATE TABLE rewards (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			match TEXT NOT NULL,
			reward TEXT NOT NULL,
			reward_type TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			CONSTRAINT match_unique UNIQUE (match)
		);
	`)
	if err != nil {
		panic(err)
	}
	return db, tearDown
}

func TestRewardSaveRepository_Save(t *testing.T) {
	db, tearDown := setupReewardSaveTestDB()
	defer tearDown()
	repo := NewRewardSaveRepository(
		db,
		func(ctx context.Context) *sql.Tx {
			return nil
		},
	)

	testCases := []struct {
		name        string
		data        map[string]any
		expectedErr bool
	}{
		{
			name:        "Save new reward",
			data:        map[string]any{"match": "match1", "reward": "reward1", "reward_type": "type1"},
			expectedErr: false,
		},
		{
			name:        "Update existing reward",
			data:        map[string]any{"match": "match1", "reward": "reward_updated", "reward_type": "type_updated"},
			expectedErr: false,
		},
		{
			name:        "Save duplicate match",
			data:        map[string]any{"match": "match2", "reward": "reward2", "reward_type": "type2"},
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.Save(context.Background(), tc.data)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
