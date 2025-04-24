package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type RewardExistsRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) *sqlx.Tx
}

func NewRewardExistsRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
) *RewardExistsRepository {
	return &RewardExistsRepository{db: db, txProvider: txProvider}
}

func (r *RewardExistsRepository) Exists(ctx context.Context, match string) (bool, error) {
	var exists bool
	err := query(ctx, r.db, r.txProvider, &exists, rewardExistsByIDQuery, match)
	if err != nil {
		return false, err
	}
	return exists, nil
}

const rewardExistsByIDQuery = `SELECT EXISTS (SELECT 1	FROM rewards WHERE match = ?)`
