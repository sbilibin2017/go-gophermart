package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderGetService interface {
	Get(ctx context.Context, req *types.OrderGetRequest) (*types.OrderGetResponse, *types.APIStatus, error)
}

func OrderGetHandler(svc OrderGetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		number := getPathParam(r, "number")

		resp, status, err := svc.Get(r.Context(), &types.OrderGetRequest{Number: number})
		if err != nil {
			http.Error(w, status.Message, status.Status)
			return
		}

		encodeResponse(w, resp)
	}
}
