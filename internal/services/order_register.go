package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderRegisterOrderFilterOneRepository interface {
	FilterOne(ctx context.Context, number string) (*types.OrderDB, error)
}

type OrderRegisterOrderSaveRepository interface {
	Save(ctx context.Context, number string, status string, accrual int64) error
}

type OrderRegisterRewardFilterOneILikeRepository interface {
	FilterOneILike(ctx context.Context, description string) (*types.RewardDB, error)
}

type OrderRegisterOrderValidator interface {
	Struct(v any) error
}

type OrderRegisterService struct {
	v  OrderRegisterOrderValidator
	of OrderRegisterOrderFilterOneRepository
	os OrderRegisterOrderSaveRepository
	rf OrderRegisterRewardFilterOneILikeRepository
}

func NewOrderRegisterService(
	v OrderRegisterOrderValidator,
	of OrderRegisterOrderFilterOneRepository,
	os OrderRegisterOrderSaveRepository,
	rf OrderRegisterRewardFilterOneILikeRepository,
) *OrderRegisterService {
	return &OrderRegisterService{
		v:  v,
		of: of,
		os: os,
		rf: rf,
	}
}

func (svc *OrderRegisterService) Register(
	ctx context.Context, req *types.OrderRequest,
) (*types.APIStatus, *types.APIStatus) {
	if err := svc.v.Struct(req); err != nil {
		valErr := formatValidationError(err)
		return nil, &types.APIStatus{
			StatusCode: http.StatusBadRequest,
			Message:    valErr.Message,
		}
	}

	existingOrder, err := svc.of.FilterOne(ctx, req.Order)
	if err != nil {
		return nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}
	if existingOrder != nil {
		return nil, &types.APIStatus{
			StatusCode: http.StatusConflict,
			Message:    "Order number already exists",
		}
	}

	var accrual int64

	for _, good := range req.Goods {
		reward, err := svc.rf.FilterOneILike(ctx, good.Description)
		if err != nil {
			return nil, &types.APIStatus{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal server error while checking reward",
			}
		}
		if reward == nil {
			return nil, &types.APIStatus{
				StatusCode: http.StatusNotFound,
				Message:    "No reward found for the product",
			}
		}

		goodAccrual, err := calculateAccrual(reward, good)
		if err != nil {
			return nil, &types.APIStatus{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid reward type",
			}
		}

		accrual += goodAccrual
	}

	if err := svc.os.Save(ctx, req.Order, types.OrderAccrualStatusRegistered, accrual); err != nil {
		return nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save order with reward",
		}
	}

	return &types.APIStatus{
		StatusCode: http.StatusOK,
		Message:    "Reward successfully registered",
	}, nil
}

func calculateAccrual(reward *types.RewardDB, good types.Good) (int64, error) {
	switch reward.RewardType {
	case types.RewardTypePercent:
		return int64(float64(good.Price) * (float64(reward.Reward) / 100)), nil
	case types.RewardTypePoint:
		return reward.Reward, nil
	default:
		return 0, errors.New("invalid reward type")
	}
}
