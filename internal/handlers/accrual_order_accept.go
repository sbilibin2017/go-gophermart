package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderAcceptService interface {
	Accept(ctx context.Context, req *services.OrderAcceptRequest) (*types.APIStatus, *types.APIStatus)
}

func OrderAcceptHandler(svc OrderAcceptService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req services.OrderAcceptRequest
		if err := decodeJSONRequest(w, r, &req); err != nil {
			return
		}
		resp, err := svc.Accept(r.Context(), &req)
		if err != nil {
			handleErrorResponse(w, err.Message, err.Status)
			return
		}
		writeTextPlainResponse(w, resp.Message, resp.Status)
	}
}
