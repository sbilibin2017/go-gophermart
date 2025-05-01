package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserOrderUpdloadService interface {
	Upload(ctx context.Context, req *types.UserOrderUploadRequest) (*types.APISuccessStatus, *types.APIErrorStatus)
}

func UserOrderUpdloadHandler(svc UserOrderUpdloadService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, err := getLoginFromContext(w, r)
		if err != nil {
			return
		}

		req := types.UserOrderUploadRequest{
			Login:  login,
			Number: getURLParam(r, "number"),
		}

		successStatus, errorStatus := svc.Upload(r.Context(), &req)
		if errorStatus != nil {
			sendTextPlainResponse(w, errorStatus.Message, errorStatus.StatusCode)
			return
		}

		sendTextPlainResponse(w, successStatus.Message, successStatus.StatusCode)
	}
}
