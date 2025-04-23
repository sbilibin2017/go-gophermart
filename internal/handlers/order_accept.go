package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderAcceptService interface {
	Accept(ctx context.Context, order *types.OrderAcceptRequest) (*types.APIStatus, *types.APIStatus, error)
}

func OrderAcceptHandler(svc OrderAcceptService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrderAcceptRequest
		if err := decodeRequest(w, r, &req); err != nil {
			return
		}
		resp, status, err := svc.Accept(r.Context(), &req)
		if err != nil {
			http.Error(w, status.Message, status.Status)
			return
		}
		writeTextPlainResponse(w, resp.Status, resp.Message)
	}
}
