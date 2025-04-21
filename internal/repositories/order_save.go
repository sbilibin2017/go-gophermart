package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/repositories/helpers"
)

type OrderSaveRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) *sqlx.Tx
}

func NewOrderSaveRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
) *OrderSaveRepository {
	return &OrderSaveRepository{db: db, txProvider: txProvider}
}

func (r *OrderSaveRepository) Save(
	ctx context.Context, data map[string]any,
) error {
	_, err := helpers.ExecContext(ctx, r.db, r.txProvider, orderSaveQuery, data)
	if err != nil {
		return err
	}
	return nil
}

const orderSaveQuery = `
	INSERT INTO orders (order_id, price, created_at, updated_at)
	VALUES (:order_id, :price, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	ON CONFLICT (order_id) DO UPDATE
	SET price = EXCLUDED.price,
		updated_at = CURRENT_TIMESTAMP
`
