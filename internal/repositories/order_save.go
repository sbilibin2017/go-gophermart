package repositories

import (
	"context"
)

type OrderExecutor interface {
	Execute(
		ctx context.Context,
		query string,
		args any,
	) error
}

type OrderSaveRepository struct {
	e OrderExecutor
}

func NewOrderSaveRepository(
	e OrderExecutor,
) *OrderSaveRepository {
	return &OrderSaveRepository{e: e}
}

func (r *OrderSaveRepository) Save(
	ctx context.Context, order *OrderSave,
) error {
	err := r.e.Execute(ctx, orderSaveQuery, order)
	if err != nil {
		return err
	}
	return nil
}

type OrderSave struct {
	OrderID string `db:"order_id"`
	Status  string `db:"status"`
	Accrual int64  `db:"accrual"`
}

const orderSaveQuery = `
	INSERT INTO orders (order_id, status, accrual, created_at, updated_at)
	VALUES (:order_id, :status, :accrual, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	ON CONFLICT (order_id) DO UPDATE
	SET status = EXCLUDED.status,
		accrual = EXCLUDED.accrual,
		updated_at = CURRENT_TIMESTAMP
`
