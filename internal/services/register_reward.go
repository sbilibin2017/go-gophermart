package services

import (
	"context"
	"errors"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type RewardExistsRepository interface {
	Exists(ctx context.Context, match string) (bool, error)
}

type RewardSaveRepository interface {
	Save(ctx context.Context, match string, reward uint, rewardType string) error
}

type RewardService struct {
	re RewardExistsRepository
	rs RewardSaveRepository
}

func NewRewardService(
	re RewardExistsRepository,
	rs RewardSaveRepository,
) *RewardService {
	return &RewardService{re: re, rs: rs}
}

func (svc *RewardService) Register(
	ctx context.Context, reward *types.RegisterRewardRequest,
) error {
	exists, err := svc.re.Exists(ctx, reward.Match)
	if err != nil {
		return ErrRegisterRewardInternal
	}
	if exists {
		return ErrGoodRewardAlreadyExists
	}
	err = svc.rs.Save(ctx, reward.Match, reward.Reward, reward.RewardType)
	if err != nil {
		return ErrRegisterRewardInternal
	}
	return nil
}

var (
	ErrGoodRewardAlreadyExists = errors.New("вознаграждение для указанного товара уже существует")
	ErrRegisterRewardInternal  = errors.New("внутренняя ошибка при регистрации вознаграждения")
)
