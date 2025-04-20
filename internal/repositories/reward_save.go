package repositories

import (
	"context"
	"database/sql"
)

type RewardSaveRepository struct {
	db         *sql.DB
	txProvider func(ctx context.Context) *sql.Tx
}

func NewRewardSaveRepository(
	db *sql.DB,
	txProvider func(ctx context.Context) *sql.Tx,
) *RewardSaveRepository {
	return &RewardSaveRepository{db: db, txProvider: txProvider}
}

func (r *RewardSaveRepository) Save(
	ctx context.Context, data map[string]any,
) error {
	_, err := execContext(
		ctx, r.db, r.txProvider, rewardSaveQuery,
		data["match"], data["reward"], data["reward_type"],
	)
	return err
}

const rewardSaveQuery = `
	INSERT INTO rewards (match, reward, reward_type, created_at, updated_at)
	VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	ON CONFLICT (match) DO UPDATE
	SET reward = EXCLUDED.reward,
		reward_type = EXCLUDED.reward_type,
		updated_at = CURRENT_TIMESTAMP
`
