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
		args any,
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
	ctx context.Context, filter *RewardFilterILike, // принимаем указатель на фильтр
) (*types.RewardDB, error) {
	query, argMap := buildGoodRewardFilterILikeQuery(filter)

	var result *types.RewardDB
	err := r.q.Query(ctx, &result, query, argMap)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type RewardFilterILike struct {
	RewardID string   `db:"reward_id"`
	Fields   []string // поля для запроса
}

func buildGoodRewardFilterILikeQuery(filter *RewardFilterILike) (string, map[string]any) {
	argMap := make(map[string]any)
	argMap["reward_id"] = "%" + filter.RewardID + "%" // используем фильтр как указатель
	fieldsQuery := "*"
	if len(filter.Fields) > 0 {
		fieldsQuery = strings.Join(filter.Fields, ", ")
	}
	query := fmt.Sprintf(goodRewardFilterILikeQuery, fieldsQuery)
	return query, argMap
}

const goodRewardFilterILikeQuery = `
	SELECT %s 
	FROM rewards 
	WHERE reward_id ILIKE :reward_id
`
