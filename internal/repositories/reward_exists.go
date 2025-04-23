package repositories

import (
	"context"
)

type RewardExistsQuerier interface {
	Query(
		ctx context.Context,
		dest any,
		query string,
		argMap map[string]any,
	) error
}

type RewardExistsRepository struct {
	q RewardExistsQuerier
}

func NewRewardExistsRepository(
	q RewardExistsQuerier,
) *RewardExistsRepository {
	return &RewardExistsRepository{q: q}
}

func (r *RewardExistsRepository) Exists(
	ctx context.Context, rewardID string,
) (bool, error) {
	argMap := map[string]any{
		"reward_id": rewardID,
	}
	var exists bool
	err := r.q.Query(ctx, &exists, rewardExistsByIDQuery, argMap)
	if err != nil {
		return false, err
	}
	return exists, nil
}

var rewardExistsByIDQuery = `SELECT EXISTS(SELECT 1 FROM rewards WHERE reward_id = :reward_id)`
