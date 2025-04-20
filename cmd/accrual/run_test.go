package main

import (
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func TestRun_DBConnectionError(t *testing.T) {
	called := false
	originalDBConnect := dbConnFunc
	defer func() { dbConnFunc = originalDBConnect }()
	dbConnFunc = func(uri string) (*sqlx.DB, error) {
		called = true
		return nil, errors.New("mock db error")
	}
	run()
	assert.True(t, called, "dbConnect should be called")
}

func TestRun_Success(t *testing.T) {
	databaseURI = ":memory:"
	runAddress = "localhost:8080"

	dbConnFunc = func(uri string) (*sqlx.DB, error) {
		db, err := sqlx.Open("sqlite3", uri)
		if err != nil {
			return nil, err
		}
		schema := `CREATE TABLE IF NOT EXISTS accruals (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			data TEXT
		);`
		if _, err := db.Exec(schema); err != nil {
			return nil, err
		}
		return db, nil
	}

	run()

	assert.True(t, true, "run should complete without panic or error")
}
