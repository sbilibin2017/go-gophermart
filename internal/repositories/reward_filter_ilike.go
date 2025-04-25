package repositories

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type RewardFilterILikeRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) *sqlx.Tx
}

func NewRewardFilterILikeRepository(db *sqlx.DB, txProvider func(ctx context.Context) *sqlx.Tx) *RewardFilterILikeRepository {
	return &RewardFilterILikeRepository{db: db, txProvider: txProvider}
}

func (r *RewardFilterILikeRepository) FilterILike(
	ctx context.Context,
	description *RewardFilterILike,
) (*RewardFilterILikeDB, error) {
	var result RewardFilterILikeDB
	err := query(ctx, r.db, r.txProvider, &result, rewardFilterILikeQuery, description)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type RewardFilterILike struct {
	Description string `db:"match"`
}

type RewardFilterILikeDB struct {
	Match      string    `db:"match"`
	Reward     int64     `db:"reward"`
	RewardType string    `db:"reward_type"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

const rewardFilterILikeQuery = `
	SELECT * 
	FROM rewards 
	WHERE match ILIKE :description
	LIMIT 1
`
