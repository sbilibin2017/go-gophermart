package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type RewardSaveRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) *sqlx.Tx
}

func NewRewardSaveRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
) *RewardSaveRepository {
	return &RewardSaveRepository{db: db, txProvider: txProvider}
}

func (r *RewardSaveRepository) Save(ctx context.Context, reward *RewardSave) error {
	return command(ctx, r.db, r.txProvider, rewardSaveQuery, reward)
}

type RewardSave struct {
	Match      string `db:"match"`
	Reward     int64  `db:"reward"`
	RewardType string `db:"reward_type"`
}

const rewardSaveQuery = `
	INSERT INTO rewards (match, reward, reward_type, created_at, updated_at)
	VALUES (:match, :reward, :reward_type, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	ON CONFLICT (match) DO UPDATE
	SET reward = EXCLUDED.reward,
		reward_type = EXCLUDED.reward_type,
		updated_at = CURRENT_TIMESTAMP
`
