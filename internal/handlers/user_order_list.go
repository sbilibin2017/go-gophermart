package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserOrderListService interface {
	FilterByLoginAndSort(ctx context.Context, login string) ([]*types.UserOrder, error)
}

type UserOrderListHTTPErrorRegistry interface {
	Get(err error) *types.HTTPError
}

type UserOrderListHandler struct {
	httpErrRegistry UserOrderListHTTPErrorRegistry
	svc             UserOrderListService
}

func (h *UserOrderListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	login, err := contextutils.GetLogin(r.Context())
	if err != nil {
		sendTextPlainResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	orders, err := h.svc.FilterByLoginAndSort(r.Context(), login)
	if handleServiceError(w, err, h.httpErrRegistry) {
		return
	}

	if len(orders) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	sendJSONResponse(w, http.StatusOK, orders)
}
