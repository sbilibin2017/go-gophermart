package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/services"
)

type AccrualOrderGetService interface {
	Get(ctx context.Context, req *services.AccrualOrderGetRequest) (*services.AccrualOrderGetResponse, *services.APIStatus, *services.APIStatus)
}

func AccrualOrderGetHandler(svc AccrualOrderGetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		number := getURLParam(r, "number")
		req := &services.AccrualOrderGetRequest{Order: number}
		resp, status, err := svc.Get(r.Context(), req)
		if err != nil {
			handleError(w, err)
			return
		}
		encodeJSONResponse(w, resp, status)
	}
}
