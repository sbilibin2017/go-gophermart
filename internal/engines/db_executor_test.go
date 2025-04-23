package engines

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/glebarez/sqlite"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	db, err := sqlx.Open("sqlite", ":memory:")
	require.NoError(t, err)

	logger.Init()

	schema := `
	CREATE TABLE test_table (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT
	);`
	_, err = db.Exec(schema)
	require.NoError(t, err)

	return db
}

func TestDBExecutor_Execute_WithoutTransaction(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	executor := DBExecutor{
		db: db,
		txProvider: func(ctx context.Context) (*sqlx.Tx, bool) {
			return nil, false
		},
	}

	query := `INSERT INTO test_table (name) VALUES (:name)`
	args := map[string]any{"name": "test_name"}

	err := executor.Execute(context.Background(), query, args)
	assert.NoError(t, err)

	var name string
	err = db.Get(&name, "SELECT name FROM test_table WHERE name = ?", "test_name")
	assert.NoError(t, err)
	assert.Equal(t, "test_name", name)
}

func TestDBExecutor_Execute_WithTransaction(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	tx, err := db.BeginTxx(context.Background(), nil)
	require.NoError(t, err)

	executor := DBExecutor{
		db: db,
		txProvider: func(ctx context.Context) (*sqlx.Tx, bool) {
			return tx, true
		},
	}

	query := `INSERT INTO test_table (name) VALUES (:name)`
	args := map[string]any{"name": "tx_name"}

	err = executor.Execute(context.Background(), query, args)
	assert.NoError(t, err)

	err = tx.Commit()
	assert.NoError(t, err)

	var name string
	err = db.Get(&name, "SELECT name FROM test_table WHERE name = ?", "tx_name")
	assert.NoError(t, err)
	assert.Equal(t, "tx_name", name)
}

func setupMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "sqlite3")

	logger.Init()

	return sqlxDB, mock
}

func TestDBExecutor_Execute_ErrorInNamedExecContext_WithoutTransaction(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	mock.ExpectExec("INSERT INTO test_table").WillReturnError(fmt.Errorf("execution error"))

	executor := DBExecutor{
		db: db,
		txProvider: func(ctx context.Context) (*sqlx.Tx, bool) {
			return nil, false
		},
	}

	query := `INSERT INTO test_table (name) VALUES (:name)`
	args := map[string]any{"name": "test_name"}

	err := executor.Execute(context.Background(), query, args)
	assert.Error(t, err)
	assert.Equal(t, "execution error", err.Error())

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDBExecutor_Execute_ErrorInNamedExecContext_WithTransaction(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO test_table").
		WithArgs("tx_name").
		WillReturnError(fmt.Errorf("transaction execution error"))
	mock.ExpectRollback()
	executor := DBExecutor{
		db: db,
		txProvider: func(ctx context.Context) (*sqlx.Tx, bool) {
			tx, _ := db.Beginx()
			return tx, true
		},
	}

	query := `INSERT INTO test_table (name) VALUES (:name)`
	args := map[string]any{"name": "tx_name"}

	err := executor.Execute(context.Background(), query, args)
	assert.Error(t, err)
	assert.Equal(t, "transaction execution error", err.Error())

}
