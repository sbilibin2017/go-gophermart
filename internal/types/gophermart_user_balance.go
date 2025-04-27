package types

type GophermartUserCurrentBalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn int64   `json:"withdrawn"`
}
