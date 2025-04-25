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

func (r *OrderExistsRepository) Exists(ctx context.Context, order *OrderExistsNumber) (bool, error) {
	var exists bool
	err := query(ctx, r.db, r.txProvider, &exists, orderExistsByIDQuery, order)
	if err != nil {
		return false, err
	}
	return exists, nil
}

type OrderExistsNumber struct {
	Number string `db:"number"`
}

const orderExistsByIDQuery = `SELECT EXISTS (SELECT 1 FROM orders WHERE number = :number)`
