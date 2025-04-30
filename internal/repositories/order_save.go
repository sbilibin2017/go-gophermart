package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderSaveRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, error)
}

func NewOrderSaveRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
) *OrderSaveRepository {
	return &OrderSaveRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *OrderSaveRepository) Save(
	ctx context.Context, order *types.Order,
) error {
	e := getExecutor(ctx, r.db, r.txProvider)
	_, err := sqlx.NamedExecContext(ctx, e, orderUpsertQuery, order)

	return err
}

const orderUpsertQuery = `
INSERT INTO orders (number, accrual, status)
VALUES (:number, :accrual, :status)
ON CONFLICT (number)
DO UPDATE SET
	accrual = EXCLUDED.accrual,
	status = EXCLUDED.status
`
