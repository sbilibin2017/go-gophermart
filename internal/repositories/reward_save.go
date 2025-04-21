package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/repositories/helpers"
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

func (r *RewardSaveRepository) Save(
	ctx context.Context, data map[string]any,
) error {
	_, err := helpers.ExecContext(ctx, r.db, r.txProvider, rewardSaveQuery, data)
	return err
}

const rewardSaveQuery = `
	INSERT INTO rewards (reward_id, reward, reward_type, created_at, updated_at)
	VALUES (:reward_id, :reward, :reward_type, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	ON CONFLICT (reward_id) DO UPDATE
	SET reward = EXCLUDED.reward,
		reward_type = EXCLUDED.reward_type,
		updated_at = CURRENT_TIMESTAMP
`
