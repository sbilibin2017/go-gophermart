package engines

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestName struct {
	Name string `db:"name"`
}

func TestDBQuerier_Query_WithoutTransaction(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()
	mock.ExpectQuery("SELECT name FROM test_table WHERE name = ?").
		WithArgs("test_name").
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("test_name"))
	querier := DBQuerier{
		db: db,
		txProvider: func(ctx context.Context) (*sqlx.Tx, bool) {
			return nil, false
		},
	}
	var result TestName
	query := `SELECT name FROM test_table WHERE name = :name`
	args := map[string]any{"name": "test_name"}
	err := querier.Query(context.Background(), &result, query, args)
	assert.NoError(t, err)
	assert.Equal(t, "test_name", result.Name)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDBQuerier_Query_WithTransaction(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT name FROM test_table WHERE name = ?").
		WithArgs("tx_name").
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("tx_name"))
	mock.ExpectCommit()
	querier := DBQuerier{
		db: db,
		txProvider: func(ctx context.Context) (*sqlx.Tx, bool) {
			tx, _ := db.Beginx()
			return tx, true
		},
	}
	var result TestName
	query := `SELECT name FROM test_table WHERE name = :name`
	args := map[string]any{"name": "tx_name"}
	err := querier.Query(context.Background(), &result, query, args)
	assert.NoError(t, err)
	assert.Equal(t, "tx_name", result.Name)
}

func TestDBQuerier_Query_ErrorInNamedQueryContext_WithoutTransaction(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()
	mock.ExpectQuery("SELECT name FROM test_table WHERE name = ?").
		WithArgs("test_name").
		WillReturnError(fmt.Errorf("execution error"))
	querier := DBQuerier{
		db: db,
		txProvider: func(ctx context.Context) (*sqlx.Tx, bool) {
			return nil, false
		},
	}
	var name string
	query := `SELECT name FROM test_table WHERE name = :name`
	args := map[string]any{"name": "test_name"}
	err := querier.Query(context.Background(), &name, query, args)
	assert.Error(t, err)
	assert.Equal(t, "execution error", err.Error())
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDBQuerier_Query_ErrorInNamedQueryContext_WithTransaction(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT name FROM test_table WHERE name = ?").
		WithArgs("tx_name").
		WillReturnError(fmt.Errorf("transaction query error"))
	mock.ExpectRollback()
	querier := DBQuerier{
		db: db,
		txProvider: func(ctx context.Context) (*sqlx.Tx, bool) {
			tx, _ := db.Beginx()
			return tx, true
		},
	}
	var name string
	query := `SELECT name FROM test_table WHERE name = :name`
	args := map[string]any{"name": "tx_name"}
	err := querier.Query(context.Background(), &name, query, args)
	assert.Error(t, err)
	assert.Equal(t, "transaction query error", err.Error())
}

func TestDBQuerier_Query_ScanError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()
	mock.ExpectQuery("SELECT name FROM test_table WHERE name = ?").
		WithArgs("tx_name").
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("tx_name"))
	querier := DBQuerier{
		db: db,
		txProvider: func(ctx context.Context) (*sqlx.Tx, bool) {
			return nil, false
		},
	}
	var result struct {
		Age int `db:"name"`
	}
	query := `SELECT name FROM test_table WHERE name = :name`
	args := map[string]any{"name": "tx_name"}
	err := querier.Query(context.Background(), &result, query, args)
	assert.Error(t, err)
}

func TestDBQuerier_Query_ScanError_InTransaction(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()
	mock.ExpectBegin()
	type Result struct {
		Name int `db:"name"`
	}
	mock.ExpectQuery("SELECT name FROM test_table WHERE name = ?").
		WithArgs("tx_name").
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("not_an_int"))
	mock.ExpectCommit()
	querier := DBQuerier{
		db: db,
		txProvider: func(ctx context.Context) (*sqlx.Tx, bool) {
			tx, _ := db.Beginx()
			return tx, true
		},
	}
	var result Result
	query := `SELECT name FROM test_table WHERE name = :name`
	args := map[string]any{"name": "tx_name"}
	err := querier.Query(context.Background(), &result, query, args)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "converting driver.Value type string")
}
