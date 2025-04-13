package repositories

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type OrderExistsRepository struct {
	db *sqlx.DB
}

func NewOrderExistsRepository(db *sqlx.DB) *OrderExistsRepository {
	return &OrderExistsRepository{db: db}
}

// Обновлённый запрос, учитывающий правильные имена столбцов
var orderExistsQuery = `SELECT EXISTS(SELECT 1 FROM orders WHERE order_id = $1)`

// Структура для передачи ID заказа
type OrderExistsID struct {
	OrderID uint64 `db:"order_id"`
}

// Метод для проверки существования заказа
func (r *OrderExistsRepository) Exists(ctx context.Context, orderId *OrderExistsID) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, orderExistsQuery, orderId.OrderID)
	if err != nil {
		return false, err
	}
	return exists, nil
}
