package repositories

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"
)

type OrderGetByIDRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, bool)
}

func NewOrderGetByIDRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, bool),
) *OrderGetByIDRepository {
	return &OrderGetByIDRepository{db: db, txProvider: txProvider}
}

func (r *OrderGetByIDRepository) GetByID(
	ctx context.Context, orderID string, fields []string,
) (map[string]any, error) {
	query := buildGetOrderInfoByIDQuery(fields)
	argMap := map[string]any{
		"order_id": orderID,
	}
	var result map[string]any
	err := getContextNamed(ctx, r.db, r.txProvider, &result, query, argMap)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func buildGetOrderInfoByIDQuery(fields []string) string {
	var sb strings.Builder
	if len(fields) > 0 {
		sb.WriteString("SELECT ")
		sb.WriteString(strings.Join(fields, ", "))
	} else {
		sb.WriteString("SELECT *")
	}
	sb.WriteString(" FROM orders WHERE order_id = :order_id")
	return sb.String()
}
