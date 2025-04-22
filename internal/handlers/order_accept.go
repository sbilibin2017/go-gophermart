package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderAcceptService interface {
	Accept(ctx context.Context, order *types.OrderAcceptRequest) (*types.APIResponse[any], *types.APIError)
}

func OrderAcceptHandler(svc OrderAcceptService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrderAcceptRequest
		if err := decodeJSONRequest(w, r, &req); err != nil {
			return
		}

		resp, err := svc.Accept(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Message, err.Status)
			return
		}

		sendTextResponse(w, resp)
	}
}
