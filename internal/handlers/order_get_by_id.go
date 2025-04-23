package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderGetService interface {
	GetByID(ctx context.Context, req *types.OrderGetByIDRequest) (*types.OrderGetByIDResponse, *types.APIStatus, error)
}

func OrderGetHandler(svc OrderGetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		number := getPathParam(r, "number")

		resp, status, err := svc.GetByID(r.Context(), &types.OrderGetByIDRequest{Number: number})
		if err != nil {
			http.Error(w, status.Message, status.Status)
			return
		}

		encodeResponse(w, resp)
	}
}
