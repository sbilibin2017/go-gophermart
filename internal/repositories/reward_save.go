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
	return &RewardSaveRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *RewardSaveRepository) Save(ctx context.Context, match map[string]any) error {
	query := rewardSaveQuery
	args := []any{match["match"], match["reward"], match["reward_type"]}

	if tx, ok := r.txProvider(ctx); ok {
		_, err := tx.ExecContext(ctx, query, args...)
		return err
	}

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

var rewardSaveQuery = `
INSERT INTO rewards (match, reward, reward_type, created_at, updated_at) 
VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (match) 
DO UPDATE SET 
    reward = EXCLUDED.reward,
    reward_type = EXCLUDED.reward_type,
    updated_at = CURRENT_TIMESTAMP
`
