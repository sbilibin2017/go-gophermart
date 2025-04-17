package services

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/dto"
)

type RewardExistsRepository interface {
	Exists(ctx context.Context, filter *dto.RewardExistsFilterDB) (bool, error)
}

type RewardSaveRepository interface {
	Save(ctx context.Context, reward *dto.RewardDB) error
}

type Tx interface {
	Do(ctx context.Context, fn func(tx *sqlx.Tx) error) error
}

type RegisterRewardService struct {
	re RewardExistsRepository
	rs RewardSaveRepository
	tx Tx
}

func NewRegisterRewardService(
	re RewardExistsRepository,
	rs RewardSaveRepository,
	tx Tx,
) *RegisterRewardService {
	return &RegisterRewardService{re: re, rs: rs, tx: tx}
}

func (svc *RegisterRewardService) Register(
	ctx context.Context, reward *domain.Reward,
) error {
	return svc.tx.Do(ctx, func(tx *sqlx.Tx) error {
		filter := &dto.RewardExistsFilterDB{
			Match: reward.Match,
		}

		exists, err := svc.re.Exists(ctx, filter)
		if err != nil {
			return err
		}
		if exists {
			return ErrGoodRewardAlreadyExists
		}

		reward := &dto.RewardDB{
			Match:      reward.Match,
			Reward:     reward.Reward,
			RewardType: string(reward.RewardType),
		}

		err = svc.rs.Save(ctx, reward)
		if err != nil {
			return err
		}

		return nil
	})
}

var (
	ErrGoodRewardAlreadyExists = errors.New("вознаграждение для указанного товара уже существует")
)
