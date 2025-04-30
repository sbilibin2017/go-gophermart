package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type OrderExistsRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, error)
}

func NewOrderExistsRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
) *OrderExistsRepository {
	return &OrderExistsRepository{db: db, txProvider: txProvider}
}

func (r *OrderExistsRepository) Exists(ctx context.Context, number string) (bool, error) {
	e := getExecutor(ctx, r.db, r.txProvider)
	var exists bool
	err := sqlx.GetContext(ctx, e, &exists, orderExistsQuery, number)

	return exists, err
}

const orderExistsQuery = `
SELECT EXISTS (
	SELECT 1 FROM orders WHERE number = $1
)
`
