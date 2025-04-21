package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

var (
	ErrRewardAlreadyExists   = errors.New("reward already exists")
	ErrRewardIsNotRegistered = errors.New("reward is not registered")
)

type RegisterRewardExistsRepository interface {
	Exists(ctx context.Context, filter map[string]any) (bool, error)
}

type RegisterRewardSaveRepository interface {
	Save(ctx context.Context, data map[string]any) error
}

type RegisterRewardValidator interface {
	Struct(s any) error
}

type RegisterRewardService struct {
	v  RegisterRewardValidator
	re RegisterRewardExistsRepository
	rs RegisterRewardSaveRepository
}

func NewRegisterRewardService(
	v RegisterRewardValidator,
	re RegisterRewardExistsRepository,
	rs RegisterRewardSaveRepository,
) *RegisterRewardService {
	return &RegisterRewardService{
		v:  v,
		re: re,
		rs: rs,
	}
}

func (svc *RegisterRewardService) Register(
	ctx context.Context, reward *types.RegisterRewardRequest,
) (*string, *types.APIError) {
	if err := svc.v.Struct(reward); err != nil {
		return nil, types.NewValidationErrorResponse(err)
	}

	exists, err := svc.re.Exists(ctx, map[string]any{"reward_id": reward.Match})
	if err != nil {
		return nil, types.NewInternalError()
	}
	if exists {
		return nil, types.NewAPIError(ErrRewardAlreadyExists.Error(), http.StatusConflict)
	}

	err = svc.rs.Save(ctx, map[string]any{
		"reward_id":   reward.Match,
		"reward":      reward.Reward,
		"reward_type": reward.RewardType,
	})
	if err != nil {
		return nil, types.NewInternalError()
	}

	s := "Reward registered successfully"
	return &s, nil
}
