package handlers

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/handlers/utils"
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
	Register(ctx context.Context, u *domain.User) (*domain.UserToken, error)
}

func UserRegisterHandler(svc UserRegisterService) httprouter.Handle {
	validate := validator.New()
	validate.RegisterValidation("login", validation.ValidateLogin)
	validate.RegisterValidation("password", validation.ValidatePassword)
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var req UserRegisterRequest
		if err := utils.DecodeJSON(w, r, &req); err != nil {
			return
		}
		err := validate.Struct(req)
		if err != nil {
			utils.HandleValidationError(w, err)
			return
		}
		user := &domain.User{Login: req.Login, Password: req.Password}
		token, err := svc.Register(r.Context(), user)
		if err != nil {
			handleUserRegisterError(w, err)
			return
		}
		resp := &UserRegisterResponse{
			AccessToken: token.Access,
		}
		utils.EncodeJSON(w, http.StatusOK, resp)
	}
}

func handleUserRegisterError(w http.ResponseWriter, err error) {
	switch err {
	case services.ErrUserAlreadyExists:
		http.Error(w, err.Error(), http.StatusConflict)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
