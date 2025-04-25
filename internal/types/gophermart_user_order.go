package types

import "time"

// GophermartUserOrderRequest - структура для запроса на загрузку номера заказа
type GophermartUserOrderRequest struct {
	Login  string `json:"login"`  // Логин пользователя
	Number string `json:"number"` // Номер заказа
}

// GophermartUserOrderResponse - структура для ответа на запрос загрузки номера заказа
type GophermartUserOrderResponse struct {
	StatusCode int    `json:"status_code"` // HTTP статус код (например, 200)
	Message    string `json:"message"`     // Сообщение о результате
}

// GophermartUserOrdersRequest - структура для запроса на получение списка загруженных номеров заказов
type GophermartUserOrdersRequest struct {
	Login string `json:"login"` // Логин пользователя
}

// GophermartUserOrdersItem - элемент ответа на запрос о заказах (один заказ)
type GophermartUserOrdersItem struct {
	Number     string    `json:"number"`      // Номер заказа
	Status     string    `json:"status"`      // Статус обработки: NEW, PROCESSING, INVALID, PROCESSED
	Accrual    *int64    `json:"accrual"`     // Начисленные баллы (может быть nil, если расчёт не завершён)
	UploadedAt time.Time `json:"uploaded_at"` // Время загрузки заказа в формате RFC3339
}

// GophermartUserOrdersResponse - структура для ответа на запрос списка заказов
type GophermartUserOrdersResponse struct {
	StatusCode int                        `json:"status_code"` // HTTP статус код (например, 200)
	Message    string                     `json:"message"`     // Сообщение о результате
	Orders     []GophermartUserOrdersItem `json:"orders"`      // Список заказов
}

// GophermartUserOrder — информация о заказах пользователя, загруженных для расчёта начислений.
type GophermartUserOrderDB struct {
	Login      string    `db:"user_id"`      // Уникальный идентификатор пользователя
	Number     string    `db:"order_number"` // Номер заказа (уникален для пользователя)
	Status     string    `db:"status"`       // Статус расчёта начислений: REGISTERED, PROCESSING, INVALID, PROCESSED
	Accrual    *int64    `db:"accrual"`      // Начисленные баллы (может быть nil, если расчёт не завершён)
	UploadedAt time.Time `db:"uploaded_at"`  // Время загрузки заказа
	CreatedAt  time.Time `db:"created_at"`   // Время регистрации записи
	UpdatedAt  time.Time `db:"updated_at"`   // Время последнего обновления записи
}
