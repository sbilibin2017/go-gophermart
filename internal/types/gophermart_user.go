package types

import "time"

type GophermartUserUserRegisterRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type GophermartUserUserRegisterResponse struct {
	Token string `json:"token"`
}

type GophermartUserUserLoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type GophermartUserUserLoginResponse struct {
	Token string `json:"token"`
}

type GophermartUserDB struct {
	Login     string    `db:"login"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
