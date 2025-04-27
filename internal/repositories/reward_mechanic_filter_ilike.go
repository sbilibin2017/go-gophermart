package repositories

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type RewardMechanicFilterILikeRepository struct {
	db *sqlx.DB
}

func NewRewardMechanicFilterILikeRepository(db *sqlx.DB) *RewardMechanicFilterILikeRepository {
	return &RewardMechanicFilterILikeRepository{db: db}
}

func (repo *RewardMechanicFilterILikeRepository) FilterILikeByDescription(
	ctx context.Context,
	description string,
	fields []string,
) (map[string]any, error) {
	q, args := buildILikeQuery(fields, description)
	var result map[string]any
	err := query(ctx, repo.db, q, &result, args)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func buildILikeQuery(fields []string, description string) (string, map[string]any) {
	columns := buildColumnsString(fields)
	query := fmt.Sprintf(
		rewardMechanicFilterILikeQuery,
		columns,
	)
	args := map[string]any{
		"description": "%" + description + "%",
	}
	return query, args
}

const rewardMechanicFilterILikeQuery = `
	SELECT %s FROM accrual_reward_mechanic WHERE match ILIKE :description
`
