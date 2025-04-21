package helpers

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func TestQueryRowContext(t *testing.T) {
	db, err := sqlx.Open("sqlite", "file::memory:?cache=shared")
	assert.NoError(t, err)
	_, err = db.Exec(`DROP TABLE IF EXISTS test_table`)
	assert.NoError(t, err)
	_, err = db.Exec(`CREATE TABLE test_table (id INTEGER PRIMARY KEY, name TEXT)`)
	assert.NoError(t, err)
	_, err = db.Exec(`INSERT INTO test_table (name) VALUES ('test_name')`)
	assert.NoError(t, err)
	txProvider := func(ctx context.Context) *sqlx.Tx {
		return nil
	}
	query := "SELECT name FROM test_table WHERE id = :id"
	args := map[string]any{"id": 1}
	row, err := QueryRowContext(context.Background(), db, txProvider, query, args)
	assert.NoError(t, err)
	assert.NotNil(t, row)
	var name string
	err = row.Scan(&name)
	assert.NoError(t, err)
	assert.Equal(t, "test_name", name)
}

func TestQueryRowContextWithTransaction(t *testing.T) {
	db, err := sqlx.Open("sqlite", "file::memory:?cache=shared")
	assert.NoError(t, err)
	_, err = db.Exec(`DROP TABLE IF EXISTS test_table`)
	assert.NoError(t, err)
	_, err = db.Exec(`CREATE TABLE test_table (id INTEGER PRIMARY KEY, name TEXT)`)
	assert.NoError(t, err)
	_, err = db.Exec(`INSERT INTO test_table (name) VALUES ('transaction_test')`)
	assert.NoError(t, err)
	txProvider := func(ctx context.Context) *sqlx.Tx {
		tx, err := db.BeginTxx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		return tx
	}
	query := "SELECT name FROM test_table WHERE id = :id"
	args := map[string]any{"id": 1}
	tx := txProvider(context.Background())
	ctxWithTx := TxToContext(context.Background(), tx)
	row, err := QueryRowContext(ctxWithTx, db, func(ctx context.Context) *sqlx.Tx {
		return TxFromContext(ctx)
	}, query, args)
	assert.NoError(t, err)
	assert.NotNil(t, row)
	var name string
	err = row.Scan(&name)
	assert.NoError(t, err)
	assert.Equal(t, "transaction_test", name)
	err = tx.Commit()
	assert.NoError(t, err)
}

func TestQueryRowContextWithNoResult(t *testing.T) {
	db, err := sqlx.Open("sqlite", "file::memory:?cache=shared")
	assert.NoError(t, err)
	_, err = db.Exec(`DROP TABLE IF EXISTS test_table`)
	assert.NoError(t, err)
	_, err = db.Exec(`CREATE TABLE test_table (id INTEGER PRIMARY KEY, name TEXT)`)
	assert.NoError(t, err)
	query := "SELECT name FROM test_table WHERE id = :id"
	args := map[string]any{"id": 999} // ID не существует
	txProvider := func(ctx context.Context) *sqlx.Tx {
		return nil
	}
	row, err := QueryRowContext(context.Background(), db, txProvider, query, args)
	assert.NoError(t, err)
	assert.NotNil(t, row)
	var name string
	err = row.Scan(&name)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "sql: no rows in result set")
}

func TestQueryRowContextWithTransactionNoResult(t *testing.T) {
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
	query := "SELECT name FROM test_table WHERE id = :id"
	args := map[string]any{"id": 999} // ID не существует
	tx := txProvider(context.Background())
	ctxWithTx := TxToContext(context.Background(), tx)
	row, err := QueryRowContext(ctxWithTx, db, func(ctx context.Context) *sqlx.Tx {
		return TxFromContext(ctx)
	}, query, args)
	assert.NoError(t, err)
	assert.NotNil(t, row)
	var name string
	err = row.Scan(&name)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "sql: no rows in result set")
	err = tx.Commit()
	assert.NoError(t, err)
}
