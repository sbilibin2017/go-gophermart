package types

type GophermartUserRegisterRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type GophermartUserRegisterResponse struct {
	Token string `json:"token"`
}

type GophermartUserLoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type GophermartUserLoginResponse struct {
	Token string `json:"token"`
}
