package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type AccrualRewardMechanicRegisterExistsRepository interface {
	Exists(ctx context.Context, match string) (bool, error)
}

type AccrualRewardMechanicRegisterSaveRepository interface {
	Save(ctx context.Context, match string, reward int64, rewardType string) error
}

type AccrualRewardMechanicRegisterValidator interface {
	Struct(v any) error
}

type AccrualRewardMechanicRegisterService struct {
	v  AccrualRewardMechanicRegisterValidator
	re AccrualRewardMechanicRegisterExistsRepository
	rs AccrualRewardMechanicRegisterSaveRepository
}

func NewAccrualRewardMechanicRegisterService(
	v AccrualRewardMechanicRegisterValidator,
	re AccrualRewardMechanicRegisterExistsRepository,
	rs AccrualRewardMechanicRegisterSaveRepository,
) *AccrualRewardMechanicRegisterService {
	return &AccrualRewardMechanicRegisterService{
		v:  v,
		re: re,
		rs: rs,
	}
}

func (s *AccrualRewardMechanicRegisterService) Register(ctx context.Context, req *types.AccrualRewardMechanicRegisterRequest) (*types.APIStatus, *types.APIStatus) {
	if err := s.v.Struct(req); err != nil {
		return nil, &types.APIStatus{
			StatusCode: http.StatusBadRequest,
			Message:    "Reward mechanic data is invalid",
		}
	}
	exists, err := s.re.Exists(ctx, req.Match)
	if err != nil {
		return nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Reward mechanic is not registered",
		}
	}
	if exists {
		return nil, &types.APIStatus{
			StatusCode: http.StatusConflict,
			Message:    "Reward mechanic with this match key already exists",
		}
	}
	if err := s.rs.Save(ctx, req.Match, req.Reward, req.RewardType); err != nil {
		return nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Reward mechanic is not registered",
		}
	}
	return &types.APIStatus{
		StatusCode: http.StatusOK,
		Message:    "Reward mechanic registered successfully",
	}, nil
}
