package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type RewardRegisterMechanicFilterOneRepository interface {
	FilterOne(ctx context.Context, match string) (*types.RewardMechanicDB, error)
}

type RewardRegisterMechanicSaveRepository interface {
	Save(ctx context.Context, match string, reward int64, rewardType string) error
}

type RewardRegisterMechanicValidator interface {
	Struct(v any) error
}

type RewardRegisterMechanicService struct {
	v  RewardRegisterMechanicValidator
	rf RewardRegisterMechanicFilterOneRepository
	rs RewardRegisterMechanicSaveRepository
}

func NewRewardRegisterMechanicService(
	v RewardRegisterMechanicValidator,
	rf RewardRegisterMechanicFilterOneRepository,
	rs RewardRegisterMechanicSaveRepository,
) *RewardRegisterMechanicService {
	return &RewardRegisterMechanicService{
		v:  v,
		rf: rf,
		rs: rs,
	}
}

func (svc *RewardRegisterMechanicService) Register(
	ctx context.Context, req *types.RewardRegisterMechanicRequest,
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
