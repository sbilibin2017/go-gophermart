package services

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/db"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
)

type RewardExistsRepository interface {
	Exists(ctx context.Context, match string) (bool, error)
}

type RewardSaveRepository interface {
	Save(ctx context.Context, match string, reward uint, rewardType string) error
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

func (svc *RegisterRewardService) Register(
	ctx context.Context, reward *domain.Reward,
) error {
	return db.WithTx(ctx, svc.db, func(tx *sqlx.Tx) error {
		exists, err := svc.re.Exists(ctx, reward.Match)
		if err != nil {
			return domain.ErrFailedToRegisterReward
		}
		if exists {
			return domain.ErrRewardKeyAlreadyRegistered
		}

		err = svc.rs.Save(ctx, reward.Match, reward.Reward, string(reward.RewardType))
		if err != nil {
			return domain.ErrFailedToRegisterReward
		}

		return nil
	})
}
