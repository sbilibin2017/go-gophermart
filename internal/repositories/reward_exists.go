package repositories

import (
	"context"
	"database/sql"
)

type RewardExistsRepository struct {
	db         *sql.DB
	txProvider func(ctx context.Context) *sql.Tx
}

func NewRewardExistsRepository(
	db *sql.DB,
	txProvider func(ctx context.Context) *sql.Tx,
) *RewardExistsRepository {
	return &RewardExistsRepository{db: db, txProvider: txProvider}
}

func (r *RewardExistsRepository) Exists(
	ctx context.Context, match map[string]any,
) (bool, error) {
	var exists bool
	row, err := queryRowContext(ctx, r.db, r.txProvider, rewardExistsQuery, match["match"])
	if err != nil {
		return false, err
	}
	if err := scanRow(row, &exists); err != nil {
		return false, err
	}
	return exists, nil
}



const rewardExistsQuery = "SELECT EXISTS(SELECT 1 FROM rewards WHERE match = $1)"
