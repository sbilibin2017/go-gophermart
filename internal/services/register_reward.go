package services

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/storage"
)

type RewardExistsRepository interface {
	Exists(ctx context.Context, tx *sqlx.Tx, filter *repositories.RewardExistsFilter) (bool, error)
}

type RewardSaveRepository interface {
	Save(ctx context.Context, tx *sqlx.Tx, reward *repositories.RewardSave) error
}

type RewardService struct {
	re RewardExistsRepository
	rs RewardSaveRepository
	db *sqlx.DB
}

func NewRewardService(
	re RewardExistsRepository,
	rs RewardSaveRepository,
	db *sqlx.DB,
) *RewardService {
	return &RewardService{re: re, rs: rs, db: db}
}

func (svc *RewardService) Register(
	ctx context.Context, reward *domain.Reward,
) error {
	return storage.WithTx(ctx, svc.db, func(tx *sqlx.Tx) error {
		filter := &repositories.RewardExistsFilter{
			Match: reward.Match,
		}

		exists, err := svc.re.Exists(ctx, tx, filter)
		if err != nil {
			return ErrRegisterRewardInternal
		}
		if exists {
			return ErrGoodRewardAlreadyExists
		}

		save := &repositories.RewardSave{
			Match:      reward.Match,
			Reward:     reward.Reward,
			RewardType: string(reward.RewardType),
		}

		err = svc.rs.Save(ctx, tx, save)
		if err != nil {
			return ErrRegisterRewardInternal
		}
		return nil
	})
}

var (
	ErrGoodRewardAlreadyExists = errors.New("вознаграждение для указанного товара уже существует")
	ErrRegisterRewardInternal  = errors.New("внутренняя ошибка при регистрации вознаграждения")
)
