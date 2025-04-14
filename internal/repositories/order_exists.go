package repositories

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/sbilibin2017/go-gophermart/internal/queries"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderExistRepository struct {
	db *sqlx.DB
}

func NewOrderExistRepository(db *sqlx.DB) *OrderExistRepository {
	return &OrderExistRepository{db: db}
}

// Метод для проверки существования заказа
func (r *OrderExistRepository) Exists(
	ctx context.Context,
	filter *types.OrderExistsFilter,
) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, queries.OrderExistsQuery, filter.Number)
	if err != nil {
		return false, err
	}
	return exists, nil
}
