package services

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
)

type RewardExistsRepository interface {
	Exists(ctx context.Context, filter map[string]any) (bool, error)
}

type RewardSaveRepository interface {
	Save(ctx context.Context, match map[string]any) error
}

type RegisterRewardService struct {
	re RewardExistsRepository
	rs RewardSaveRepository
	db *sqlx.DB
}

func NewRegisterRewardService(
	re RewardExistsRepository,
	rs RewardSaveRepository,
	db *sqlx.DB,
) *RegisterRewardService {
	return &RegisterRewardService{re: re, rs: rs, db: db}
}

var (
	ErrGoodRewardAlreadyExists = errors.New("reward already exists")
)

func (svc *RegisterRewardService) Register(ctx context.Context, reward *domain.Reward) error {
	exists, err := svc.re.Exists(ctx, map[string]any{"match": reward.Match})
	if err != nil {
		return domain.ErrFailedToRegisterReward
	}
	if exists {
		return domain.ErrRewardKeyAlreadyRegistered
	}

	err = svc.rs.Save(ctx, map[string]any{
		"match":       reward.Match,
		"reward":      reward.Reward,
		"reward_type": string(reward.RewardType),
	})
	if err != nil {
		return domain.ErrFailedToRegisterReward
	}

	return nil
}
