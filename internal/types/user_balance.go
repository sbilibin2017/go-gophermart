package types

type UserBalanceCurrentRequest string

type UserBalanceCurrentResponse struct {
	Current   float64 `json:"current"`
	Withdrawn int64   `json:"withdrawn"`
}
