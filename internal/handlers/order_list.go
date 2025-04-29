package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
)

type OrderListService interface {
	List(ctx context.Context, login string) ([]*domain.Order, error)
}

func OrderListHandler(
	svc OrderListService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, ok := middlewares.GetLoginFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		orders, err := svc.List(r.Context(), login)
		if err != nil {

			http.Error(w, "Internal server error", http.StatusInternalServerError) // 500
			return
		}

		if orders == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(orders); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}
