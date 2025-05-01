package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/types"
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
	ctx context.Context, reward *types.RewardDB,
) error {
	return execContext(
		ctx,
		r.db,
		r.txProvider,
		rewardSaveQuery,
		reward.Match, reward.Reward, reward.RewardType,
	)
}

const rewardSaveQuery = `
	INSERT INTO reward (match, reward, reward_type)
	VALUES ($1, $2, $3)
	ON CONFLICT (match)
	DO UPDATE SET
		reward = EXCLUDED.reward,
		reward_type = EXCLUDED.reward_type
`
