package domain

type UserBalance struct {
	Current   float64 `json:"current" db:"current"`     // Текущий баланс пользователя
	Withdrawn int64   `json:"withdrawn" db:"withdrawn"` // Сумма использованных баллов за весь период
}
