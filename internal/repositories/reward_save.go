package repositories

import (
	"context"
)

type RewardExecutor interface {
	Execute(
		ctx context.Context,
		query string,
		arg any,
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
	ctx context.Context, rewardSave *RewardSave,
) error {
	err := r.e.Execute(ctx, rewardSaveQuery, rewardSave) // Передаем структуру напрямую
	if err != nil {
		return err
	}

	return nil
}

type RewardSave struct {
	RewardID   string `db:"reward_id"`
	Reward     int64  `db:"reward"`
	RewardType string `db:"reward_type"`
}

const rewardSaveQuery = `
	INSERT INTO rewards (reward_id, reward, reward_type, created_at, updated_at)
	VALUES (:reward_id, :reward, :reward_type, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	ON CONFLICT (reward_id) DO UPDATE
	SET reward = EXCLUDED.reward,
		reward_type = EXCLUDED.reward_type,
		updated_at = CURRENT_TIMESTAMP
`
