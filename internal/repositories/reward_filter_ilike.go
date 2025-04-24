package repositories

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type RewardFilterILikeRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) *sqlx.Tx
}

func NewRewardFilterILikeRepository(db *sqlx.DB, txProvider func(ctx context.Context) *sqlx.Tx) *RewardFilterILikeRepository {
	return &RewardFilterILikeRepository{db: db, txProvider: txProvider}
}

func (r *RewardFilterILikeRepository) FilterILike(
	ctx context.Context,
	description string,
	fields []string,
) (map[string]any, error) {
	columns := getColumns(fields)
	params := map[string]any{
		"description": "%" + description + "%",
	}
	query := fmt.Sprintf(rewardFilterILikeQueryTemplate, columns)
	result, err := queryNamed(ctx, r.db, r.txProvider, query, params)
	if err != nil {
		return nil, err
	}
	return result, nil
}

const rewardFilterILikeQueryTemplate = `
	SELECT %s 
	FROM rewards 
	WHERE match ILIKE :description
	LIMIT 1
`
