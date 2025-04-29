package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
)

type OrderListRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, bool)
}

func NewOrderListRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, bool),
) *OrderListRepository {
	return &OrderListRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *OrderListRepository) ListByLoginOrderedDesc(
	ctx context.Context, login string,
) ([]*domain.Order, error) {
	params := map[string]any{"login": login}
	var orders []*domain.Order

	// Передаем txProvider в queryStructs для обработки транзакции
	err := queryStructs(ctx, r.db, r.txProvider, orderListQuery, params, &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

const orderListQuery = `
SELECT number, login, status, accrual, uploaded_at
FROM orders
WHERE login = :login
ORDER BY uploaded_at DESC
`
