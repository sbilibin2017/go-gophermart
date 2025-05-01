package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserAuthService interface {
	Authenticate(ctx context.Context, req *types.UserAuthRequest) (string, *types.APISuccessStatus, *types.APIErrorStatus)
}

func UserAuthHandler(svc UserAuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserAuthRequest
		err := decodeRequestBody(w, r, &req)
		if err != nil {
			return
		}

		tokenString, successStatus, errorStatus := svc.Authenticate(r.Context(), &req)
		if errorStatus != nil {
			sendTextPlainResponse(w, errorStatus.Message, errorStatus.StatusCode)
			return
		}

		setTokenHeader(w, tokenString)
		sendTextPlainResponse(w, successStatus.Message, successStatus.StatusCode)
	}
}
