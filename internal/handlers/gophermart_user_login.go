package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type GophermartUserLoginService interface {
	Login(ctx context.Context, req *types.GophermartUserRegisterRequest) (*types.GophermartUserRegisterResponse, *types.APIStatus, *types.APIStatus)
}

func GophermartUserLoginHandler(svc GophermartUserLoginService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GophermartUserRegisterRequest
		if err := decodeJSONRequest(w, r, &req); err != nil {
			return
		}
		resp, status, err := svc.Login(r.Context(), &req)
		if err != nil {
			handleError(w, err)
			return
		}
		setAuthorizationHeader(w, resp.Token)
		sendTextPlainResponse(w, status)
	}
}
