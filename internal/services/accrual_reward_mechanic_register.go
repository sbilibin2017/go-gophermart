package services

import (
	"context"
	"net/http"
)

type AccrualRewardMechanicRegisterService struct {
	v  StructValidator
	re ExistsRepository
	rs SaveRepository
}

func NewAccrualRewardMechanicRegisterService(
	v StructValidator,
	re ExistsRepository,
	rs SaveRepository,
) *AccrualRewardMechanicRegisterService {
	return &AccrualRewardMechanicRegisterService{
		v:  v,
		re: re,
		rs: rs,
	}
}

func (s *AccrualRewardMechanicRegisterService) Register(
	ctx context.Context, req *AccrualRewardMechanicRegisterRequest,
) (*APIStatus, *APIStatus) {
	if err := s.v.Struct(req); err != nil {
		return nil, &APIStatus{
			StatusCode: http.StatusBadRequest,
			Message:    "Reward mechanic data is invalid",
		}
	}
	exists, err := s.re.Exists(ctx, map[string]any{"match": req.Match})
	if err != nil {
		return nil, &APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error checking if reward mechanic exists",
		}
	}
	if exists {
		return nil, &APIStatus{
			StatusCode: http.StatusConflict,
			Message:    "Reward mechanic with this match key already exists",
		}
	}
	if err := s.rs.Save(ctx, map[string]any{"match": req.Match, "reward": req.Reward, "reward_type": req.RewardType}); err != nil {
		return nil, &APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error saving reward mechanic",
		}
	}
	return &APIStatus{
		StatusCode: http.StatusOK,
		Message:    "Reward mechanic registered successfully",
	}, nil
}

type AccrualRewardMechanicRegisterRequest struct {
	Match      string `json:"match" validate:"required"`
	Reward     int64  `json:"reward" validate:"required,gt=0"`
	RewardType string `json:"reward_type" validate:"required"`
}
