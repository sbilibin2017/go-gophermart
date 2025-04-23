package repositories

import (
	"context"
)

type RewardExecutor interface {
	Execute(
		ctx context.Context,
		query string,
		argMap map[string]any,
	) error
}

type RewardSaveRepository struct {
	e RewardExecutor
}

func NewRewardSaveRepository(
	e RewardExecutor,
) *RewardSaveRepository {
	return &RewardSaveRepository{e: e}
}

func (r *RewardSaveRepository) Save(
	ctx context.Context, rewardID string, reward int64, rewardType string,
) error {
	argMap := map[string]any{
		"reward_id":   rewardID,
		"reward":      reward,
		"reward_type": rewardType,
	}
	err := r.e.Execute(ctx, rewardSaveQuery, argMap)
	if err != nil {
		return err
	}

	return nil
}

const rewardSaveQuery = `
	INSERT INTO rewards (reward_id, reward, reward_type, created_at, updated_at)
	VALUES (:reward_id, :reward, :reward_type, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	ON CONFLICT (reward_id) DO UPDATE
	SET reward = EXCLUDED.reward,
		reward_type = EXCLUDED.reward_type,
		updated_at = CURRENT_TIMESTAMP
`
