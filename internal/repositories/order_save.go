package repositories

import (
	"context"
)

type OrderExecutor interface {
	Execute(
		ctx context.Context,
		query string,
		argMap map[string]any,
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
	ctx context.Context, orderID string, status string, accrual float64,
) error {
	argMap := map[string]any{
		"order_id": orderID,
		"status":   status,
		"accrual":  accrual,
	}
	err := r.e.Execute(ctx, orderSaveQuery, argMap)
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
