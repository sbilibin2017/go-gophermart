package engines

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func createUsersTable2(db *sqlx.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT
	)`)
	return err
}

func TestDBQuerier_Query_NoTransaction(t *testing.T) {
	db, err := sqlx.Open("sqlite", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	err = createUsersTable(db)
	assert.NoError(t, err)

	_, err = db.Exec("INSERT INTO users (name) VALUES (?)", "John Doe")
	assert.NoError(t, err)

	querier := NewDBQuerier(db, func(ctx context.Context) *sqlx.Tx { return nil })

	query := "SELECT name FROM users WHERE name = :name"
	args := map[string]any{"name": "John Doe"}
	var result string

	err = querier.Query(context.Background(), &result, query, args)
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", result)
}

func TestDBQuerier_Query_WithTransaction(t *testing.T) {
	db, err := sqlx.Open("sqlite", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	err = createUsersTable2(db)
	assert.NoError(t, err)

	_, err = db.Exec("INSERT INTO users (name) VALUES (?)", "Jane Doe")
	assert.NoError(t, err)

	tx, err := db.Beginx()
	assert.NoError(t, err)

	querier := NewDBQuerier(db, func(ctx context.Context) *sqlx.Tx { return tx })

	query := "SELECT name FROM users WHERE name = :name"
	args := map[string]any{"name": "Jane Doe"}
	var result string

	err = querier.Query(context.Background(), &result, query, args)
	assert.NoError(t, err)

	err = tx.Commit()
	assert.NoError(t, err)
	assert.Equal(t, "Jane Doe", result)
}

func TestDBQuerier_Query_NoRowsReturned(t *testing.T) {
	db, err := sqlx.Open("sqlite", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	err = createUsersTable2(db)
	assert.NoError(t, err)

	_, err = db.Exec("INSERT INTO users (name) VALUES (?)", "Alice")
	assert.NoError(t, err)

	querier := NewDBQuerier(db, func(ctx context.Context) *sqlx.Tx { return nil })

	query := "SELECT name FROM users WHERE name = :name"
	args := map[string]any{"name": "Bob"}
	var result []string

	err = querier.Query(context.Background(), &result, query, args)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestDBQuerier_Query_ErrorDuringScan(t *testing.T) {
	db, err := sqlx.Open("sqlite", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	err = createUsersTable2(db)
	assert.NoError(t, err)

	_, err = db.Exec("INSERT INTO users (name) VALUES (?)", "Charlie")
	assert.NoError(t, err)

	querier := NewDBQuerier(db, func(ctx context.Context) *sqlx.Tx { return nil })

	query := "SELECT name FROM users WHERE name = :name"
	args := map[string]any{"name": "Charlie"}
	var result map[string]any

	err = querier.Query(context.Background(), &result, query, args)
	assert.Error(t, err)
}

func TestDBQuerier_Query_RowsErr(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	query := "SELECT name FROM users WHERE name = :name"
	args := map[string]any{"name": "Jane Doe"}

	mock.ExpectQuery(query).
		WithArgs(args["name"]).
		WillReturnRows(sqlmock.NewRows([]string{"name"}))

	mock.ExpectQuery(query).
		WithArgs(args["name"]).
		WillReturnError(errors.New("simulated rows error"))

	querier := NewDBQuerier(sqlxDB, func(ctx context.Context) *sqlx.Tx { return nil })

	var result []string

	err = querier.Query(context.Background(), &result, query, args)
	assert.Error(t, err)
}

func TestDBQuerier_Query_RowsErr2(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	querier := NewDBQuerier(sqlxDB, func(ctx context.Context) *sqlx.Tx { return nil })

	query := "SELECT name FROM users"
	args := map[string]any{}

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("John Doe").
		RowError(0, errors.New("simulated rows error"))

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WillReturnRows(rows)

	var result string
	err = querier.Query(context.Background(), &result, query, args)

	assert.Error(t, err)
	assert.Equal(t, "simulated rows error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}
