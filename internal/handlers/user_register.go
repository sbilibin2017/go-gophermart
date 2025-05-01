package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserRegisterService interface {
	Register(ctx context.Context, req *types.UserRegisterRequest) (string, *types.APISuccessStatus, *types.APIErrorStatus)
}

func UserRegisterHandler(svc UserRegisterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserRegisterRequest
		err := decodeRequestBody(w, r, &req)
		if err != nil {
			return
		}

		tokenString, successStatus, errorStatus := svc.Register(r.Context(), &req)
		if errorStatus != nil {
			sendTextPlainResponse(w, errorStatus.Message, errorStatus.StatusCode)
			return
		}

		setTokenHeader(w, tokenString)
		sendTextPlainResponse(w, successStatus.Message, successStatus.StatusCode)
	}
}
