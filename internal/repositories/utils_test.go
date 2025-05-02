package repositories

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestExecContext(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening mock database connection: %v", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	mock.ExpectExec("INSERT INTO users").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	txProvider := func(ctx context.Context) (*sqlx.Tx, error) { return nil, nil }
	err = execContext(context.Background(), sqlxDB, txProvider, "INSERT INTO users (id) VALUES ($1)", 1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetContext(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening mock database connection: %v", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	rows := sqlmock.NewRows([]string{"name"}).AddRow("John")
	mock.ExpectQuery("SELECT name FROM users").WithArgs(1).WillReturnRows(rows)
	txProvider := func(ctx context.Context) (*sqlx.Tx, error) { return nil, nil }
	var user string
	err = getContext(context.Background(), sqlxDB, txProvider, "SELECT name FROM users WHERE id = $1", &user, 1)
	assert.NoError(t, err)
	assert.Equal(t, "John", user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSelectContext(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening mock database connection: %v", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	rows := sqlmock.NewRows([]string{"name"}).AddRow("John")
	mock.ExpectQuery("SELECT name FROM users").WithArgs(1).WillReturnRows(rows)
	txProvider := func(ctx context.Context) (*sqlx.Tx, error) { return nil, nil }
	var users []string
	err = selectContext(context.Background(), sqlxDB, txProvider, "SELECT name FROM users WHERE id = $1", &users, 1)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "John", users[0])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHandleExecutorErrorNoRows(t *testing.T) {
	err := handleExecutorError(sql.ErrNoRows)
	assert.NoError(t, err)
}

func TestHandleExecutorErrorGeneric(t *testing.T) {
	err := handleExecutorError(sql.ErrConnDone)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
}

func TestGetExecutor_WithTransaction(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening mock database connection: %v", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	tx := &sqlx.Tx{}
	txProvider := func(ctx context.Context) (*sqlx.Tx, error) {
		return tx, nil
	}
	executor := getExecutor(context.Background(), sqlxDB, txProvider)
	assert.Equal(t, tx, executor)
}
