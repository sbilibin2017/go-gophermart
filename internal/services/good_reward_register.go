package services

import (
	"context"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type GoodRewardRegisterGoodRewardExists interface {
	Exists(ctx context.Context, match string) (bool, error)
}

type GoodRewardRegisterGoodRewardSave interface {
	Save(ctx context.Context, goodReward *types.GoodReward) error
}

type GoodRewardRegisterService struct {
	gre GoodRewardRegisterGoodRewardExists
	grs GoodRewardRegisterGoodRewardSave
}

func NewGoodRewardRegisterService(
	gre GoodRewardRegisterGoodRewardExists,
	grs GoodRewardRegisterGoodRewardSave,
) *GoodRewardRegisterService {
	return &GoodRewardRegisterService{
		gre: gre,
		grs: grs,
	}
}

func (s *GoodRewardRegisterService) Register(
	ctx context.Context, goodReward *types.GoodReward,
) error {
	exists, err := s.gre.Exists(ctx, goodReward.Match)
	if err != nil {
		return err
	}
	if exists {
		return types.ErrGoodRewardAlreadyExists
	}
	return s.grs.Save(ctx, goodReward)
}
