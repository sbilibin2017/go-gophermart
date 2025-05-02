package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserOrderUploadListService interface {
	List(ctx context.Context, login string) ([]*types.UserOrderUploadedListResponse, *types.APIStatus)
}

func UserOrderUploadListHandler(
	svc UserOrderUploadListService,
	loginProvider func(ctx context.Context) (string, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, err := getLoginFromContext(w, r, loginProvider)
		if err != nil {
			return
		}

		orders, apiErr := svc.List(r.Context(), login)
		if apiErr != nil {
			sendTextPlainResponse(w, apiErr.Message, apiErr.StatusCode)
			return
		}

		if err := encodeResponseBody(w, orders, http.StatusOK); err != nil {
			sendTextPlainResponse(w, errFailedToEncodeResponseBody, http.StatusInternalServerError)
		}
	}
}
