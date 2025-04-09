package requests

import (
	"github.com/go-playground/validator/v10"
	"github.com/sbilibin2017/go-gophermart/internal/validators"
)

type UserRegisterRequest struct {
	Login    string `validate:"login"`
	Password string `validate:"password"`
}

func (u *UserRegisterRequest) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("login", validators.ValidateLogin)
	validate.RegisterValidation("password", validators.ValidatePassword)
	return validate.Struct(u)
}
