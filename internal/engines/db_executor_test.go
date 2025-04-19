package engines

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func createUsersTable(db *sqlx.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT
	)`)
	return err
}

func TestDBExecutor_Execute_NoTransaction(t *testing.T) {
	db, err := sqlx.Open("sqlite", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	err = createUsersTable(db)
	assert.NoError(t, err)

	executor := NewDBExecutor(db, func(ctx context.Context) *sqlx.Tx {
		return nil
	})

	query := "INSERT INTO users (name) VALUES (:name)"
	args := map[string]any{
		"name": "John Doe",
	}

	err = executor.Execute(context.Background(), query, args)
	assert.NoError(t, err)

	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM users WHERE name = ?", "John Doe")
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestDBExecutor_Execute_WithTransaction(t *testing.T) {
	db, err := sqlx.Open("sqlite", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	err = createUsersTable(db)
	assert.NoError(t, err)

	tx, err := db.Beginx()
	assert.NoError(t, err)

	executor := NewDBExecutor(db, func(ctx context.Context) *sqlx.Tx {
		return tx
	})

	query := "INSERT INTO users (name) VALUES (:name)"
	args := map[string]any{
		"name": "John Doe",
	}

	err = executor.Execute(context.Background(), query, args)
	assert.NoError(t, err)

	err = tx.Commit()
	assert.NoError(t, err)

	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM users WHERE name = ?", "John Doe")
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestDBExecutor_CreateTable(t *testing.T) {
	db, err := sqlx.Open("sqlite", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	err = createUsersTable(db)
	assert.NoError(t, err)

	var tableName string
	err = db.Get(&tableName, "SELECT name FROM sqlite_master WHERE type='table' AND name='users'")
	assert.NoError(t, err)
	assert.Equal(t, "users", tableName)
}

func TestDBExecutor_Execute_WithRollback(t *testing.T) {
	db, err := sqlx.Open("sqlite", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	err = createUsersTable(db)
	assert.NoError(t, err)

	tx, err := db.Beginx()
	assert.NoError(t, err)

	executor := NewDBExecutor(db, func(ctx context.Context) *sqlx.Tx {
		return tx
	})

	query := "INSERT INTO users (name) VALUES (:name)"
	args := map[string]any{
		"name": "John Doe",
	}

	err = executor.Execute(context.Background(), query, args)
	assert.NoError(t, err)

	err = tx.Rollback()
	assert.NoError(t, err)

	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM users WHERE name = ?", "John Doe")
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestDBExecutor_Execute_ErrorHandling(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	query := "INSERT INTO users (name) VALUES (:name)"
	args := map[string]any{
		"name": "John Doe",
	}

	mock.ExpectBegin()
	mock.ExpectExec(query).
		WithArgs(args["name"]).
		WillReturnError(errors.New("simulated error"))

	executor := NewDBExecutor(sqlxDB, func(ctx context.Context) *sqlx.Tx {
		tx, _ := sqlxDB.Beginx()
		return tx
	})
	err = executor.Execute(context.Background(), query, args)
	assert.Error(t, err)

}
