package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type RewardSaveRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, bool)
}

func NewRewardSaveRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, bool),
) *RewardSaveRepository {
	return &RewardSaveRepository{db: db, txProvider: txProvider}
}

func (r *RewardSaveRepository) Save(ctx context.Context, reward map[string]any) error {
	err := execContextNamed(ctx, r.db, r.txProvider, rewardSaveQuery, reward)
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
