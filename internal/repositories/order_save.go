package repositories

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

type status string

const (
	StatusNew        status = "NEW"
	StatusProcessing status = "PROCESSING"
	StatusInvalid    status = "INVALID"
	StatusProcessed  status = "PROCESSED"
)

type OrderSaveDB struct {
	Number  uint64  `db:"number"`
	Status  status  `db:"status"`
	Accrual float64 `db:"accrual"`
}

type OrderSaveRepository struct {
	db *sqlx.DB
}

func NewOrderSaveRepository(db *sqlx.DB) *OrderSaveRepository {
	return &OrderSaveRepository{db: db}
}

func (r *OrderSaveRepository) Save(
	ctx context.Context, order *OrderSaveDB,
) error {
	logger.Logger.Infof("Saving order: number=%d, status=%s, accrual=%.2f", order.Number, order.Status, order.Accrual)

	_, err := r.db.NamedExecContext(ctx, orderSaveQuery, order)
	if err != nil {
		logger.Logger.Errorf("Failed to save order number=%d: %v", order.Number, err)
		return err
	}

	logger.Logger.Debugf("Order number=%d saved successfully", order.Number)
	return nil
}

var orderSaveQuery = `
	INSERT INTO orders (number, status, accrual, created_at, updated_at)
	VALUES (:number, :status, :accrual, now(), now())
	ON CONFLICT (number) DO UPDATE 
	SET 
		status = EXCLUDED.status,
		accrual = EXCLUDED.accrual,
		updated_at = now()
`
