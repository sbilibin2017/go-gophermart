package repositories

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type AccrualRewardMechanicSaveRepository struct {
	db *sqlx.DB
}

func NewAccrualRewardMechanicSaveRepository(db *sqlx.DB) *AccrualRewardMechanicSaveRepository {
	return &AccrualRewardMechanicSaveRepository{db: db}
}

func (r *AccrualRewardMechanicSaveRepository) Save(
	ctx context.Context,
	match string,
	reward int64,
	rewardType string,
) error {
	args := map[string]any{
		"match":       match,
		"reward":      reward,
		"reward_type": rewardType,
	}
	return exec(ctx, r.db, saveAccrualRewardMechanicQuery, args)
}

const saveAccrualRewardMechanicQuery = `
	INSERT INTO accrual_reward_mechanic (match, reward, reward_type, created_at, updated_at)
	VALUES (:match, :reward, :reward_type, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	ON CONFLICT (match)
	DO UPDATE SET reward = :reward, reward_type = :reward_type, updated_at = CURRENT_TIMESTAMP;
`
