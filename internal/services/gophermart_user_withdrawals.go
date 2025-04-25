package services

import "time"

// GophermartUserWithdrawalsRequest - структура для запроса информации о выводах средств пользователя
type GophermartUserWithdrawalsRequest struct {
	Login string `json:"login"` // Логин пользователя
}

// GophermartUserWithdrawalInfo - информация о выводе средств
type GophermartUserWithdrawalResponse struct {
	Order       string    `json:"order"`        // Номер заказа
	Sum         int64     `json:"sum"`          // Сумма списанных баллов
	ProcessedAt time.Time `json:"processed_at"` // Время обработки вывода средств (формат RFC3339)
}
