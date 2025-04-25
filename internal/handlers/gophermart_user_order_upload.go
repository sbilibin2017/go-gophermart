package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserOrderUploadRequest struct {
	Number string `json:"number" validate:"required,luhn"`
}

type UserOrderUploadResponse struct {
	Message string `json:"message"`
}

type UserOrderUploadService interface {
	Upload(ctx context.Context, req *UserOrderUploadRequest, login string) (*UserOrderUploadResponse, *types.APIStatus)
}

func UserOrderUploadHandler(svc UserOrderUploadService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, err := middlewares.GetLogin(r)
		if err != nil {
			handleErrorResponse(w, "User not authenticated", http.StatusUnauthorized)
			return
		}
		number := getPathParam(r, "number")
		resp, status := svc.Upload(r.Context(), &UserOrderUploadRequest{Number: number}, login)
		if status != nil {
			handleErrorResponse(w, status.Message, status.Status)
			return
		}
		encodeJSONResponse(w, resp, http.StatusOK)
	}
}
