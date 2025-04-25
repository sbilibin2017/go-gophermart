package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserLoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

type UserLoginService interface {
	Login(ctx context.Context, req *UserLoginRequest) (*UserLoginResponse, *types.APIStatus)
}

func UserLoginHandler(svc UserLoginService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UserLoginRequest
		if err := decodeJSONRequest(w, r, &req); err != nil {
			return
		}
		resp, err := svc.Login(r.Context(), &req)
		if err != nil {
			handleErrorResponse(w, err.Message, err.Status)
			return
		}
		setAuthorizationHeader(w, resp.Token)
		encodeJSONResponse(w, resp, http.StatusOK)
	}
}
