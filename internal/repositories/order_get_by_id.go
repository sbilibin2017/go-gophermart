package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderGetByIDQuerier interface {
	Query(
		ctx context.Context,
		dest any,
		query string,
		args any,
	) error
}

type OrderGetByIDRepository struct {
	q OrderGetByIDQuerier
}

func NewOrderGetByIDRepository(q OrderGetByIDQuerier) *OrderGetByIDRepository {
	return &OrderGetByIDRepository{q: q}
}

func (r *OrderGetByIDRepository) GetByID(
	ctx context.Context,
	filter *OrderGetByIDFilter, // Передаём указатель на OrderGetByIDFilter
) (*types.OrderDB, error) {
	query := buildGetOrderByIDQuery(filter.Fields)

	result := new(types.OrderDB)
	err := r.q.Query(ctx, result, query, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type OrderGetByIDFilter struct {
	OrderID string   `db:"order_id"`
	Fields  []string // оставляем как обычный срез
}

func buildGetOrderByIDQuery(fields []string) string {
	fieldsQuery := "*"
	if len(fields) > 0 {
		fieldsQuery = strings.Join(fields, ", ")
	}
	return fmt.Sprintf(orderGetByIDQueryTemplate, fieldsQuery)
}

const orderGetByIDQueryTemplate = `
	SELECT %s FROM orders WHERE order_id = :order_id
`
