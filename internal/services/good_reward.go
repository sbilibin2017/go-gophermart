package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type GoodRewardExistsRepository interface {
	ExistsByID(ctx context.Context, rewardID string) (bool, error)
}

type GoodRewardSaveRepository interface {
	Save(ctx context.Context, reward map[string]any) error
}

type GoodRewardValidator interface {
	Struct(v any) error
}

type GoodRewardService struct {
	v  GoodRewardValidator
	re GoodRewardExistsRepository
	rs GoodRewardSaveRepository
}

func NewGoodRewardService(
	v GoodRewardValidator,
	re GoodRewardExistsRepository,
	rs GoodRewardSaveRepository,
) *GoodRewardService {
	return &GoodRewardService{
		v:  v,
		re: re,
		rs: rs,
	}
}

func (svc *GoodRewardService) Register(
	ctx context.Context, req *types.RewardRegisterRequest,
) (*types.APIResponse[any], *types.APIError) {
	if err := svc.v.Struct(req); err != nil {
		return &types.APIResponse[any]{
			Status:  http.StatusBadRequest,
			Message: "Invalid reward request",
		}, nil
	}

	exists, err := svc.re.ExistsByID(ctx, req.Match)
	if err != nil {
		return &types.APIResponse[any]{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error while checking reward existence",
		}, nil
	}

	if exists {
		return &types.APIResponse[any]{
			Status:  http.StatusConflict,
			Message: "Reward with this match already exists",
		}, nil
	}

	err = svc.rs.Save(
		ctx,
		map[string]any{
			"reward_id":   req.Match,
			"reward":      req.Reward,
			"reward_type": string(req.RewardType),
		},
	)

	if err != nil {
		return &types.APIResponse[any]{
			Status:  http.StatusInternalServerError,
			Message: "Failed to save reward",
		}, nil
	}

	return &types.APIResponse[any]{
		Status:  http.StatusOK,
		Message: "Reward successfully registered",
	}, nil
}
