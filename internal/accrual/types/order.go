package types

import "time"

// Good - информация о товаре в заказе
type RegisterOrderGood struct {
	Description string  `json:"description"` // Наименование товара
	Price       float64 `json:"price"`       // Цена товара
}

// RegisterOrderRequest - модель для данных запроса при регистрации нового заказа
type RegisterOrderRequest struct {
	Order string              `json:"order"` // Номер заказа
	Goods []RegisterOrderGood `json:"goods"` // Список товаров в заказе
}

// RegisterOrderResponse - модель для данных ответа при регистрации нового заказа
type RegisterOrderResponse struct {
	StatusCode int    `json:"status_code"` // HTTP статус код (например, 200)
	Message    string `json:"message"`     // Сообщение о результате
}

// OrderAccrualRequest - модель для данных запроса при получении информации о расчёте начислений по заказу
type OrderAccrualRequest struct {
	Number string `json:"number"` // Номер заказа
}

// OrderAccrualResponse - модель для данных ответа при запросе информации о расчёте начислений
type OrderAccrualResponse struct {
	Order      string `json:"order"`             // Номер заказа
	Status     string `json:"status"`            // Статус расчёта начислений
	Accrual    *int64 `json:"accrual,omitempty"` // Начисленные баллы (может быть nil, если расчёт не завершён)
	StatusCode int    `json:"status_code"`       // HTTP статус код (например, 200)
	Message    string `json:"message"`           // Сообщение о результате
}

// OrderAccrual — информация о расчёте начислений баллов лояльности для конкретного заказа.
type OrderDB struct {
	Number    string    `db:"number"`     // Номер заказа
	Status    string    `db:"status"`     // Статус расчёта начислений (REGISTERED, INVALID, PROCESSING, PROCESSED)
	Accrual   *int64    `db:"accrual"`    // Начисленные баллы (может быть nil, если расчёт не завершён)
	CreatedAt time.Time `db:"created_at"` // Время регистрации записи
	UpdatedAt time.Time `db:"updated_at"` // Время последнего обновления записи
}
