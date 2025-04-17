package services

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/dto"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
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

var (
	ErrGoodRewardAlreadyExists = errors.New("вознаграждение для указанного товара уже существует")
)

func (svc *RegisterRewardService) Register(
	ctx context.Context, reward *domain.Reward,
) error {
	logger.Logger.Info("Начало регистрации награды для товара с match:", reward.Match)

	err := svc.tx.Do(ctx, func(tx *sqlx.Tx) error {
		logger.Logger.Info("Проверка существования награды для товара с match:", reward.Match)
		exists, err := svc.re.Exists(ctx, &dto.RewardExistsFilterDB{Match: reward.Match})
		if err != nil {
			logger.Logger.Error("Ошибка при проверке существования награды:", err)
			return err
		}
		if exists {
			logger.Logger.Info("Награда для товара с match уже существует.")
			return ErrGoodRewardAlreadyExists
		}

		rewardDTO := &dto.RewardDB{
			Match:      reward.Match,
			Reward:     reward.Reward,
			RewardType: string(reward.RewardType),
		}

		logger.Logger.Info("Сохранение награды для товара с match:", reward.Match)
		if err := svc.rs.Save(ctx, rewardDTO); err != nil {
			logger.Logger.Error("Ошибка при сохранении награды:", err)
			return err
		}

		logger.Logger.Info("Награда для товара с match успешно зарегистрирована.")
		return nil
	})

	if errors.Is(err, ErrGoodRewardAlreadyExists) {
		return ErrGoodRewardAlreadyExists
	}

	return err
}
