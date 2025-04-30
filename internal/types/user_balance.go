package types

type UserBalance struct {
	Login     string  `json:"login" db:"login"`
	Current   float64 `json:"current" db:"current"`
	Withdrawn int64   `json:"withdrawn" db:"withdrawn"`
}
