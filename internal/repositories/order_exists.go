package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/repositories/helpers"
	"go.uber.org/zap"
)

type OrderExistsRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) *sqlx.Tx
}

func NewOrderExistsRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
) *OrderExistsRepository {
	return &OrderExistsRepository{db: db, txProvider: txProvider}
}

func (r *OrderExistsRepository) Exists(
	ctx context.Context, filter map[string]any,
) (bool, error) {
	var exists bool
	row, err := helpers.QueryRowContext(ctx, r.db, r.txProvider, orderExistsQuery, filter)
	if err != nil {
		logger.Logger.Error("Error executing query to check if order exists", zap.Error(err), zap.Any("filter", filter))
		return false, err
	}
	if err := row.Scan(&exists); err != nil {
		logger.Logger.Error("Error scanning result for order existence", zap.Error(err), zap.Any("filter", filter))
		return false, err
	}
	return exists, nil
}

const orderExistsQuery = `
	SELECT EXISTS(
		SELECT 1 FROM orders WHERE order_id = :order_id
	)
`
