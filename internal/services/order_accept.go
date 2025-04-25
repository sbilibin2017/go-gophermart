package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/constants"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderAcceptOrderExistsRepository interface {
	Exists(ctx context.Context, order *repositories.OrderExistsNumber) (bool, error)
}

type OrderAcceptOrderSaveRepository interface {
	Save(ctx context.Context, order *repositories.OrderSave) error
}

type OrderAcceptGoodRewardFilterILikeRepository interface {
	FilterILike(ctx context.Context, description *repositories.RewardFilterILike) (*repositories.RewardFilterILikeDB, error)
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
		v:  v,
		oe: oe,
		os: os,
		rf: rf,
	}
}

func (svc *OrderAcceptService) Accept(
	ctx context.Context, req *OrderAcceptRequest,
) (*types.APIStatus, *types.APIStatus) {
	err := svc.v.Struct(req)
	if err != nil {
		return nil, &types.APIStatus{
			Status:  http.StatusBadRequest,
			Message: "Invalid request format",
		}
	}
	exists, err := svc.oe.Exists(ctx, &repositories.OrderExistsNumber{Number: req.Order})
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
			&repositories.RewardFilterILike{
				Description: good.Description,
			},
		)
		if err != nil {
			return nil, &types.APIStatus{
				Status:  http.StatusInternalServerError,
				Message: "Internal server error",
			}
		}
		if filtered.RewardType == "" {
			return nil, &types.APIStatus{
				Status:  http.StatusInternalServerError,
				Message: "Invalid reward data",
			}
		}
		accrualAmount := calcAccrual(good.Price, filtered.Reward, filtered.RewardType)
		if accrualAmount == nil {
			return nil, &types.APIStatus{
				Status:  http.StatusInternalServerError,
				Message: "Error in accrual calculation",
			}
		}
		accrual += *accrualAmount
	}
	orderData := &repositories.OrderSave{
		Number:  req.Order,
		Status:  constants.ORDER_STATUS_REGISTERED,
		Accrual: &accrual,
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
	case constants.REWARD_TYPE_PERCENT:
		p := price * reward / 100
		return &p
	case constants.REWARD_TYPE_POINT:
		return &reward
	default:
		return nil
	}
}

type OrderAcceptRequest struct {
	Order string `json:"order" validate:"required,luhn"`
	Goods []Good `json:"goods" validate:"required,dive,required"`
}

type Good struct {
	Description string `json:"description" validate:"required"`
	Price       int64  `json:"price" validate:"required,gt=0"`
}
