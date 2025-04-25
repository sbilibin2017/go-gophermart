package types

import "time"

// GophermartUserRegisterRequest - структура для запроса на регистрацию пользователя
type GophermartUserRegisterRequest struct {
	Login    string `json:"login"`    // Логин пользователя
	Password string `json:"password"` // Пароль пользователя
}

// GophermartUserRegisterResponse - структура для ответа на регистрацию пользователя
type GophermartUserRegisterResponse struct {
	StatusCode int    `json:"status_code"` // HTTP статус код (например, 200)
	Message    string `json:"message"`     // Сообщение о результате
}

// GophermartUserLoginRequest - структура для запроса на аутентификацию пользователя
type GophermartUserLoginRequest struct {
	Login    string `json:"login"`    // Логин пользователя
	Password string `json:"password"` // Пароль пользователя
}

// GophermartUserLoginResponse - структура для ответа на запрос аутентификации пользователя
type GophermartUserLoginResponse struct {
	StatusCode int    `json:"status_code"` // HTTP статус код (например, 200)
	Message    string `json:"message"`     // Сообщение о результате
}

// GophermartUserDB — информация о пользователе, зарегистрированном в системе.
type GophermartUserDB struct {
	Login     string    `db:"login"`      // Уникальный логин пользователя
	Password  string    `db:"password"`   // Хэшированный пароль пользователя
	CreatedAt time.Time `db:"created_at"` // Время регистрации записи
	UpdatedAt time.Time `db:"updated_at"` // Время последнего обновления записи
}
