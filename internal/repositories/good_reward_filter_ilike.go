package repositories

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"
)

type OrderRewardGetByDescriptionRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, bool)
}

func NewOrderRewardFilterILikeRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, bool),
) *OrderRewardGetByDescriptionRepository {
	return &OrderRewardGetByDescriptionRepository{db: db, txProvider: txProvider}
}

func (r *OrderRewardGetByDescriptionRepository) FilterILike(
	ctx context.Context, match string, fields []string,
) (map[string]any, error) {
	query := buildFilterQuery(fields)
	argMap := map[string]any{
		"match": "%" + match + "%",
	}
	var result map[string]any
	err := getContextNamed(ctx, r.db, r.txProvider, &result, query, argMap)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func buildFilterQuery(fields []string) string {
	var sb strings.Builder
	if len(fields) > 0 {
		sb.WriteString("SELECT ")
		sb.WriteString(strings.Join(fields, ", "))
	} else {
		sb.WriteString("SELECT *")
	}
	sb.WriteString(" FROM rewards WHERE reward_id ILIKE :match")
	return sb.String()
}
