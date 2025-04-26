package types

import "time"

type GophermartUserBalance struct {
	Login     string    `db:"login"`
	Current   float64   `db:"current"`
	Withdrawn int64     `db:"withdrawn"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
