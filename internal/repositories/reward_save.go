package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/queries"
)

type RewardSaveRepository struct {
	db *sqlx.DB
}

func NewRewardSaveRepository(db *sqlx.DB) *RewardSaveRepository {
	return &RewardSaveRepository{db: db}
}

func (r *RewardSaveRepository) Save(
	ctx context.Context, reward map[string]any,
) error {
	if len(reward) == 0 {
		return nil
	}

	tx := middlewares.TxFromContext(ctx)

	args := map[string]interface{}{
		"match":       reward["match"],
		"reward":      reward["reward"],
		"reward_type": reward["reward_type"],
	}

	if tx != nil {
		_, err := tx.NamedExecContext(ctx, queries.RewardSaveQuery, args)
		if err != nil {
			return err
		}
	} else {
		_, err := r.db.NamedExecContext(ctx, queries.RewardSaveQuery, args)
		if err != nil {
			return err
		}
	}

	return nil
}
