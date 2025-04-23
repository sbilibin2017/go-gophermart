package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type RewardFilterILikeQuerier interface {
	Query(
		ctx context.Context,
		dest any,
		query string,
		argMap map[string]any,
	) error
}

type RewardFilterILikeRepository struct {
	q RewardFilterILikeQuerier
}

func NewRewardFilterILikeRepository(
	q RewardFilterILikeQuerier,
) *RewardFilterILikeRepository {
	return &RewardFilterILikeRepository{q: q}
}

func (r *RewardFilterILikeRepository) FilterILike(
	ctx context.Context, rewardID string, fields []string,
) (*types.RewardDB, error) {
	query, argMap := buildGoodRewardFilterILikeQuery(rewardID, fields)

	var result *types.RewardDB
	err := r.q.Query(ctx, &result, query, argMap)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func buildGoodRewardFilterILikeQuery(reward_id string, fields []string) (string, map[string]any) {
	argMap := make(map[string]any)
	argMap["reward_id"] = "%" + reward_id + "%"
	fieldsQuery := strings.Join(fields, ", ")
	query := fmt.Sprintf(goodRewardFilterILikeQuery, fieldsQuery)
	return query, argMap
}

const goodRewardFilterILikeQuery = "SELECT %s FROM rewards WHERE reward_id ILIKE :reward_id"
