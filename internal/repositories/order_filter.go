package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderFilterOneRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, error)
}

func NewOrderFilterOneRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
) *OrderFilterOneRepository {
	return &OrderFilterOneRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *OrderFilterOneRepository) FilterOne(
	ctx context.Context, number string,
) (*types.OrderDB, error) {
	var order types.OrderDB
	err := getContext(
		ctx,
		r.db,
		r.txProvider,
		orderFilterOneQuery,
		&order,
		number,
	)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

const orderFilterOneQuery = `
	SELECT number, status, accrual 
	FROM orders 
	WHERE number = $1
`
