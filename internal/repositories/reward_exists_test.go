package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func setupRewardExistsTestDB() (*sql.DB, func()) {
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
			match TEXT NOT NULL
		);
		INSERT INTO rewards (match) VALUES ('match1');
		INSERT INTO rewards (match) VALUES ('match2');
	`)
	if err != nil {
		panic(err)
	}
	return db, tearDown
}

func TestRewardExistsRepository_Exists(t *testing.T) {
	db, tearDown := setupRewardExistsTestDB()
	defer tearDown()
	repo := NewRewardExistsRepository(
		db,
		func(ctx context.Context) *sql.Tx {
			return nil
		},
	)
	testCases := []struct {
		name     string
		match    string
		expected bool
	}{
		{"Match exists", "match1", true},
		{"Match does not exist", "match3", false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			exists, err := repo.Exists(context.Background(), map[string]any{"match": tc.match})
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, exists)
		})
	}
}

func TestRewardExistsRepository_Exists_QueryError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRewardExistsRepository(db, func(ctx context.Context) *sql.Tx {
		return nil
	})
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("match1").
		WillReturnError(fmt.Errorf("query error"))
	exists, err := repo.Exists(context.Background(), map[string]any{"match": "match1"})
	assert.Error(t, err)
	assert.False(t, exists)
}

func TestRewardExistsRepository_Exists_QueryContextError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock БД: %v", err)
	}
	defer db.Close()
	repo := NewRewardExistsRepository(db, func(ctx context.Context) *sql.Tx {
		return nil
	})
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("match1").
		WillReturnError(fmt.Errorf("query context error"))
	exists, err := repo.Exists(context.Background(), map[string]any{"match": "match1"})
	assert.Error(t, err)
	assert.Equal(t, "query context error", err.Error())
	assert.False(t, exists)
}

func TestRewardExistsRepository_Exists_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()
	repo := NewRewardExistsRepository(db, func(ctx context.Context) *sql.Tx {
		return nil
	})
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("match1").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(nil))
	_, err = repo.Exists(context.Background(), map[string]any{"match": "match1"})
	assert.Error(t, err)
}
