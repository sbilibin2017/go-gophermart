package services

// GophermartUserRegisterRequest - структура для запроса на регистрацию пользователя
type GophermartUserRegisterRequest struct {
	Login    string `json:"login"`    // Логин пользователя
	Password string `json:"password"` // Пароль пользователя
}
