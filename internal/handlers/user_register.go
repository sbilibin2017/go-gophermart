package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserRegisterService interface {
	Register(ctx context.Context, req *types.UserRegisterRequest) (string, *types.APIStatus, *types.APIStatus)
}

func UserRegisterHandler(svc UserRegisterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserRegisterRequest

		if err := decodeRequestBody(w, r, &req); err != nil {
			return
		}

		tokenString, status, err := svc.Register(r.Context(), &req)
		if err != nil {
			sendTextPlainResponse(w, err.Message, err.StatusCode)
			return
		}

		setTokenHeader(w, tokenString)
		sendTextPlainResponse(w, status.Message, status.StatusCode)
	}
}
