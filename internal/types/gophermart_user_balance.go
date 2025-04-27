package types

import "time"

type GophermartUserCurrentBalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn int64   `json:"withdrawn"`
}

type GophermartUserBalanceDB struct {
	Login     string    `db:"login"`
	Current   float64   `db:"current"`
	Withdrawn int64     `db:"withdrawn"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
