package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderAcceptOrderExistsByIDRepository interface {
	Exists(ctx context.Context, orderID string) (bool, error)
}

type OrderAcceptOrderSaveRepository interface {
	Save(ctx context.Context, order map[string]any) error
}

type OrderAcceptRewardFilterILikeRepository interface {
	FilterILike(ctx context.Context, match string, fields []string) (map[string]any, error)
}

type OrderAcceptValidator interface {
	Struct(v any) error
}

type OrderAcceptService struct {
	v  OrderAcceptValidator
	oe OrderAcceptOrderExistsByIDRepository
	os OrderAcceptOrderSaveRepository
	rf OrderAcceptRewardFilterILikeRepository
}

func NewOrderAcceptService(
	v OrderAcceptValidator,
	oe OrderAcceptOrderExistsByIDRepository,
	os OrderAcceptOrderSaveRepository,
	rf OrderAcceptRewardFilterILikeRepository,
) *OrderAcceptService {
	return &OrderAcceptService{
		v:  v,
		oe: oe,
		os: os,
		rf: rf,
	}
}

func (svc *OrderAcceptService) Accept(
	ctx context.Context, req *types.OrderAcceptRequest,
) (*types.APIResponse[any], *types.APIError) {
	if err := svc.v.Struct(req); err != nil {
		return &types.APIResponse[any]{
			Status:  http.StatusBadRequest,
			Message: "Invalid reward request",
		}, nil
	}

	exists, err := svc.oe.Exists(ctx, req.Order)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}
	if exists {
		return nil, &types.APIError{
			Status:  http.StatusConflict,
			Message: "Order has already been processed",
		}
	}

	var accrual int64

	for _, good := range req.Goods {
		filtered, err := svc.rf.FilterILike(ctx, good.Description, []string{"reward_type", "reward"})
		if err != nil {
			return nil, &types.APIError{
				Status:  http.StatusInternalServerError,
				Message: "Internal server error while filtering reward information",
			}
		}

		accrual += calcAccrual(
			good.Price,
			filtered["reward"].(int64),
			types.RewardType(filtered["reward_type"].(string)),
		)
	}

	err = svc.os.Save(ctx,
		map[string]any{
			"order_id": req.Order,
			"status":   string(types.OrderStatusRegistered),
			"accrual":  accrual,
		})
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to save the order",
		}
	}

	return &types.APIResponse[any]{
		Status:  http.StatusAccepted,
		Message: "Order accepted for processing",
	}, nil
}

func calcAccrual(price int64, reward int64, rewardType types.RewardType) int64 {
	switch rewardType {
	case types.RewardTypePercent:
		return int64(price * reward / 100)
	case types.RewardTypePoint:
		return reward
	default:
		return 0
	}
}
