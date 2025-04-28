package repositories

import (
	"context"
)

type OrderExistsRepository struct {
	q Querier
}

func NewOrderExistsRepository(
	q Querier,
) *OrderExistsRepository {
	return &OrderExistsRepository{q: q}
}

func (repo *OrderExistsRepository) Exists(
	ctx context.Context,
	filter map[string]any,
) (bool, error) {
	var exists bool
	err := repo.q.Query(ctx, accrualOrderRegisterExistsQuery, &exists, filter["number"])
	if err != nil {
		return false, err
	}
	return exists, nil
}

const accrualOrderRegisterExistsQuery = `
	SELECT EXISTS (
		SELECT 1
		FROM accrual_order
		WHERE number = $1
	)
`
