package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type RewardExistsRepository interface {
	Exists(ctx context.Context, match string) (bool, error)
}

type RewardSaveRepository interface {
	Save(ctx context.Context, data map[string]any) error
}

type RewardRegisterValidator interface {
	Struct(v any) error
}

type RewardRegisterService struct {
	v  RewardRegisterValidator
	re RewardExistsRepository
	rs RewardSaveRepository
}

func NewRewardRegisterService(
	v RewardRegisterValidator,
	re RewardExistsRepository,
	rs RewardSaveRepository,
) *RewardRegisterService {
	return &RewardRegisterService{
		v:  v,
		re: re,
		rs: rs,
	}
}

func (svc *RewardRegisterService) Register(
	ctx context.Context, req *types.RewardRegisterRequest,
) (*types.APIStatus, *types.APIStatus) {
	err := svc.v.Struct(req)
	if err != nil {
		return nil, &types.APIStatus{
			Status:  http.StatusBadRequest,
			Message: "Reward is not valid",
		}
	}

	exists, err := svc.re.Exists(ctx, req.Match)
	if err != nil {
		return nil, &types.APIStatus{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}
	if exists {
		return nil, &types.APIStatus{
			Status:  http.StatusConflict,
			Message: "Reward already exists",
		}
	}

	err = svc.rs.Save(ctx, convertToMap(req))
	if err != nil {
		return nil, &types.APIStatus{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}

	return &types.APIStatus{
		Status:  http.StatusOK,
		Message: "Reward registered successfully",
	}, nil
}
