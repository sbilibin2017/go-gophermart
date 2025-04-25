package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserBalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn int64   `json:"withdrawn"`
}

type UserBalanceService interface {
	GetBalance(ctx context.Context, login string) (*UserBalanceResponse, *types.APIStatus)
}

func UserBalanceHandler(svc UserBalanceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, err := middlewares.GetLogin(r)
		if err != nil {
			handleErrorResponse(w, "User not authenticated", http.StatusUnauthorized)
			return
		}
		balance, errStatus := svc.GetBalance(r.Context(), login)
		if errStatus != nil {
			handleErrorResponse(w, errStatus.Message, errStatus.Status)
			return
		}
		encodeJSONResponse(w, balance, http.StatusOK)
	}
}
