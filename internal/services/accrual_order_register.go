package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/constants"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type AccrualOrderRegisterExistsRepository interface {
	Exists(ctx context.Context, number string) (bool, error)
}

type AccrualOrderRegisterSaveRepository interface {
	Save(ctx context.Context, number string, accrual int64, status string) error
}

type AccrualOrderRegisterRewardMechanicFindRepository interface {
	FilterILikeByDescription(ctx context.Context, description string, fields []string) (map[string]any, error)
}

type AccrualOrderRegisterValidator interface {
	Struct(v any) error
}

type AccrualOrderRegisterService struct {
	v  AccrualOrderRegisterValidator
	ro AccrualOrderRegisterExistsRepository
	rs AccrualOrderRegisterSaveRepository
	rm AccrualOrderRegisterRewardMechanicFindRepository
}

func NewAccrualOrderRegisterService(
	v AccrualOrderRegisterValidator,
	ro AccrualOrderRegisterExistsRepository,
	rs AccrualOrderRegisterSaveRepository,
	rm AccrualOrderRegisterRewardMechanicFindRepository,
) *AccrualOrderRegisterService {
	return &AccrualOrderRegisterService{
		v:  v,
		ro: ro,
		rs: rs,
		rm: rm,
	}
}

func (svc *AccrualOrderRegisterService) Register(
	ctx context.Context, req *types.AccrualOrderRegisterRequest,
) (*types.APIStatus, *types.APIStatus) {
	if err := svc.v.Struct(req); err != nil {
		return nil, &types.APIStatus{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid order data",
		}
	}

	exists, err := svc.ro.Exists(ctx, req.Order)
	if err != nil {
		return nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Order is not registered",
		}
	}
	if exists {
		return nil, &types.APIStatus{
			StatusCode: http.StatusConflict,
			Message:    "Order already registered",
		}
	}

	var accrual int64
	for _, good := range req.Goods {
		mechanic, err := svc.rm.FilterILikeByDescription(
			ctx, good.Description, []string{"reward", "reward_type"},
		)
		if err != nil {
			return nil, &types.APIStatus{
				StatusCode: http.StatusInternalServerError,
				Message:    "Order is not registered",
			}
		}
		if a := calcAccrual(good.Price, mechanic); a != nil {
			accrual += *a
		}
	}
	if err := svc.rs.Save(ctx, req.Order, accrual, constants.ORDER_STATUS_REGISTERED); err != nil {
		return nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Order is not registered",
		}
	}
	return &types.APIStatus{
		StatusCode: http.StatusAccepted,
		Message:    "Order successfully registered",
	}, nil
}

func calcAccrual(price int64, mechanic map[string]any) *int64 {
	switch mechanic["reward_type"] {
	case constants.REWARD_TYPE_PERCENT:
		v := int64(price * mechanic["reward"].(int64) / 100)
		return &v
	case constants.REWARD_TYPE_POINT:
		v := int64(price)
		return &v
	default:
		return nil
	}
}
