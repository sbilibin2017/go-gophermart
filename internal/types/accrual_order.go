package types

import "time"

type AccrualOrderGetRequest struct {
	Order string `json:"order"`
}

type AccrualOrderGetResponse struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual int64  `json:"accrual"`
}

type AccrualOrderRegisterRequest struct {
	Order string `json:"order"`
	Goods []struct {
		Description string `json:"description"`
		Price       int64  `json:"price"`
	} `json:"goods"`
}

type AccrualOrderDB struct {
	Number    string    `db:"number"`
	Status    string    `db:"status"`
	Accrual   int64     `db:"accrual"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
