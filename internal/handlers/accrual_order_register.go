package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type AccrualOrderRegisterService interface {
	Register(ctx context.Context, req *types.AccrualOrderRegisterRequest) (*types.APIStatus, *types.APIStatus)
}

func AccrualOrderRegisterHandler(svc AccrualOrderRegisterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AccrualOrderRegisterRequest
		if err := decodeJSONRequest(w, r, &req); err != nil {
			return
		}
		resp, err := svc.Register(r.Context(), &req)
		if err != nil {
			handleError(w, err)
			return
		}
		sendTextPlainResponse(w, resp)
	}
}
