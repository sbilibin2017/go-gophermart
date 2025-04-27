package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type GophermartUserOrderUploadService interface {
	Upload(ctx context.Context, req *types.GophermartUserOrderUploadRequest, login string) (*types.APIStatus, *types.APIStatus)
}

func GophermartUserOrderUploadHandler(svc GophermartUserOrderUploadService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login := getUserLoginFromContext(w, r)
		if login == nil {
			return
		}
		status, err := svc.Upload(
			r.Context(),
			&types.GophermartUserOrderUploadRequest{Number: getURLParam(r, "number")},
			*login,
		)
		if err != nil {
			handleError(w, err)
			return
		}
		sendTextPlainResponse(w, status)
	}
}
