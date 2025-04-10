package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/handlers/validation"
	"github.com/sbilibin2017/go-gophermart/internal/services"
)

type UserRegisterRequest struct {
	Login    string `validate:"login"`
	Password string `validate:"password"`
}

type UserRegisterResponse struct {
	AccessToken string `json:"access_token"`
}

type UserRegisterService interface {
	Register(ctx context.Context, u *domain.User) (string, error)
}

var (
	ErrLoginValidation                = errors.New("login must contain only letters and digits, minimum 3 characters")
	ErrPasswordValidation             = errors.New("password must be at least 8 characters long and include upper/lowercase letters, a number, and a special character")
	ErrUserRegisterInvalidJSONPayload = errors.New("invalid JSON payload")
)

func UserRegisterHandler(svc UserRegisterService) httprouter.Handle {
	validate := validator.New()
	validate.RegisterValidation("login", validation.ValidateLogin)
	validate.RegisterValidation("password", validation.ValidatePassword)

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var req UserRegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, ErrUserRegisterInvalidJSONPayload.Error(), http.StatusBadRequest)
			return
		}
		if err := validate.Struct(req); err != nil {
			if ve, ok := err.(validator.ValidationErrors); ok {
				for _, e := range ve {
					switch e.Field() {
					case "Login":
						http.Error(w, ErrLoginValidation.Error(), http.StatusBadRequest)
						return
					case "Password":
						http.Error(w, ErrPasswordValidation.Error(), http.StatusBadRequest)
						return
					}
				}
			}
		}
		user := &domain.User{Login: req.Login, Password: req.Password}
		token, err := svc.Register(r.Context(), user)
		if err != nil {
			switch err {
			case services.ErrUserAlreadyExists:
				http.Error(w, err.Error(), http.StatusConflict)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(UserRegisterResponse{AccessToken: token})
	}
}
