package repositories

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type OrderExistRepository struct {
	db *sqlx.DB
}

func NewOrderExistRepository(db *sqlx.DB) *OrderExistRepository {
	return &OrderExistRepository{db: db}
}

// Обновлённый запрос, учитывающий правильные имена столбцов
var orderExistsQuery = `SELECT EXISTS(SELECT 1 FROM orders WHERE order_id = $1)`

// Структура для передачи ID заказа
type OrderExistsID struct {
	OrderID uint64 `db:"order_id"`
}

// Метод для проверки существования заказа
func (r *OrderExistRepository) Exists(ctx context.Context, orderID *OrderExistsID) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, orderExistsQuery, orderID.OrderID)
	if err != nil {
		return false, err
	}
	return exists, nil
}
