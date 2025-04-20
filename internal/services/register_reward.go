package services

import (
	"context"
	"errors"
)

var (
	ErrRewardAlreadyExists   = errors.New("reward already exists")
	ErrRewardIsNotRegistered = errors.New("reward is not registered")
)

type RewardExistsRepository interface {
	Exists(ctx context.Context, filter map[string]any) (bool, error)
}

type RewardSaveRepository interface {
	Save(ctx context.Context, data map[string]any) error
}

type RegisterRewardSaveService struct {
	re RewardExistsRepository
	rs RewardSaveRepository
}

func NewRegisterRewardSaveService(
	re RewardExistsRepository,
	rs RewardSaveRepository,
) *RegisterRewardSaveService {
	return &RegisterRewardSaveService{
		re: re,
		rs: rs,
	}
}

func (svc *RegisterRewardSaveService) Register(
	ctx context.Context, match string, reward uint64, rewardType string,
) error {
	exists, err := svc.re.Exists(ctx, map[string]any{"match": match})
	if err != nil {
		return ErrRewardIsNotRegistered
	}
	if exists {
		return ErrRewardAlreadyExists
	}

	err = svc.rs.Save(ctx, map[string]any{
		"match":       match,
		"reward":      reward,
		"reward_type": rewardType,
	})
	if err != nil {
		return ErrRewardIsNotRegistered
	}

	return nil
}
