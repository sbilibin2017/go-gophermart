package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/repositories/helpers"
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

func (r *RewardExistsRepository) Exists(
	ctx context.Context, filter map[string]any,
) (bool, error) {
	row, err := helpers.QueryRowContext(ctx, r.db, r.txProvider, rewardExistsQuery, filter)
	if err != nil {
		return false, err
	}
	exists, err := helpers.Scan[bool](row)
	if err != nil {
		return false, err
	}
	return exists, nil
}

const rewardExistsQuery = `
	SELECT EXISTS(
		SELECT 1 FROM rewards WHERE reward_id = :reward_id
	)
`
