package repositories

import (
	"context"
)

type RewardExistsQuerier interface {
	Query(
		ctx context.Context,
		dest any,
		query string,
		args any,
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
	ctx context.Context, filter *RewardExistsFilter, // принимаем указатель на фильтр
) (bool, error) {
	var exists bool
	err := r.q.Query(ctx, &exists, rewardExistsByIDQuery, filter)
	if err != nil {
		return false, err
	}
	return exists, nil
}

type RewardExistsFilter struct {
	RewardID string `db:"reward_id"`
}

var rewardExistsByIDQuery = `
	SELECT EXISTS(
		SELECT 1
		FROM rewards
		WHERE reward_id = :reward_id
	)
`
