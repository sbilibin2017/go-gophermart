package repositories

import (
	"context"
	"testing"
	"time"

	_ "github.com/glebarez/sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setupOrderGetTest(t *testing.T) *OrderGetRepository {
	db, err := sqlx.Open("sqlite", ":memory:")
	assert.NoError(t, err)
	_, err = db.Exec(`
		CREATE TABLE orders (
			number TEXT PRIMARY KEY,
			status TEXT,
			accrual INTEGER,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		);
	`)
	assert.NoError(t, err)
	_, err = db.Exec(`
		INSERT INTO orders (number, status, accrual, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?);
	`, "12345", "NEW", nil, time.Now(), time.Now())
	assert.NoError(t, err)
	return NewOrderGetRepository(db, nil)
}

func TestOrderGet(t *testing.T) {
	repo := setupOrderGetTest(t)

	t.Run("Get specific fields", func(t *testing.T) {
		order, err := repo.Get(context.Background(), "12345", []string{"number", "status"})
		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, "12345", order["number"])
		assert.Equal(t, "NEW", order["status"])
	})

	t.Run("Get all fields", func(t *testing.T) {
		order, err := repo.Get(context.Background(), "12345", []string{})
		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, "12345", order["number"])
		assert.Equal(t, "NEW", order["status"])
		assert.Contains(t, order, "created_at")
		assert.Contains(t, order, "updated_at")
	})

	t.Run("Order does not exist", func(t *testing.T) {
		order, err := repo.Get(context.Background(), "nonexistent_order", []string{})
		assert.NoError(t, err)
		assert.Empty(t, order)
	})

	t.Run("Empty database", func(t *testing.T) {
		db, err := sqlx.Open("sqlite", ":memory:")
		assert.NoError(t, err)
		_, err = db.Exec(`
			CREATE TABLE orders (
				number TEXT PRIMARY KEY,
				status TEXT,
				accrual INTEGER,
				created_at TIMESTAMP,
				updated_at TIMESTAMP
			);
		`)
		assert.NoError(t, err)
		repo := NewOrderGetRepository(db, nil)
		order, err := repo.Get(context.Background(), "12345", []string{})
		assert.NoError(t, err)
		assert.Empty(t, order)
	})
}
