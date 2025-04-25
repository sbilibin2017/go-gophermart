package services

// GophermartUserBalanceRequest - структура для запроса на получение текущего баланса пользователя
type GophermartUserBalanceRequest struct {
	Login string `json:"login"` // Логин пользователя
}

// GophermartUserBalanceResponse - структура для ответа на запрос текущего баланса
type GophermartUserBalanceResponse struct {
	Current   float64 `json:"current"`   // Текущий баланс баллов лояльности
	Withdrawn int64   `json:"withdrawn"` // Сумма использованных баллов
}
