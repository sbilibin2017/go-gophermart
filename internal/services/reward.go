package services

import (
	"context"

	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/errors"
	"github.com/sbilibin2017/go-gophermart/internal/models"
)

type RewardExistsRepository interface {
	Exists(ctx context.Context, match string) (bool, error)
}

type RewardSaveRepository interface {
	Save(ctx context.Context, reward *models.RewardDB) error
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
	ctx context.Context, reward *domain.Reward,
) error {
	exists, err := svc.re.Exists(ctx, reward.Match)
	if err != nil {
		return errors.ErrInternal
	}
	if exists {
		return errors.ErrGoodRewardAlreadyExists
	}

	dbReward := models.NewRewardDB(
		reward.Match,
		reward.Reward,
		string(reward.RewardType),
	)

	err = svc.rs.Save(ctx, dbReward)
	if err != nil {
		return errors.ErrInternal
	}

	return nil
}
