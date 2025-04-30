package services

import (
	"context"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderRegisterOrderExists interface {
	Exists(ctx context.Context, number string) (bool, error)
}

type OrderRegisterGoodRewardFindILike interface {
	FindILike(ctx context.Context, match string) (*types.GoodReward, error)
}

type OrderRegisterOrderSave interface {
	Save(ctx context.Context, order *types.Order) error
}

type OrderRegisterService struct {
	oe OrderRegisterOrderExists
	os OrderRegisterOrderSave
	of OrderRegisterGoodRewardFindILike
}

func NewOrderRegisterService(
	oe OrderRegisterOrderExists,
	os OrderRegisterOrderSave,
	of OrderRegisterGoodRewardFindILike,
) *OrderRegisterService {
	return &OrderRegisterService{
		oe: oe,
		os: os,
		of: of,
	}
}

func (s *OrderRegisterService) Register(
	ctx context.Context, req *types.OrderRequest,
) error {
	exists, err := s.oe.Exists(ctx, req.Order)
	if err != nil {
		return err
	}
	if exists {
		return types.ErrOrderAlreadyExists
	}

	order := types.Order{
		Number: req.Order,
		Status: types.ORDER_STATUS_REGISTERED,
	}

	var accrual int64
	for _, good := range req.Goods {
		goodReward, err := s.of.FindILike(ctx, good.Description)
		if err != nil {
			return types.ErrInternal
		}

		switch goodReward.RewardType {
		case types.REWARD_TYPE_PERCENT:
			accrual += (good.Price * goodReward.Reward) / 100
		case types.REWARD_TYPE_POINT:
			accrual += goodReward.Reward
		}
	}

	order.Accrual = &accrual

	if err := s.os.Save(ctx, &order); err != nil {
		return err
	}

	return nil
}
