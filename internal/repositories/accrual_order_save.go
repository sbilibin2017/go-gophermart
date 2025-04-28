package repositories

import (
	"context"
)

type OrderSaveRepository struct {
	e Executor
}

func NewOrderSaveRepository(
	e Executor,
) *OrderSaveRepository {
	return &OrderSaveRepository{e: e}
}

func (repo *OrderSaveRepository) Save(
	ctx context.Context,
	data map[string]any,
) error {
	return repo.e.Exec(ctx, saveAccrualOrderQuery, data)
}

const saveAccrualOrderQuery = `
	INSERT INTO accrual_order (number, accrual, status, created_at, updated_at)
	VALUES (:number, :accrual, :status, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	ON CONFLICT (number)
	DO UPDATE SET accrual = :accrual, status = :status, updated_at = CURRENT_TIMESTAMP;
`
