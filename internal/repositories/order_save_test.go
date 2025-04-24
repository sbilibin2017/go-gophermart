package repositories

import (
	"context"
	"testing"
	"time"

	_ "github.com/glebarez/sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupOrderSaveTest(t *testing.T) *OrderSaveRepository {
	db, err := sqlx.Open("sqlite", ":memory:")
	require.NoError(t, err)
	_, err = db.Exec(`
		CREATE TABLE orders (
			number TEXT PRIMARY KEY,
			status TEXT,
			accrual INTEGER,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		);
	`)
	require.NoError(t, err)
	return NewOrderSaveRepository(db, nil)
}

func TestOrderSaveRepository_Save(t *testing.T) {
	repo := setupOrderSaveTest(t)

	t.Run("Save new order", func(t *testing.T) {
		order := map[string]any{
			"number":  "order123",
			"status":  "NEW",
			"accrual": int64(100),
		}
		err := repo.Save(context.Background(), order)
		require.NoError(t, err)
		var number, status string
		var accrual int64
		var createdAt, updatedAt time.Time
		err = repo.db.QueryRowx("SELECT number, status, accrual, created_at, updated_at FROM orders WHERE number = ?", "order123").
			Scan(&number, &status, &accrual, &createdAt, &updatedAt)
		require.NoError(t, err)
		assert.Equal(t, "order123", number)
		assert.Equal(t, "NEW", status)
		assert.Equal(t, int64(100), accrual)
		assert.False(t, createdAt.IsZero())
		assert.False(t, updatedAt.IsZero())
	})

	t.Run("Update existing order", func(t *testing.T) {
		order := map[string]any{
			"number":  "order123",
			"status":  "NEW",
			"accrual": int64(100),
		}
		err := repo.Save(context.Background(), order)
		require.NoError(t, err)
		order["status"] = "UPDATED"
		order["accrual"] = int64(200)
		err = repo.Save(context.Background(), order)
		require.NoError(t, err)
		var status string
		var accrual int64
		err = repo.db.QueryRowx("SELECT status, accrual FROM orders WHERE number = ?", "order123").
			Scan(&status, &accrual)
		require.NoError(t, err)
		assert.Equal(t, "UPDATED", status)
		assert.Equal(t, int64(200), accrual)
	})

	t.Run("Update order with partial fields", func(t *testing.T) {
		initial := map[string]any{
			"number":  "order456",
			"status":  "NEW",
			"accrual": int64(50),
		}
		err := repo.Save(context.Background(), initial)
		require.NoError(t, err)
		partial := map[string]any{
			"number":  "order456",
			"status":  "PARTIALLY_UPDATED",
			"accrual": nil, // Partial update: accrual set to nil
		}
		err = repo.Save(context.Background(), partial)
		require.NoError(t, err)
		var status string
		var accrual interface{}
		err = repo.db.QueryRowx("SELECT status, accrual FROM orders WHERE number = ?", "order456").
			Scan(&status, &accrual)
		require.NoError(t, err)
		assert.Equal(t, "PARTIALLY_UPDATED", status)
		assert.Nil(t, accrual)
	})
}
