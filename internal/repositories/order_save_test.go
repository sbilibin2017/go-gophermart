package repositories

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderSaveRepository_Save(t *testing.T) {
	setupQuery := `CREATE TABLE orders (
		number BIGINT PRIMARY KEY,
		status TEXT NOT NULL,
		accrual DOUBLE PRECISION NOT NULL,		
		updated_at TIMESTAMP NOT NULL DEFAULT now()
	)`
	db, cleanup := setupDB(t, setupQuery)
	defer cleanup()

	repo := NewOrderSaveRepository(db)

	order := &types.OrderDB{
		Number:  12345,
		Status:  "new",
		Accrual: 100.5,
	}

	t.Run("Insert new order", func(t *testing.T) {
		err := repo.Save(context.Background(), order)
		require.NoError(t, err)

		var accrual float64
		err = db.Get(&accrual, `SELECT accrual FROM orders WHERE number = $1`, order.Number)
		require.NoError(t, err)
		assert.Equal(t, 100.5, accrual)
	})

	t.Run("Update existing order", func(t *testing.T) {
		order.Status = "processed"
		order.Accrual = 200.75
		err := repo.Save(context.Background(), order)
		require.NoError(t, err)

		var accrual float64
		err = db.Get(&accrual, `SELECT accrual FROM orders WHERE number = $1`, order.Number)
		require.NoError(t, err)
		assert.Equal(t, 200.75, accrual)
	})
}

func TestOrderSaveRepository_Save_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewOrderSaveRepository(sqlxDB)

	order := &types.OrderDB{
		Number:  42,
		Status:  "new",
		Accrual: 999.99,
	}

	mock.ExpectExec(`INSERT INTO orders`).WillReturnError(errors.New("mocked db error"))

	err = repo.Save(context.Background(), order)

	require.Error(t, err)
	assert.EqualError(t, err, "mocked db error")
	require.NoError(t, mock.ExpectationsWereMet())
}
