package services

// GophermartUserLoginRequest - структура для запроса на аутентификацию пользователя
type GophermartUserLoginRequest struct {
	Login    string `json:"login"`    // Логин пользователя
	Password string `json:"password"` // Пароль пользователя
}
