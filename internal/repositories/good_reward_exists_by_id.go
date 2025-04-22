package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type RewardExistsRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, bool)
}

func NewRewardExistsRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, bool),
) *RewardExistsRepository {
	return &RewardExistsRepository{db: db, txProvider: txProvider}
}

func (r *RewardExistsRepository) ExistsByID(ctx context.Context, rewardID string) (bool, error) {
	argMap := map[string]any{
		"reward_id": rewardID,
	}
	var exists bool
	err := getContextNamed(ctx, r.db, r.txProvider, &exists, rewardExistsByIDQuery, argMap)
	if err != nil {
		return false, err
	}
	return exists, nil
}

var rewardExistsByIDQuery = `SELECT EXISTS(SELECT 1 FROM rewards WHERE reward_id = :reward_id)`
