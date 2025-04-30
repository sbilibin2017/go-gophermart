package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderFindRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, error)
}

func NewOrderFindRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
) *OrderFindRepository {
	return &OrderFindRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *OrderFindRepository) Find(ctx context.Context, number string) (*types.Order, error) {
	e := getExecutor(ctx, r.db, r.txProvider)
	var order types.Order
	err := sqlx.GetContext(ctx, e, &order, orderFindQuery, number)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

const orderFindQuery = `
SELECT number, accrual, status
FROM orders
WHERE number = $1
LIMIT 1
`
