package types

type UserWithdrawal struct {
	Login       string `json:"login" db:"login"`
	Number      string `json:"number" db:"number"`
	Sum         int64  `json:"withdrawn" db:"withdrawn"`
	ProcessedAt string `json:"processed_at" db:"processed_at"`
}
