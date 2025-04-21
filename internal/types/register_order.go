package types

type RegisterOrderRequest struct {
	Order string              `json:"order" validate:"required,luhn"`
	Goods []RegisterOrderGood `json:"goods" validate:"required"`
}

type RegisterOrderGood struct {
	Description string `json:"description" validate:"required"`
	Price       int64  `json:"price" validate:"required,gt=0"`
}
