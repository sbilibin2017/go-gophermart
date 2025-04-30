package types

type Order struct {
	Number string `json:"number" db:"number"`
	Goods  []Good `json:"goods"`
}

type Good struct {
	Description string `json:"description" db:"description"`
	Price       int64  `json:"price" db:"price"`
}
