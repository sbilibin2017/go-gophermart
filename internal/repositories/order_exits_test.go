package repositories

import (
	"context"
	"testing"
	"time"

	_ "github.com/glebarez/sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setupOrderExistsTest(t *testing.T) *OrderExistsRepository {
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
	return NewOrderExistsRepository(db, nil)
}

func TestOrderExists(t *testing.T) {
	repo := setupOrderExistsTest(t)

	t.Run("Order exists", func(t *testing.T) {
		exists, err := repo.Exists(context.Background(), "12345")
		assert.NoError(t, err)
		assert.True(t, exists, "Order should exist")
	})

	t.Run("Order does not exist", func(t *testing.T) {
		exists, err := repo.Exists(context.Background(), "nonexistent_order_id")
		assert.NoError(t, err)
		assert.False(t, exists, "Order should not exist")
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
		repo := NewOrderExistsRepository(db, nil)
		exists, err := repo.Exists(context.Background(), "12345")
		assert.NoError(t, err)
		assert.False(t, exists, "Order should not exist in an empty database")
	})
}
