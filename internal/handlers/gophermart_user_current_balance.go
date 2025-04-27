package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type GophermartUserCurrentBalanceService interface {
	Get(ctx context.Context, login string) (*types.GophermartUserCurrentBalanceResponse, *types.APIStatus, *types.APIStatus)
}

func GophermartUserCurrentBalanceHandler(svc GophermartUserCurrentBalanceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login := getUserLoginFromContext(w, r)
		if login == nil {
			return
		}
		resp, status, err := svc.Get(r.Context(), *login)
		if err != nil {
			handleError(w, status)
			return
		}
		encodeJSONResponse(w, resp, status)
	}
}
