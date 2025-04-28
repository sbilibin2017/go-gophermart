package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

type GophermartUserOrderUploadService struct {
	v  StructValidator
	ro ExistsRepository
	rs SaveRepository
}

func NewGophermartUserOrderUploadService(
	v StructValidator,
	ro ExistsRepository,
	rs SaveRepository,
) *GophermartUserOrderUploadService {
	return &GophermartUserOrderUploadService{
		v:  v,
		ro: ro,
		rs: rs,
	}
}

func (svc *GophermartUserOrderUploadService) Upload(
	ctx context.Context, req *GophermartUserOrderUploadRequest, login string,
) (*APIStatus, *APIStatus) {
	if err := svc.v.Struct(req); err != nil {
		return nil, &APIStatus{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid order number",
		}
	}
	exists, err := svc.ro.Exists(ctx, map[string]any{"login": login, "order": req.Number})
	if err != nil {
		logger.Logger.Errorf("Error checking if order number exists: %v", err)
		return nil, &APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "User order number is not accepted for processing",
		}
	}
	if exists {
		return nil, &APIStatus{
			StatusCode: http.StatusOK,
			Message:    "Order number already uploaded by this user",
		}
	}
	otherExists, err := svc.ro.Exists(ctx, map[string]any{"order": req.Number})
	if err != nil {
		logger.Logger.Errorf("Error checking if order number exists for other users: %v", err)
		return nil, &APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "User order number is not accepted for processing",
		}
	}
	if otherExists {
		return nil, &APIStatus{
			StatusCode: http.StatusConflict,
			Message:    "Order number already uploaded by another user",
		}
	}
	if err := svc.rs.Save(ctx, map[string]any{"login": login, "order": req.Number}); err != nil {
		logger.Logger.Errorf("Error saving order number: %v", err)
		return nil, &APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "User order number is not accepted for processing",
		}
	}
	return &APIStatus{
		StatusCode: http.StatusAccepted,
		Message:    "Order number accepted for processing",
	}, nil
}

type GophermartUserOrderUploadRequest struct {
	Number string `json:"number" validate:"required,luhn"`
}
