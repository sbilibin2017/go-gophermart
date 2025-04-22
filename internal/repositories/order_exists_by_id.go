package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type OrderExistsRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, bool)
}

func NewOrderExistsRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, bool),
) *OrderExistsRepository {
	return &OrderExistsRepository{db: db, txProvider: txProvider}
}

func (r *OrderExistsRepository) Exists(ctx context.Context, orderID string) (bool, error) {
	argMap := map[string]any{
		"order_id": orderID,
	}
	var exists bool
	err := getContextNamed(ctx, r.db, r.txProvider, &exists, orderExistsByIDQuery, argMap)
	if err != nil {
		return false, err
	}
	return exists, nil
}

var orderExistsByIDQuery = `SELECT EXISTS(SELECT 1 FROM orders WHERE order_id = :order_id)`
