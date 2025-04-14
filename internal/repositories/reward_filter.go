package repositories

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/queries"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type RewardFilterRepository struct {
	db *sqlx.DB
}

func NewRewardFilterRepository(db *sqlx.DB) *RewardFilterRepository {
	return &RewardFilterRepository{db: db}
}

func getRewardFilterQueryArgs(filter []*types.RewardFilter) []interface{} {
	matches := make([]string, 0, len(filter))
	for _, f := range filter {
		matches = append(matches, "%"+f.Description+"%")
	}
	return []interface{}{matches}
}

func (r *RewardFilterRepository) Filter(
	ctx context.Context, filter []*types.RewardFilter,
) ([]*types.RewardDB, error) {
	args := getRewardFilterQueryArgs(filter)
	var rewards []*types.RewardDB
	err := r.db.SelectContext(ctx, &rewards, queries.RewardFilterQuery, args...)
	if err != nil {
		return nil, err
	}
	return rewards, nil
}
