package helpers

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func TestExecContext(t *testing.T) {
	db, err := sqlx.Open("sqlite", "file::memory:?cache=shared")
	assert.NoError(t, err)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS test_table (id INTEGER PRIMARY KEY, name TEXT)`)
	assert.NoError(t, err)
	txProvider := func(ctx context.Context) *sqlx.Tx {
		return nil
	}
	query := "INSERT INTO test_table (name) VALUES (:name)"
	args := map[string]any{"name": "test_name"}
	result, err := ExecContext(context.Background(), db, txProvider, query, args)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	var name string
	err = db.Get(&name, "SELECT name FROM test_table WHERE id = 1")
	assert.NoError(t, err)
	assert.Equal(t, "test_name", name)
	rowsAffected, _ := result.RowsAffected()
	assert.Equal(t, int64(1), rowsAffected)
}

func TestExecContextWithTransaction(t *testing.T) {
	db, err := sqlx.Open("sqlite", "file::memory:?cache=shared")
	assert.NoError(t, err)
	_, err = db.Exec(`DROP TABLE IF EXISTS test_table`)
	assert.NoError(t, err)
	_, err = db.Exec(`CREATE TABLE test_table (id INTEGER PRIMARY KEY, name TEXT)`)
	assert.NoError(t, err)
	txProvider := func(ctx context.Context) *sqlx.Tx {
		tx, err := db.BeginTxx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		return tx
	}
	query := "INSERT INTO test_table (name) VALUES (:name)"
	args := map[string]any{"name": "transaction_test"}
	tx := txProvider(context.Background())
	ctxWithTx := contextutils.TxToContext(context.Background(), tx)
	result, err := ExecContext(ctxWithTx, db, func(ctx context.Context) *sqlx.Tx {
		return contextutils.TxFromContext(ctx)
	}, query, args)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	err = tx.Commit()
	assert.NoError(t, err)
	var name string
	err = db.Get(&name, "SELECT name FROM test_table WHERE id = 1")
	assert.NoError(t, err)
	assert.Equal(t, "transaction_test", name)
	rowsAffected, _ := result.RowsAffected()
	assert.Equal(t, int64(1), rowsAffected)
	query2 := "INSERT INTO test_table (name) VALUES (:name)"
	args2 := map[string]any{"name": "test_name"}
	result2, err := ExecContext(context.Background(), db, func(ctx context.Context) *sqlx.Tx {
		return nil
	}, query2, args2)
	assert.NoError(t, err)
	assert.NotNil(t, result2)
	var name2 string
	err = db.Get(&name2, "SELECT name FROM test_table WHERE id = 2")
	assert.NoError(t, err)
	assert.Equal(t, "test_name", name2)
	rowsAffected2, _ := result2.RowsAffected()
	assert.Equal(t, int64(1), rowsAffected2)
}
