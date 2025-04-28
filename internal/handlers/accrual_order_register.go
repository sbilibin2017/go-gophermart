package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/services"
)

type AccrualOrderRegisterService interface {
	Register(ctx context.Context, req *services.AccrualOrderRegisterRequest) (*services.APIStatus, *services.APIStatus)
}

func AccrualOrderRegisterHandler(svc AccrualOrderRegisterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req services.AccrualOrderRegisterRequest
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
