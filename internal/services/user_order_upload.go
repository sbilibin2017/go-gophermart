package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserOrderUploadRepository interface {
	Exists(ctx context.Context, filter *repositories.UserOrderExistsFilter) (bool, error)
}

type UserOrderUploadValidator interface {
	Struct(v any) error
}

type UserOrderUploadService struct {
	v    UserOrderUploadValidator
	repo UserOrderUploadRepository
}

func NewUserOrderUploadService(
	v UserOrderUploadValidator,
	repo UserOrderUploadRepository,
) *UserOrderUploadService {
	return &UserOrderUploadService{
		v:    v,
		repo: repo,
	}
}

func (svc *UserOrderUploadService) Upload(
	ctx context.Context, req *UserOrderUploadRequest,
) *types.APIStatus {
	// Используем валидатор для проверки формата запроса
	if err := svc.v.Struct(req); err != nil {
		return &types.APIStatus{
			Status:  http.StatusBadRequest,
			Message: "Invalid order number format",
		}
	}

	// Проверка, был ли этот заказ уже загружен пользователем
	exists, err := svc.repo.Exists(ctx, &repositories.UserOrderExistsFilter{
		Number: req.Number,
		Login:  req.Login,
	})
	if err != nil {
		return &types.APIStatus{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}
	if exists {
		return &types.APIStatus{
			Status:  http.StatusOK,
			Message: "Order already uploaded by this user",
		}
	}

	// Проверка, был ли этот номер загружен другим пользователем
	existsByOther, err := svc.repo.Exists(ctx, &repositories.UserOrderExistsFilter{
		Number: req.Number,
	})
	if err != nil {
		return &types.APIStatus{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}
	if existsByOther {
		return &types.APIStatus{
			Status:  http.StatusConflict,
			Message: "Order already uploaded by another user",
		}
	}

	// Если заказ не был загружен этим пользователем и другим, принимаем его в обработку
	return &types.APIStatus{
		Status:  http.StatusAccepted,
		Message: "Order accepted for processing",
	}
}

type UserOrderUploadRequest struct {
	Number string `json:"number"`
	Login  string `json:"login"`
}
