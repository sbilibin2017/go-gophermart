package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderGetService interface {
	Get(ctx context.Context, number string) (*types.Order, error)
}

type OrderGetValidator interface {
	Struct(v any) error
}

func OrderGetHandler(
	val OrderGetValidator,
	svc OrderGetService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		number := chi.URLParam(r, "number")

		req := struct {
			Number string `json:"number" validate:"required,luhn"`
		}{
			Number: number,
		}

		if err := val.Struct(&req); err != nil {
			valErr := formatValidationError(err)
			http.Error(w, valErr.Message, http.StatusBadRequest)
			return
		}

		order, err := svc.Get(r.Context(), req.Number)
		if err != nil {
			if errors.Is(err, types.ErrOrderNotFound) {
				http.Error(w, capitalize(err.Error()), http.StatusNoContent)
				return
			}
			http.Error(w, capitalize(types.ErrInternal.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(order)
	}
}
