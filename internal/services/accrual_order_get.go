package services

import (
	"context"
	"net/http"
)

type AccrualOrderGetService struct {
	v  StructValidator
	ro FilterRepository
}

func NewAccrualOrderGetService(
	v StructValidator,
	ro FilterRepository,
) *AccrualOrderGetService {
	return &AccrualOrderGetService{
		v:  v,
		ro: ro,
	}
}

func (svc *AccrualOrderGetService) Get(
	ctx context.Context, req *AccrualOrderGetRequest,
) (*AccrualOrderGetResponse, *APIStatus, *APIStatus) {
	if err := svc.v.Struct(req); err != nil {
		return nil, nil, &APIStatus{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid order number format",
		}
	}
	order, err := svc.ro.Filter(ctx, map[string]any{"number": req.Order}, []string{"number", "status", "accrual"})
	if err != nil {
		return nil, nil, &APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error retrieving accrual data",
		}
	}
	if order == nil {
		return nil, nil, &APIStatus{
			StatusCode: http.StatusNotFound,
			Message:    "Order is not registered",
		}
	}
	var response *AccrualOrderGetResponse
	err = mapToStruct(response, order)
	if err != nil {
		return nil, nil, &APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error mapping accrual data",
		}
	}
	return response, &APIStatus{
		StatusCode: http.StatusOK,
		Message:    "Success",
	}, nil
}

type AccrualOrderGetRequest struct {
	Order string `json:"order" validate:"required,luhn"`
}

type AccrualOrderGetResponse struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual *int64 `json:"accrual,omitempty"`
}
