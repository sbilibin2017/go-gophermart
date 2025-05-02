package repositories

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func setupUserSaveDB() (*sqlx.DB, func(), error) {
	db, err := sqlx.Connect("sqlite", ":memory:")
	if err != nil {
		return nil, nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			login TEXT PRIMARY KEY,
			password TEXT NOT NULL
		)
	`)
	if err != nil {
		return nil, nil, err
	}
	tearDown := func() {
		db.Close()
	}
	return db, tearDown, nil
}

func TestUserSaveRepository_Save_Insert(t *testing.T) {
	db, tearDown, err := setupUserSaveDB()
	if err != nil {
		t.Fatalf("Error setting up SQLite DB: %v", err)
	}
	defer tearDown()
	txProvider := func(ctx context.Context) (*sqlx.Tx, error) {
		return nil, nil
	}
	repo := NewUserSaveRepository(db, txProvider)
	err = repo.Save(context.Background(), "user1", "password1")
	assert.NoError(t, err)
	var password string
	err = db.Get(&password, "SELECT password FROM users WHERE login = ?", "user1")
	assert.NoError(t, err)
	assert.Equal(t, "password1", password)
}

func TestUserSaveRepository_Save_Update(t *testing.T) {
	db, tearDown, err := setupUserSaveDB()
	if err != nil {
		t.Fatalf("Error setting up SQLite DB: %v", err)
	}
	defer tearDown()
	txProvider := func(ctx context.Context) (*sqlx.Tx, error) {
		return nil, nil
	}
	repo := NewUserSaveRepository(db, txProvider)
	err = repo.Save(context.Background(), "user1", "password1")
	assert.NoError(t, err)
	err = repo.Save(context.Background(), "user1", "newpassword")
	assert.NoError(t, err)
	var password string
	err = db.Get(&password, "SELECT password FROM users WHERE login = ?", "user1")
	assert.NoError(t, err)
	assert.Equal(t, "newpassword", password)
}

func TestUserSaveRepository_Save_Conflict(t *testing.T) {
	db, tearDown, err := setupUserSaveDB()
	if err != nil {
		t.Fatalf("Error setting up SQLite DB: %v", err)
	}
	defer tearDown()
	txProvider := func(ctx context.Context) (*sqlx.Tx, error) {
		return nil, nil
	}
	repo := NewUserSaveRepository(db, txProvider)
	err = repo.Save(context.Background(), "user1", "password1")
	assert.NoError(t, err)
	err = repo.Save(context.Background(), "user1", "updatedpassword")
	assert.NoError(t, err)
	var password string
	err = db.Get(&password, "SELECT password FROM users WHERE login = ?", "user1")
	assert.NoError(t, err)
	assert.Equal(t, "updatedpassword", password)
}
