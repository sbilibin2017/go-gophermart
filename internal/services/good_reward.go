package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

const (
	RewardRegisteredSuccessMessage = "Reward successfully registered"
	InvalidRewardRequestMessage    = "Invalid reward request"
	InternalServerErrorMessage     = "Internal server error while checking reward existence"
	RewardExistsMessage            = "Reward with this match already exists"
	FailedToSaveRewardMessage      = "Failed to save reward"
)

type GoodRewardExistsRepository interface {
	Exists(ctx context.Context, rewardID string) (bool, error)
}

type GoodRewardSaveRepository interface {
	Save(ctx context.Context, rewardID string, reward int64, rewardType string) error
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
	ctx context.Context, req *types.GoodRewardRegisterRequest,
) (*types.APIStatus, error) {
	if err := svc.v.Struct(req); err != nil {
		return &types.APIStatus{
			Status:  http.StatusBadRequest,
			Message: InvalidRewardRequestMessage,
		}, err
	}

	exists, err := svc.re.Exists(ctx, req.Match)
	if err != nil {
		return &types.APIStatus{
			Status:  http.StatusInternalServerError,
			Message: InternalServerErrorMessage,
		}, err
	}

	if exists {
		return &types.APIStatus{
			Status:  http.StatusConflict,
			Message: RewardExistsMessage,
		}, nil
	}

	err = svc.rs.Save(
		ctx,
		req.Match,
		req.Reward,
		string(req.RewardType),
	)

	if err != nil {
		return &types.APIStatus{
			Status:  http.StatusInternalServerError,
			Message: FailedToSaveRewardMessage,
		}, err
	}

	return &types.APIStatus{
		Status:  http.StatusOK,
		Message: RewardRegisteredSuccessMessage,
	}, nil
}
