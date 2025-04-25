package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
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

func (r *OrderSaveRepository) Save(ctx context.Context, order *OrderSave) error {
	return command(ctx, r.db, r.txProvider, orderSaveQuery, order)
}

type OrderSave struct {
	Number  string `db:"number"`
	Status  string `db:"status"`
	Accrual *int64 `db:"accrual"`
}

const orderSaveQuery = `
	INSERT INTO orders ("number", "status", "accrual", created_at, updated_at)
	VALUES (:number, :status, :accrual, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	ON CONFLICT ("number") DO UPDATE 
	SET status = EXCLUDED.status,
		accrual = EXCLUDED.accrual,
		updated_at = CURRENT_TIMESTAMP
`
