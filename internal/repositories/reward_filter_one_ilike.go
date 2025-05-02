package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type RewardFilterOneILikeRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, error)
}

func NewRewardFilterOneILikeRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
) *RewardFilterOneILikeRepository {
	return &RewardFilterOneILikeRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *RewardFilterOneILikeRepository) FilterOneILike(
	ctx context.Context, description string,
) (*types.RewardDB, error) {
	var reward types.RewardDB
	query, args := buildRewardFilterOneILikeQuery(description)
	err := getContext(
		ctx,
		r.db,
		r.txProvider,
		query,
		&reward,
		args,
	)
	if err != nil {
		return nil, err
	}
	return &reward, nil
}

func buildRewardFilterOneILikeQuery(description string) (string, []any) {
	return rewardFilterOneILikeQuery, []interface{}{"%" + description + "%"}
}

const rewardFilterOneILikeQuery = `
	SELECT *
	FROM reward 
	WHERE match ILIKE $1
`
