package repositories

import (
	"context"
)

type AccrualRewardMechanicSaveRepository struct {
	e Executor
}

func NewAccrualRewardMechanicSaveRepository(
	e Executor,
) *AccrualRewardMechanicSaveRepository {
	return &AccrualRewardMechanicSaveRepository{e: e}
}

func (repo *AccrualRewardMechanicSaveRepository) Save(
	ctx context.Context,
	data map[string]any,
) error {
	return repo.e.Exec(ctx, saveAccrualRewardMechanicQuery, data)
}

const saveAccrualRewardMechanicQuery = `
	INSERT INTO accrual_reward_mechanic (match, reward, reward_type, created_at, updated_at)
	VALUES (:match, :reward, :reward_type, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	ON CONFLICT (match)
	DO UPDATE SET reward = :reward, reward_type = :reward_type, updated_at = CURRENT_TIMESTAMP;
`
