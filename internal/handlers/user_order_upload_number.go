package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserOrderUploadNumberService interface {
	Upload(ctx context.Context, req *types.UserOrderUploadNumberRequest) (*types.APIStatus, *types.APIStatus)
}

func UserOrderUploadNumberHandler(
	svc UserOrderUploadNumberService,
	loginProvider func(ctx context.Context) (string, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, err := getLoginFromContext(w, r, loginProvider)
		if err != nil {
			return
		}

		req := types.UserOrderUploadNumberRequest{
			Login:  login,
			Number: getURLParam(r, "number"),
		}

		req.Login = login

		status, errStatus := svc.Upload(r.Context(), &req)
		if errStatus != nil {
			sendTextPlainResponse(w, errStatus.Message, errStatus.StatusCode)
			return
		}

		sendTextPlainResponse(w, status.Message, status.StatusCode)
	}
}
