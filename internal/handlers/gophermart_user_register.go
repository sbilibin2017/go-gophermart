package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/services"
)

type GophermartUserRegisterService interface {
	Register(ctx context.Context, req *services.GophermartUserRegisterRequest) (*services.GophermartUserRegisterResponse, *services.APIStatus, *services.APIStatus)
}

func GophermartUserRegisterHandler(svc GophermartUserRegisterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req services.GophermartUserRegisterRequest
		if err := decodeJSONRequest(w, r, &req); err != nil {
			return
		}
		resp, status, err := svc.Register(r.Context(), &req)
		if err != nil {
			handleError(w, err)
			return
		}
		setAuthorizationHeader(w, resp.Token)
		sendTextPlainResponse(w, status)
	}
}
