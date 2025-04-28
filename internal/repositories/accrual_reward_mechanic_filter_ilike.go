package repositories

import (
	"context"
	"fmt"
)

type AccrualRewardMechanicFilterILikeRepository struct {
	q Querier
}

func NewAccrualRewardMechanicFilterILikeRepository(
	q Querier,
) *AccrualRewardMechanicFilterILikeRepository {
	return &AccrualRewardMechanicFilterILikeRepository{q: q}
}

func (repo *AccrualRewardMechanicFilterILikeRepository) FilterILike(
	ctx context.Context,
	s string,
	fields []string,
) (map[string]any, error) {
	q, args := buildILikeQuery(fields, s)
	var result map[string]any
	err := repo.q.Query(ctx, q, &result, args)
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
