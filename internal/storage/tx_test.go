package storage

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestTx_Do_Success(t *testing.T) {
	db, _, _, cleanup := Setup(t)
	defer cleanup()
	tx := NewTx(db)
	_, err := db.Exec("CREATE TABLE users (id SERIAL PRIMARY KEY, username VARCHAR(255) NOT NULL);")
	require.NoError(t, err)
	operation := func(tx *sql.Tx) error {
		_, err := tx.Exec("INSERT INTO users (username) VALUES ('testuser')")
		return err
	}
	err = tx.Do(context.Background(), operation)
	require.NoError(t, err)
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'testuser'").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestTx_Do_NilDB(t *testing.T) {
	tx := NewTx(nil)
	operation := func(tx *sql.Tx) error {
		return nil
	}
	err := tx.Do(context.Background(), operation)
	assert.NoError(t, err)
}

func TestTx_Do_OperationError_Rollback(t *testing.T) {
	db, _, _, cleanup := Setup(t)
	defer cleanup()
	_, err := db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, username TEXT NOT NULL);")
	require.NoError(t, err)
	tx := NewTx(db)
	operation := func(tx *sql.Tx) error {
		_, err := tx.Exec("INSERT INTO users (username) VALUES (?)", "testuser")
		if err != nil {
			return err
		}
		return errors.New("operation failed")
	}
	err = tx.Do(context.Background(), operation)
	require.Error(t, err)

}
