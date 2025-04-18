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

type RewardExistsRepository interface {
	Exists(ctx context.Context, filter map[string]any) (bool, error)
}

type RewardSaveRepository interface {
	Save(ctx context.Context, reward map[string]any) error
}

type RewardSaveService struct {
	re RewardExistsRepository
	rs RewardSaveRepository
}

func NewRewardSaveService(
	re RewardExistsRepository,
	rs RewardSaveRepository,
) *RewardSaveService {
	return &RewardSaveService{
		re: re,
		rs: rs,
	}
}

func (svc *RewardSaveService) Register(
	ctx context.Context, reward *types.Reward,
) error {
	exists, err := svc.re.Exists(ctx, map[string]any{"match": reward.Match})
	if err != nil {
		return ErrRewardIsNotRegistered
	}
	if exists {
		return ErrRewardAlreadyExists
	}

	err = svc.rs.Save(ctx, map[string]any{
		"match":       reward.Match,
		"reward":      reward.Reward,
		"reward_type": reward.RewardType,
	})
	if err != nil {
		return ErrRewardIsNotRegistered
	}

	return nil
}
