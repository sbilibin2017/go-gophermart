package services

import (
	"context"
	"errors"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

var (
	ErrRegisterOrderAlreadyExists   = errors.New("order already exists")
	ErrRegisterOrderIsNotRegistered = errors.New("order is not registered")
)

type RegisterOrderExistsRepository interface {
	Exists(ctx context.Context, filter map[string]any) (bool, error)
}

type RegisterOrderRewardFilterRepository interface {
	Filter(ctx context.Context, filter map[string]any, fields []string) (map[string]any, error)
}

type RegisterOrderSaveRepository interface {
	Save(ctx context.Context, data map[string]any) error
}

type RegisterOrderService struct {
	rf RegisterOrderRewardFilterRepository
	oe RegisterOrderExistsRepository
	os RegisterOrderSaveRepository
}

func NewRegisterOrderService(
	rf RegisterOrderRewardFilterRepository,
	oe RegisterOrderExistsRepository,
	os RegisterOrderSaveRepository,
) *RegisterOrderService {
	return &RegisterOrderService{
		rf: rf,
		oe: oe,
		os: os,
	}
}

func (svc *RegisterOrderService) Register(
	ctx context.Context, order *types.RegisterOrderRequest,
) error {
	orderExists, err := svc.oe.Exists(ctx, map[string]any{"order_id": order.Order})
	if err != nil {
		return ErrRegisterOrderIsNotRegistered
	}
	if orderExists {
		return ErrRegisterOrderAlreadyExists
	}

	price := int64(0)
	for _, good := range order.Goods {
		reward, err := svc.rf.Filter(
			ctx, map[string]any{
				"reward_id": good.Description,
			},
			[]string{"reward", "reward_type"},
		)
		if err != nil {
			return ErrRegisterOrderIsNotRegistered
		}
		if reward == nil {
			return ErrRegisterOrderIsNotRegistered
		}
		price += applyReward(good.Price, reward["reward_type"], reward["reward"])
	}

	orderData := map[string]any{
		"order_id": order.Order,
		"price":    price,
	}
	err = svc.os.Save(ctx, orderData)
	if err != nil {
		return ErrRegisterOrderIsNotRegistered
	}

	return nil
}

func applyReward(goodPrice int64, rewardType any, rewardValue any) int64 {
	var price int64

	switch rewardType {
	case "%":
		price = int64(goodPrice * rewardValue.(int64) / 100)
	case "pt":
		price = rewardValue.(int64)
	}

	return price
}
