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

func TestOrderExistsRepository_Exists(t *testing.T) {
	setupQuery := `CREATE TABLE orders (
		number BIGINT PRIMARY KEY
	)`
	db, cleanup := setupDB(t, setupQuery)
	defer cleanup()

	repo := NewOrderExistRepository(db)

	existingorderID := uint64(12345)
	_, err := db.Exec(`INSERT INTO orders (number) VALUES ($1)`, existingorderID)
	require.NoError(t, err)

	testCases := []struct {
		name       string
		orderID    uint64
		wantExists bool
	}{
		{
			name:       "Order exists",
			orderID:    existingorderID,
			wantExists: true,
		},
		{
			name:       "Order does not exist",
			orderID:    99999,
			wantExists: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			exists, err := repo.Exists(context.Background(), &types.OrderExistsFilter{Number: tc.orderID})
			require.NoError(t, err)
			assert.Equal(t, tc.wantExists, exists)
		})
	}
}

func TestOrderExistsRepository_Exists_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewOrderExistRepository(sqlxDB)

	orderID := &types.OrderExistsFilter{Number: 42}

	mock.ExpectQuery(``).
		WithArgs(orderID.Number).
		WillReturnError(errors.New("mocked db error"))

	exists, err := repo.Exists(context.Background(), orderID)

	require.Error(t, err)
	require.False(t, exists)
	require.EqualError(t, err, "mocked db error")

	require.NoError(t, mock.ExpectationsWereMet())
}
