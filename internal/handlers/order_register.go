package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderRegisterService interface {
	Register(ctx context.Context, order *types.OrderRequest) error
}

type OrderRegisterValidator interface {
	Struct(v any) error
}

func OrderRegisterHandler(
	val OrderRegisterValidator,
	svc OrderRegisterService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrderRequest

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, errInvalidJSONFormat, http.StatusBadRequest)
			return
		}

		if err := val.Struct(&req); err != nil {
			valErr := formatValidationError(err)
			http.Error(w, valErr.Message, http.StatusBadRequest)
			return
		}

		err := svc.Register(r.Context(), &req)
		if err != nil {
			switch {
			case errors.Is(err, types.ErrOrderAlreadyExists):
				http.Error(w, capitalize(err.Error()), http.StatusConflict)
			default:
				http.Error(w, capitalize(types.ErrInternal.Error()), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Order successfully received"))
	}
}
