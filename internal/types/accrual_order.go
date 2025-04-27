package types

type AccrualOrderGetRequest struct {
	Order string `json:"order" validate:"required,luhn"`
}

type AccrualOrderGetResponse struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual int64  `json:"accrual"`
}

type AccrualOrderRegisterRequest struct {
	Order string `json:"order" validate:"required,luhn"`
	Goods []struct {
		Description string `json:"description" validate:"required"`
		Price       int64  `json:"price" validate:"required,gt=0"`
	} `json:"goods" validate:"required,min=1"`
}
