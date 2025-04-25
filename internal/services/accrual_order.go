package services

// AccrualOrderGetRequest - модель для данных запроса при получении информации о расчёте начислений по заказу
type AccrualOrderGetRequest struct {
	Number string `json:"number"` // Номер заказа
}

// AccrualOrderGetResponse - модель для данных ответа при запросе информации о расчёте начислений
type AccrualOrderGetResponse struct {
	Order   string `json:"order"`             // Номер заказа
	Status  string `json:"status"`            // Статус расчёта начислений
	Accrual *int64 `json:"accrual,omitempty"` // Начисленные баллы (может быть nil, если расчёт не завершён)
}
