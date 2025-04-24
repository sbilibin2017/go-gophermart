package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
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

func (r *OrderExistsRepository) Exists(ctx context.Context, number string) (bool, error) {
	var exists bool
	err := query(ctx, r.db, r.txProvider, &exists, orderExistsByIDQuery, number)
	if err != nil {
		return false, err
	}
	return exists, nil
}

const orderExistsByIDQuery = `SELECT EXISTS (SELECT 1 FROM orders WHERE number = ?)`
