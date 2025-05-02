package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type RewardRegisterRewardFilterOneRepository interface {
	FilterOne(ctx context.Context, match string) (*types.RewardDB, error)
}

type RewardRegisterRewardSaveRepository interface {
	Save(ctx context.Context, match string, reward int64, rewardType string) error
}

type RewardRegisterValidator interface {
	Struct(v any) error
}

type RewardRegisterService struct {
	v  RewardRegisterValidator
	rf RewardRegisterRewardFilterOneRepository
	rs RewardRegisterRewardSaveRepository
}

func NewRewardRegisterService(
	v RewardRegisterValidator,
	rf RewardRegisterRewardFilterOneRepository,
	rs RewardRegisterRewardSaveRepository,
) *RewardRegisterService {
	return &RewardRegisterService{
		v:  v,
		rf: rf,
		rs: rs,
	}
}

func (svc *RewardRegisterService) Register(
	ctx context.Context, req *types.RewardRequest,
) (*types.APIStatus, *types.APIStatus) {
	if err := svc.v.Struct(req); err != nil {
		valErr := formatValidationError(err)
		return nil, &types.APIStatus{
			StatusCode: http.StatusBadRequest,
			Message:    valErr.Message,
		}
	}

	existingReward, err := svc.rf.FilterOne(ctx, req.Match)
	if err != nil {
		return nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}
	if existingReward != nil {
		return nil, &types.APIStatus{
			StatusCode: http.StatusConflict,
			Message:    "Reward match key already exists",
		}
	}

	if err := svc.rs.Save(ctx, req.Match, req.Reward, req.RewardType); err != nil {
		return nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save reward",
		}
	}

	return &types.APIStatus{
		StatusCode: http.StatusOK,
		Message:    "Reward successfully registered",
	}, nil
}
