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

func TestUnitOfWork_Do_Success(t *testing.T) {
	db, _, _, cleanup := SetupPostgresContainer(t)
	defer cleanup()
	uow := NewTransaction(db)
	_, err := db.Exec("CREATE TABLE users (id SERIAL PRIMARY KEY, username VARCHAR(255) NOT NULL);")
	require.NoError(t, err)
	operation := func(tx *sql.Tx) error {
		_, err := tx.Exec("INSERT INTO users (username) VALUES ('testuser')")
		return err
	}
	err = uow.Do(context.Background(), operation)
	require.NoError(t, err)
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'testuser'").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestUnitOfWork_Do_NilDB(t *testing.T) {
	uow := NewTransaction(nil)
	operation := func(tx *sql.Tx) error {
		return nil
	}
	err := uow.Do(context.Background(), operation)
	assert.NoError(t, err)
}

func TestUnitOfWork_Do_OperationError_Rollback(t *testing.T) {
	db, _, _, cleanup := SetupPostgresContainer(t)
	defer cleanup()
	_, err := db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, username TEXT NOT NULL);")
	require.NoError(t, err)
	uow := NewTransaction(db)
	operation := func(tx *sql.Tx) error {
		_, err := tx.Exec("INSERT INTO users (username) VALUES (?)", "testuser")
		if err != nil {
			return err
		}
		return errors.New("operation failed")
	}
	err = uow.Do(context.Background(), operation)
	require.Error(t, err)

}
