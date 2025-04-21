package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/repositories/helpers"
)

type OrderExistsRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) *sqlx.Tx
}

func NewOrderExistsRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
) *OrderExistsRepository {
	return &OrderExistsRepository{db: db, txProvider: txProvider}
}

func (r *OrderExistsRepository) Exists(
	ctx context.Context, filter map[string]any,
) (bool, error) {
	row, err := helpers.QueryRowContext(ctx, r.db, r.txProvider, orderExistsQuery, filter)
	if err != nil {
		return false, err
	}
	exists, err := helpers.Scan[bool](row)
	if err != nil {
		return false, err
	}
	return exists, nil
}

const orderExistsQuery = `
	SELECT EXISTS(
		SELECT 1 FROM orders WHERE order_id = :order_id
	)
`
