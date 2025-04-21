package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/repositories/helpers"
)

type OrderFilterRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) *sqlx.Tx
}

func NewOrderFilterRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
) *OrderFilterRepository {
	return &OrderFilterRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *OrderFilterRepository) Filter(
	ctx context.Context, filter map[string]any, fields []string,
) (map[string]any, error) {
	query := buildOrderFilterQuery(filter, fields)
	row, err := helpers.QueryRowContext(ctx, r.db, r.txProvider, query, filter)
	if err != nil {
		return nil, err
	}
	return helpers.MapScan(row), nil
}

func buildOrderFilterQuery(filter map[string]any, fields []string) string {
	fs := "*"
	if len(fields) > 0 {
		fs = strings.Join(fields, ", ")
	}
	query := fmt.Sprintf("SELECT %s FROM order WHERE 1=1", fs)
	for key := range filter {
		query += fmt.Sprintf(" AND %s = :%s", key, key)
	}
	return query
}
