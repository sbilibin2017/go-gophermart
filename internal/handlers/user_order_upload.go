package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserOrderUpdloadService interface {
	Upload(ctx context.Context, req *types.UserOrderUploadRequest) (*types.APISuccessStatus, *types.APIErrorStatus)
}

func UserOrderUpdloadHandler(svc UserOrderUpdloadService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := contextutils.GetClaims(r.Context())
		if err != nil {
			sendTextPlainResponse(w, types.ErrUnauthorized.Error(), http.StatusUnauthorized)
			return

		}

		req := types.UserOrderUploadRequest{
			Login:  claims.Login,
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
