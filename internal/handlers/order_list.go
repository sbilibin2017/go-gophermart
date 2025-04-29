package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

type OrderListService interface {
	List(ctx context.Context, login string) ([]*domain.Order, error)
}

func OrderListHandler(
	svc OrderListService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, ok := contextutils.GetLogin(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Fetch the orders using the service
		orders, err := svc.List(r.Context(), login)
		if err != nil {
			logger.Logger.Errorf("Error fetching orders for user %s: %v", login, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError) // 500
			return
		}

		// If no orders are found, return a No Content response
		if orders == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Prepare the response with the appropriate content type and status
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Ensure that orders is a slice of pointers, and use `json.NewEncoder(w).Encode` for output
		if err := json.NewEncoder(w).Encode(orders); err != nil {
			logger.Logger.Errorf("Error encoding response: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}
