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
	Save(ctx context.Context, reward *types.RewardDB) error
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
) (*types.APISuccessStatus, *types.APIErrorStatus) {
	if err := svc.v.Struct(req); err != nil {
		valErr := formatValidationError(err)
		return nil, &types.APIErrorStatus{
			StatusCode: http.StatusBadRequest,
			Message:    valErr.Message,
		}
	}

	existingReward, err := svc.rf.FilterOne(ctx, req.Match)
	if err != nil {
		return nil, &types.APIErrorStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}
	if existingReward != nil {
		return nil, &types.APIErrorStatus{
			StatusCode: http.StatusConflict,
			Message:    "Reward match key already exists",
		}
	}

	reward := &types.RewardDB{
		Match:      req.Match,
		Reward:     req.Reward,
		RewardType: req.RewardType,
	}
	if err := svc.rs.Save(ctx, reward); err != nil {
		return nil, &types.APIErrorStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save reward",
		}
	}

	return &types.APISuccessStatus{
		StatusCode: http.StatusOK,
		Message:    "Reward successfully registered",
	}, nil
}
