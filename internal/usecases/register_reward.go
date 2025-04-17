package usecases

import (
	"context"

	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/dto"
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
	err := uc.val.Struct(req)

	if err != nil {
		return nil, err
	}

	reward := &domain.Reward{
		Match:      req.Match,
		Reward:     req.Reward,
		RewardType: domain.RewardType(req.RewardType),
	}

	err = uc.svc.Register(ctx, reward)

	if err != nil {
		return nil, err
	}

	return &dto.RegisterRewardResponse{
		Message: "Информация о вознаграждении за товар зарегистрирована",
	}, nil

}
