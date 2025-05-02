package types

type UserBalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn int64   `json:"withdrawn"`
}

type UserBalanceDB struct {
	Login     string  `db:"login"`
	Current   float64 `db:"current"`
	Withdrawn int64   `db:"withdrawn"`
}
