package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserRegisterRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserRegisterResponse struct {
	Token string `json:"token"`
}

type UserRegisterService interface {
	Register(ctx context.Context, req *UserRegisterRequest) (*UserRegisterResponse, *types.APIStatus)
}

func UserRegisterHandler(svc UserRegisterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UserRegisterRequest
		if err := decodeJSONRequest(w, r, &req); err != nil {
			return
		}
		resp, err := svc.Register(r.Context(), &req)
		if err != nil {
			handleErrorResponse(w, err.Message, err.Status)
			return
		}
		setAuthorizationHeader(w, resp.Token)
		encodeJSONResponse(w, resp, http.StatusOK)
	}
}
