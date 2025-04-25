package repositories

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type OrderGetRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) *sqlx.Tx
}

func NewOrderGetRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
) *OrderGetRepository {
	return &OrderGetRepository{db: db, txProvider: txProvider}
}

func (r *OrderGetRepository) Get(
	ctx context.Context, filter *OrderGetFilter,
) (*OrderGetDB, error) {
	var order OrderGetDB
	err := query(ctx, r.db, r.txProvider, &order, orderGetQuery, filter)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

type OrderGetFilter struct {
	Number string `db:"number"`
}

type OrderGetDB struct {
	Number    string    `db:"number"`
	Status    string    `db:"status"`
	Accrual   *int64    `db:"accrual"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

const orderGetQuery = "SELECT * FROM orders WHERE number = :number"
