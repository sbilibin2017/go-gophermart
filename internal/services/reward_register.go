package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type RewardExistsRepository interface {
	Exists(ctx context.Context, match *repositories.RewardExistsMatch) (bool, error)
}

type RewardSaveRepository interface {
	Save(ctx context.Context, reward *repositories.RewardSave) error
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
	ctx context.Context, req *RewardRegisterRequest,
) (*types.APIStatus, *types.APIStatus) {
	err := svc.v.Struct(req)
	if err != nil {
		return nil, &types.APIStatus{
			Status:  http.StatusBadRequest,
			Message: "Reward is not valid",
		}
	}
	exists, err := svc.re.Exists(ctx, &repositories.RewardExistsMatch{
		Match: req.Match,
	})
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
	err = svc.rs.Save(ctx, &repositories.RewardSave{
		Match:      req.Match,
		Reward:     req.Reward,
		RewardType: req.RewardType,
	})
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

type RewardRegisterRequest struct {
	Match      string `json:"match" validate:"required"`
	Reward     int64  `json:"reward" validate:"required,gte=0"`
	RewardType string `json:"reward_type" validate:"required,reward_type"`
}
