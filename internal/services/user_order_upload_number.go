package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/services/validation"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserOrderUploadNumberUserOrderFilterOneRepository interface {
	FilterOne(ctx context.Context, number string, login *string) (*types.UserOrderDB, error)
}

type UserOrderUploadNumberUserOrderSaveRepository interface {
	Save(ctx context.Context, number string, login string, status string) error
}

type UserOrderUploadNumberValidator interface {
	Struct(v any) error
}

type UserOrderUploadNumberService struct {
	v   UserOrderUploadNumberValidator
	uof UserOrderUploadNumberUserOrderFilterOneRepository
	uos UserOrderUploadNumberUserOrderSaveRepository
}

func NewUserOrderUploadNumberService(
	v UserOrderUploadNumberValidator,
	uof UserOrderUploadNumberUserOrderFilterOneRepository,
	uos UserOrderUploadNumberUserOrderSaveRepository,
) *UserOrderUploadNumberService {
	return &UserOrderUploadNumberService{
		v:   v,
		uof: uof,
		uos: uos,
	}
}

func (svc *UserOrderUploadNumberService) Upload(
	ctx context.Context, req *types.UserOrderUploadNumberRequest,
) (*types.APIStatus, *types.APIStatus) {
	if err := svc.v.Struct(req); err != nil {
		valErr := formatValidationError(err)
		if valErr != nil {
			if validation.IsLuhnValidationError(valErr) {
				return nil, &types.APIStatus{
					StatusCode: http.StatusUnprocessableEntity,
					Message:    valErr.Message,
				}
			}
			return nil, &types.APIStatus{
				StatusCode: http.StatusBadRequest,
				Message:    valErr.Message,
			}
		}
		return nil, &types.APIStatus{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request structure",
		}
	}

	existingByUser, err := svc.uof.FilterOne(ctx, req.Number, &req.Login)
	if err != nil {
		return nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to check order ownership",
		}
	}
	if existingByUser != nil {
		return &types.APIStatus{
			StatusCode: http.StatusOK,
			Message:    "Order number already uploaded by this user",
		}, nil
	}

	existingByAny, err := svc.uof.FilterOne(ctx, req.Number, nil)
	if err != nil {
		return nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to check global order existence",
		}
	}
	if existingByAny != nil {
		return nil, &types.APIStatus{
			StatusCode: http.StatusConflict,
			Message:    "Order number already uploaded by another user",
		}
	}

	err = svc.uos.Save(ctx, req.Number, req.Login, types.GOPHERMART_USER_ORDER_STATUS_NEW)
	if err != nil {
		return nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save new order number",
		}
	}

	return &types.APIStatus{
		StatusCode: http.StatusAccepted,
		Message:    "Order number accepted for processing",
	}, nil
}
