package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type GophermartUserRegisterService interface {
	Register(ctx context.Context, req *types.GophermartUserUserRegisterRequest) (*types.GophermartUserUserRegisterResponse, *types.APIStatus, *types.APIStatus)
}

func GophermartUserRegisterHandler(svc GophermartUserRegisterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GophermartUserUserRegisterRequest
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
