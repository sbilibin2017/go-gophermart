package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserWithdrawListService interface {
	ListWithdrawalsByLoginAndOrdered(ctx context.Context, login string) ([]*types.UserWithdrawal, error)
}

type UserWithdrawListHTTPErrorRegistry interface {
	Get(err error) *types.HTTPError
}

type UserWithdrawalsHandler struct {
	httpErrRegistry UserWithdrawListHTTPErrorRegistry
	svc             UserWithdrawListService
}

func (h *UserWithdrawalsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	login, err := contextutils.GetLogin(r.Context())
	if err != nil {
		sendTextPlainResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	withdrawals, err := h.svc.ListWithdrawalsByLoginAndOrdered(r.Context(), login)
	if handled := handleServiceError(w, err, h.httpErrRegistry); handled {
		return
	}

	if len(withdrawals) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	sendJSONResponse(w, http.StatusOK, convertAllWithdrawalsToResponse(withdrawals))
}

func convertAllWithdrawalsToResponse(withdrawals []*types.UserWithdrawal) []any {
	var response []any
	for _, withdrawal := range withdrawals {
		response = append(response, convertUserWithdrawalToEResponse(withdrawal))
	}
	return response
}

func convertUserWithdrawalToEResponse(withdrawal *types.UserWithdrawal) *struct {
	Order       string `json:"order"`
	Sum         int64  `json:"sum"`
	ProcessedAt string `json:"processed_at"`
} {
	return &struct {
		Order       string `json:"order"`
		Sum         int64  `json:"sum"`
		ProcessedAt string `json:"processed_at"`
	}{
		Order:       withdrawal.Number,
		Sum:         withdrawal.Sum,
		ProcessedAt: withdrawal.ProcessedAt,
	}
}
