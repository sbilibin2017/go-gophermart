package types

type Withdrawn struct {
	Login  string `json:"login" db:"login"`
	Number string `json:"number" db:"number"`
	Sum    int64  `json:"withdrawn" db:"withdrawn"`
}
