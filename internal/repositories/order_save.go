package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type OrderSaveRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, bool)
}

func NewOrderSaveRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, bool),
) *OrderSaveRepository {
	return &OrderSaveRepository{db: db, txProvider: txProvider}
}

func (r *OrderSaveRepository) Save(ctx context.Context, order map[string]any) error {
	err := execContextNamed(ctx, r.db, r.txProvider, orderSaveQuery, order)
	if err != nil {
		return err
	}
	return nil
}

const orderSaveQuery = `
	INSERT INTO orders (order_id, status, accrual, created_at, updated_at)
	VALUES (:order_id, :status, :accrual, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	ON CONFLICT (order_id) DO UPDATE
	SET status = EXCLUDED.status,
		accrual = EXCLUDED.accrual,
		updated_at = CURRENT_TIMESTAMP
`
