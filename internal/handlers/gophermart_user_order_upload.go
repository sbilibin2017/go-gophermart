package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/jwt"
	"github.com/sbilibin2017/go-gophermart/internal/services"
)

type GophermartUserOrderUploadService interface {
	Upload(ctx context.Context, req *services.GophermartUserOrderUploadRequest, login string) (*services.APIStatus, *services.APIStatus)
}

func GophermartUserOrderUploadHandler(svc GophermartUserOrderUploadService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login := getUserLoginFromContext(w, r, jwt.GetJWTPayload)
		if login == nil {
			return
		}
		status, err := svc.Upload(
			r.Context(),
			&services.GophermartUserOrderUploadRequest{Number: getURLParam(r, "number")},
			*login,
		)
		if err != nil {
			handleError(w, err)
			return
		}
		sendTextPlainResponse(w, status)
	}
}
