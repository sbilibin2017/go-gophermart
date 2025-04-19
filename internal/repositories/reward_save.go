package repositories

import (
	"context"
)

type RewardSaveExecutor interface {
	Execute(ctx context.Context, query string, args map[string]any) error
}

type RewardSaveRepository struct {
	e RewardSaveExecutor
}

func NewRewardSaveRepository(e RewardSaveExecutor) *RewardSaveRepository {
	return &RewardSaveRepository{e: e}
}

func (r *RewardSaveRepository) Save(
	ctx context.Context, args map[string]any,
) error {
	return r.e.Execute(ctx, rewardSaveQuery, args)
}

const rewardSaveQuery = `
	INSERT INTO rewards (match, reward, reward_type, created_at, updated_at)
	VALUES (:match, :reward, :reward_type, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	ON CONFLICT (match) DO UPDATE
	SET reward = EXCLUDED.reward,
		reward_type = EXCLUDED.reward_type,
		updated_at = CURRENT_TIMESTAMP
`
