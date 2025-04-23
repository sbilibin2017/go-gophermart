package repositories

import (
	"context"
)

type OrderExistsQuerier interface {
	Query(
		ctx context.Context,
		dest any,
		query string,
		args any,
	) error
}

type OrderExistsRepository struct {
	q OrderExistsQuerier
}

func NewOrderExistsRepository(q OrderExistsQuerier) *OrderExistsRepository {
	return &OrderExistsRepository{q: q}
}

func (r *OrderExistsRepository) Exists(
	ctx context.Context,
	filter *OrderExistsFilter, // Передаём указатель на OrderExistsFilter
) (bool, error) {
	var exists bool
	err := r.q.Query(ctx, &exists, orderExistsByIDQuery, filter)
	if err != nil {
		return false, err
	}
	return exists, nil
}

type OrderExistsFilter struct {
	OrderID string `db:"order_id"`
}

const orderExistsByIDQuery = `
	SELECT EXISTS (
		SELECT 1
		FROM orders
		WHERE order_id = :order_id
	)
`
