package repositories

import (
	"context"
)

type RewardExistsQuerier interface {
	Query(ctx context.Context, dest any, query string, args map[string]any) error
}

type RewardExistsRepository struct {
	q RewardExistsQuerier
}

func NewRewardExistsRepository(q RewardExistsQuerier) *RewardExistsRepository {
	return &RewardExistsRepository{q: q}
}

func (r *RewardExistsRepository) Exists(
	ctx context.Context, match map[string]any,
) (bool, error) {
	var exists bool
	err := r.q.Query(ctx, &exists, rewardExistsQuery, match)
	if err != nil {
		return false, err
	}
	return exists, nil
}

const rewardExistsQuery = "SELECT EXISTS(SELECT 1 FROM rewards WHERE match = :match)"
