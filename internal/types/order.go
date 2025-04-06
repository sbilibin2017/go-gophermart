package types

import "time"

type Order struct {
	Order     string    `json:"order"`
	Status    string    `json:"status"`
	Goods     []Good    `json:"goods"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
