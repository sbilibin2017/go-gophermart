package db_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"

	"github.com/sbilibin2017/go-gophermart/internal/db"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	dbConn, err := sqlx.Open("sqlite", ":memory:")
	require.NoError(t, err)
	schema := `CREATE TABLE IF NOT EXISTS test (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT);`
	_, err = dbConn.Exec(schema)
	require.NoError(t, err)
	return dbConn
}

func TestWithTx_Success(t *testing.T) {
	dbConn := setupTestDB(t)
	defer dbConn.Close()
	ctx := context.Background()
	err := db.WithTx(ctx, dbConn, func(tx *sqlx.Tx) error {
		_, err := tx.Exec("INSERT INTO test(name) VALUES (?)", "John")
		return err
	})
	require.NoError(t, err)
	var count int
	err = dbConn.Get(&count, "SELECT COUNT(*) FROM test WHERE name = ?", "John")
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestWithTx_Rollback(t *testing.T) {
	dbConn := setupTestDB(t)
	defer dbConn.Close()
	ctx := context.Background()
	customErr := errors.New("force rollback")
	err := db.WithTx(ctx, dbConn, func(tx *sqlx.Tx) error {
		_, execErr := tx.Exec("INSERT INTO test(name) VALUES (?)", "Jane")
		require.NoError(t, execErr)
		return customErr
	})
	require.ErrorIs(t, err, customErr)
	var count int
	err = dbConn.Get(&count, "SELECT COUNT(*) FROM test WHERE name = ?", "Jane")
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestTxFromContext(t *testing.T) {
	dbConn := setupTestDB(t)
	defer dbConn.Close()
	ctx := context.Background()
	err := db.WithTx(ctx, dbConn, func(tx *sqlx.Tx) error {
		ctxWithTx := context.WithValue(ctx, db.TxKey, tx)
		txFromCtx := db.TxFromContext(ctxWithTx)
		assert.Equal(t, tx, txFromCtx)
		return nil
	})
	require.NoError(t, err)
}

func TestWithTx_BeginTxError(t *testing.T) {
	dbConn, err := sqlx.Open("sqlite", ":memory:")
	require.NoError(t, err)
	dbConn.Close()
	ctx := context.Background()
	err = db.WithTx(ctx, dbConn, func(tx *sqlx.Tx) error {
		t.Fatal("This should not be called if BeginTxx fails")
		return nil
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to begin transaction")
}
