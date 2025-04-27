package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type GophermartUserOrderExistsRepository interface {
	ExistsByLoginAndNumber(ctx context.Context, login *string, number string) (bool, error)
}

type GophermartUserOrderSaveRepository interface {
	Save(ctx context.Context, login, number string) error
}

type GophermartUserOrderUploadValidator interface {
	Struct(v any) error
}

type GophermartUserOrderUploadService struct {
	v  GophermartUserOrderUploadValidator
	ro GophermartUserOrderExistsRepository
	rs GophermartUserOrderSaveRepository
}

func NewGophermartUserOrderUploadService(
	v GophermartUserOrderUploadValidator,
	ro GophermartUserOrderExistsRepository,
	rs GophermartUserOrderSaveRepository,
) *GophermartUserOrderUploadService {
	return &GophermartUserOrderUploadService{
		v:  v,
		ro: ro,
		rs: rs,
	}
}

func (svc *GophermartUserOrderUploadService) Upload(
	ctx context.Context, req *types.GophermartUserOrderUploadRequest, login string,
) (*types.APIStatus, *types.APIStatus) {
	if err := svc.v.Struct(req); err != nil {
		return nil, &types.APIStatus{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid order number",
		}
	}
	exists, err := svc.ro.ExistsByLoginAndNumber(ctx, &login, req.Number)
	if err != nil {
		logger.Logger.Errorf("Error checking if order number exists: %v", err)
		return nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "User order number is not accepted for processing",
		}
	}
	if exists {
		return nil, &types.APIStatus{
			StatusCode: http.StatusOK,
			Message:    "Order number already uploaded by this user",
		}
	}
	otherExists, err := svc.ro.ExistsByLoginAndNumber(ctx, nil, req.Number)
	if err != nil {
		logger.Logger.Errorf("Error checking if order number exists for other users: %v", err)
		return nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "User order number is not accepted for processing",
		}
	}
	if otherExists {
		return nil, &types.APIStatus{
			StatusCode: http.StatusConflict,
			Message:    "Order number already uploaded by another user",
		}
	}
	if err := svc.rs.Save(ctx, login, req.Number); err != nil {
		logger.Logger.Errorf("Error saving order number: %v", err)
		return nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "User order number is not accepted for processing",
		}
	}
	return &types.APIStatus{
		StatusCode: http.StatusAccepted,
		Message:    "Order number accepted for processing",
	}, nil
}
