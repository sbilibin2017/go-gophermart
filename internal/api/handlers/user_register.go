package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/engines/json"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/usecases"
	"github.com/sbilibin2017/go-gophermart/internal/usecases/validators"
)

type UserRegisterUsecase interface {
	Execute(ctx context.Context, req *usecases.UserRegisterRequest) (*usecases.UserRegisterResponse, error)
}

type Decoder interface {
	Decode(r *http.Request, v any) error
}

func UserRegisterHandler(
	uc UserRegisterUsecase,
	decoder Decoder,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req usecases.UserRegisterRequest

		if err := decoder.Decode(r, &req); err != nil {
			http.Error(w, json.ErrRequestDecoderUnprocessableJson.Error(), http.StatusBadRequest)
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
	case services.ErrUserAlreadyExists:
		http.Error(w, err.Error(), http.StatusConflict)
	case validators.ErrInvalidLogin, validators.ErrInvalidPassword:
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, errors.New("internal error").Error(), http.StatusInternalServerError)
	}
}
