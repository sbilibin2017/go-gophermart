package services

import "time"

// GophermartUserOrdersRequest - структура для запроса на получение списка загруженных номеров заказов
type GophermartUserOrdersRequest struct {
	Login string `json:"login"` // Логин пользователя
}

// GophermartUserOrdersResponse - элемент ответа на запрос о заказах (один заказ)
type GophermartUserOrdersResponse struct {
	Number     string    `json:"number"`      // Номер заказа
	Status     string    `json:"status"`      // Статус обработки: NEW, PROCESSING, INVALID, PROCESSED
	Accrual    *int64    `json:"accrual"`     // Начисленные баллы (может быть nil, если расчёт не завершён)
	UploadedAt time.Time `json:"uploaded_at"` // Время загрузки заказа в формате RFC3339
}
