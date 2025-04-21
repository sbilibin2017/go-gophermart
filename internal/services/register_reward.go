package services

import (
	"context"
	"errors"

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

type RegisterRewardService struct {
	re RegisterRewardExistsRepository
	rs RegisterRewardSaveRepository
}

func NewRegisterRewardService(
	re RegisterRewardExistsRepository,
	rs RegisterRewardSaveRepository,
) *RegisterRewardService {
	return &RegisterRewardService{
		re: re,
		rs: rs,
	}
}

func (svc *RegisterRewardService) Register(
	ctx context.Context, reward *types.RegisterRewardRequest,
) error {
	exists, err := svc.re.Exists(ctx, map[string]any{"reward_id": reward.Match})
	if err != nil {
		return ErrRewardIsNotRegistered
	}
	if exists {
		return ErrRewardAlreadyExists
	}

	err = svc.rs.Save(ctx, map[string]any{
		"reward_id":   reward.Match,
		"reward":      reward.Reward,
		"reward_type": reward.RewardType,
	})
	if err != nil {
		return ErrRewardIsNotRegistered
	}

	return nil
}
