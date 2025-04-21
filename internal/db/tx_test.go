package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func setupTestDB2(t *testing.T) *sql.DB {
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

func TestSetTxAndGetTx(t *testing.T) {
	db := setupTestDB2(t)
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	ctx := context.Background()
	ctxWithTx := SetTx(ctx, tx)
	retrievedTx, ok := GetTx(ctxWithTx)
	assert.True(t, ok)
	assert.Equal(t, tx, retrievedTx)
}
