package services

// AccrualRegisterOrderGood - информация о товаре в заказе
type AccrualRegisterOrderGood struct {
	Description string  `json:"description"` // Наименование товара
	Price       float64 `json:"price"`       // Цена товара
}

// AccrualRegisterOrderRequest - модель для данных запроса при регистрации нового заказа
type AccrualRegisterOrderRequest struct {
	Order string                     `json:"order"` // Номер заказа
	Goods []AccrualRegisterOrderGood `json:"goods"` // Список товаров в заказе
}
