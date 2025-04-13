package repositories

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type OrderSaveRepository struct {
	db *sqlx.DB
}

func NewOrderSaveRepository(db *sqlx.DB) *OrderSaveRepository {
	return &OrderSaveRepository{db: db}
}

var saveOrderQuery = `
	INSERT INTO orders (order_id, status, accrual, created_at, updated_at)
	VALUES (:order_id, :status, :accrual, now(), now())
	ON CONFLICT (order_id) DO UPDATE 
	SET 
		status = EXCLUDED.status,
		accrual = EXCLUDED.accrual,
		updated_at = now()
`

type Status string

const (
	StatusNew        Status = "NEW"        // заказ загружен, но не попал в обработку
	StatusProcessing Status = "PROCESSING" // вознаграждение рассчитывается
	StatusInvalid    Status = "INVALID"    // система отказала в расчёте
	StatusProcessed  Status = "PROCESSED"  // расчёт успешно завершён
)

type OrderSave struct {
	OrderID uint64  `db:"order_id"`
	Status  Status  `db:"status"`
	Accrual float64 `db:"accrual"`
}

func (r *OrderSaveRepository) Save(ctx context.Context, order *OrderSave) error {
	_, err := r.db.NamedExecContext(ctx, saveOrderQuery, order)
	return err
}
