package repositories

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderExistsRepository_Exists(t *testing.T) {
	setupQuery := `CREATE TABLE orders (
		order_id BIGINT PRIMARY KEY
	)`
	db, cleanup := setupTestPostgres(t, setupQuery)
	defer cleanup()

	repo := NewOrderExistRepository(db)

	existingorderID := uint64(12345)
	_, err := db.Exec(`INSERT INTO orders (order_id) VALUES ($1)`, existingorderID)
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
			exists, err := repo.Exists(context.Background(), &OrderExistsID{OrderID: tc.orderID})
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

	orderID := &OrderExistsID{OrderID: 42}

	mock.ExpectQuery(``).
		WithArgs(orderID.OrderID).
		WillReturnError(errors.New("mocked db error"))

	exists, err := repo.Exists(context.Background(), orderID)

	require.Error(t, err)
	require.False(t, exists)
	require.EqualError(t, err, "mocked db error")

	require.NoError(t, mock.ExpectationsWereMet())
}
