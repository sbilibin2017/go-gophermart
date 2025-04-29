package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
)

type OrderSaveRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, bool)
}

func NewOrderSaveRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, bool),
) *OrderSaveRepository {
	return &OrderSaveRepository{db: db, txProvider: txProvider}
}

func (repo *OrderSaveRepository) Save(
	ctx context.Context, order *domain.Order,
) error {
	// Передаем txProvider в exec для обработки транзакции
	return exec(ctx, repo.db, repo.txProvider, orderSaveQuery, order)
}

const orderSaveQuery = `
INSERT INTO orders (number, login, status)
	VALUES (:number, :login, :status)
ON CONFLICT (number) DO UPDATE
	SET login = EXCLUDED.login,
	    status = EXCLUDED.status
`
