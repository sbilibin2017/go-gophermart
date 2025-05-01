package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type RewardFilterOneRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, error)
}

func NewRewardFilterOneRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
) *RewardFilterOneRepository {
	return &RewardFilterOneRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *RewardFilterOneRepository) FilterOne(
	ctx context.Context, match string,
) (*types.RewardDB, error) {
	var reward types.RewardDB
	err := getContext(
		ctx,
		r.db,
		r.txProvider,
		rewardFilterOneQuery,
		&reward,
		match,
	)
	if err != nil {
		return nil, err
	}
	return &reward, nil
}

const rewardFilterOneQuery = `
	SELECT *
	FROM reward 
	WHERE match = &1
`
