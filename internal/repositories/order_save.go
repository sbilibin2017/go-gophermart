package repositories

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/queries"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderSaveRepository struct {
	db *sqlx.DB
}

func NewOrderSaveRepository(db *sqlx.DB) *OrderSaveRepository {
	return &OrderSaveRepository{db: db}
}

func (r *OrderSaveRepository) Save(ctx context.Context, order *types.OrderDB) error {
	_, err := r.db.NamedExecContext(ctx, queries.OrderSaveQuery, order)
	return err
}
