package types

type OrderAcceptRequest struct {
	Order string `json:"order" validate:"required,luhn"`
	Goods []struct {
		Description string `json:"description" validate:"required"`
		Price       int64  `json:"price" validate:"required,gt=0"`
	} `json:"goods" validate:"required,dive,required"`
}
