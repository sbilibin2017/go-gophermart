package repositories

import (
	"context"
)

type AccrualRewardMechanicExistsRepository struct {
	q Querier
}

func NewAccrualRewardMechanicExistsRepository(
	q Querier,
) *AccrualRewardMechanicExistsRepository {
	return &AccrualRewardMechanicExistsRepository{q: q}
}

func (repo *AccrualRewardMechanicExistsRepository) Exists(
	ctx context.Context,
	filter map[string]any,
) (bool, error) {
	var exists bool
	err := repo.q.Query(ctx, accrualRewardMechanicExistsQuery, &exists, filter["match"])
	if err != nil {
		return false, err
	}
	return exists, nil
}

const accrualRewardMechanicExistsQuery = `
	SELECT EXISTS (
		SELECT 1
		FROM accrual_reward_mechanic
		WHERE match = $1
	)
`
