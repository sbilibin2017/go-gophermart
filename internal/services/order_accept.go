package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

// Константы для сообщений об ошибках
const (
	ErrInvalidRewardRequest           = "Invalid reward request"
	ErrInternalServerError            = "Internal server error"
	ErrInternalServerErrorFilter      = "Internal server error while filtering reward information"
	ErrOrderAlreadyProcessed          = "Order has already been processed"
	ErrFailedToSaveOrder              = "Failed to save the order"
	SuccessOrderAcceptedForProcessing = "Order accepted for processing"
)

// Интерфейсы
type OrderAcceptOrderExistsByIDRepository interface {
	Exists(ctx context.Context, orderID string) (bool, error)
}

type OrderAcceptOrderSaveRepository interface {
	Save(ctx context.Context, orderID string, status string, accrual int64) error
}

type OrderAcceptGoodRewardFilterILikeRepository interface {
	FilterILike(ctx context.Context, match string, fields []string) (*types.GoodRewardDB, error)
}

type OrderAcceptValidator interface {
	Struct(v any) error
}

type OrderAcceptService struct {
	v  OrderAcceptValidator
	oe OrderAcceptOrderExistsByIDRepository
	os OrderAcceptOrderSaveRepository
	rf OrderAcceptGoodRewardFilterILikeRepository
}

func NewOrderAcceptService(
	v OrderAcceptValidator,
	oe OrderAcceptOrderExistsByIDRepository,
	os OrderAcceptOrderSaveRepository,
	rf OrderAcceptGoodRewardFilterILikeRepository,
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
) (*types.APIStatus, error) {
	if err := svc.v.Struct(req); err != nil {
		return &types.APIStatus{
			Status:  http.StatusBadRequest,
			Message: ErrInvalidRewardRequest,
		}, nil
	}

	exists, err := svc.oe.Exists(ctx, req.Order)
	if err != nil {
		return &types.APIStatus{
			Status:  http.StatusInternalServerError,
			Message: ErrInternalServerError,
		}, nil
	}
	if exists {
		return &types.APIStatus{
			Status:  http.StatusConflict,
			Message: ErrOrderAlreadyProcessed,
		}, nil
	}

	var accrual int64

	for _, good := range req.Goods {
		filtered, err := svc.rf.FilterILike(ctx, good.Description, []string{"reward_type", "reward"})
		if err != nil {
			return &types.APIStatus{
				Status:  http.StatusInternalServerError,
				Message: ErrInternalServerErrorFilter,
			}, nil
		}

		accrual += calcAccrual(
			good.Price,
			filtered.Reward,                       // Assuming filtered is a pointer to GoodRewardDB and contains Reward
			types.RewardType(filtered.RewardType), // Assuming filtered.RewardType is a string
		)
	}

	err = svc.os.Save(ctx,
		req.Order,
		string(types.OrderStatusRegistered),
		accrual,
	)
	if err != nil {
		return &types.APIStatus{
			Status:  http.StatusInternalServerError,
			Message: ErrFailedToSaveOrder,
		}, nil
	}

	return &types.APIStatus{
		Status:  http.StatusAccepted,
		Message: SuccessOrderAcceptedForProcessing,
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
