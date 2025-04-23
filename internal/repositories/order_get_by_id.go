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
		argMap map[string]any,
	) error
}

type OrderGetByIDRepository struct {
	q OrderGetByIDQuerier
}

func NewOrderGetByIDRepository(
	q OrderGetByIDQuerier,
) *OrderGetByIDRepository {
	return &OrderGetByIDRepository{q: q}
}

func (r *OrderGetByIDRepository) GetByID(
	ctx context.Context, orderID string, fields []string,
) (*types.OrderDB, error) {
	query := buildGetOrderInfoByIDQuery(fields)
	argMap := map[string]any{
		"order_id": orderID,
	}
	result := new(types.OrderDB)
	err := r.q.Query(ctx, result, query, argMap)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func buildGetOrderInfoByIDQuery(fields []string) string {
	fieldsQuery := "*"
	if len(fields) > 0 {
		fieldsQuery = strings.Join(fields, ", ")
	}
	query := fmt.Sprintf(orderGetByIDQuery, fieldsQuery)
	return query
}

const orderGetByIDQuery = "SELECT %s FROM orders WHERE order_id = :order_id"
