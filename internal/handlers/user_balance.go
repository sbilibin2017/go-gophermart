package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserBalanceService interface {
	GetBalanceByLogin(ctx context.Context, login string) (*types.UserBalance, error)
}

type UserBalanceHTTPErrorRegistry interface {
	Get(err error) *types.HTTPError
}

type UserBalanceHandler struct {
	svc             UserBalanceService
	httpErrRegistry UserBalanceHTTPErrorRegistry
}

func (h *UserBalanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	login, err := contextutils.GetLogin(r.Context())
	if err != nil {
		sendTextPlainResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	balance, err := h.svc.GetBalanceByLogin(r.Context(), login)
	if handleServiceError(w, err, h.httpErrRegistry) {
		return
	}

	sendJSONResponse(w, http.StatusOK, balance)
}
