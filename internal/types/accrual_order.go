package types

import "time"

// AccrualRegisterOrderGood - информация о товаре в заказе
type AccrualRegisterOrderGood struct {
	Description string  `json:"description"` // Наименование товара
	Price       float64 `json:"price"`       // Цена товара
}

// AccrualRegisterOrderRequest - модель для данных запроса при регистрации нового заказа
type AccrualRegisterOrderRequest struct {
	Order string              `json:"order"` // Номер заказа
	Goods []AccrualRegisterOrderGood `json:"goods"` // Список товаров в заказе
}

// AccrualRegisterOrderResponse - модель для данных ответа при регистрации нового заказа
type AccrualRegisterOrderResponse struct {
	StatusCode int    `json:"status_code"` // HTTP статус код (например, 200)
	Message    string `json:"message"`     // Сообщение о результате
}

// AccrualOrderAccrualRequest - модель для данных запроса при получении информации о расчёте начислений по заказу
type AccrualOrderAccrualRequest struct {
	Number string `json:"number"` // Номер заказа
}

// AccrualOrderAccrualResponse - модель для данных ответа при запросе информации о расчёте начислений
type AccrualOrderAccrualResponse struct {
	Order      string `json:"order"`             // Номер заказа
	Status     string `json:"status"`            // Статус расчёта начислений
	Accrual    *int64 `json:"accrual,omitempty"` // Начисленные баллы (может быть nil, если расчёт не завершён)
	StatusCode int    `json:"status_code"`       // HTTP статус код (например, 200)
	Message    string `json:"message"`           // Сообщение о результате
}

// AccrualOrderDB — информация о расчёте начислений баллов лояльности для конкретного заказа.
type AccrualOrderDB struct {
	Number    string    `db:"number"`     // Номер заказа
	Status    string    `db:"status"`     // Статус расчёта начислений (REGISTERED, INVALID, PROCESSING, PROCESSED)
	Accrual   *int64    `db:"accrual"`    // Начисленные баллы (может быть nil, если расчёт не завершён)
	CreatedAt time.Time `db:"created_at"` // Время регистрации записи
	UpdatedAt time.Time `db:"updated_at"` // Время последнего обновления записи
}
