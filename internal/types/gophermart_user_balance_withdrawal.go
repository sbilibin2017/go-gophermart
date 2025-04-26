package types

import "time"

// GophermartUserBalanceWithdrawal - модель для таблицы gophermart_user_balance_withdrawal
type GophermartUserBalanceWithdrawal struct {
	Login       string    `db:"login"`        // references gophermart_user(login)
	Number      string    `db:"number"`       // primary key part 1
	Sum         int64     `db:"sum"`          // not null
	ProcessedAt time.Time `db:"processed_at"` // not null
	CreatedAt   time.Time `db:"created_at"`   // not null
	UpdatedAt   time.Time `db:"updated_at"`   // not null
}
