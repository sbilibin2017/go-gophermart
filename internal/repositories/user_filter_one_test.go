package repositories

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func setupUserFilterOneDB() (*sqlx.DB, func(), error) {
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

func TestUserFilterOneRepository_FilterOne_Found(t *testing.T) {
	db, tearDown, err := setupUserFilterOneDB()
	if err != nil {
		t.Fatalf("Error setting up SQLite DB: %v", err)
	}
	defer tearDown()

	txProvider := func(ctx context.Context) (*sqlx.Tx, error) {
		return nil, nil
	}

	repo := NewUserFilterOneRepository(db, txProvider)

	_, err = db.Exec("INSERT INTO users (login, password) VALUES ($1, $2)", "user1", "password1")
	assert.NoError(t, err)

	user, err := repo.FilterOne(context.Background(), "user1")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "user1", user.Login)
	assert.Equal(t, "password1", user.Password)
}

func TestUserFilterOneRepository_FilterOne_Error(t *testing.T) {
	db, tearDown, err := setupUserFilterOneDB()
	if err != nil {
		t.Fatalf("Error setting up SQLite DB: %v", err)
	}
	defer tearDown()

	txProvider := func(ctx context.Context) (*sqlx.Tx, error) {
		return nil, nil
	}

	repo := NewUserFilterOneRepository(db, txProvider)

	// Ошибка при выполнении запроса
	_, err = db.Exec("DROP TABLE users") // Удаляем таблицу, чтобы возникла ошибка
	assert.NoError(t, err)

	user, err := repo.FilterOne(context.Background(), "user1")
	assert.Error(t, err)
	assert.Nil(t, user)
}
