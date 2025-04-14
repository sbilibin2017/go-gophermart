package types

// Order представляет структуру заказа.
// Он содержит уникальный номер заказа и список товаров, включенных в этот заказ.
type Order struct {
	Number uint64 // Номер заказа, уникальный идентификатор заказа
	Goods  []Good // Список товаров, включенных в данный заказ
}

// Good описывает товар в заказе.
// Включает описание товара и его цену.
type Good struct {
	Description string // Описание товара, например, название или характеристика
	Price       uint64 // Цена товара в наименьшей денежной единице (например, в копейках)
}

// OrderRegisterRequest представляет запрос на регистрацию заказа.
// Он включает номер заказа и список товаров, которые должны быть зарегистрированы.
type OrderRegisterRequest struct {
	Number uint64 `json:"order" validate:"required"`      // Номер заказа, обязательное поле
	Goods  []Good `json:"goods" validate:"dive,required"` // Список товаров в заказе, обязательное поле
}

// OrderRegisterResponse представляет ответ на запрос регистрации заказа.
// Он содержит сообщение о результате регистрации.
type OrderRegisterResponse struct {
	Message string `json:"message"` // Сообщение, которое возвращается в ответ на запрос
}

type OrderExistsFilter struct {
	Number uint64 `db:"number"`
}

type Status string

const (
	StatusNew        Status = "NEW"        // заказ загружен, но не попал в обработку
	StatusProcessing Status = "PROCESSING" // вознаграждение рассчитывается
	StatusInvalid    Status = "INVALID"    // система отказала в расчёте
	StatusProcessed  Status = "PROCESSED"  // расчёт успешно завершён
)

type OrderDB struct {
	Number  uint64  `db:"number"`
	Status  Status  `db:"status"`
	Accrual float64 `db:"accrual"`
}
