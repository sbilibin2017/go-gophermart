package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	_, err = db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}
	return db
}

func TestGetExecutorWithTx(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	ctx := SetTx(context.Background(), tx)
	executor := GetExecutor(ctx, db)
	txExecutor, ok := executor.(*txExecutor)
	assert.True(t, ok)
	assert.Equal(t, tx, txExecutor.tx)
}

func TestGetExecutorWithDB(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	ctx := context.Background()
	executor := GetExecutor(ctx, db)
	dbExecutor, ok := executor.(*dbExecutor)
	assert.True(t, ok)
	assert.Equal(t, db, dbExecutor.db)
}

func TestDBExecutorMethods(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	dbExecutor := &dbExecutor{db: db}
	_, err := dbExecutor.ExecContext(context.Background(), "INSERT INTO users (name) VALUES (?)", "John")
	assert.NoError(t, err)
	var name string
	err = dbExecutor.QueryRowContext(context.Background(), "SELECT name FROM users WHERE id = ?", 1).Scan(&name)
	assert.NoError(t, err)
	assert.Equal(t, "John", name)
	rows, err := dbExecutor.QueryContext(context.Background(), "SELECT id, name FROM users")
	assert.NoError(t, err)
	defer rows.Close()
	var id int
	for rows.Next() {
		err := rows.Scan(&id, &name)
		assert.NoError(t, err)
		assert.Equal(t, "John", name)
	}
}

func TestTxExecutorMethods(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	txExecutor := &txExecutor{tx: tx}
	_, err = txExecutor.ExecContext(context.Background(), "INSERT INTO users (name) VALUES (?)", "Alice")
	assert.NoError(t, err)
	var name string
	err = txExecutor.QueryRowContext(context.Background(), "SELECT name FROM users WHERE id = ?", 1).Scan(&name)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", name)
	rows, err := txExecutor.QueryContext(context.Background(), "SELECT id, name FROM users")
	assert.NoError(t, err)
	defer rows.Close()
	var id int
	for rows.Next() {
		err := rows.Scan(&id, &name)
		assert.NoError(t, err)
		assert.Equal(t, "Alice", name)
	}
	err = tx.Commit()
	assert.NoError(t, err)
}
