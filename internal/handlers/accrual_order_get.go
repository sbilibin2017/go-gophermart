package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type AccrualOrderGetService interface {
	Get(ctx context.Context, req *types.AccrualOrderGetRequest) (*types.AccrualOrderGetResponse, *types.APIStatus, *types.APIStatus)
}

func AccrualOrderGetHandler(svc AccrualOrderGetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		number := getURLParam(r, "number")
		req := &types.AccrualOrderGetRequest{Order: number}
		resp, status, err := svc.Get(r.Context(), req)
		if err != nil {
			handleError(w, err)
			return
		}
		encodeJSONResponse(w, resp, status)
	}
}
