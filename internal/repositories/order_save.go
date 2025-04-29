package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type OrderSaveRepository struct {
	db *sqlx.DB
}

func NewOrderSaveRepository(db *sqlx.DB) *OrderSaveRepository {
	return &OrderSaveRepository{db: db}
}

func (repo *OrderSaveRepository) Save(
	ctx context.Context, number string, login string,
) error {
	args := map[string]any{
		"number": number,
		"login":  login,
	}
	return exec(ctx, repo.db, orderSaveQuery, args)
}

const orderSaveQuery = `
INSERT INTO orders (number, login)
	VALUES (:number, :login)
ON CONFLICT (number) DO UPDATE
	SET login = EXCLUDED.login
`
