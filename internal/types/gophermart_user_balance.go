package types

import "time"

// GophermartUserBalanceRequest - структура для запроса на получение текущего баланса пользователя
type GophermartUserBalanceRequest struct {
	Login string `json:"login"` // Логин пользователя
}

// GophermartUserBalanceResponse - структура для ответа на запрос текущего баланса
type GophermartUserBalanceResponse struct {
	StatusCode int     `json:"status_code"` // HTTP статус код (например, 200)
	Message    string  `json:"message"`     // Сообщение о результате
	Current    float64 `json:"current"`     // Текущий баланс баллов лояльности
	Withdrawn  int64   `json:"withdrawn"`   // Сумма использованных баллов
}

// GophermartUserBalanceDB - информация о балансе пользователя
type GophermartUserBalanceDB struct {
	Login     string    `db:"login"`      // Уникальный идентификатор пользователя
	Current   float64   `db:"current"`    // Текущий баланс баллов лояльности
	Withdrawn int64     `db:"withdrawn"`  // Сумма использованных баллов
	CreatedAt time.Time `db:"created_at"` // Время регистрации записи
	UpdatedAt time.Time `db:"updated_at"` // Время последнего обновления записи
}
