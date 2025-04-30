package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type GoodRewardRegisterService interface {
	Register(ctx context.Context, goodReward *types.GoodReward) error
}

type GoodRewardRegisterValidator interface {
	Struct(v any) error
}

func GoodRewardRegisterHandler(
	val GoodRewardRegisterValidator,
	svc GoodRewardRegisterService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GoodReward

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, errInternalServerError, http.StatusBadRequest)
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
			case errors.Is(err, types.ErrGoodRewardAlreadyExists):
				http.Error(w, capitalize(err.Error()), http.StatusConflict)
			default:
				http.Error(w, errInternalServerError, http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Good reward successfully registered"))
	}
}
