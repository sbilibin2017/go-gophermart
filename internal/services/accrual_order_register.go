package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/constants"
)

type AccrualOrderRegisterService struct {
	v  StructValidator
	ro ExistsRepository
	rs SaveRepository
	rm FilterILikeRepository
}

func NewAccrualOrderRegisterService(
	v StructValidator,
	ro ExistsRepository,
	rs SaveRepository,
	rm FilterILikeRepository,
) *AccrualOrderRegisterService {
	return &AccrualOrderRegisterService{
		v:  v,
		ro: ro,
		rs: rs,
		rm: rm,
	}
}

func (svc *AccrualOrderRegisterService) Register(
	ctx context.Context, req *AccrualOrderRegisterRequest,
) (*APIStatus, *APIStatus) {
	if err := svc.v.Struct(req); err != nil {
		return nil, &APIStatus{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid order data",
		}
	}
	exists, err := svc.ro.Exists(ctx, map[string]any{"number": req.Order})
	if err != nil {
		return nil, &APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error checking if order exists",
		}
	}
	if exists {
		return nil, &APIStatus{
			StatusCode: http.StatusConflict,
			Message:    "Order already registered",
		}
	}
	var accrual int64
	for _, good := range req.Goods {
		mechanic, err := svc.rm.FilterILike(ctx, good.Description, []string{"reward", "reward_type"})
		if err != nil {
			return nil, &APIStatus{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error retrieving reward mechanic data",
			}
		}
		if a := calcAccrual(good.Price, mechanic); a != nil {
			accrual += *a
		}
	}
	if err := svc.rs.Save(ctx, map[string]any{"number": req.Order, "accrual": accrual, "status": constants.ORDER_STATUS_REGISTERED}); err != nil {
		return nil, &APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error saving order",
		}
	}
	return &APIStatus{
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

type AccrualOrderRegisterRequest struct {
	Order string `json:"order" validate:"required,luhn"`
	Goods []struct {
		Description string `json:"description" validate:"required"`
		Price       int64  `json:"price" validate:"required,gt=0"`
	} `json:"goods" validate:"required,min=1"`
}
