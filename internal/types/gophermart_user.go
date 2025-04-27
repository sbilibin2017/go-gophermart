package types

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
