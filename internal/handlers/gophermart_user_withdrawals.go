package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserWithdrawalsResponse struct {
	Order       string    `json:"order"`
	Sum         int64     `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type UserWithdrawalsService interface {
	List(ctx context.Context, login string) ([]*UserWithdrawalsResponse, *types.APIStatus)
}

func UserWithdrawalsHandler(svc UserWithdrawalsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, err := middlewares.GetLogin(r)
		if err != nil {
			handleErrorResponse(w, "User not authenticated", http.StatusUnauthorized)
			return
		}
		withdrawals, errStatus := svc.List(r.Context(), login)
		if errStatus != nil {
			handleErrorResponse(w, errStatus.Message, errStatus.Status)
			return
		}
		if len(withdrawals) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		encodeJSONResponse(w, withdrawals, http.StatusOK)
	}
}
