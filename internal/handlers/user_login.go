package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserLoginService interface {
	Login(ctx context.Context, req *types.UserLoginRequest) (string, *types.APIStatus, *types.APIStatus)
}

func UserLoginHandler(svc UserLoginService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserLoginRequest

		if err := decodeRequestBody(w, r, &req); err != nil {
			return
		}

		tokenString, status, err := svc.Login(r.Context(), &req)
		if err != nil {
			sendTextPlainResponse(w, err.Message, err.StatusCode)
			return
		}

		setTokenHeader(w, tokenString)
		sendTextPlainResponse(w, status.Message, status.StatusCode)
	}
}
