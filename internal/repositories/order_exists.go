package repositories

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

type OrderExistsFilterDB struct {
	Number uint64 `db:"number"`
}

type OrderExistRepository struct {
	db *sqlx.DB
}

func NewOrderExistRepository(db *sqlx.DB) *OrderExistRepository {
	return &OrderExistRepository{db: db}
}

func (r *OrderExistRepository) Exists(
	ctx context.Context,
	filter *OrderExistsFilterDB,
) (bool, error) {
	logger.Logger.Debugf("Checking if order exists: number=%d", filter.Number)

	var exists bool
	err := r.db.GetContext(ctx, &exists, orderExistsQuery, filter.Number)
	if err != nil {
		logger.Logger.Errorf("Failed to check order existence for number=%d: %v", filter.Number, err)
		return false, err
	}

	logger.Logger.Debugf("Order existence result for number=%d: %t", filter.Number, exists)
	return exists, nil
}

var orderExistsQuery = `SELECT EXISTS(SELECT 1 FROM orders WHERE number = $1)`
