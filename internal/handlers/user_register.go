package handlers

import (
	"context"

	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/errors"
	"github.com/sbilibin2017/go-gophermart/internal/usecases"
)

type UserRegisterUsecase interface {
	Execute(ctx context.Context, req *usecases.UserRegisterRequest) (*usecases.UserRegisterResponse, error)
}

type Decoder interface {
	Decode(v any) error
}

func UserRegisterHandler(uc UserRegisterUsecase, d Decoder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req usecases.UserRegisterRequest
		if err := d.Decode(&req); err != nil {
			http.Error(w, errors.ErrUnprocessableJson.Error(), http.StatusBadRequest)
			return
		}
		resp, err := uc.Execute(r.Context(), &req)
		if err != nil {
			handleUserRegisterError(w, err)
			return
		}
		w.Header().Set("Authorization", "Bearer "+resp.AccessToken)
		w.WriteHeader(http.StatusOK)
	}
}

func handleUserRegisterError(w http.ResponseWriter, err error) {
	switch err {
	case errors.ErrUserAlreadyExists:
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.ErrInvalidLogin:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.ErrInvalidPassword:
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, errors.ErrInternal.Error(), http.StatusInternalServerError)
	}
}
