package repositories

import (
	"context"
)

type OrderExistsQuerier interface {
	Query(
		ctx context.Context,
		dest any,
		query string,
		argMap map[string]any,
	) error
}

type OrderExistsRepository struct {
	q OrderExistsQuerier
}

func NewOrderExistsRepository(
	q OrderExistsQuerier,
) *OrderExistsRepository {
	return &OrderExistsRepository{q: q}
}

func (r *OrderExistsRepository) Exists(
	ctx context.Context, orderID string,
) (bool, error) {
	argMap := map[string]any{
		"order_id": orderID,
	}
	var exists bool
	err := r.q.Query(ctx, &exists, orderExistsByIDQuery, argMap)
	if err != nil {
		return false, err
	}
	return exists, nil
}

var orderExistsByIDQuery = `SELECT EXISTS(SELECT 1 FROM orders WHERE order_id = :order_id)`
