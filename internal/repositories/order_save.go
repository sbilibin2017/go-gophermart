package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
)

type OrderSaveRepository struct {
	db *sqlx.DB
}

func NewOrderSaveRepository(db *sqlx.DB) *OrderSaveRepository {
	return &OrderSaveRepository{db: db}
}

func (repo *OrderSaveRepository) Save(
	ctx context.Context, order *domain.Order,
) error {
	return exec(ctx, repo.db, orderSaveQuery, order)
}

const orderSaveQuery = `
INSERT INTO orders (number, login, status)
	VALUES (:number, :login, :status)
ON CONFLICT (number) DO UPDATE
	SET login = EXCLUDED.login,
	    status = EXCLUDED.status
`
