package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type AccrualOrderSaveRepository struct {
	db *sqlx.DB
}

func NewAccrualOrderSaveRepository(db *sqlx.DB) *AccrualOrderSaveRepository {
	return &AccrualOrderSaveRepository{db: db}
}

func (repo *AccrualOrderSaveRepository) Save(
	ctx context.Context,
	number string,
	accrual int64,
	status string,
) error {
	args := map[string]any{
		"number":  number,
		"accrual": accrual,
		"status":  status,
	}
	return exec(ctx, repo.db, saveAccrualOrderQuery, args)
}

const saveAccrualOrderQuery = `
	INSERT INTO accrual_order (number, accrual, status, created_at, updated_at)
	VALUES (:number, :accrual, :status, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	ON CONFLICT (number)
	DO UPDATE SET accrual = :accrual, status = :status, updated_at = CURRENT_TIMESTAMP;
`
