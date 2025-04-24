package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderAcceptOrderExistsRepository interface {
	Exists(ctx context.Context, number string) (bool, error)
}

type OrderAcceptOrderSaveRepository interface {
	Save(ctx context.Context, order map[string]any) error
}

type OrderAcceptGoodRewardFilterILikeRepository interface {
	FilterILike(ctx context.Context, description string, fields []string) (map[string]any, error)
}

type OrderAcceptValidator interface {
	Struct(v any) error
}

type OrderAcceptService struct {
	v  OrderAcceptValidator
	oe OrderAcceptOrderExistsRepository
	os OrderAcceptOrderSaveRepository
	rf OrderAcceptGoodRewardFilterILikeRepository
}

func NewOrderAcceptService(
	v OrderAcceptValidator,
	oe OrderAcceptOrderExistsRepository,
	os OrderAcceptOrderSaveRepository,
	rf OrderAcceptGoodRewardFilterILikeRepository,

) *OrderAcceptService {
	return &OrderAcceptService{
		oe: oe,
		os: os,
		rf: rf,
		v:  v,
	}
}

func (svc *OrderAcceptService) Accept(
	ctx context.Context, req *types.OrderAcceptRequest,
) (*types.APIStatus, *types.APIStatus) {
	err := svc.v.Struct(req)
	if err != nil {
		return nil, &types.APIStatus{
			Status:  http.StatusBadRequest,
			Message: "Invalid request format",
		}
	}

	exists, err := svc.oe.Exists(ctx, req.Order)
	if err != nil {
		return nil, &types.APIStatus{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}
	if exists {
		return nil, &types.APIStatus{
			Status:  http.StatusConflict,
			Message: "Order already accepted",
		}
	}

	var accrual int64
	for _, good := range req.Goods {
		filtered, err := svc.rf.FilterILike(
			ctx,
			good.Description,
			[]string{"reward", "reward_type"},
		)
		if err != nil {
			return nil, &types.APIStatus{
				Status:  http.StatusInternalServerError,
				Message: "Internal server error",
			}
		}

		a := calcAccrual(good.Price, filtered["reward"].(int64), filtered["reward_type"].(string))
		if a == nil {
			return nil, &types.APIStatus{
				Status:  http.StatusInternalServerError,
				Message: "Internal server error",
			}
		}
		accrual += *a
	}

	orderData := map[string]any{
		"order":   req.Order,
		"status":  string(types.OrderStatusRegistered),
		"accrual": accrual,
	}

	err = svc.os.Save(ctx, orderData)
	if err != nil {
		return nil, &types.APIStatus{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}

	return &types.APIStatus{
		Status:  http.StatusAccepted,
		Message: "Order successfully accepted for processing",
	}, nil
}

func calcAccrual(price int64, reward int64, rewardType string) *int64 {
	switch rewardType {
	case string(types.RewardTypePercent):
		p := price * reward / 100
		return &p
	case string(types.RewardTypePoint):
		p := reward
		return &p
	default:
		return nil
	}
}
