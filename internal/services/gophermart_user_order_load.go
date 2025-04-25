package services

// GophermartUserOrderRequest - структура для запроса на загрузку номера заказа
type GophermartUserOrderRequest struct {
	Login  string `json:"login"`  // Логин пользователя
	Number string `json:"number"` // Номер заказа
}
