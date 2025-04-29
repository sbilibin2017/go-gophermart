package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/domain"
)

type UserRegisterService interface {
	Register(ctx context.Context, user *domain.User) (*string, error)
}

type UserRegisterValidator interface {
	Struct(v any) error
}

func UserRegisterHandler(
	val UserRegisterValidator,
	svc UserRegisterService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData struct {
			Login    string `json:"login" validate:"required"`
			Password string `json:"password" validate:"required"`
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&requestData); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if err := val.Struct(requestData); err != nil {
			http.Error(w, "validation failed", http.StatusBadRequest)
			return
		}

		user := &domain.User{
			Login:    requestData.Login,
			Password: requestData.Password,
		}

		token, err := svc.Register(r.Context(), user)
		if err != nil {
			if errors.Is(err, domain.ErrLoginAlreadyTaken) {
				http.Error(w, "login already taken", http.StatusConflict)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Authorization", "Bearer "+*token)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User successfully registered"))
	}
}
