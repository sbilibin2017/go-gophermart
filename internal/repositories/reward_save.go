package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type RewardSaveRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, error)
}

func NewRewardSaveRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
) *RewardSaveRepository {
	return &RewardSaveRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *RewardSaveRepository) Save(
	ctx context.Context, match string, reward int64, rewardType string,
) error {
	return execContext(
		ctx,
		r.db,
		r.txProvider,
		rewardSaveQuery,
		match, reward, rewardType,
	)
}

const rewardSaveQuery = `
	INSERT INTO rewards (match, reward, reward_type)
	VALUES ($1, $2, $3)
	ON CONFLICT (match)
	DO UPDATE SET
		reward = EXCLUDED.reward,
		reward_type = EXCLUDED.reward_type
`
