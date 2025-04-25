package types

import "time"

// GophermartUserBalanceWithdrawRequest - структура для запроса на списание баллов с баланса
type GophermartUserBalanceWithdrawRequest struct {
	Order string `json:"order"` // Номер заказа, для которого списываются баллы
	Sum   int64  `json:"sum"`   // Сумма баллов, которую нужно списать
}

// GophermartUserBalanceWithdrawResponse - структура для ответа на запрос списания баллов
type GophermartUserBalanceWithdrawResponse struct {
	StatusCode int    `json:"status_code"` // HTTP статус код (например, 200)
	Message    string `json:"message"`     // Сообщение о результате
}

// GophermartUserWithdrawalsRequest - структура для запроса информации о выводах средств пользователя
type GophermartUserWithdrawalsRequest struct {
	Login string `json:"login"` // Логин пользователя
}

// GophermartUserWithdrawalsResponse - структура для ответа с информацией о выводах средств
type GophermartUserWithdrawalsResponse struct {
	StatusCode  int                        `json:"status_code"` // HTTP статус код (например, 200)
	Message     string                     `json:"message"`     // Сообщение о результате
	Withdrawals []GophermartUserWithdrawal `json:"withdrawals"` // Список выводов средств
}

// GophermartUserWithdrawalInfo - информация о выводе средств
type GophermartUserWithdrawal struct {
	Order       string    `json:"order"`        // Номер заказа
	Sum         int64     `json:"sum"`          // Сумма списанных баллов
	ProcessedAt time.Time `json:"processed_at"` // Время обработки вывода средств (формат RFC3339)
}

// GophermartUserBalanceWithdrawalDB - информация о выводе средств с накопительного счёта (для базы данных)
type GophermartUserBalanceWithdrawalDB struct {
	Login     string    `db:"login"`        // Уникальный идентификатор пользователя
	Number    string    `db:"number"`       // Номер заказа, на который были списаны баллы
	Sum       int64     `db:"sum"`          // Сумма списанных баллов
	Processed time.Time `db:"processed_at"` // Время списания баллов
	CreatedAt time.Time `db:"created_at"`   // Время регистрации записи
	UpdatedAt time.Time `db:"updated_at"`   // Время последнего обновления записи
}
