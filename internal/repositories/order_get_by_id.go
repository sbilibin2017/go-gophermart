package repositories

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type OrderGetRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) *sqlx.Tx
}

func NewOrderGetRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
) *OrderGetRepository {
	return &OrderGetRepository{db: db, txProvider: txProvider}
}

func (r *OrderGetRepository) Get(
	ctx context.Context, number string, fields []string,
) (map[string]any, error) {
	columns := getColumns(fields)
	q := fmt.Sprintf(orderGetQuery, columns)
	var order map[string]any
	err := query(ctx, r.db, r.txProvider, &order, q, number)
	if err != nil {
		return nil, err
	}
	return order, nil
}

const orderGetQuery = "SELECT %s FROM orders WHERE number = ?"
