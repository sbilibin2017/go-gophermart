package usecases

import (
	"context"

	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/dto"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

type RegisterRewardService interface {
	Register(ctx context.Context, reward *domain.Reward) error
}

type RegisterRewardValidator interface {
	Struct(i interface{}) error
}

type RegisterRewardUsecase struct {
	svc RegisterRewardService
	val RegisterRewardValidator
}

func NewRegisterRewardUsecase(
	svc RegisterRewardService,
	val RegisterRewardValidator,
) *RegisterRewardUsecase {
	return &RegisterRewardUsecase{
		svc: svc,
		val: val,
	}
}

func (uc *RegisterRewardUsecase) Execute(
	ctx context.Context, req *dto.RegisterRewardRequest,
) (*dto.RegisterRewardResponse, error) {
	logger.Logger.Info("Запуск процесса регистрации награды для товара:", req.Match)

	err := uc.val.Struct(req)
	if err != nil {
		logger.Logger.Info("Ошибка валидации данных для товара:", req.Match, "Ошибка:", err)
		return nil, err
	}
	logger.Logger.Info("Данные успешно валидированы для товара:", req.Match)

	reward := &domain.Reward{
		Match:      req.Match,
		Reward:     req.Reward,
		RewardType: domain.RewardType(req.RewardType),
	}
	logger.Logger.Info("Преобразование данных в доменную модель для товара:", req.Match)

	err = uc.svc.Register(ctx, reward)
	if err != nil {
		logger.Logger.Info("Ошибка при регистрации награды для товара:", req.Match, "Ошибка:", err)
		return nil, err
	}
	logger.Logger.Info("Награда для товара успешно зарегистрирована:", req.Match)

	return &dto.RegisterRewardResponse{
		Message: "Информация о вознаграждении за товар зарегистрирована",
	}, nil
}
